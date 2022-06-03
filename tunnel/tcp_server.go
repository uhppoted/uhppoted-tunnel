package tunnel

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type tcpServer struct {
	addr        *net.TCPAddr
	connections map[net.Conn]struct{}
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
	}

	return &out, nil
}

func (tcp *tcpServer) Close() {
}

func (tcp *tcpServer) Run(relay relay) error {
	router := Switch{
		relay: relay,
	}

	return tcp.listen(&router)
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

func (tcp *tcpServer) listen(router *Switch) error {
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
				tcp.connections[socket] = struct{}{}
				tcp.Unlock()

				go func(socket *net.TCPConn) {
					for {
						buffer := make([]byte, 2048) // buffer is handed off to router
						if N, err := socket.Read(buffer); err != nil {
							if err == io.EOF {
								infof("TCP  client connection %v closed ", socket.RemoteAddr())
								break
							}
							warnf("TCP  error reading from socket (%v)", err)
						} else {
							tcp.received(buffer[:N], router, socket)
						}
					}
				}(socket)
			}
		}
	}
}

func (tcp *tcpServer) received(packet []byte, router *Switch, socket net.Conn) {
	hex := dump(packet, "                                ")
	debugf("TCP  received %v bytes from %v\n%s\n", len(packet), socket.RemoteAddr(), hex)

	ix := 0
	for ix < len(packet) {
		size := uint(packet[ix])
		size <<= 8
		size += uint(packet[ix+1])

		id, message := depacketize(packet[ix:])

		router.reply(id, message)

		// if reply := relay(id, message); reply != nil && len(reply) > 0 {
		// 	packet := packetize(id, reply)

		// 	if N, err := socket.Write(packet); err != nil {
		// 		warnf("error relaying reply to %v (%v)", socket.RemoteAddr(), err)
		// 	} else if N != len(packet) {
		// 		warnf("relayed reply with %v of %v bytes to %v", N, len(reply), socket.RemoteAddr())
		// 	} else {
		// 		infof("relayed reply with %v bytes to %v", len(reply), socket.RemoteAddr())
		// 	}
		// }

		ix += 6 + int(size)
	}
}

func (tcp *tcpServer) send(conn *net.TCPConn, id uint32, message []byte) {
	packet := packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		warnf("error sending message to %v (%v)", conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf("TCP  sent %v of %v bytes to %v", N, len(message), conn.RemoteAddr())
	} else {
		infof("TCP  sent %v bytes to %v", len(message), conn.RemoteAddr())
	}
}
