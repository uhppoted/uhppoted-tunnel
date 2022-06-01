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
	return tcp.listen(relay)
}

func (tcp *tcpServer) Send(message []byte) []byte {
	for c, _ := range tcp.connections {
		if reply := tcp.send(c, message); reply != nil && len(reply) > 0 {
			return reply
		}
	}

	return nil
}

func (tcp *tcpServer) listen(relay func([]byte) []byte) error {
	socket, err := net.Listen("tcp", fmt.Sprintf("%v", tcp.addr))
	if err != nil {
		return err
	}

	infof("TCP  listening on %v", socket.Addr())

	for {
		if client, err := socket.Accept(); err != nil {
			errorf("%v", err)
		} else {
			infof("TCP  incoming connection (%v)", client.RemoteAddr())

			if socket, ok := client.(*net.TCPConn); !ok {
				errorf("TCP  %v", "invalid TCP socket")
			} else {
				tcp.Lock()
				tcp.connections[socket] = nil
				tcp.Unlock()

				go func(socket *net.TCPConn) {
					buffer := make([]byte, 2048)

					for {
						if N, err := socket.Read(buffer); err != nil {
							if err == io.EOF {
								infof("TCP  client connection %v closed ", socket.RemoteAddr())
								break
							}
							warnf("TCP  error reading from socket (%v)", err)

						} else if ch, ok := tcp.connections[socket]; !ok || ch == nil {
							tcp.received(buffer[:N], relay, socket)
						} else {
							select {
							case ch <- buffer[:N]:
								debugf("TCP  dispatched packet")
							default:
								debugf("TCP  dropped packet")
							}
						}
					}
				}(socket)
			}
		}
	}
}

func (tcp *tcpServer) received(packet []byte, relay func([]byte) []byte, socket net.Conn) {
	hex := dump(packet, "                                ")
	debugf("TCP  received %v bytes from %v\n%s\n", len(packet), socket.RemoteAddr(), hex)

	ix := 0
	for ix < len(packet) {
		size := uint(packet[ix])
		size <<= 8
		size += uint(packet[ix+1])

		message := depacketize(packet[ix : ix+2+int(size)])

		if reply := relay(message); reply != nil && len(reply) > 0 {
			packet := packetize(reply)

			if N, err := socket.Write(packet); err != nil {
				warnf("error relaying reply to %v (%v)", socket.RemoteAddr(), err)
			} else if N != len(packet) {
				warnf("relayed reply with %v of %v bytes to %v", N, len(reply), socket.RemoteAddr())
			} else {
				infof("relayed reply with %v bytes to %v", len(reply), socket.RemoteAddr())
			}
		}

		ix += 2 + int(size)
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
		warnf("TCP  sent %v of %v bytes to %v", N, len(message), conn.RemoteAddr())
	} else {
		infof("TCP  sent %v bytes to %v", len(message), conn.RemoteAddr())

		select {
		case <-time.After(tcp.timeout):
			infof("TCP  timeout waiting for reply from %v", conn.RemoteAddr())
			return nil

		case buffer := <-ch:
			hex := dump(buffer, "                           ")
			debugf("TCP  received %v bytes from %v\n%s\n", N, conn.RemoteAddr(), hex)

			return depacketize(buffer)
		}
	}

	return nil
}
