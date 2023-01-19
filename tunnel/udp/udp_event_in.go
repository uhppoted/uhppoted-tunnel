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

type udpEventIn struct {
	conn.Conn
	hwif   string
	addr   *net.UDPAddr
	retry  conn.Backoff
	ctx    context.Context
	closed chan struct{}
}

func NewUDPEventIn(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*udpEventIn, error) {
	addr, err := net.ResolveUDPAddr("udp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve UDP address '%v'", spec)
	}

	if addr.Port == 0 {
		return nil, fmt.Errorf("UDP event requires a non-zero port")
	}

	udp := udpEventIn{
		Conn: conn.Conn{
			Tag: "UDP",
		},
		hwif:   hwif,
		addr:   addr,
		retry:  retry,
		ctx:    ctx,
		closed: make(chan struct{}),
	}

	udp.Infof("connector::udp-event-in")

	return &udp, nil
}

func (udp *udpEventIn) Close() {
	udp.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-udp.closed:
		udp.Infof("closed")

	case <-timeout.C:
		udp.Infof("close timeout")
	}
}

func (udp *udpEventIn) Run(router *router.Switch) (err error) {
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
				udp.Warnf("Failed to create UDP event socket (%v)", socket)
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

func (udp *udpEventIn) Send(id uint32, message []byte) {
}

func (udp *udpEventIn) listen(socket net.PacketConn, router *router.Switch) {
	udp.Infof("listening on %v", udp.addr)

	defer socket.Close()

	for {
		buffer := make([]byte, 2048) // NTS: buffer is handed off to router

		N, remote, err := socket.ReadFrom(buffer)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			udp.Warnf("%v", err)
		}

		if err != nil {
			return
		}

		id := protocol.NextID()
		udp.Dumpf(buffer[:N], "event %v  %v bytes from %v", id, N, remote)

		router.Received(id, buffer[:N], nil)
	}
}
