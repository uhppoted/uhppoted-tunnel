package tunnel

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/uhppoted/uhppoted-tunnel/router"
)

type tcpServer struct {
	addr        *net.TCPAddr
	connections map[net.Conn]struct{}
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
		addr:        addr,
		connections: map[net.Conn]struct{}{},
		closed:      make(chan struct{}),
	}

	return &out, nil
}

func (tcp *tcpServer) Close() {
}

func (tcp *tcpServer) Run(router *router.Switch) error {
	socket, err := net.Listen("tcp", fmt.Sprintf("%v", tcp.addr))
	if err != nil {
		return err
	}

	defer socket.Close()

	go func() {
		tcp.listen(socket, router)
	}()

	<-tcp.closed

	return nil
}

func (tcp *tcpServer) Send(id uint32, message []byte) {
	for c, _ := range tcp.connections {
		if socket, ok := c.(*net.TCPConn); ok && socket != nil {
			go func() {
				tcp.send(socket, id, message)
			}()
		}
	}
}

func (tcp *tcpServer) listen(socket net.Listener, router *router.Switch) {
	infof("TCP", "listening on %v", socket.Addr())

	for {
		if client, err := socket.Accept(); err != nil {
			errorf("TCP", "%v", err)
		} else {
			infof("TCP", "incoming connection (%v)", client.RemoteAddr())

			if socket, ok := client.(*net.TCPConn); !ok {
				errorf("TCP", "%v", "invalid TCP socket")
			} else {
				tcp.Lock()
				tcp.connections[socket] = struct{}{}
				tcp.Unlock()

				go func(socket *net.TCPConn) {
					for {
						buffer := make([]byte, 2048) // buffer is handed off to router
						if N, err := socket.Read(buffer); err != nil {
							if err == io.EOF {
								infof("TCP", "client connection %v closed ", socket.RemoteAddr())
								break
							}
							warnf("TCP", "error reading from socket (%v)", err)
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
}

func (tcp *tcpServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
	hex := dump(buffer, "                                ")
	debugf("TCP", "received %v bytes from %v\n%s\n", len(buffer), socket.RemoteAddr(), hex)

	for len(buffer) > 0 {
		id, msg, remaining := depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			tcp.send(socket, id, message)
		})
	}
}

func (tcp *tcpServer) send(conn net.Conn, id uint32, message []byte) {
	packet := packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		warnf("TCP", "msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf("TCP", "msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		infof("TCP", "msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
