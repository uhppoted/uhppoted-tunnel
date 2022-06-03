package tunnel

import (
	"fmt"
	"net"
	"time"
)

type udpBroadcast struct {
	addr *net.UDPAddr
	ch   chan message
}

type message struct {
	id      uint32
	message []byte
}

func NewUDPBroadcast(spec string) (*udpBroadcast, error) {
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

	out := udpBroadcast{
		addr: addr,
		ch:   make(chan message),
	}

	return &out, nil
}

func (udp *udpBroadcast) Close() {
}

func (udp *udpBroadcast) Run(relay relay) error {
	router := Switch{
		relay: relay,
	}

	for {
		select {
		case msg := <-udp.ch:
			router.reply(msg.id, msg.message)
		}
	}
}

func (udp *udpBroadcast) Send(id uint32, msg []byte) {
	go func() {
		if reply := udp.send(id, msg); reply != nil {
			udp.ch <- message{
				id:      id,
				message: reply,
			}
		}
	}()
}

func (udp *udpBroadcast) send(id uint32, message []byte) []byte {
	hex := dump(message, "                                ")
	debugf("UDP  broadcast%v\n%s\n", "", hex)

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
			debugf("UDP  sent %v bytes to %v\n", N, udp.addr)

			reply := make([]byte, 2048)

			if N, remote, err := socket.ReadFromUDP(reply); err != nil {
				warnf("%v", err)
			} else {
				hex := dump(reply[:N], "                                ")
				debugf("UDP  received %v bytes from %v\n%s", N, remote, hex)

				return reply[:N]
			}
		}
	}

	return nil
}
