package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
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

func main() {
	var sciondAddr string
	var localAddr snet.UDPAddr
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Parse()

	ctx := context.Background()
	var err error

	err = tsp.StartHandler(
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

	err = tsp.StartOriginator(newSciondConnector(sciondAddr, ctx), ctx)
	if err != nil {
		log.Fatal("Failed to start TSP originator:", err)
	}

	ntpHosts := []string{"time.apple.com", "time.facebook.com", "time.google.com",
		"0.pool.ntp.org", "1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org"} 

	for {
		for _, ntpHost := range ntpHosts {
			err = ets.GetNTPTime(ntpHost);
			if err != nil {
				log.Printf("Failed to get NTP time: %v\n", err)
			}
		}

		time.Sleep(30 * time.Second +
			time.Duration(rand.Int63n(250)) * time.Millisecond)
	}
}
