package tsp

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/daemon"
	"github.com/scionproto/scion/go/lib/snet"
)

type PathInfo struct {
	LocalIA addr.IA
	CoreASes map[addr.IA][]snet.Path
	LocalTSHosts []*net.UDPAddr
}

var patherLog = log.New(ioutil.Discard, "[tsp/pather] ", log.LstdFlags)

func StartPather(c daemon.Connector, ctx context.Context) (<-chan PathInfo, error) {
	localIA, err := c.LocalIA(ctx)
	if err != nil {
		return nil, err
	}

	pathInfos := make(chan PathInfo)

	go func() {
		ticker := time.NewTicker(15 * time.Second)

		for {
			select {
			case <-ticker.C:
				patherLog.Printf("Looking up TSP broadcast paths\n")

				corePaths, err := c.Paths(ctx,
					addr.IA{I: 0, A: 0}, localIA, daemon.PathReqFlags{Refresh: true})
				if err != nil {
					patherLog.Printf("Failed to lookup core paths: %v\n", err)
				}
				coreASes := make(map[addr.IA][]snet.Path)
				if corePaths != nil {
					for _, p := range corePaths {
						coreASes[p.Destination()] = append(coreASes[p.Destination()], p)
					}
				}

				patherLog.Printf("Reachable core ASes:\n")
				for coreAS := range coreASes {
					patherLog.Printf("%v", coreAS)
					for _, p := range coreASes[coreAS] {
						patherLog.Printf("\t%v\n", p)
					}
				}

				localSvcInfo, err := c.SVCInfo(ctx, []addr.HostSVC{addr.SvcTS})
				if err != nil {
					patherLog.Printf("Failed to lookup local TS service info: %v\n", err)
				}
				var localTSHosts []*net.UDPAddr
				localSvcAddr, ok := localSvcInfo[addr.SvcTS]
				if ok {
					localSvcUdpAddr, err := net.ResolveUDPAddr("udp", localSvcAddr)
					if err != nil {
						patherLog.Printf("Failed to resolve local TS service addr: %v\n", err)
					}
					localTSHosts = append(localTSHosts, localSvcUdpAddr)
				}

				patherLog.Printf("Reachable local time services:\n")
				for _, localTSHost := range localTSHosts {
					patherLog.Printf("\t%v\n", localTSHost)
				}

				pathInfos <- PathInfo{
					LocalIA: localIA,
					CoreASes: coreASes,
					LocalTSHosts: localTSHosts,
				}
			}
		}
	}()

	return pathInfos, nil
}
