package ip

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	// "net/netip"
	// "strings"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

// type address interface {
// 	*net.UDPAddr | *net.TCPAddr
// }

type ipOut struct {
	conn.Conn
	hwif          string
	broadcastAddr *net.UDPAddr
	timeout       time.Duration
	controllers   map[uint32]any
	ctx           context.Context
	ch            chan protocol.Message
	closed        chan struct{}
}

func NewIPOut(hwif string, spec string, timeout time.Duration, ctx context.Context) (*ipOut, error) {
	broadcast, err := net.ResolveUDPAddr("udp", spec)
	if err != nil {
		return nil, err
	}

	if broadcast == nil {
		return nil, fmt.Errorf("unable to resolve UDP broadcast address '%v'", spec)
	}

	if broadcast.Port == 0 {
		return nil, fmt.Errorf("UDP broadcast address requires a non-zero port")
	}

	ip := ipOut{
		Conn: conn.Conn{
			Tag: "IP",
		},
		hwif:          hwif,
		broadcastAddr: broadcast,
		timeout:       timeout,
		controllers:   map[uint32]any{},
		ctx:           ctx,
		ch:            make(chan protocol.Message),
		closed:        make(chan struct{}),
	}

	// for k, v := range map[uint32]string{405419896: "tcp:192.168.1.100:60000"} {
	// 	if addr, err := resolve(v); err != nil {
	// 		ip.Warnf("invalid controller address '%v' (%v)", v, err)
	// 	} else {
	// 		ip.controllers[k] = addr
	// 	}
	// }

	ip.Infof("connector::ip-out")

	return &ip, nil
}

func (ip *ipOut) Close() {
	ip.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-ip.closed:
		ip.Infof("closed")

	case <-timeout.C:
		ip.Infof("close timeout")
	}
}

func (ip *ipOut) Run(router *router.Switch) error {
loop:
	for {
		select {
		case msg := <-ip.ch:
			router.Received(msg.ID, msg.Message, nil)

		case <-ip.ctx.Done():
			break loop
		}
	}

	close(ip.closed)

	return nil
}

func (ip *ipOut) Send(id uint32, msg []byte) {
	go func() {
		ip.send(id, msg)
	}()
}

func (ip *ipOut) send(id uint32, message []byte) {
	if len(message) == 64 && message[0] == 0x17 {
		controller := binary.LittleEndian.Uint32(message[4:])

		if v, ok := ip.controllers[controller]; ok {
			switch addr := v.(type) {
			case *net.UDPAddr:
				ip.udpSendto(id, message, addr)
				return

			case *net.TCPAddr:
				ip.tcpSendto(id, message, addr)
				return
			}
		}
	}

	ip.broadcast(id, message)
}

func (ip *ipOut) udpSendto(id uint32, message []byte, addr *net.UDPAddr) {
	fmt.Printf(">>>>> UDP/SENDTO %v\n", addr)
}

func (ip *ipOut) tcpSendto(id uint32, message []byte, addr *net.TCPAddr) {
	fmt.Printf(">>>>> TCP/SENDTO %v\n", addr)
}

func (ip *ipOut) broadcast(id uint32, message []byte) {
	ip.Dumpf(message, "broadcast (%v bytes)", len(message))

	listener := net.ListenConfig{
		Control: func(network, address string, connection syscall.RawConn) error {
			if ip.hwif != "" {
				return conn.BindToDevice(connection, ip.hwif, conn.IsIPv4(ip.broadcastAddr.IP), ip.Conn)
			} else {
				return nil
			}
		},
	}

	if bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0"); err != nil {
		ip.Warnf("%v", err)
	} else if socket, err := listener.ListenPacket(context.Background(), "udp4", fmt.Sprintf("%v", bind)); err != nil {
		ip.Warnf("%v", err)
	} else if socket == nil {
		ip.Warnf("invalid UDP socket (%v)", socket)
	} else {
		defer socket.Close()

		if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
			ip.Warnf("%v", err)
		}

		if err := socket.SetReadDeadline(time.Now().Add(5*time.Second + ip.timeout)); err != nil {
			ip.Warnf("%v", err)
		}

		if N, err := socket.WriteTo(message, ip.broadcastAddr); err != nil {
			ip.Warnf("%v", err)
		} else {
			ip.Debugf("sent %v bytes to %v\n", N, ip.broadcastAddr)

			ctx, cancel := context.WithTimeout(ip.ctx, ip.timeout+5*time.Second)

			defer cancel()

			go func() {
				for {
					reply := make([]byte, 2048)

					if N, remote, err := socket.ReadFrom(reply); err != nil && !errors.Is(err, net.ErrClosed) {
						ip.Warnf("%v", err)
						return
					} else if err != nil {
						return
					} else {
						ip.Dumpf(reply[0:N], "received %v bytes from %v", N, remote)

						ip.ch <- protocol.Message{
							ID:      id,
							Message: reply[:N],
						}
					}
				}
			}()

			select {
			case <-time.After(ip.timeout):
				// Ok

			case <-ctx.Done():
				ip.Warnf("%v", ctx.Err())
			}
		}
	}
}

// func resolve(addr string) (any, error) {
// 	if strings.HasPrefix(addr, "tcp:") {
// 		if v, err := netip.ParseAddrPort(addr[4:]); err != nil {
// 			return nil, err
// 		} else {
// 			return net.TCPAddrFromAddrPort(v), nil
// 		}
// 	}
//
// 	if strings.HasPrefix(addr, "udp:") {
// 		if v, err := netip.ParseAddrPort(addr[4:]); err != nil {
// 			return nil, err
// 		} else {
// 			return net.UDPAddrFromAddrPort(v), nil
// 		}
// 	}
//
// 	if v, err := netip.ParseAddrPort(addr); err != nil {
// 		return nil, err
// 	} else {
// 		return net.UDPAddrFromAddrPort(v), nil
// 	}
// }
