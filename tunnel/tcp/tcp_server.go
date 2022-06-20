package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
)

type tcpServer struct {
	tag         string
	addr        *net.TCPAddr
	retryDelay  time.Duration
	connections map[net.Conn]struct{}
	closing     chan struct{}
	closed      chan struct{}
	sync.RWMutex
}

func NewTCPServer(spec string) (*tcpServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)

	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	} else if addr.Port == 0 {
		return nil, fmt.Errorf("TCP host requires a non-zero port")
	}

	out := tcpServer{
		tag:         "TCP",
		addr:        addr,
		retryDelay:  15 * time.Second,
		connections: map[net.Conn]struct{}{},
		closing:     make(chan struct{}),
		closed:      make(chan struct{}),
	}

	return &out, nil
}

func (tcp *tcpServer) Close() {
	infof(tcp.tag, "closing")
	close(tcp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		infof(tcp.tag, "closed")

	case <-timeout.C:
		infof(tcp.tag, "close timeout")
	}
}

func (tcp *tcpServer) Run(router *router.Switch) (err error) {
	var socket net.Listener
	var closing = false
	var delay = 0 * time.Second

	go func() {
		for !closing {
			time.Sleep(delay)

			socket, err = net.Listen("tcp", fmt.Sprintf("%v", tcp.addr))
			if err != nil {
				warnf(tcp.tag, "%v", err)
			} else if socket == nil {
				warnf(tcp.tag, "%v", fmt.Errorf("Failed to create TCP listen socket (%v)", socket))
			}

			delay = tcp.retryDelay

			tcp.listen(socket, router)
		}

		for k, _ := range tcp.connections {
			k.Close()
		}

		tcp.closed <- struct{}{}
	}()

	<-tcp.closing

	closing = true
	socket.Close()

	return nil
}

func (tcp *tcpServer) Send(id uint32, message []byte) {
	for c, _ := range tcp.connections {
		go func() {
			tcp.send(c, id, message)
		}()
	}
}

func (tcp *tcpServer) listen(socket net.Listener, router *router.Switch) {
	infof(tcp.tag, "listening on %v", socket.Addr())

	defer socket.Close()

	for {
		client, err := socket.Accept()
		if err != nil {
			errorf(tcp.tag, "%v", err)
			return
		}

		infof(tcp.tag, "incoming connection (%v)", client.RemoteAddr())

		if socket, ok := client.(*net.TCPConn); !ok {
			warnf(tcp.tag, "invalid TCP socket (%v)", socket)
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
							infof(tcp.tag, "client connection %v closed ", socket.RemoteAddr())
						} else {
							warnf(tcp.tag, "%v", err)
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
	dumpf(tcp.tag, buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

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
		warnf(tcp.tag, "msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf(tcp.tag, "msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		infof(tcp.tag, "msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
