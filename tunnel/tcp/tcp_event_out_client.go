package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tcpEventOutClient struct {
	tcpEventClient
}

func NewTCPEventOutClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpEventOutClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	tcp := tcpEventOutClient{
		tcpEventClient{
			Conn: conn.Conn{
				Tag: "TCP",
			},
			hwif:    hwif,
			addr:    addr,
			retry:   retry,
			timeout: 5 * time.Second,
			ch:      make(chan protocol.Message, 16),
			ctx:     ctx,
			closed:  make(chan struct{}),
		},
	}

	tcp.tcpEventClient.received = tcp.received
	tcp.tcpEventClient.send = tcp.send

	tcp.Infof("connector::tcp-event-out-client")

	return &tcp, nil
}

func (tcp *tcpEventOutClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	// tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())
	//
	// for len(buffer) > 0 {
	// 	id, msg, remaining := protocol.Depacketize(buffer)
	// 	buffer = remaining
	//
	// 	router.Received(id, msg, nil)
	// }
}

func (tcp *tcpEventOutClient) send(conn net.Conn, id uint32, msg []byte) {
	packet := protocol.Packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}
}
