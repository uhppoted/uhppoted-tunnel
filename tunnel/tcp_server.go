package tunnel

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type tcpServer struct {
	addr        *net.TCPAddr
	timeout     time.Duration
	connections map[net.Conn]chan []byte
	sync.RWMutex
}

func NewTCPServer(spec string) (*tcpServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	if addr.Port == 0 {
		return nil, fmt.Errorf("TCP host requires a non-zero port")
	}

	out := tcpServer{
		addr:        addr,
		timeout:     5 * time.Second,
		connections: map[net.Conn]chan []byte{},
	}

	return &out, nil
}

func (tcp *tcpServer) Close() {
}

func (tcp *tcpServer) Run(relay func([]byte) []byte) error {
	return tcp.listen()
}

func (tcp *tcpServer) Send(message []byte) []byte {
	for c, _ := range tcp.connections {
		if reply := tcp.send(c, message); reply != nil && len(reply) > 0 {
			return reply
		}
	}

	return nil
}

func (tcp *tcpServer) listen() error {
	socket, err := net.Listen("tcp", fmt.Sprintf("%v", tcp.addr))
	if err != nil {
		return err
	}

	infof("TCP/out  listening on %v", socket.Addr())

	for {
		if client, err := socket.Accept(); err != nil {
			errorf("%v", err)
		} else {
			infof("TCP/out  incoming connection (%v)", client.RemoteAddr())

			if socket, ok := client.(*net.TCPConn); !ok {
				errorf("TCP/out  %v", "invalid TCP socket")
			} else {
				tcp.Lock()
				tcp.connections[socket] = nil
				tcp.Unlock()

				go func(socket *net.TCPConn) {
					buffer := make([]byte, 2048)

					for {
						if N, err := socket.Read(buffer); err != nil {
							if err == io.EOF {
								infof("TCP/out  client connection %v closed ", socket.RemoteAddr())
								break
							}
							warnf("TCP/out  error reading from socket (%v)", err)

						} else if ch, ok := tcp.connections[socket]; !ok || ch == nil {
							warnf("TCP/out  discarding %v byte packet from %v", N, socket.RemoteAddr())
						} else {
							select {
							case ch <- buffer[:N]:
								debugf("TCP/out  dispatched packet")
							default:
								debugf("TCP/out  dropped packet")
							}
						}
					}
				}(socket)
			}
		}
	}
}

func (tcp *tcpServer) send(conn net.Conn, message []byte) []byte {
	ch := make(chan []byte)

	tcp.Lock()
	tcp.connections[conn] = ch
	tcp.Unlock()

	defer func() {
		tcp.Lock()
		tcp.connections[conn] = nil
		tcp.Unlock()

		close(ch)
	}()

	packet := packetize(message)

	if N, err := conn.Write(packet); err != nil {
		warnf("error sending message to %v (%v)", conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf("TCP/out  sent %v of %v bytes to %v", N, len(message), conn.RemoteAddr())
	} else {
		infof("TCP/out  sent %v bytes to %v", len(message), conn.RemoteAddr())

		select {
		case <-time.After(tcp.timeout):
			infof("TCP/out  timeout waiting for reply from %v", conn.RemoteAddr())
			return nil

		case buffer := <-ch:
			hex := dump(buffer, "                           ")
			debugf("TCP/out  received %v bytes from %v\n%s\n", N, conn.RemoteAddr(), hex)

			return depacketize(buffer)
		}
	}

	return nil
}
