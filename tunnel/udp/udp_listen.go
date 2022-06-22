package udp

import (
	"fmt"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type udpListen struct {
	conn.Conn
	addr    *net.UDPAddr
	retry   conn.Backoff
	closing chan struct{}
	closed  chan struct{}
}

func NewUDPListen(spec string, retry conn.Backoff) (*udpListen, error) {
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
		Conn: conn.Conn{
			Tag: "UDP",
		},
		addr:    addr,
		retry:   retry,
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	return &udp, nil
}

func (udp *udpListen) Close() {
	udp.Infof("closing")
	close(udp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-udp.closed:
		udp.Infof("closed")

	case <-timeout.C:
		udp.Infof("close timeout")
	}
}

func (udp *udpListen) Run(router *router.Switch) (err error) {
	var socket *net.UDPConn
	var closing = false

	go func() {
	loop:
		for !closing {
			socket, err = net.ListenUDP("udp", udp.addr)
			if err != nil {
				udp.Warnf("%v", err)
			} else if socket == nil {
				udp.Warnf("Failed to create UDP listen socket (%v)", socket)
			} else {
				udp.retry.Reset()
				udp.listen(socket, router)
			}

			if !udp.retry.Wait(udp.Tag, udp.closing) {
				break loop
			}
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
	udp.Infof("listening on %v", udp.addr)

	defer socket.Close()

	for {
		buffer := make([]byte, 2048) // NTS buffer is handed off to router

		N, remote, err := socket.ReadFromUDP(buffer)
		if err != nil {
			udp.Warnf("%v", err)
			return
		}

		id := protocol.NextID()
		udp.Dumpf(buffer[:N], "request %v  %v bytes from %v", id, N, remote)

		h := func(reply []byte) {
			udp.Dumpf(reply, "reply %v  %v bytes for %v", id, len(reply), remote)

			if N, err := socket.WriteToUDP(reply, remote); err != nil {
				udp.Warnf("%v", err)
			} else {
				udp.Debugf("sent %v bytes to %v\n", N, remote)
			}
		}

		router.Received(id, buffer[:N], h)
	}
}
