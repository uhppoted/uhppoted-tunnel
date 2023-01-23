package tls

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tlsEventClient struct {
	conn.Conn
	hwif    string
	addr    *net.TCPAddr
	config  *tls.Config
	retry   conn.Backoff
	timeout time.Duration
	ch      chan protocol.Message
	ctx     context.Context
	closed  chan struct{}

	received func([]byte, *router.Switch, net.Conn)
	send     func(net.Conn, uint32, []byte)
}

func (tcp *tlsEventClient) Close() {
	tcp.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		tcp.Infof("closed")

	case <-timeout.C:
		tcp.Infof("close timeout")
	}
}

func (tcp *tlsEventClient) Run(router *router.Switch) error {
	tcp.connect(router)
	tcp.closed <- struct{}{}

	return nil
}

func (tcp *tlsEventClient) Send(id uint32, msg []byte) {
	select {
	case tcp.ch <- protocol.Message{ID: id, Message: msg}:
	default:
	}
}

func (tcp *tlsEventClient) connect(router *router.Switch) {
	for {
		tcp.Infof("connecting to %v", tcp.addr)

		dialer := &net.Dialer{
			Control: func(network, address string, connection syscall.RawConn) error {
				if tcp.hwif != "" {
					return conn.BindToDevice(connection, tcp.hwif, conn.IsIPv4(tcp.addr.IP), tcp.Conn)
				} else {
					return nil
				}
			},
		}

		if socket, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%v", tcp.addr), tcp.config); err != nil {
			tcp.Warnf("%v", err)
		} else if socket == nil {
			tcp.Warnf("connect %v failed (%v)", tcp.addr, socket)
		} else {
			tcp.retry.Reset()
			eof := make(chan struct{})

			go func() {
				for {
					select {
					case msg := <-tcp.ch:
						tcp.Infof("msg %v  relaying to %v", msg.ID, socket.RemoteAddr())
						tcp.send(socket, msg.ID, msg.Message)

					case <-eof:
						return

					case <-tcp.ctx.Done():
						socket.Close()
						return
					}
				}
			}()

			if err := tcp.listen(socket, router); err != nil && !errors.Is(err, net.ErrClosed) {
				tcp.Warnf("%v", err)
			}

			close(eof)
		}

		if !tcp.retry.Wait(tcp.Tag) {
			return
		}
	}
}

func (tcp *tlsEventClient) listen(socket net.Conn, router *router.Switch) error {
	tcp.Infof("connected  to %v", socket.RemoteAddr())

	defer socket.Close()

	for {
		buffer := make([]byte, 2048)
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		go func() {
			tcp.received(buffer[:N], router, socket)
		}()
	}
}
