package tunnel

import (
	"fmt"
	"net"
)

type tcpInClient struct {
	addr *net.TCPAddr
}

func NewTCPIn(spec string) (*tcpInClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	in := tcpInClient{
		addr: addr,
	}

	return &in, nil
}

func (tcp *tcpInClient) Listen(relay func([]byte) []byte) error {
	socket, err := net.Dial("tcp", fmt.Sprintf("%v", tcp.addr))
	if err != nil {
		return fmt.Errorf("Error connecting to  %v (%v)", tcp.addr, err)
	} else if socket == nil {
		return fmt.Errorf("Failed to create TCP connection to %v (%v)", tcp.addr, socket)
	}

	defer socket.Close()

	infof("TCP  connected to %v", tcp.addr)

	buffer := make([]byte, 2048)

	for {
		if N, err := socket.Read(buffer); err != nil {
			warnf("%v", err)
			break
		} else {
			hex := dump(buffer[:N], "                           ")
			debugf("TCP  received %v bytes from %v\n%s\n", N, socket.RemoteAddr(), hex)

			ix := 0
			for ix < N {
				size := uint(buffer[ix])
				size <<= 8
				size += uint(buffer[ix+1])

				message := depacketize(buffer[ix : ix+2+int(size)])

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
	}

	return nil
}

func (tcp *tcpInClient) Close() {

}
