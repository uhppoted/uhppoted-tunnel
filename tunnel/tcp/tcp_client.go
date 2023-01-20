package tcp

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

type tcpClient struct {
	conn.Conn
	hwif    string
	addr    *net.TCPAddr
	retry   conn.Backoff
	timeout time.Duration
	ch      chan protocol.Message
	ctx     context.Context
	closed  chan struct{}
}

func NewTCPInClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpClient, error) {
	client, err := makeTCPClient(hwif, spec, retry, ctx)

	if err == nil {
		client.Infof("connector::tcp-client-in")
	}

	return client, err
}

func NewTCPOutClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpClient, error) {
	client, err := makeTCPClient(hwif, spec, retry, ctx)

	if err == nil {
		client.Infof("connector::tcp-client-out")
	}

	return client, err
}

func makeTCPClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	in := tcpClient{
		Conn: conn.Conn{
			Tag: "TCP",
		},
		hwif:    hwif,
		addr:    addr,
		retry:   retry,
		timeout: 5 * time.Second,
		ch:      make(chan protocol.Message, 16),
		ctx:     ctx,
		closed:  make(chan struct{}),
	}

	return &in, nil
}

func (tcp *tcpClient) Close() {
	tcp.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		tcp.Infof("closed")

	case <-timeout.C:
		tcp.Infof("close timeout")
	}
}

func (tcp *tcpClient) Run(router *router.Switch) error {
	tcp.connect(router)
	tcp.closed <- struct{}{}

	return nil
}

func (tcp *tcpClient) Send(id uint32, msg []byte) {
	select {
	case tcp.ch <- protocol.Message{ID: id, Message: msg}:
	default:
	}
}

func (tcp *tcpClient) connect(router *router.Switch) {
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

		if socket, err := dialer.Dial("tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
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

func (tcp *tcpClient) listen(socket net.Conn, router *router.Switch) error {
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

func (tcp *tcpClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			tcp.send(socket, id, message)
		})
	}
}

func (tcp *tcpClient) send(conn net.Conn, id uint32, msg []byte) []byte {
	packet := protocol.Packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}

	return nil
}
