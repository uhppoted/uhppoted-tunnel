package udp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type udpBroadcast struct {
	conn.Conn
	hwif    string
	addr    *net.UDPAddr
	timeout time.Duration
	ctx     context.Context
	ch      chan protocol.Message
	closed  chan struct{}
}

func NewUDPBroadcast(hwif string, spec string, timeout time.Duration, ctx context.Context) (*udpBroadcast, error) {
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
		hwif:    hwif,
		addr:    addr,
		timeout: timeout,
		ctx:     ctx,
		ch:      make(chan protocol.Message),
		closed:  make(chan struct{}),
	}

	return &out, nil
}

func (udp *udpBroadcast) Close() {
	udp.Infof("closing")

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

		case <-udp.ctx.Done():
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

	listener := net.ListenConfig{
		Control: func(network, address string, connection syscall.RawConn) error {
			if udp.hwif != "" {
				return conn.BindToDevice(connection, udp.hwif, conn.IsIPv4(udp.addr.IP), udp.Conn)
			} else {
				return nil
			}
		},
	}

	if bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0"); err != nil {
		udp.Warnf("%v", err)
	} else if socket, err := listener.ListenPacket(context.Background(), "udp4", fmt.Sprintf("%v", bind)); err != nil {
		udp.Warnf("%v", err)
	} else if socket == nil {
		udp.Warnf("invalid UDP socket (%v)", socket)
	} else {
		defer socket.Close()

		if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
			udp.Warnf("%v", err)
		}

		if err := socket.SetReadDeadline(time.Now().Add(5*time.Second + udp.timeout)); err != nil {
			udp.Warnf("%v", err)
		}

		if N, err := socket.WriteTo(message, udp.addr); err != nil {
			udp.Warnf("%v", err)
		} else {
			udp.Debugf("sent %v bytes to %v\n", N, udp.addr)

			ctx, cancel := context.WithTimeout(udp.ctx, udp.timeout+5*time.Second)

			defer cancel()

			go func() {
				for {
					reply := make([]byte, 2048)

					if N, remote, err := socket.ReadFrom(reply); err != nil && !errors.Is(err, net.ErrClosed) {
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

			select {
			case <-time.After(udp.timeout):
				// Ok

			case <-ctx.Done():
				udp.Warnf("%v", ctx.Err())
			}
		}
	}
}
