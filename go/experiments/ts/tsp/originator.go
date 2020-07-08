package tsp

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/topology"
	"github.com/scionproto/scion/go/proto"
)

var originatorLog = log.New(ioutil.Discard, "[tsp/originator] ", log.LstdFlags)

func newPacket(localIA addr.IA, remoteIA addr.IA, path snet.Path) *snet.Packet {
	var dpPath *spath.Path
	if path != nil {
		dpPath = path.Path()
	}
	return &snet.Packet{
		PacketInfo: snet.PacketInfo{
			Destination: snet.SCIONAddress{
				IA: remoteIA,
				Host: addr.SvcTS | addr.SVCMcast,
			},
			Path: dpPath,
			Payload: common.RawBytes(
				[]byte(localIA.String() + " > " + remoteIA.String()),
			),
		},
	}
}

func StartOriginator(c sciond.Connector, ctx context.Context) error {
	localIA, err := c.LocalIA(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			originatorLog.Printf("Initiating TSP broadcast\n")

			corePaths, err := c.Paths(ctx,
				addr.IA{I: 0, A: 0}, localIA, sciond.PathReqFlags{Refresh: true})
			if err != nil {
				originatorLog.Printf("Failed to lookup core paths: %v\n", err)
			}
			coreASes := make(map[addr.IA][]snet.Path)
			if corePaths != nil {
				for _, p := range corePaths {
					coreASes[p.Destination()] = append(coreASes[p.Destination()], p)
				}
			}

			originatorLog.Printf("Reachable core ASes:\n")
			for coreAS := range coreASes {
				originatorLog.Printf("%v", coreAS)
				for _, p := range coreASes[coreAS] {
					originatorLog.Printf("\t%v\n", p)
				}
			}

			svcInfoReply, err := c.SVCInfo(ctx, []proto.ServiceType{proto.ServiceType_ts})
			if err != nil {
				originatorLog.Printf("Failed to lookup local TS service info: %v\n", err)
			}
			localTSHosts := make(map[string]bool)
			if svcInfoReply != nil {
				for _, i := range svcInfoReply.Entries[0].HostInfos {
					localTSHosts[i.Host().IP().String()] = true
				}
			}

			originatorLog.Printf("Reachable local time services:\n")
			for localTSHost := range localTSHosts {
				originatorLog.Printf("%v", localTSHost)
			}

			for remoteIA := range coreASes {
				ps := coreASes[remoteIA]
				sp := ps[rand.Intn(len(ps))]
				propagatePacketTo(
					newPacket(localIA, remoteIA, sp),
					sp.UnderlayNextHop())
			}
			for localTSHost := range localTSHosts {
				propagatePacketTo(
					newPacket(localIA, localIA, /* path: */ nil),
					&net.UDPAddr{IP: net.ParseIP(localTSHost), Port: topology.EndhostPort})
			}

			time.Sleep(15 * time.Second)
		}
	}()

	return nil
}
