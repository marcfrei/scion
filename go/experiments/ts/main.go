package main

import (
	"context"
	"flag"
	"log"
	"net"
	"math/rand"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/topology"
	"github.com/scionproto/scion/go/proto"

	"github.com/scionproto/scion/go/experiments/ts/tsp"
)

const (
	ModeClient = "client"
	ModeServer = "server"
)

func newPacket(localIA addr.IA, localHost *net.UDPAddr, remoteIA addr.IA, path snet.Path) *snet.Packet {
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
			Source: snet.SCIONAddress{
				IA: localIA,
				Host: addr.HostFromIP(localHost.IP),
			},
			Path: dpPath,
			Payload: common.RawBytes(
				[]byte(localIA.String() + " > " + remoteIA.String()),
			),
		},
	}
}

func main() {
	var mode string
	var sciondAddr string
	var localAddr snet.UDPAddr
	flag.StringVar(&mode, "mode", ModeClient, "Run in "+ModeClient+" or "+ModeServer+" mode")
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

	if mode == ModeServer {
		packetConn, localPort, err := packetDispatcher.Register(context.TODO(),
			localAddr.IA, localAddr.Host, addr.SvcTS)
		if err != nil {
			log.Fatal("Failed to register listener:", err)
		}

		log.Printf("Listening on %v:%d - %v\n", localAddr.Host.IP, localPort, addr.SvcTS)

		for {
			var packet snet.Packet
			var lastHop net.UDPAddr
			err := packetConn.ReadFrom(&packet, &lastHop)
			if err != nil {
				log.Println("Failed to read packet:", err)
				continue
			}

			log.Printf("Packet: %v\n", string(packet.Payload.(common.RawBytes)))
		}
	} else {
		err := tsp.StartPropagator(packetDispatcher, context.TODO(),
			localAddr.IA, localAddr.Host, addr.SvcNone)
		if err != nil {
			log.Fatal("Failed to start TSP propagatot:", err)
		}

		corePaths, err := sciondConnector.Paths(context.TODO(),
			addr.IA{I: 0, A: 0}, localAddr.IA, sciond.PathReqFlags{Refresh: true})
		if err != nil {
			log.Fatal("Failed to lookup core paths:", err)
		}
		coreASes := make(map[addr.IA][]snet.Path)
		for _, p := range corePaths {
			coreASes[p.Destination()] = append(coreASes[p.Destination()], p)
		}

		log.Printf("Reachable core ASes:\n")
		for coreAS := range coreASes {
			log.Printf("%v", coreAS)
			for _, p := range coreASes[coreAS] {
				log.Printf("\t%v\n", p)
			}
		}

		svcInfoReply, err := sciondConnector.SVCInfo(context.TODO(),
			[]proto.ServiceType{proto.ServiceType_ts})
		if err != nil {
			log.Fatal("Failed to lookup local TS service info:", err)
		}
		localTSHosts := make(map[string]bool)
		for _, i := range svcInfoReply.Entries[0].HostInfos {
			localTSHosts[i.Host().IP().String()] = true
		}

		log.Printf("Reachable local time services:\n")
		for localTSHost := range localTSHosts {
			log.Printf("%v", localTSHost)
		}

		for remoteIA := range coreASes {
			ps := coreASes[remoteIA]
			sp := ps[rand.Intn(len(ps))]
			tsp.PropagatePacketTo(
				newPacket(localAddr.IA, localAddr.Host, remoteIA, sp),
				sp.UnderlayNextHop())
		}
		for localTSHost := range localTSHosts {
			tsp.PropagatePacketTo(
				newPacket(localAddr.IA, localAddr.Host, localAddr.IA, /* path: */ nil),
				&net.UDPAddr{IP: net.ParseIP(localTSHost), Port: topology.EndhostPort})
		}

		select{}
	}
}
