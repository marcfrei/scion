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

const (
	flagStartRound = 0
	flagBroadcast = 1
	flagUpdate = 2

	roundPeriod = 30 * time.Second
	roundDuration = 20 * time.Second
)

type syncEntry struct {
	ia addr.IA
	syncInfos []tsp.SyncInfo
}

var ntpHosts = []string{
	"time.apple.com",
	"time.facebook.com",
	"time.google.com",
	// "0.pool.ntp.org",
	// "1.pool.ntp.org",
	// "2.pool.ntp.org",
	// "3.pool.ntp.org",
}

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

func ntpClockOffset() time.Duration {
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

	return clockOffset
}

func syncEntryForIA(syncEntries []syncEntry, ia addr.IA) *syncEntry {
	for i, x := range syncEntries {
		if x.ia == ia {
			return &syncEntries[i]
		}
	}
	return nil
}

func syncInfoForHost(syncInfos []tsp.SyncInfo, host addr.HostAddr) *tsp.SyncInfo {
	for i, x := range syncInfos {
		if x.Source.Host.Equal(host) {
			return &syncInfos[i]
		}
	}
	return nil
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

	var syncEntries []syncEntry

	flag := flagStartRound
	now := time.Now().UTC()
	syncTime := now.Add(roundPeriod).Truncate(roundPeriod)
	log.Printf("Next sync time: %v\n", syncTime)
	timer := time.NewTimer(syncTime.Add(-(roundDuration / 2)).Sub(now))

	var clockOffset time.Duration

	for {
		select {
		case pathInfo = <-pathInfos:
			log.Printf("Received new path info.\n")
		case syncInfo := <-syncInfos:
			log.Printf("Received new sync info: %v\n", syncInfo)
			se := syncEntryForIA(syncEntries, syncInfo.Source.IA)
			if se == nil {
				syncEntries = append(syncEntries, syncEntry{
					ia: syncInfo.Source.IA,
				})
				se = &syncEntries[len(syncEntries) - 1]
			}
			si := syncInfoForHost(se.syncInfos, syncInfo.Source.Host)
			if si != nil {
				*si = syncInfo
			} else {
				se.syncInfos = append(se.syncInfos, syncInfo)
			}
		case now = <-timer.C:
			log.Printf("Received new timer signal: %v\n", now.UTC())
			switch flag {
			case flagStartRound:
				log.Printf("START ROUND\n")

				syncEntries = nil

				flag = flagBroadcast
				now = time.Now().UTC()
				timer.Reset(syncTime.Sub(now))
			case flagBroadcast:
				log.Printf("BROADCAST\n")

				clockOffset = ntpClockOffset()

				for remoteIA, ps := range pathInfo.CoreASes {
					if syncEntryForIA(syncEntries, remoteIA) == nil {
						syncEntries = append(syncEntries, syncEntry{
							ia: remoteIA,
						})
					}
					sp := ps[rand.Intn(len(ps))]
					tsp.PropagatePacketTo(
						newPacket(pathInfo.LocalIA, remoteIA, sp, clockOffset),
						sp.UnderlayNextHop())
				}
				if syncEntryForIA(syncEntries, pathInfo.LocalIA) == nil {
					syncEntries = append(syncEntries, syncEntry{
						ia: pathInfo.LocalIA,
					})
				}
				for _, ip := range pathInfo.LocalTSHosts {
					tsp.PropagatePacketTo(
						newPacket(pathInfo.LocalIA, pathInfo.LocalIA, /* path: */ nil, clockOffset),
						&net.UDPAddr{IP: ip, Port: topology.EndhostPort})
				}

				flag = flagUpdate
				now = time.Now().UTC()
				timer.Reset(syncTime.Add(roundDuration / 2).Sub(now))
			case flagUpdate:
				log.Printf("UPDATE\n")

				for _, se := range syncEntries {
					log.Printf("\t%v:\n", se.ia)
					for _, sie := range se.syncInfos {
						log.Printf("\t\t%v\n", sie)
					}
				}

				flag = flagStartRound
				now = time.Now().UTC()
				syncTime = now.Add(roundPeriod).Truncate(roundPeriod)
				log.Printf("Next sync time: %v\n", syncTime)
				timer.Reset(syncTime.Add(-(roundDuration / 2)).Sub(now))
			}
		}
	}
}
