package tunnel

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/uhppoted/uhppoted-tunnel/router"
)

type udpListen struct {
	addr   *net.UDPAddr
	closed chan struct{}
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
		addr:   addr,
		closed: make(chan struct{}),
	}

	return &udp, nil
}

func (udp *udpListen) Close() {
	close(udp.closed)
}

func (udp *udpListen) Run(router *router.Switch) error {
	socket, err := net.ListenUDP("udp", udp.addr)
	if err != nil {
		return fmt.Errorf("Error creating UDP listen socket (%v)", err)
	} else if socket == nil {
		return fmt.Errorf("Failed to create UDP listen socket (%v)", socket)
	}

	defer socket.Close()

	go func() {
		udp.listen(socket, router)
	}()

	<-udp.closed

	return nil
}

func (udp *udpListen) Send(id uint32, message []byte) {
}

func (udp *udpListen) listen(socket *net.UDPConn, router *router.Switch) {
	infof("UDP", "listening on %v", udp.addr)

	for {
		buffer := make([]byte, 2048) // NTS buffer is handed off to router

		if N, remote, err := socket.ReadFromUDP(buffer); err != nil {
			if errors.Is(err, io.EOF) {
				infof("UDP", "listen socket %v closed ", socket)
			} else {
				warnf("UDP", "error reading from socket (%v)", err)
			}

			return
		} else {
			id := nextID()
			dumpf(buffer[:N], "UDP  request %v  %v bytes from %v", id, N, remote)

			h := func(reply []byte) {
				dumpf(reply, "UDP  reply %v  %v bytes for %v", id, len(reply), remote)

				if N, err := socket.WriteToUDP(reply, remote); err != nil {
					warnf("UDP", "%v", err)
				} else {
					debugf("UDP  sent %v bytes to %v\n", N, remote)
				}
			}

			router.Received(id, buffer[:N], h)
		}
	}
}
