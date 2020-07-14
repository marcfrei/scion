package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net"
	"sort"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/topology"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"

	"github.com/scionproto/scion/go/experiments/ts/ets"
	"github.com/scionproto/scion/go/experiments/ts/tsp"
)

func newSciondConnector(addr string, ctx context.Context) sciond.Connector {
	c, err := sciond.NewService(addr).Connect(ctx)
	if err != nil {
		log.Fatal("Failed to create SCION connector:", err)
	}
	return c
}

func newPacketDispatcher(c sciond.Connector) snet.PacketDispatcherService {
	return &snet.DefaultPacketDispatcherService{
		Dispatcher: reconnect.NewDispatcherService(reliable.NewDispatcher("")),
		SCMPHandler: snet.NewSCMPHandler(sciond.RevHandler{Connector: c}),
	}
}

func newPacket(localIA addr.IA, remoteIA addr.IA, path snet.Path,
	payload time.Duration) *snet.Packet {
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
				[]byte(payload.String()),
			),
		},
	}
}

func main() {
	var sciondAddr string
	var localAddr snet.UDPAddr
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Parse()

	ctx := context.Background()
	var err error

	pathInfos, err := tsp.StartPather(newSciondConnector(sciondAddr, ctx), ctx)
	if err != nil {
		log.Fatal("Failed to start TSP originator:", err)
	}
	var pathInfo tsp.PathInfo

	syncInfos, err := tsp.StartHandler(
		newPacketDispatcher(newSciondConnector(sciondAddr, ctx)), ctx,
		localAddr.IA, localAddr.Host)
	if err != nil {
		log.Fatal("Failed to start TSP handler:", err)
	}

	localAddr.Host.Port = 0
	err = tsp.StartPropagator(
		newPacketDispatcher(newSciondConnector(sciondAddr, ctx)), ctx,
		localAddr.IA, localAddr.Host)
	if err != nil {
		log.Fatal("Failed to start TSP propagator:", err)
	}

	ticker := time.NewTicker(30 * time.Second)

	ntpHosts := []string{"time.apple.com", "time.facebook.com", "time.google.com",
		"0.pool.ntp.org", "1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org"}

	for {
		select {
		case pathInfo = <-pathInfos:
			log.Printf("Received new path info.\n")
		case syncInfo := <-syncInfos:
			log.Printf("Received new sync info: %v\n", syncInfo)
		case <-ticker.C:
			log.Printf("Received new timer tick.\n")

			var clockOffsets []time.Duration
			for _, ntpHost := range ntpHosts {
				off, err := ets.GetNTPClockOffset(ntpHost);
				if err != nil {
					log.Printf("Failed to get NTP clock offset: %v\n", err)
				}
				clockOffsets = append(clockOffsets, off)
			}
			sort.Slice(clockOffsets, func(i, j int) bool {
				return clockOffsets[i] < clockOffsets[j]
			})

			var clockOffset time.Duration
			if len(clockOffsets) == 0 {
				clockOffset = 0
			} else {
				i := len(clockOffsets) / 2
				if len(clockOffsets) % 2 != 0 {
					clockOffset = clockOffsets[i]
				} else {
					clockOffset = (clockOffsets[i] + clockOffsets[i - 1]) / 2
				}
			}
			log.Printf("Median NTP clock offset %v\n", clockOffset)

			for remoteIA, ps := range pathInfo.CoreASes {
				sp := ps[rand.Intn(len(ps))]
				tsp.PropagatePacketTo(
					newPacket(pathInfo.LocalIA, remoteIA, sp, clockOffset),
					sp.UnderlayNextHop())
			}
			for _, ip := range pathInfo.LocalTSHosts {
				tsp.PropagatePacketTo(
					newPacket(pathInfo.LocalIA, pathInfo.LocalIA, /* path: */ nil, clockOffset),
					&net.UDPAddr{IP: ip, Port: topology.EndhostPort})
			}
		}
	}
}
