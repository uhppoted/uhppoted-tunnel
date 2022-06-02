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

func (udp *udpListen) Run(relay relay) error {
	router := Switch{
		relay: relay,
	}

	return udp.listen(&router)
}

func (udp *udpListen) Send(message []byte) []byte {
	return nil
}

func (udp *udpListen) listen(router *Switch) error {
	socket, err := net.ListenUDP("udp", udp.addr)
	if err != nil {
		return fmt.Errorf("Error creating UDP listen socket (%v)", err)
	} else if socket == nil {
		return fmt.Errorf("Failed to create UDP listen socket (%v)", socket)
	}

	defer socket.Close()

	infof("UDP  listening on %v", udp.addr)

	for {
		buffer := make([]byte, 2048) // NTS buffer is handed off to router
		N, remote, err := socket.ReadFromUDP(buffer)
		if err != nil {
			debugf("%v", err)
			break
		}

		udp.dump(buffer[:N], "request  %v bytes from %v", N, remote)

		h := func(reply []byte) {
			udp.dump(reply, "reply  %v bytes for %v", N, remote)

			if N, err := socket.WriteToUDP(reply, remote); err != nil {
				warnf("%v", err)
			} else {
				debugf("UDP sent %v bytes to %v\n", N, remote)
			}
		}

		router.received(nextID(), buffer[:N], h)
	}

	return nil
}

func (udp *udpListen) dump(message []byte, format string, args ...any) {
	hex := dump(message, "                                ")
	preamble := fmt.Sprintf(format, args...)

	debugf("UDP  %v\n%s", preamble, hex)
}
