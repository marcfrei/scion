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
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/topology"

	"github.com/scionproto/scion/go/experiments/ts/ets"
	"github.com/scionproto/scion/go/experiments/ts/tsp"
)

const (
	flagStartRound = 0
	flagBroadcast = 1
	flagUpdate = 2

	roundPeriod = 30 * time.Second
	roundDuration = 10 * time.Second
)

type syncEntry struct {
	ia addr.IA
	syncInfos []tsp.SyncInfo
}

var ntpHosts = []string{
	"time.apple.com",
	"time.facebook.com",
	"time.google.com",
	"0.pool.ntp.org",
	"1.pool.ntp.org",
	"2.pool.ntp.org",
	"3.pool.ntp.org",
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

func medianTimeDuration(ds []time.Duration) time.Duration {
	sort.Slice(ds, func(i, j int) bool {
		return ds[i] < ds[j]
	})
	var m time.Duration
	n := len(ds)
	if n == 0 {
		m = 0
	} else {
		i := n / 2
		if n % 2 != 0 {
			m = ds[i]
		} else {
			m = (ds[i] + ds[i - 1]) / 2
		}
	}
	return m
}

func midpointTimeDuration(ds []time.Duration, f int) time.Duration {
	if f == 0 {
		return 0
	}
	n := len(ds)
	if n < 3 * f + 1 {
		return 0
	}
	sort.Slice(ds, func(i, j int) bool {
		return ds[i] < ds[j]
	})
	return (ds[f] + ds[n - f - 1]) / 2
}

func ntpClockOffset() time.Duration {
	var clockOffsets []time.Duration
	for _, ntpHost := range ntpHosts {
		off, err := ets.GetNTPClockOffset(ntpHost);
		if err != nil {
			log.Printf("Failed to get NTP clock offset from %v: %v\n", ntpHost, err)
		}
		clockOffsets = append(clockOffsets, off)
	}
	return medianTimeDuration(clockOffsets)
}

func syncEntryClockOffset(syncInfos []tsp.SyncInfo) time.Duration {
	var clockOffsets []time.Duration
	for _, syncInfo := range syncInfos {
		clockOffsets = append(clockOffsets, syncInfo.ClockOffset)
	}
	return medianTimeDuration(clockOffsets)
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

	var clockCorrection time.Duration
	var syncEntries []syncEntry

	flag := flagStartRound
	now := time.Now().UTC().Add(clockCorrection)
	syncTime := now.Add(roundPeriod).Truncate(roundPeriod)
	timer := time.NewTimer(syncTime.Add(-(roundDuration / 2)).Sub(now))

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
				now = time.Now().UTC().Add(clockCorrection)
				timer.Reset(syncTime.Sub(now))
			case flagBroadcast:
				log.Printf("BROADCAST\n")

				clockOffset := ntpClockOffset()

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
				now = time.Now().UTC().Add(clockCorrection)
				timer.Reset(syncTime.Add(roundDuration / 2).Sub(now))
			case flagUpdate:
				log.Printf("UPDATE\n")

				var clockOffsets []time.Duration
				for _, se := range syncEntries {
					clockOffsets = append(clockOffsets, syncEntryClockOffset(se.syncInfos))
					log.Printf("\t%v:\n", se.ia)
					for _, sie := range se.syncInfos {
						log.Printf("\t\t%v\n", sie)
					}
				}

				var f int
				if len(clockOffsets) != 0 {
					f = (len(clockOffsets) - 1) / 3
				}
				d := midpointTimeDuration(clockOffsets, f);
				log.Printf("Delta correction: %v\n", d)

				clockCorrection += d
				log.Printf("Overall correction: %v\n", clockCorrection)

				flag = flagStartRound
				now = time.Now().UTC().Add(clockCorrection)
				syncTime = now.Add(roundPeriod).Truncate(roundPeriod)
				timer.Reset(syncTime.Add(-(roundDuration / 2)).Sub(now))
			}
		}
	}
}
