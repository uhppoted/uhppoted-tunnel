package tunnel

import (
	"fmt"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/router"
)

type udpListen struct {
	addr       *net.UDPAddr
	retryDelay time.Duration
	closing    chan struct{}
	closed     chan struct{}
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
		addr:       addr,
		retryDelay: 15 * time.Second,
		closing:    make(chan struct{}),
		closed:     make(chan struct{}),
	}

	return &udp, nil
}

func (udp *udpListen) Close() {
	infof("UDP", "closing")
	close(udp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-udp.closed:
		infof("UDP", "closed")

	case <-timeout.C:
		infof("UDP", "close timeout")
	}
}

func (udp *udpListen) Run(router *router.Switch) (err error) {
	var socket *net.UDPConn
	var closing = false
	var delay = 0 * time.Second

	go func() {
		for !closing {
			time.Sleep(delay)

			socket, err = net.ListenUDP("udp", udp.addr)
			if err != nil {
				return
			} else if socket == nil {
				err = fmt.Errorf("Failed to create UDP listen socket (%v)", socket)
				return
			}

			delay = udp.retryDelay

			udp.listen(socket, router)
		}

		udp.closed <- struct{}{}
	}()

	<-udp.closing

	closing = true
	socket.Close()

	return nil
}

func (udp *udpListen) Send(id uint32, message []byte) {
}

func (udp *udpListen) listen(socket *net.UDPConn, router *router.Switch) {
	infof("UDP", "listening on %v", udp.addr)

	defer socket.Close()

	for {
		buffer := make([]byte, 2048) // NTS buffer is handed off to router

		N, remote, err := socket.ReadFromUDP(buffer)
		if err != nil {
			warnf("UDP", "%v", err)
			return
		}

		id := nextID()
		dumpf("UDP", buffer[:N], "request %v  %v bytes from %v", id, N, remote)

		h := func(reply []byte) {
			dumpf("UDP", reply, "reply %v  %v bytes for %v", id, len(reply), remote)

			if N, err := socket.WriteToUDP(reply, remote); err != nil {
				warnf("UDP", "%v", err)
			} else {
				debugf("UDP", "sent %v bytes to %v\n", N, remote)
			}
		}

		router.Received(id, buffer[:N], h)
	}
}
