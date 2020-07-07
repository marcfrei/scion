package tsp

import (
	"context"
	"log"
	"net"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/snet"
)

func StartHandler(s snet.PacketDispatcherService, ctx context.Context,
	localIA addr.IA, localHost *net.UDPAddr) error {
	packetConn, localPort, err := s.Register(ctx, localIA, localHost, addr.SvcTS)
	if err != nil {
		return err
	}

	log.Printf("[TSP handler] Listening in %v on %v:%d - %v\n",
		localIA, localHost.IP, localPort, addr.SvcTS)

	go func() {
		for {
			var packet snet.Packet
			var lastHop net.UDPAddr
			err := packetConn.ReadFrom(&packet, &lastHop)
			if err != nil {
				log.Printf("[TSP handler] Failed to read packet: %v\n", err)
				continue
			}

			log.Printf("[TSP handler] Received packet: %v\n", string(packet.Payload.(common.RawBytes)))
		}
	}()

	return nil
}
