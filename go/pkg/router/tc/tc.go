// Copyright 2021 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tc adds traffic control capabilities (packet prioritization) to the
// border router. Traffic control can be enabled/disabled in 'scion/go/posix-router/main.go'.
package tc

import (
	"net"

	"golang.org/x/net/ipv4"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/underlay/conn"
)

type TrafficClass int
type SchedulerId int

// Identify the possible traffic classes. The class "ClsOthers" must be zero. Packets where the
// traffic class is not set will be assigned to this class by default. Every class must be
// registered inside the 'NewQueues()' function.
const (
	ClsOthers TrafficClass = iota
	ClsColibri
	ClsEpic
	ClsOhpEmpty
	ClsSNC
)

// Identify the possible scheduling algorithms.
const (
	SchedOthersOnly SchedulerId = iota
	SchedStrictPriority
)

// Number of packets to schedule in a single scheduling round.
const outputBatchCnt = 2

// Queues describes the queues (one for each traffic class) for a certain router interface.
// The 'mapping' translates traffic classes to their respective queue. The 'nonempty' channel is
// used to signal to the scheduler that packets are ready, which is necessary to omit
// computationally expensive spinning over the queues. The 'scheduler' denotes the scheduling
// algorithm that will be used, and 'writeBuffer' is a pre-allocated buffer to speed up the
// WriteTo() function.
type Queues struct {
	mapping     map[TrafficClass]*ZeroAllocQueue
	nonempty    chan struct{}
	scheduler   Scheduler
	writeBuffer conn.Messages
}

type Scheduler interface {
	Schedule(qs *Queues) ([]ipv4.Message, error)
}

// NewQueues creates new queues. If 'scheduling' is set to true, a queue is allocated for each
// traffic class, otherwise only one large queue for 'ClsOthers' will be created.
func NewQueues(scheduling bool, maxPacketLength int) *Queues {
	qs := &Queues{}
	qs.nonempty = make(chan struct{}, 1)
	qs.mapping = make(map[TrafficClass]*ZeroAllocQueue)

	qs.mapping[ClsOthers] = newZeroAllocQueue(128, maxPacketLength)
	if scheduling {
		qs.mapping[ClsColibri] = newZeroAllocQueue(32, maxPacketLength)
		qs.mapping[ClsEpic] = newZeroAllocQueue(32, maxPacketLength)
		qs.mapping[ClsOhpEmpty] = newZeroAllocQueue(32, maxPacketLength)
		qs.mapping[ClsSNC] = newZeroAllocQueue(32, maxPacketLength)
	}

	qs.writeBuffer = conn.NewReadMessages(outputBatchCnt)
	return qs
}

// SetScheduler assigns the provided scheduler to the queues.
func (qs *Queues) SetScheduler(s SchedulerId, sncPacketsScheduled prometheus.Counter) {
	switch s {
	case SchedOthersOnly:
		qs.scheduler = &OthersOnlyScheduler{}
	case SchedStrictPriority:
		qs.scheduler = &StrictPriorityScheduler{
			SNCPacketsScheduled: sncPacketsScheduled,
		}
	default:
		qs.scheduler = &StrictPriorityScheduler{
			SNCPacketsScheduled: sncPacketsScheduled,
		}
	}
}

// Enqueue adds the message 'm' (destined to 'outAddr') to the queue corresponding to traffic
// class 'tc'.
func (qs *Queues) Enqueue(tc TrafficClass, m []byte, outAddr *net.UDPAddr) error {
	// Get the correct queue
	q, ok := qs.mapping[tc]
	if !ok {
		return serrors.New("unknown traffic class", "tc", tc)
	}

	err := q.enqueue(m, outAddr)
	if err != nil {
		return serrors.WithCtx(err, "tc", tc)
	}
	qs.setToNonempty()
	return nil
}

// setToNonempty signals to the scheduler that new messages are ready to be scheduled.
func (qs *Queues) setToNonempty() {
	select {
	case qs.nonempty <- struct{}{}:
	default:
	}
}

// WaitUntilNonempty blocks until new messages are ready to be scheduled.
func (qs *Queues) WaitUntilNonempty() {
	<-qs.nonempty
}

// dequeue reads up to 'batchSize' number of messages from the queue corresponding to the traffic
// class 'tc' and stores them in the message buffer 'ms'.
func (qs *Queues) dequeue(tc TrafficClass, batchSize int, ms []ipv4.Message) (int, error) {
	q, ok := qs.mapping[tc]
	if !ok {
		return 0, serrors.New("unknown traffic class")
	}
	return q.dequeue(batchSize, ms), nil
}

// Schedule dequeues messages from the queues and prioritizes them according the the specified
// Scheduler. It returns the messages that will be forwarded.
func (qs *Queues) Schedule() ([]ipv4.Message, error) {
	return qs.scheduler.Schedule(qs)
}

// ReturnBuffers returns the buffers to the queues where they were taken from.
func (qs *Queues) ReturnBuffers(ms []ipv4.Message) error {
	counter := 0
	for tc, q := range qs.mapping {
		if counter+q.borrowed > len(ms) {
			return serrors.New("too many packets to return", "number", len(ms))
		}

		for i := 0; i < q.borrowed; i++ {
			// Reset buffer
			ms[counter+i].Buffers[0] = ms[counter+i].Buffers[0][:0]

			// Return packet to queue buffer
			select {
			case q.emptyPackets <- ms[counter+i]:
			default:
				return serrors.New("queueBuffer full", "traffic class", tc)
			}
		}

		counter += q.borrowed
		q.borrowed = 0
	}
	return nil
}

// ZeroAllocQueue is a thread-safe queue for IP packets that operates on a fixed number of
// preallocated packets. The memory overhead is therefore constant.
type ZeroAllocQueue struct {
	emptyPackets  chan ipv4.Message
	filledPackets chan ipv4.Message
	borrowed      int
}

// newZeroAllocQueue creates a new queue with 'queueSize' IP packets, where the maximal size of a
// packet is given by 'maxPacketLength'.
func newZeroAllocQueue(queueSize, maxPacketLength int) *ZeroAllocQueue {
	q := &ZeroAllocQueue{}
	q.emptyPackets = make(chan ipv4.Message, queueSize)
	q.filledPackets = make(chan ipv4.Message, queueSize)
	q.borrowed = 0

	msgs := conn.NewReadMessages(queueSize)
	for _, msg := range msgs {
		msg.Buffers[0] = make([]byte, maxPacketLength)
		msg.Addr = &net.UDPAddr{
			IP: make(net.IP, 16),
		}
		q.emptyPackets <- msg
	}
	return q
}

func (q *ZeroAllocQueue) enqueue(m []byte, outAddr *net.UDPAddr) error {
	// Retrieve free buffer if available
	var p ipv4.Message
	select {
	case p = <-q.emptyPackets:
	default:
		return serrors.New("no empty packets available")
	}

	// Copy the message into the buffer
	p.Buffers[0] = p.Buffers[0][:len(m)]
	copy(p.Buffers[0], m)

	addr := p.Addr.(*net.UDPAddr)
	if outAddr != nil {
		addr.IP = addr.IP[:len(outAddr.IP)]
		copy(addr.IP, outAddr.IP)
		addr.Port = outAddr.Port
	} else {
		addr.IP = addr.IP[:0]
	}

	// Put the packet into the queue
	select {
	case q.filledPackets <- p:
		return nil
	default:
		return serrors.New("queue full")
	}
}

func (q *ZeroAllocQueue) dequeue(batchSize int, ms []ipv4.Message) int {
	var counter int
L:
	for counter = 0; counter < batchSize; counter++ {
		select {
		case m := <-q.filledPackets:
			ms[counter] = m
		default:
			break L // this breaks out of the loop (not above it, but after)
		}
	}

	q.borrowed = q.borrowed + counter
	return counter
}

func (tc TrafficClass) String() string {
	switch tc {
	case ClsOthers:
		return "Others"
	case ClsColibri:
		return "COLIBRI"
	case ClsEpic:
		return "EPIC"
	case ClsOhpEmpty:
		return "OHP/Empty"
	case ClsSNC:
		return "SNC (Time Sync)"
	default:
		return "Unknown traffic class"
	}
}
