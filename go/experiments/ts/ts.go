package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/l4"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sock/reliable"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/topology"
)

const (
	ModeClient = "client"
	ModeServer = "server"
)

func main() {
	var mode string
	var sciondAddr string
	var localAddr snet.UDPAddr
	var remoteIA addr.IA
	flag.StringVar(&mode, "mode", ModeClient, "Run in "+ModeClient+" or "+ModeServer+" mode")
	flag.StringVar(&sciondAddr, "sciond", "", "SCIOND address")
	flag.Var(&localAddr, "local", "Local address")
	flag.Var(&remoteIA, "remote", "Remote ISD-AS number")
	flag.Parse()

	sciondConnector, err := sciond.NewService(sciondAddr).Connect(context.TODO())
	if err != nil {
		log.Fatal("Failed to create SCION connector:", err)
	}
	defer sciondConnector.Close(context.TODO())

	dispatcher := reliable.NewDispatcher("")
 	// TODO: wrap with reconnect.NewDispatcherService?

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
		packetConn, localPort, err := packetDispatcher.Register(context.TODO(),
			localAddr.IA, localAddr.Host, addr.SvcNone)
		if err != nil {
			log.Fatal("Failed to register listener:", err)
		}

		corePaths, err := sciondConnector.Paths(context.TODO(),
			addr.IA{I: 0, A: 0}, localAddr.IA, sciond.PathReqFlags{Refresh: true})
		if err != nil {
			log.Fatal("Failed to lookup core paths:", err)
		}
		coreASes := make(map[addr.IA]bool)
		for _, p := range corePaths {
			coreASes[p.Destination()] = true
		}

		log.Printf("Reachable core ASes from %v:\n", localAddr.IA)
		for coreAS := range coreASes {
			log.Printf("%v", coreAS)
		}
		log.Printf("Paths to reachable core ASes from %v:\n", localAddr.IA)
		for _, p := range corePaths {
			log.Printf("%v\n", p)
		}

		var selectedPath snet.Path
		for _, p := range corePaths {
			if p.Destination() == remoteIA {
				selectedPath = p
				break
			}
		}
		log.Printf("Selected path from %v to %v:\n", localAddr.IA, remoteIA)
		log.Printf("%v\n", selectedPath)

		var path *spath.Path
		var nextHop *net.UDPAddr
		if selectedPath != nil {
			path = selectedPath.Path()
			nextHop = selectedPath.UnderlayNextHop()
		} else {
			path = nil
			nextHop = &net.UDPAddr{
				IP: localAddr.Host.IP,
				Port: topology.EndhostPort,
			}
		}

		err = packetConn.WriteTo(
			&snet.Packet{
				PacketInfo: snet.PacketInfo{
					Destination: snet.SCIONAddress{
						IA: remoteIA,
						Host: addr.SvcTS | addr.SVCMcast,
					},
					Source: snet.SCIONAddress{
						IA: localAddr.IA,
						Host: addr.HostFromIP(localAddr.Host.IP),
					},
					Path: path,
					L4Header: &l4.UDP{
						SrcPort: localPort,
					},
					Payload: common.RawBytes(
						[]byte(localAddr.IA.String() + " > " + remoteIA.String()),
					),
				},
			},
			nextHop,
		)
		if err != nil {
			log.Println("Failed to write packet:", err)
		}
	}
}