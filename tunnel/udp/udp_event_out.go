package udp

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type udpEventOut struct {
	conn.Conn
	hwif   string
	addr   *net.UDPAddr
	ctx    context.Context
	ch     chan protocol.Message
	closed chan struct{}
}

func NewUDPEventOut(hwif string, spec string, ctx context.Context) (*udpEventOut, error) {
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

	udp := udpEventOut{
		Conn: conn.Conn{
			Tag: "UDP",
		},
		hwif:   hwif,
		addr:   addr,
		ctx:    ctx,
		ch:     make(chan protocol.Message),
		closed: make(chan struct{}),
	}

	udp.Infof("connector::udp-event-out")

	return &udp, nil
}

func (udp *udpEventOut) Close() {
	udp.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-udp.closed:
		udp.Infof("closed")

	case <-timeout.C:
		udp.Infof("close timeout")
	}
}

func (udp *udpEventOut) Run(router *router.Switch) error {
loop:
	for {
		select {
		case msg := <-udp.ch:
			router.Received(msg.ID, msg.Message, nil)

		case <-udp.ctx.Done():
			break loop
		}
	}

	close(udp.closed)

	return nil
}

func (udp *udpEventOut) Send(id uint32, msg []byte) {
	go func() {
		udp.send(id, msg)
	}()
}

func (udp *udpEventOut) send(id uint32, message []byte) {
	udp.Dumpf(message, "event/out (%v bytes)", len(message))

	dialer := &net.Dialer{
		Control: func(network, address string, connection syscall.RawConn) error {
			if udp.hwif != "" {
				return conn.BindToDevice(connection, udp.hwif, conn.IsIPv4(udp.addr.IP), udp.Conn)
			} else {
				return nil
			}
		},
	}

	if socket, err := dialer.Dial("udp", fmt.Sprintf("%v", udp.addr)); err != nil {
		udp.Warnf("%v", err)
	} else if socket == nil {
		udp.Warnf("invalid UDP socket (%v)", socket)
	} else {
		defer socket.Close()

		if err := socket.SetWriteDeadline(time.Now().Add(1000 * time.Millisecond)); err != nil {
			udp.Warnf("%v", err)
		}

		if N, err := socket.Write(message); err != nil {
			udp.Warnf("%v", err)
		} else {
			udp.Debugf("sent %v bytes to %v\n", N, udp.addr)
		}
	}
}
