package tsp

import (
	"context"
	"log"
	"net"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/l4"
	"github.com/scionproto/scion/go/lib/snet"
)

const (
	nPropagators = 16
	nPropagateRequests = 128
)

type propagateRequest struct {
	pkt *snet.Packet
	nextHop *net.UDPAddr
}

type propagator struct {
	id int
	packetConn snet.PacketConn
	localIA addr.IA
	localHost addr.HostAddr
	localPort uint16
	propagateRequests chan propagateRequest
}

var (
	localHost net.UDPAddr
	propagators chan *propagator
	propagateRequests chan propagateRequest
)

func newPropagator(id int, packetConn snet.PacketConn,
	localIA addr.IA, localHost addr.HostAddr, localPort uint16) propagator {
	return propagator{
		id: id,
		packetConn: packetConn,
		localIA: localIA,
		localHost: localHost,
		localPort: localPort,
		propagateRequests: make(chan propagateRequest),
	}
}

func (p *propagator) start() {
	go func() {
		for {
			propagators <- p
			log.Printf("[TSP propagator, %d] Awaiting requests\n", p.id)
			select {
			case r := <-p.propagateRequests:
				log.Printf("[TSP propagator, %d] Received request %v: %v, %v\n", p.id, r, r.pkt, r.nextHop)
				r.pkt.Source = snet.SCIONAddress{IA: p.localIA, Host: p.localHost};
				r.pkt.PacketInfo.L4Header = &l4.UDP{SrcPort: p.localPort};
				err := p.packetConn.WriteTo(r.pkt, r.nextHop)
				if err != nil {
					log.Printf("[TSP propagator, %d] Failed to write packet: %v\n", p.id, err)
				}
				log.Printf("[TSP propagator, %d] Handled request\n", p.id)
			}
		}
	}()
}

func StartPropagator(s snet.PacketDispatcherService, ctx context.Context,
	localIA addr.IA, localHost *net.UDPAddr) error {
	propagators = make(chan *propagator, nPropagators)
	propagateRequests = make(chan propagateRequest, nPropagateRequests)

	for i := 0; i < cap(propagators); i++ {
		packetConn, localPort, err := s.Register(ctx, localIA, localHost, addr.SvcNone)
		if err != nil {
			// TODO: stop already started propagators
			return err
		}
		p := newPropagator(i, packetConn, localIA, addr.HostFromIP(localHost.IP), localPort)
		p.start()
	}

	go func() {
		for {
			select {
			case r := <-propagateRequests:
				log.Printf("[TSP propagator] Received request %v\n", r)
				p := <-propagators
				p.propagateRequests <- r
				log.Printf("[TSP propagator] Handled request %v\n", r)
			}
		}
	}()

	return nil
}

func propagatePacketTo(pkt *snet.Packet, nextHop *net.UDPAddr) {
	propagateRequests <- propagateRequest{pkt, nextHop}
}


