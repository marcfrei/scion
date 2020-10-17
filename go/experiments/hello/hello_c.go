package main

import (
	"context"
	"flag"
	"log"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"
)

func main() {
	var sciondAddr string
	var localAddr snet.UDPAddr
	var remoteAddr snet.UDPAddr
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Var(&remoteAddr, "remote", "Remote address")
	flag.Parse()

	var err error
	ctx := context.Background()

	sdc, err := sciond.NewService(sciondAddr).Connect(ctx)
	if err != nil {
		log.Fatal("Failed to create SCION connector:", err)
	}
	pds := &snet.DefaultPacketDispatcherService{
		Dispatcher: reconnect.NewDispatcherService(reliable.NewDispatcher("")),
		SCMPHandler: snet.DefaultSCMPHandler{
			RevocationHandler: sciond.RevHandler{Connector: sdc},
		},
	}

	localIA, err := sdc.LocalIA(ctx)
	if err != nil {
		log.Fatal("Failed to get local IA:", err)
	}

	corePaths, err := sdc.Paths(ctx,
		addr.IA{I: 0, A: 0}, localIA, sciond.PathReqFlags{Refresh: true})
	if err != nil {
		log.Fatal("Failed to lookup core paths: %v:", err)
	}
	coreASes := make(map[addr.IA][]snet.Path)
	if corePaths != nil {
		for _, p := range corePaths {
			ifcs := p.Interfaces()
			dest := ifcs[len(ifcs) - 1].IA
			coreASes[dest] = append(coreASes[dest], p)
		}
	}

	log.Printf("Reachable core ASes:\n")
	for coreAS := range coreASes {
		log.Printf("%v", coreAS)
		for _, p := range coreASes[coreAS] {
			log.Printf("\t%v\n", p)
		}
	}

	p := coreASes[remoteAddr.IA][0];
	log.Printf("Selected path to %v: %v\n", remoteAddr.IA, p)

	localAddr.Host.Port = 0
	conn, localPort, err := pds.Register(ctx, localAddr.IA, localAddr.Host, addr.SvcNone)
	if err != nil {
		log.Fatal("Failed to register client socket:", err)
	}

	log.Printf("Sending in %v on %v:%d - %v\n", localAddr.IA, localAddr.Host.IP, localPort, addr.SvcNone)

	pkt := &snet.Packet{
		PacketInfo: snet.PacketInfo{
			Source: snet.SCIONAddress{
				IA: localAddr.IA,
				Host: addr.HostFromIP(localAddr.Host.IP),
			},
			Destination: snet.SCIONAddress{
				IA: remoteAddr.IA,
				Host: addr.SvcTS | addr.SVCMcast,
			},
			Path: p.Path(),
			Payload: snet.UDPPayload{
				SrcPort: localPort,
				Payload: []byte("Hello, world!"),
			},
		},
	}

	err = conn.WriteTo(pkt, p.UnderlayNextHop())
	if err != nil {
		log.Printf("[%d] Failed to write packet: %v\n", err)
	}
}
