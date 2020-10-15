package tsp

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
)

type PathInfo struct {
	LocalIA addr.IA
	CoreASes map[addr.IA][]snet.Path
	LocalTSHosts map[string]net.IP
}

var patherLog = log.New(ioutil.Discard, "[tsp/pather] ", log.LstdFlags)

func StartPather(c sciond.Connector, ctx context.Context) (<-chan PathInfo, error) {
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
					addr.IA{I: 0, A: 0}, localIA, sciond.PathReqFlags{Refresh: true})
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

				// svcInfoReply, err := c.SVCInfo(ctx,
				// 	[]proto.ServiceType{proto.ServiceType_ts})
				// if err != nil {
				// 	patherLog.Printf("Failed to lookup local TS service info: %v\n", err)
				// }
				localTSHosts := make(map[string]net.IP)
				// if svcInfoReply != nil {
				// 	for _, i := range svcInfoReply.Entries[0].HostInfos {
				// 		ip := i.Host().IP()
				// 		localTSHosts[ip.String()] = ip
				// 	}
				// }

				patherLog.Printf("Reachable local time services:\n")
				for localTSHost := range localTSHosts {
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
