package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"sort"
	"time"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/daemon"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"
	"github.com/scionproto/scion/go/lib/topology"

	"github.com/scionproto/scion/go/experiments/ts/ets"
	"github.com/scionproto/scion/go/experiments/ts/etd"
	"github.com/scionproto/scion/go/experiments/ts/tsp"
)

const (
	flagStartRound = 0
	flagBroadcast = 1
	flagUpdate = 2

	roundPeriod = 5 * time.Second
	roundDuration = 2 * time.Second
)

type timeSource interface {
	fetchClockOffset() (time.Duration, error)
}

type mbgTimeSource string
type ntpTimeSource string

type syncEntry struct {
	ia addr.IA
	syncInfos []tsp.SyncInfo
}

var (
	timeSources = []timeSource{
		// mbgTimeSource("/dev/mbgclock0"),
		ntpTimeSource("time.apple.com"),
		ntpTimeSource("time.facebook.com"),
		ntpTimeSource("time.google.com"),
		ntpTimeSource("time.windows.com"),
	}

	clockCorrection time.Duration
)

func newSciondConnector(addr string, ctx context.Context) daemon.Connector {
	c, err := daemon.NewService(addr).Connect(ctx)
	if err != nil {
		log.Fatal("Failed to create SCION connector:", err)
	}
	return c
}

func newPacketDispatcher(c daemon.Connector) snet.PacketDispatcherService {
	return &snet.DefaultPacketDispatcherService{
		Dispatcher: reconnect.NewDispatcherService(reliable.NewDispatcher("")),
		SCMPHandler: snet.DefaultSCMPHandler{
			RevocationHandler: daemon.RevHandler{Connector: c},
		},
	}
}

func median(ds []time.Duration) time.Duration {
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

func midpoint(ds []time.Duration, f int) time.Duration {
	sort.Slice(ds, func(i, j int) bool {
		return ds[i] < ds[j]
	})
	return (ds[f] + ds[len(ds) - 1 - f]) / 2
}

func (s mbgTimeSource) fetchClockOffset() (time.Duration, error) {
	return ets.FetchMBGClockOffset(string(s))
}

func (s ntpTimeSource) fetchClockOffset() (time.Duration, error) {
	return ets.FetchNTPClockOffset(string(s))
}

func fetchClockOffset() <-chan time.Duration {
	clockOffset := make(chan time.Duration, 1)
	go func () {
		var clockOffsets []time.Duration
		ch := make(chan time.Duration)
		for _, h := range timeSources {
				go func(s timeSource) {
						off, err := s.fetchClockOffset()
						if err != nil {
							log.Printf("Failed to fetch clock offset from %v: %v\n", s, err)
						}
						ch <- off
				}(h)
		}
		for i := 0; i != len(timeSources); i++ {
			clockOffsets = append(clockOffsets, <-ch)
		}
		m := median(clockOffsets)
		log.Printf("Fetched local clock offset: %v\n", m)
		clockOffset <- m
	}()
	return clockOffset
}

func syncEntryClockOffset(syncInfos []tsp.SyncInfo) time.Duration {
	var clockOffsets []time.Duration
	for _, syncInfo := range syncInfos {
		clockOffsets = append(clockOffsets, syncInfo.ClockOffset)
	}
	return median(clockOffsets)
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

func syncedTime() time.Time {
	return time.Now().UTC().Add(clockCorrection)
}

func newTimer() *time.Timer {
	t := time.NewTimer(0)
	<-t.C
	return t
}

func scheduleNextRound(t *time.Timer, syncTime *time.Time) {
	now := syncedTime()
	*syncTime = now.Add(roundPeriod).Truncate(roundPeriod)
	t.Reset(syncTime.Add(-(roundDuration / 2)).Sub(now))
}

func scheduleBroadcast(t *time.Timer, syncTime time.Time) {
	now := syncedTime()
	t.Reset(syncTime.Sub(now))
}

func scheduleUpdate(t *time.Timer, syncTime time.Time) {
	now := syncedTime()
	t.Reset(syncTime.Add(roundDuration / 2).Sub(now))
}

func main() {
	var sciondAddr string
	var localAddr snet.UDPAddr
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Parse()

	var err error
	ctx := context.Background()

	topo, err := topology.FromJSONFile(filepath.Join(".", "gen",
		fmt.Sprintf("AS%s/topology.json", localAddr.IA.A.FileFmt())))
	if err != nil {
		log.Fatal("Failed to load topology:", err)
	}
	tsInfos, err := topo.TimeServices()
	if err != nil {
		log.Fatal("Failed to load local time service infos:", err)
	}
	var peers []addr.IA
	for _, tsi := range tsInfos {
		a := tsi.Addr.SCIONAddress
		if a.IP.Equal(localAddr.Host.IP) &&
			a.Port == localAddr.Host.Port &&
			a.Zone == localAddr.Host.Zone {
			peers = tsi.Peers
		}
	}

	pathInfos, err := tsp.StartPather(newSciondConnector(sciondAddr, ctx), ctx, peers)
	if err != nil {
		log.Fatal("Failed to start TSP pather:", err)
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

	var clockOffset struct {
		d time.Duration
		c <-chan time.Duration
	}

	var syncEntries []syncEntry
	var syncTime time.Time

	clockCorrection = <-fetchClockOffset()

	flag := flagStartRound
	syncTimer := newTimer()
	scheduleNextRound(syncTimer, &syncTime)

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
		case now := <-syncTimer.C:
			now = now.UTC()
			log.Printf("Received new timer signal: %v\n", now)
			switch flag {
			case flagStartRound:
				log.Printf("START ROUND\n")

				syncEntries = nil
				clockOffset.d = 0
				clockOffset.c = fetchClockOffset()

				flag = flagBroadcast
				scheduleBroadcast(syncTimer, syncTime)
			case flagBroadcast:
				log.Printf("BROADCAST\n")

				if len(clockOffset.c) == 0 {
					clockOffset.d = <-clockOffset.c - clockCorrection

					clockCorrection += clockOffset.d
					log.Printf("Clock correction: %v\n", clockCorrection)

					flag = flagStartRound
					scheduleNextRound(syncTimer, &syncTime)
				} else {
					clockOffset.d = <-clockOffset.c - clockCorrection

					for remoteIA, ps := range pathInfo.PeerASes {
						if syncEntryForIA(syncEntries, remoteIA) == nil {
							syncEntries = append(syncEntries, syncEntry{
								ia: remoteIA,
							})
						}
						sp := ps[rand.Intn(len(ps))]
						tsp.PropagatePacketTo(
							&snet.Packet{
								PacketInfo: snet.PacketInfo{
									Destination: snet.SCIONAddress{
										IA: remoteIA,
										Host: addr.SvcTS | addr.SVCMcast,
									},
									Path: sp.Path(),
									Payload: snet.UDPPayload{
										Payload: []byte(clockOffset.d.String()),
									},
								},
							},
							sp.UnderlayNextHop())
					}

					localSyncInfo := tsp.SyncInfo{
						Source: snet.SCIONAddress{
							IA: localAddr.IA,
							Host: addr.HostFromIP(localAddr.Host.IP),
						},
						ClockOffset: clockOffset.d,
					};
					se := syncEntryForIA(syncEntries, localAddr.IA);
					if se == nil {
						syncEntries = append(syncEntries, syncEntry{
							ia: localAddr.IA,
						})
						se = &syncEntries[len(syncEntries) - 1]
					}
					si := syncInfoForHost(se.syncInfos, addr.HostFromIP(localAddr.Host.IP))
					if si != nil {
						*si = localSyncInfo
					} else {
						se.syncInfos = append(se.syncInfos, localSyncInfo)
					}

					flag = flagUpdate
					scheduleUpdate(syncTimer, syncTime)
				}
			case flagUpdate:
				log.Printf("UPDATE\n")

				var clockOffsets []time.Duration
				for _, se := range syncEntries {
					var d time.Duration
					if len(se.syncInfos) != 0 {
						d = syncEntryClockOffset(se.syncInfos)
					} else {
						d = clockOffset.d
					}
					clockOffsets = append(clockOffsets, d)
					log.Printf("\t%v:\n", se.ia)
					for _, si := range se.syncInfos {
						log.Printf("\t\t%v\n", si)
					}
				}

				loff := clockOffset.d + clockCorrection
				log.Printf("Local clock correction: %v -> %v\n", loff, now.Add(loff))

				if len(clockOffsets) == 0 {
					clockCorrection += clockOffset.d
				} else {
					f := (len(clockOffsets) - 1) / 3
					clockCorrection += midpoint(clockOffsets, f)
				}
				// log.Printf("Clock correction: %v -> %v\n", clockCorrection, syncedTime())

				goff := clockCorrection
				log.Printf("Global clock correction: %v -> %v\n", goff, now.Add(goff))
				log.Printf("Global/local clock diff: %v\n", now.Add(goff).Sub(now.Add(loff)))

				etd.StoreClockSample(now.Add(goff), now)

				flag = flagStartRound
				scheduleNextRound(syncTimer, &syncTime)
			}
		}
	}
}
