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
	localPort uint16
	propagateRequests chan propagateRequest
}

var (
	propagators chan *propagator
	propagateRequests chan propagateRequest
)

func newPropagator(id int, packetConn snet.PacketConn, localPort uint16) propagator {
	return propagator{
		id: id,
		packetConn: packetConn,
		localPort: localPort,
		propagateRequests: make(chan propagateRequest),
	}
}

func (p *propagator) start() {
		go func() {
			for {
				propagators <- p
				log.Printf("propagator[%d]: awaiting requests\n", p.id)
				select {
				case r := <-p.propagateRequests:
					log.Printf("propagator[%d]: received request %v: %v, %v\n", p.id, r, r.pkt, r.nextHop)
					r.pkt.PacketInfo.L4Header = &l4.UDP{SrcPort: p.localPort};
					err := p.packetConn.WriteTo(r.pkt, r.nextHop)
					if err != nil {
						log.Printf("propagator[%d]: failed to write packet: %v\n", p.id, err)
					}
					log.Printf("propagator[%d]: handled request\n", p.id)
				}
			}
		}()
}

func StartPropagator(s snet.PacketDispatcherService, ctx context.Context,
	ia addr.IA, registration *net.UDPAddr, svc addr.HostSVC) error {
	propagators = make(chan *propagator, nPropagators)
	propagateRequests = make(chan propagateRequest, nPropagateRequests)

	for i := 0; i < cap(propagators); i++ {
		packetConn, localPort, err := s.Register(ctx, ia, registration, svc)
		if err != nil {
			return err
		}
		p := newPropagator(i, packetConn, localPort)
		p.start()
	}

	go func() {
		for {
			select {
			case r := <-propagateRequests:
				log.Printf("propagator[main]: received request %v\n", r)
				p := <-propagators
				p.propagateRequests <- r
				log.Printf("propagator[main]: handled request %v\n", r)
			}
		}
	}()

	return nil
}

func PropagatePacketTo(pkt *snet.Packet, nextHop *net.UDPAddr) {
	propagateRequests <- propagateRequest{pkt, nextHop}
}


