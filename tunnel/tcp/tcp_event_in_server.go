package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tcpEventIn struct {
	tcpEventServer
}

func NewTCPEventInServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpEventIn, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)

	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	} else if addr.Port == 0 {
		return nil, fmt.Errorf("TCP host requires a non-zero port")
	}

	tcp := tcpEventIn{
		tcpEventServer{
			Conn: conn.Conn{
				Tag: "TCP",
			},
			hwif:        hwif,
			addr:        addr,
			retry:       retry,
			connections: map[net.Conn]struct{}{},
			ctx:         ctx,
			closed:      make(chan struct{}),
		},
	}

	tcp.tcpEventServer.received = tcp.received

	tcp.Infof("connector::tcp-event-in-server")

	return &tcp, nil
}

func (tcp *tcpEventIn) Send(id uint32, message []byte) {
}

func (tcp *tcpEventIn) received(buffer []byte, router *router.Switch, socket net.Conn) {
	tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, nil)
	}
}
