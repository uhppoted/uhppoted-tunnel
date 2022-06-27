package udp

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type udpBroadcast struct {
	conn.Conn
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
		Conn: conn.Conn{
			Tag: "UDP",
		},
		addr:    addr,
		timeout: timeout,
		ch:      make(chan protocol.Message),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	return &out, nil
}

func (udp *udpBroadcast) Close() {
	udp.Infof("closing")
	close(udp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-udp.closed:
		udp.Infof("closed")

	case <-timeout.C:
		udp.Infof("close timeout")
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
		udp.send(id, msg)
	}()
}

func (udp *udpBroadcast) send(id uint32, message []byte) {
	udp.Dumpf(message, "broadcast (%v bytes)", len(message))

	if bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0"); err != nil {
		udp.Warnf("%v", err)
	} else if socket, err := net.ListenUDP("udp", bind); err != nil {
		udp.Warnf("%v", err)
	} else if socket == nil {
		udp.Warnf("invalid UDP socket (%v)", socket)
	} else {
		defer socket.Close()

		if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
			udp.Warnf("%v", err)
		}

		if err := socket.SetReadDeadline(time.Now().Add(2 * udp.timeout)); err != nil {
			udp.Warnf("%v", err)
		}

		if N, err := socket.WriteToUDP(message, udp.addr); err != nil {
			udp.Warnf("%v", err)
		} else {
			udp.Debugf("sent %v bytes to %v\n", N, udp.addr)

			go func() {
				for {
					reply := make([]byte, 2048)

					if N, remote, err := socket.ReadFromUDP(reply); err != nil && !errors.Is(err, net.ErrClosed) {
						udp.Warnf("%v", err)
						return
					} else if err != nil {
						return
					} else {
						udp.Dumpf(reply[0:N], "received %v bytes from %v", N, remote)

						udp.ch <- protocol.Message{
							ID:      id,
							Message: reply[:N],
						}
					}
				}
			}()

			<-time.After(udp.timeout)
		}
	}
}
