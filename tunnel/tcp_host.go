package tunnel

import (
	"fmt"
	"net"
)

type tcpOutHost struct {
	addr        *net.TCPAddr
	connections map[net.Conn]struct{}
}

func NewTCPOutHost(spec string) (*tcpOutHost, error) {
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

	out := tcpOutHost{
		addr:        addr,
		connections: map[net.Conn]struct{}{},
	}

	return &out, nil
}

func (tcp *tcpOutHost) Listen() error {
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
				errorf("%v", "invalid TCP socket")
			} else {
				tcp.connections[socket] = struct{}{}
			}
		}
	}
}

func (tcp *tcpOutHost) Close() {
}

func (tcp *tcpOutHost) Send(message []byte) []byte {
	packet := packetize(message)

	for c, _ := range tcp.connections {
		if N, err := c.Write(packet); err != nil {
			warnf("error sending message to %v (%v)", c.RemoteAddr(), err)
		} else if N != len(packet) {
			warnf("TCP/out  sent %v of %v bytes to %v", N, len(message), c.RemoteAddr())
		} else {
			infof("TCP/out  sent %v bytes to %v", len(message), c.RemoteAddr())
			buffer := make([]byte, 2048)

			if N, err := c.Read(buffer); err != nil {
				warnf("%v", err)
			} else {
				hex := dump(buffer[:N], "                           ")
				debugf("TCP/out  received %v bytes from %v\n%s\n", N, c.RemoteAddr(), hex)

				size := uint(buffer[0])
				size <<= 8
				size += uint(buffer[1])

				return depacketize(buffer)
			}
		}
	}

	return nil
}
