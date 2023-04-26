package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tcpServer struct {
	conn.Conn
	hwif        string
	addr        *net.TCPAddr
	retry       conn.Backoff
	connections map[net.Conn]struct{}
	ctx         context.Context
	closed      chan struct{}
	sync.RWMutex
}

func NewTCPInServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpServer, error) {
	server, err := makeTCPServer(hwif, spec, retry, ctx)

	if err == nil {
		server.Infof("connector::tcp-server-in")
	}

	return server, err
}

func NewTCPOutServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpServer, error) {
	server, err := makeTCPServer(hwif, spec, retry, ctx)

	if err == nil {
		server.Infof("connector::tcp-server-out")
	}

	return server, err
}

func makeTCPServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)

	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	} else if addr.Port == 0 {
		return nil, fmt.Errorf("TCP host requires a non-zero port")
	}

	tcp := tcpServer{
		Conn: conn.Conn{
			Tag: "TCP",
		},
		hwif:        hwif,
		addr:        addr,
		retry:       retry,
		connections: map[net.Conn]struct{}{},
		ctx:         ctx,
		closed:      make(chan struct{}),
	}

	return &tcp, nil
}

func (tcp *tcpServer) Close() {
	tcp.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		tcp.Infof("closed")

	case <-timeout.C:
		tcp.Infof("close timeout")
	}
}

func (tcp *tcpServer) Run(router *router.Switch) (err error) {
	closing := false
	sockets := conn.NewSocketList()

	defer sockets.CloseAll()

	go func() {
	loop:
		for {

			listener := net.ListenConfig{
				Control: func(network, address string, connection syscall.RawConn) error {
					if tcp.hwif != "" {
						return conn.BindToDevice(connection, tcp.hwif, conn.IsIPv4(tcp.addr.IP), tcp.Conn)
					} else {
						return nil
					}
				},
			}

			if socket, err := listener.Listen(context.Background(), "tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
				tcp.Warnf("%v", err)
			} else if socket == nil {
				tcp.Warnf("%v", fmt.Errorf("failed to create TCP listen socket (%v)", socket))
			} else {
				sockets.Add(socket)
				tcp.retry.Reset()
				tcp.listen(socket, router)
				sockets.Close(socket)
			}

			if closing || !tcp.retry.Wait(tcp.Tag) {
				break loop
			}
		}

		for k := range tcp.connections {
			k.Close()
		}

		tcp.closed <- struct{}{}
	}()

	<-tcp.ctx.Done()

	closing = true

	return nil
}

func (tcp *tcpServer) Send(id uint32, message []byte) {
	for c := range tcp.connections {
		go func(conn net.Conn) {
			tcp.send(conn, id, message)
		}(c)
	}
}

func (tcp *tcpServer) listen(socket net.Listener, router *router.Switch) {
	tcp.Infof("listening on %v", socket.Addr())

	defer socket.Close()

	for {
		client, err := socket.Accept()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			tcp.Errorf("%v %v", err, errors.Is(err, net.ErrClosed))
		}

		if err != nil {
			return
		}

		tcp.Infof("incoming connection (%v)", client.RemoteAddr())

		if socket, ok := client.(*net.TCPConn); !ok {
			tcp.Warnf("invalid TCP socket (%v)", socket)
			client.Close()
		} else {
			tcp.Lock()
			tcp.connections[socket] = struct{}{}
			tcp.Unlock()

			go func(socket *net.TCPConn) {
				for {
					buffer := make([]byte, 2048) // buffer is handed off to router
					if N, err := socket.Read(buffer); err != nil {
						if err == io.EOF {
							tcp.Infof("client connection %v closed ", socket.RemoteAddr())
						} else {
							tcp.Warnf("%v", err)
						}
						break
					} else {
						tcp.received(buffer[:N], router, socket)
					}
				}

				tcp.Lock()
				delete(tcp.connections, socket)
				tcp.Unlock()
			}(socket)
		}
	}
}

func (tcp *tcpServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
	tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			tcp.send(socket, id, message)
		})
	}
}

func (tcp *tcpServer) send(conn net.Conn, id uint32, message []byte) {
	packet := protocol.Packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
