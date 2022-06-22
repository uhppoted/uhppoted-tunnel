package tcp

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tcpClient struct {
	tag     string
	addr    *net.TCPAddr
	retry   conn.Backoff
	timeout time.Duration
	ch      chan protocol.Message
	closing chan struct{}
	closed  chan struct{}
}

func NewTCPClient(spec string, retry conn.Backoff) (*tcpClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	in := tcpClient{
		tag:     "TCP",
		addr:    addr,
		retry:   retry,
		timeout: 5 * time.Second,
		ch:      make(chan protocol.Message, 16),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	return &in, nil
}

func (tcp *tcpClient) Close() {
	infof(tcp.tag, "closing")
	close(tcp.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		infof(tcp.tag, "closed")

	case <-timeout.C:
		infof(tcp.tag, "close timeout")
	}
}

func (tcp *tcpClient) Run(router *router.Switch) error {
	tcp.connect(router)
	tcp.closed <- struct{}{}

	return nil
}

func (tcp *tcpClient) Send(id uint32, msg []byte) {
	select {
	case tcp.ch <- protocol.Message{ID: id, Message: msg}:
	default:
	}
}

func (tcp *tcpClient) connect(router *router.Switch) {
	for {
		infof(tcp.tag, "connecting to %v", tcp.addr)

		if socket, err := net.Dial("tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
			warnf(tcp.tag, "%v", err)
		} else if socket == nil {
			warnf(tcp.tag, "connect %v failed (%v)", tcp.addr, socket)
		} else {
			tcp.retry.Reset()
			eof := make(chan struct{})

			go func() {
				for {
					select {
					case msg := <-tcp.ch:
						infof(tcp.tag, "msg %v  relaying to %v", msg.ID, socket.RemoteAddr())
						tcp.send(socket, msg.ID, msg.Message)

					case <-eof:
						return

					case <-tcp.closing:
						socket.Close()
						return
					}
				}
			}()

			if err := tcp.listen(socket, router); err != nil {
				if err == io.EOF {
					warnf(tcp.tag, "connection to %v closed ", socket.RemoteAddr())
				} else {
					warnf(tcp.tag, "connection to %v error (%v)", tcp.addr, err)
				}
			}

			close(eof)
		}

		if !tcp.retry.Wait(tcp.tag, tcp.closing) {
			return
		}
	}
}

func (tcp *tcpClient) listen(socket net.Conn, router *router.Switch) error {
	infof(tcp.tag, "connected  to %v", socket.RemoteAddr())

	defer socket.Close()

	for {
		buffer := make([]byte, 2048)
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		go func() {
			tcp.received(buffer[:N], router, socket)
		}()
	}
}

func (tcp *tcpClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	dumpf(tcp.tag, buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			tcp.send(socket, id, message)
		})
	}
}

func (tcp *tcpClient) send(conn net.Conn, id uint32, msg []byte) []byte {
	packet := protocol.Packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		warnf(tcp.tag, "msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf(tcp.tag, "msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		infof(tcp.tag, "msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}

	return nil
}
