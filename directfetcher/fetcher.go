package directfetcher

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/resolver"

	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/grpc"
	"github.com/scionproto/scion/pkg/private/serrors"
	seg "github.com/scionproto/scion/pkg/segment"
	"github.com/scionproto/scion/pkg/snet"
	"github.com/scionproto/scion/pkg/snet/path"
	"github.com/scionproto/scion/private/path/combinator"
	"github.com/scionproto/scion/private/segment/segfetcher"
	grpcfetcher "github.com/scionproto/scion/private/segment/segfetcher/grpc"
	"github.com/scionproto/scion/private/topology"
)

type dstProvider struct{}

func (d *dstProvider) Dst(context.Context, segfetcher.Request) (net.Addr, error) {
	return &snet.SVCAddr{SVC: addr.SvcCS}, nil
}

type Fetcher struct {
	topo      *topology.Loader
	requester *segfetcher.DefaultRequester
}

func New(topo *topology.Loader) *Fetcher {
	dialer := &grpc.TCPDialer{
		SvcResolver: func(dst addr.SVC) []resolver.Address {
			if base := dst.Base(); base != addr.SvcCS {
				panic("unexpected address type")
			}
			targets := []resolver.Address{}
			for _, entry := range topo.ControlServiceAddresses() {
				targets = append(targets, resolver.Address{Addr: entry.String()})
			}
			return targets
		},
	}
	return &Fetcher{
		topo: topo,
		requester: &segfetcher.DefaultRequester{
			RPC: &grpcfetcher.Requester{
				Dialer: dialer,
			},
			DstProvider: &dstProvider{},
		},
	}
}

func (f *Fetcher) FetchPaths(ctx context.Context, dst addr.IA) ([]snet.Path, error) {
	src := f.topo.IA()
	srcCore := f.topo.Core()
	dstCore, err := f.fetchASType(ctx, dst)
	if err != nil {
		return nil, err
	}

	requests := f.createRequests(src, dst, srcCore, dstCore)
	ups, cores, downs, err := f.fetchSegments(ctx, requests)
	if err != nil {
		return nil, err
	}

	paths := combinator.Combine(src, dst, ups, cores, downs, false /* findAllIdentical */)
	return f.translatePaths(paths)
}

func toWildcard(ia addr.IA) addr.IA {
	return addr.MustParseIA(fmt.Sprintf("%d-0", ia.ISD()))
}

func (f *Fetcher) fetchASType(ctx context.Context, dst addr.IA) (bool, error) {
	if dst.IsWildcard() {
		return true, nil
	}

	src := f.topo.IA()
	singleISD := src.ISD() == dst.ISD()

	var req segfetcher.Request
	if singleISD {
		req = segfetcher.Request{
			Src: toWildcard(src), Dst: dst, SegType: seg.TypeDown}
	} else {
		req = segfetcher.Request{
			Src: toWildcard(src), Dst: dst, SegType: seg.TypeCore}
	}

	replies := f.requester.Request(ctx, segfetcher.Requests{req})
	for reply := range replies {
		if singleISD {
			if reply.Err != nil {
				return false, reply.Err
			}
			if len(reply.Segments) > 0 {
				return false, nil
			} else {
				return true, nil
			}
		} else {
			if reply.Err != nil {
				return false, nil
			} else {
				return true, nil
			}
		}
	}
	return false, serrors.New("failed to fetch AS type")
}

func (f *Fetcher) createRequests(src, dst addr.IA, srcCore, dstCore bool) segfetcher.Requests {
	switch {
	case !srcCore && !dstCore:
		return segfetcher.Requests{
			{Src: src, Dst: toWildcard(src), SegType: seg.TypeUp},
			{Src: toWildcard(src), Dst: toWildcard(dst), SegType: seg.TypeCore},
			{Src: toWildcard(dst), Dst: dst, SegType: seg.TypeDown},
		}
	case !srcCore && dstCore:
		return segfetcher.Requests{
			{Src: src, Dst: toWildcard(src), SegType: seg.TypeUp},
			{Src: toWildcard(src), Dst: dst, SegType: seg.TypeCore},
		}
	case srcCore && !dstCore:
		return segfetcher.Requests{
			{Src: src, Dst: toWildcard(dst), SegType: seg.TypeCore},
			{Src: toWildcard(dst), Dst: dst, SegType: seg.TypeDown},
		}
	default:
		return segfetcher.Requests{{Src: src, Dst: dst, SegType: seg.TypeCore}}
	}
}

func (f *Fetcher) fetchSegments(ctx context.Context, requests segfetcher.Requests) (
	ups, cores, downs []*seg.PathSegment, err error) {

	if len(requests) == 0 {
		return nil, nil, nil, nil
	}

	replies := f.requester.Request(ctx, requests)
	for reply := range replies {
		if reply.Err != nil {
			return nil, nil, nil, reply.Err
		}
		for _, segMeta := range reply.Segments {
			switch reply.Req.SegType {
			case seg.TypeUp:
				ups = append(ups, segMeta.Segment)
			case seg.TypeCore:
				cores = append(cores, segMeta.Segment)
			case seg.TypeDown:
				downs = append(downs, segMeta.Segment)
			}
		}
	}

	return ups, cores, downs, nil
}

func (f *Fetcher) translatePaths(cpaths []combinator.Path) ([]snet.Path, error) {
	var paths []snet.Path
	for _, cpath := range cpaths {
		path, err := f.translatePath(cpath)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	if len(paths) == 0 {
		return nil, serrors.New("no paths after translation")
	}
	return paths, nil
}

func (f *Fetcher) translatePath(cpath combinator.Path) (snet.Path, error) {
	if len(cpath.Metadata.Interfaces) == 0 {
		return nil, serrors.New("path has no interfaces")
	}

	firstIF := cpath.Metadata.Interfaces[0]

	nextHop := f.topo.UnderlayNextHop(uint16(firstIF.ID))
	if nextHop == nil {
		return nil, serrors.New("unable to find next hop for first interface", "ifID", firstIF.ID)
	}

	return path.Path{
		Src:           cpath.Metadata.Interfaces[0].IA,
		Dst:           cpath.Metadata.Interfaces[len(cpath.Metadata.Interfaces)-1].IA,
		DataplanePath: cpath.SCIONPath,
		NextHop:       nextHop,
		Meta:          cpath.Metadata,
	}, nil
}
