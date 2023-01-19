package udp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type udpListen struct {
	conn.Conn
	hwif   string
	addr   *net.UDPAddr
	retry  conn.Backoff
	ctx    context.Context
	closed chan struct{}
}

func NewUDPListen(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*udpListen, error) {
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
		hwif:   hwif,
		addr:   addr,
		retry:  retry,
		ctx:    ctx,
		closed: make(chan struct{}),
	}

	udp.Infof("connector::udp-listen")

	return &udp, nil
}

func (udp *udpListen) Close() {
	udp.Infof("closing")

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

	listener := net.ListenConfig{
		Control: func(network, address string, connection syscall.RawConn) error {
			if udp.hwif != "" {
				return conn.BindToDevice(connection, udp.hwif, conn.IsIPv4(udp.addr.IP), udp.Conn)
			} else {
				return nil
			}
		},
	}

	go func() {
	loop:
		for {
			socket, err := listener.ListenPacket(context.Background(), "udp4", fmt.Sprintf("%v", udp.addr))
			if err != nil {
				udp.Warnf("%v", err)
			} else if socket == nil {
				udp.Warnf("Failed to create UDP listen socket (%v)", socket)
			} else {
				udp.retry.Reset()
				udp.listen(socket, router)
			}

			if closing || !udp.retry.Wait(udp.Tag) {
				break loop
			}
		}

		udp.closed <- struct{}{}
	}()

	<-udp.ctx.Done()

	closing = true
	socket.Close()

	return nil
}

func (udp *udpListen) Send(id uint32, message []byte) {
}

func (udp *udpListen) listen(socket net.PacketConn, router *router.Switch) {
	udp.Infof("listening on %v", udp.addr)

	defer socket.Close()

	for {
		buffer := make([]byte, 2048) // NTS buffer is handed off to router

		N, remote, err := socket.ReadFrom(buffer)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			udp.Warnf("%v", err)
		}

		if err != nil {
			return
		}

		id := protocol.NextID()
		udp.Dumpf(buffer[:N], "request %v  %v bytes from %v", id, N, remote)

		h := func(reply []byte) {
			udp.Dumpf(reply, "reply %v  %v bytes for %v", id, len(reply), remote)

			if N, err := socket.WriteTo(reply, remote); err != nil {
				udp.Warnf("%v", err)
			} else {
				udp.Debugf("sent %v bytes to %v\n", N, remote)
			}
		}

		router.Received(id, buffer[:N], h)
	}
}
