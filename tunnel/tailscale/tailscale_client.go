package tailscale

import (
	"context"
	"errors"
	"fmt"
	"net"
	// "syscall"
	"time"

	"tailscale.com/tsnet"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tailscaleClient struct {
	conn.Conn
	hwif    string
	addr    string
	retry   conn.Backoff
	timeout time.Duration
	ch      chan protocol.Message
	ctx     context.Context
	closed  chan struct{}
}

// func NewTCPInClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpClient, error) {
//     client, err := makeTCPClient(hwif, spec, retry, ctx)
//
//     if err == nil {
//         client.Infof("connector::tcp-client-in")
//     }
//
//     return client, err
// }

func NewTailscaleOutClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tailscaleClient, error) {
	client, err := makeTailscaleClient(hwif, spec, retry, ctx)

	if err == nil {
		client.Infof("connector::tailscale-client-out")
	}

	return client, err
}

func makeTailscaleClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tailscaleClient, error) {
	// addr, err := net.ResolveTCPAddr("tcp", spec)
	// if err != nil {
	//     return nil, err
	// } else if addr == nil {
	//     return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	// }

	in := tailscaleClient{
		Conn: conn.Conn{
			Tag: "tailscale",
		},
		hwif:    hwif,
		addr:    "uhppoted:12345",
		retry:   retry,
		timeout: 5 * time.Second,
		ch:      make(chan protocol.Message, 16),
		ctx:     ctx,
		closed:  make(chan struct{}),
	}

	return &in, nil
}

func (ts *tailscaleClient) Close() {
	ts.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-ts.closed:
		ts.Infof("closed")

	case <-timeout.C:
		ts.Infof("close timeout")
	}
}

func (ts *tailscaleClient) Run(router *router.Switch) error {
	ts.connect(router)
	ts.closed <- struct{}{}

	return nil
}

func (ts *tailscaleClient) Send(id uint32, msg []byte) {
	select {
	case ts.ch <- protocol.Message{ID: id, Message: msg}:
	default:
	}
}

func (ts *tailscaleClient) connect(router *router.Switch) {
	server := &tsnet.Server{
		Logf: func(string, ...any) {},
		Dir:  "../runtime/uhppoted-tunnel/tailscale/client",
	}

	server.Hostname = "uhppoted-client"
	server.Dir = "../runtime/uhppoted-tunnel/tailscale/client"

	defer server.Close()
	println("woot")

	for {
		ts.Infof("connecting to %v", ts.addr)

		// 	dialer := &net.Dialer{
		// 		Timeout: ts.timeout,
		// 		Control: func(network, address string, connection syscall.RawConn) error {
		// 			if ts.hwif != "" {
		// 				return conn.BindToDevice(connection, ts.hwif, conn.IsIPv4(ts.addr.IP), ts.Conn)
		// 			} else {
		// 				return nil
		// 			}
		// 		},
		// 	}
		//
		// 	if socket, err := dialer.Dial("tcp", fmt.Sprintf("%v", ts.addr)); err != nil {

		if socket, err := server.Dial(context.Background(), "tcp", fmt.Sprintf("%v", ts.addr)); err != nil {
			ts.Warnf("%v", err)
		} else if socket == nil {
			ts.Warnf("connect %v failed (%v)", ts.addr, socket)
		} else {
			ts.retry.Reset()
			eof := make(chan struct{})

			go func() {
				for {
					select {
					case msg := <-ts.ch:
						ts.Infof("msg %v  relaying to %v", msg.ID, socket.RemoteAddr())
						ts.send(socket, msg.ID, msg.Message)

					case <-eof:
						return

					case <-ts.ctx.Done():
						socket.Close()
						return
					}

					println("woot/sleep")
					time.Sleep(1000)
				}
			}()

			if err := ts.listen(socket, router); err != nil && !errors.Is(err, net.ErrClosed) {
				ts.Warnf("%v", err)
			}

			close(eof)

			time.Sleep(5000)
		}

		if !ts.retry.Wait(ts.Tag) {
			return
		}
	}
}

func (ts *tailscaleClient) listen(socket net.Conn, router *router.Switch) error {
	ts.Infof("connected  to %v", socket.RemoteAddr())

	defer socket.Close()

	for {
		buffer := make([]byte, 2048)
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		go func() {
			ts.received(buffer[:N], router, socket)
		}()
	}
}

func (ts *tailscaleClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	ts.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			ts.send(socket, id, message)
		})
	}
}

func (ts *tailscaleClient) send(conn net.Conn, id uint32, msg []byte) []byte {
	packet := protocol.Packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		ts.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		ts.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		ts.Infof("msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}

	return nil
}
