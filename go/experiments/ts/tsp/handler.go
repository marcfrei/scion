package tsp

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/snet"
)

var handlerLog = log.New(os.Stderr, "[tsp/handler] ", log.LstdFlags) 

func StartHandler(s snet.PacketDispatcherService, ctx context.Context,
	localIA addr.IA, localHost *net.UDPAddr) error {
	conn, localPort, err := s.Register(ctx, localIA, localHost, addr.SvcTS)
	if err != nil {
		return err
	}

	handlerLog.Printf("Listening in %v on %v:%d - %v\n",
		localIA, localHost.IP, localPort, addr.SvcTS)

	go func() {
		for {
			var packet snet.Packet
			var lastHop net.UDPAddr
			err := conn.ReadFrom(&packet, &lastHop)
			if err != nil {
				handlerLog.Printf("Failed to read packet: %v\n", err)
				continue
			}

			handlerLog.Printf("Received packet: %v\n", string(
				packet.Payload.(common.RawBytes)))
		}
	}()

	return nil
}
