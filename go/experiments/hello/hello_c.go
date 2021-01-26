package main

import "C"

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/daemon"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/sock/reliable/reconnect"
	"github.com/scionproto/scion/go/lib/topology/underlay"
)

func sendHello(sciondAddr string, localAddr snet.UDPAddr, remoteAddr snet.UDPAddr) {
	var err error
	ctx := context.Background()

	sdc, err := daemon.NewService(sciondAddr).Connect(ctx)
	if err != nil {
		log.Fatal("Failed to create SCION connector:", err)
	}
	pds := &snet.DefaultPacketDispatcherService{
		Dispatcher: reconnect.NewDispatcherService(reliable.NewDispatcher("")),
		SCMPHandler: snet.DefaultSCMPHandler{
			RevocationHandler: daemon.RevHandler{Connector: sdc},
		},
	}

	ps, err := sdc.Paths(ctx, remoteAddr.IA, localAddr.IA, daemon.PathReqFlags{Refresh: true})
	if err != nil {
		log.Fatal("Failed to lookup core paths: %v:", err)
	}

	log.Printf("Available paths to %v:\n", remoteAddr.IA)
	for _, p := range ps {
		log.Printf("\t%v\n", p)
	}

	sp := ps[0]
	log.Printf("Selected path to %v: %v\n", remoteAddr.IA, sp)

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
				Host: addr.HostFromIP(remoteAddr.Host.IP),
			},
			Path: sp.Path(),
			Payload: snet.UDPPayload{
				SrcPort: localPort,
				DstPort: uint16(remoteAddr.Host.Port),
				Payload: []byte("Hello, world!"),
			},
		},
	}

	nextHop := sp.UnderlayNextHop()
	if nextHop == nil && remoteAddr.IA.Equal(localAddr.IA) {
		nextHop = &net.UDPAddr{
			IP: remoteAddr.Host.IP,
			Port: underlay.EndhostPort,
			Zone: remoteAddr.Host.Zone,
		}
	}

	err = conn.WriteTo(pkt, nextHop)
	if err != nil {
		log.Printf("[%d] Failed to write packet: %v\n", err)
	}
}

//export SendHello
func SendHello(sciondAddr, local, remote *C.char) int {
	localAddr, err := snet.ParseUDPAddr(C.GoString(local))
	if err != nil {
		return 1
	}
	remoteAddr, err := snet.ParseUDPAddr(C.GoString(remote))
	if err != nil {
		return 1
	}

	sendHello(C.GoString(sciondAddr), *localAddr, *remoteAddr)
	return 0
}

func main() {
	var sciondAddr string
	var localAddr snet.UDPAddr
	var remoteAddr snet.UDPAddr
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Var(&remoteAddr, "remote", "Remote address")
	flag.Parse()

	sendHello(sciondAddr, localAddr, remoteAddr)
}
