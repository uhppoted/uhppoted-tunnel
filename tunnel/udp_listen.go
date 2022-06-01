package tunnel

import (
	"fmt"
	"net"
)

type udpListen struct {
	addr *net.UDPAddr
}

func NewUDPListen(spec string) (*udpListen, error) {
	addr, err := net.ResolveUDPAddr("udp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve UDP address '%v'", spec)
	}

	if addr.Port == 0 {
		return nil, fmt.Errorf("UDP listen requires a non-zero port")
	}

	udp := udpListen{
		addr: addr,
	}

	return &udp, nil
}

func (udp *udpListen) Close() {
}

func (udp *udpListen) Run(relay func([]byte) []byte) error {
	return udp.listen(relay)
}

func (udp *udpListen) Send(message []byte) []byte {
	return nil
}

func (udp *udpListen) listen(relay func([]byte) []byte) error {
	socket, err := net.ListenUDP("udp", udp.addr)
	if err != nil {
		return fmt.Errorf("Error creating UDP listen socket (%v)", err)
	} else if socket == nil {
		return fmt.Errorf("Failed to create UDP listen socket (%v)", socket)
	}

	defer socket.Close()

	infof("UDP  listening on %v", udp.addr)

	buffer := make([]byte, 2048)

	for {
		N, remote, err := socket.ReadFromUDP(buffer)
		if err != nil {
			debugf("%v", err)
			break
		}

		hex := dump(buffer[:N], "                                ")
		debugf("UDP  received %v bytes from %v\n%s\n", N, remote, hex)

		if reply := relay(buffer[:N]); reply != nil {
			if N, err := socket.WriteToUDP(reply, remote); err != nil {
				warnf("%v", err)
			} else {
				debugf("UDP sent %v bytes to %v\n", N, remote)
			}
		}
	}

	return nil
}
