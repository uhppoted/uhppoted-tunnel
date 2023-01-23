package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tcpEventOutServer struct {
	tcpEventServer
}

func NewTCPEventOutServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpEventOutServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)

	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	} else if addr.Port == 0 {
		return nil, fmt.Errorf("TCP host requires a non-zero port")
	}

	tcp := tcpEventOutServer{
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

	tcp.Infof("connector::tcp-event-out-client")

	return &tcp, nil
}

func (tcp *tcpEventOutServer) Send(id uint32, message []byte) {
	for c, _ := range tcp.connections {
		go func(conn net.Conn) {
			tcp.send(conn, id, message)
		}(c)
	}
}

func (tcp *tcpEventOutServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
}

func (tcp *tcpEventOutServer) send(conn net.Conn, id uint32, message []byte) {
	packet := protocol.Packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
