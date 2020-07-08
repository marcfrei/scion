package main

import (
	"context"
	"flag"
	"log"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"

	"github.com/scionproto/scion/go/experiments/ts/tsp"
)

func main() {
	var sciondAddr string
	var localAddr snet.UDPAddr
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Parse()

	sciondConnector, err := sciond.NewService(sciondAddr).Connect(context.TODO())
	if err != nil {
		log.Fatal("Failed to create SCION connector:", err)
	}
	defer sciondConnector.Close(context.TODO())

	dispatcher := reconnect.NewDispatcherService(reliable.NewDispatcher(""))

	packetDispatcher := &snet.DefaultPacketDispatcherService{
		Dispatcher: dispatcher,
		SCMPHandler: snet.NewSCMPHandler(
			sciond.RevHandler{Connector: sciondConnector},
		),
	}

	ctx := context.TODO()

	err = tsp.StartHandler(packetDispatcher, ctx, localAddr.IA, localAddr.Host)
	if err != nil {
		log.Fatal("Failed to start TSP handler:", err)
	}

	localAddr.Host.Port = 0
	err = tsp.StartPropagator(packetDispatcher, ctx, localAddr.IA, localAddr.Host)
	if err != nil {
		log.Fatal("Failed to start TSP propagator:", err)
	}

	err = tsp.StartOriginator(sciondConnector, ctx)
	if err != nil {
		log.Fatal("Failed to start TSP originator:", err)
	}

	select{}
}
