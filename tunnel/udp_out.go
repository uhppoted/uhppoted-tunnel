package tunnel

import (
	"fmt"
	"net"
	"time"
)

type udpOut struct {
	addr *net.UDPAddr
}

func NewUDPOut(spec string) (*udpOut, error) {
	addr, err := net.ResolveUDPAddr("udp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve UDP address '%v'", spec)
	}

	if addr.Port == 0 {
		return nil, fmt.Errorf("UDP requires a non-zero port")
	}

	out := udpOut{
		addr: addr,
	}

	return &out, nil
}

func (udp *udpOut) Listen() error {
	ch := make(chan bool)

	<-ch

	return nil
}

func (udp *udpOut) Close() {
}

func (udp *udpOut) Send(message []byte) []byte {
	hex := dump(message, "                           ")
	debugf("broadcast%v\n%s\n", "", hex)

	if bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:0"); err != nil {
		warnf("%v", err)
	} else if socket, err := net.ListenUDP("udp", bind); err != nil {
		warnf("%v", err)
	} else if socket == nil {
		warnf("invalid UDP socket (%v)", socket)
	} else {
		defer socket.Close()

		if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
			warnf("%v", err)
		}

		if err := socket.SetReadDeadline(time.Now().Add(5000 * time.Millisecond)); err != nil {
			warnf("%v", err)
		}

		if N, err := socket.WriteToUDP(message, udp.addr); err != nil {
			warnf("%v", err)
		} else {
			debugf(" ... sent %v bytes to %v\n", N, udp.addr)

			reply := make([]byte, 2048)

			if N, remote, err := socket.ReadFromUDP(reply); err != nil {
				warnf("%v", err)
			} else {
				debugf(" ... received %v bytes from %v\n%s", N, remote, dump(reply[:N], " ...          "))

				return reply[:N]
			}
		}
	}

	return nil
}
