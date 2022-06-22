package udp

import (
	"fmt"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
)

type udpBroadcast struct {
	tag     string
	addr    *net.UDPAddr
	timeout time.Duration
	ch      chan protocol.Message
	closing chan struct{}
	closed  chan struct{}
}

func NewUDPBroadcast(spec string, timeout time.Duration) (*udpBroadcast, error) {
	addr, err := net.ResolveUDPAddr("udp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve UDP address '%v'", spec)
	}

	if addr.Port == 0 {
		return nil, fmt.Errorf("UDP requires a non-zero port")
	}

	out := udpBroadcast{
		tag:     "UDP",
		addr:    addr,
		timeout: timeout,
		ch:      make(chan protocol.Message),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	return &out, nil
}

func (udp *udpBroadcast) Close() {
	infof(udp.tag, "closing")
	close(udp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-udp.closed:
		infof(udp.tag, "closed")

	case <-timeout.C:
		infof(udp.tag, "close timeout")
	}
}

func (udp *udpBroadcast) Run(router *router.Switch) error {
loop:
	for {
		select {
		case msg := <-udp.ch:
			router.Received(msg.ID, msg.Message, nil)

		case <-udp.closing:
			break loop
		}
	}

	close(udp.closed)

	return nil
}

func (udp *udpBroadcast) Send(id uint32, msg []byte) {
	go func() {
		if reply := udp.send(id, msg); reply != nil {
			udp.ch <- protocol.Message{
				ID:      id,
				Message: reply,
			}
		}
	}()
}

func (udp *udpBroadcast) send(id uint32, message []byte) []byte {
	dumpf(udp.tag, message, "broadcast (%v bytes)", len(message))

	if bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0"); err != nil {
		warnf(udp.tag, "%v", err)
	} else if socket, err := net.ListenUDP("udp", bind); err != nil {
		warnf(udp.tag, "%v", err)
	} else if socket == nil {
		warnf(udp.tag, "invalid UDP socket (%v)", socket)
	} else {
		defer socket.Close()

		if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
			warnf(udp.tag, "%v", err)
		}

		if err := socket.SetReadDeadline(time.Now().Add(udp.timeout)); err != nil {
			warnf(udp.tag, "%v", err)
		}

		if N, err := socket.WriteToUDP(message, udp.addr); err != nil {
			warnf(udp.tag, "%v", err)
		} else {
			debugf(udp.tag, "sent %v bytes to %v\n", N, udp.addr)

			reply := make([]byte, 2048)

			if N, remote, err := socket.ReadFromUDP(reply); err != nil {
				warnf(udp.tag, "%v", err)
			} else {
				dumpf(udp.tag, reply[0:N], "received %v bytes from %v", N, remote)

				return reply[:N]
			}
		}
	}

	return nil
}
