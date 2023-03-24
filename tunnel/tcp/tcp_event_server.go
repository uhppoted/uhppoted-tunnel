package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	// "github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tcpEventServer struct {
	conn.Conn
	hwif        string
	addr        *net.TCPAddr
	retry       conn.Backoff
	connections map[net.Conn]struct{}
	ctx         context.Context
	closed      chan struct{}

	received func([]byte, *router.Switch, net.Conn)

	sync.RWMutex
}

func (tcp *tcpEventServer) Close() {
	tcp.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		tcp.Infof("closed")

	case <-timeout.C:
		tcp.Infof("close timeout")
	}
}

func (tcp *tcpEventServer) Run(router *router.Switch) (err error) {
	var socket net.Listener
	var closing = false

	go func() {
	loop:
		for {

			listener := net.ListenConfig{
				Control: func(network, address string, connection syscall.RawConn) error {
					if tcp.hwif != "" {
						return conn.BindToDevice(connection, tcp.hwif, conn.IsIPv4(tcp.addr.IP), tcp.Conn)
					} else {
						return nil
					}
				},
			}

			socket, err = listener.Listen(context.Background(), "tcp", fmt.Sprintf("%v", tcp.addr))
			if err != nil {
				tcp.Warnf("%v", err)
			} else if socket == nil {
				tcp.Warnf("%v", fmt.Errorf("failed to create TCP listen socket (%v)", socket))
			} else {
				tcp.retry.Reset()
				tcp.listen(socket, router)
			}

			if closing || !tcp.retry.Wait(tcp.Tag) {
				break loop
			}
		}

		for k := range tcp.connections {
			k.Close()
		}

		tcp.closed <- struct{}{}
	}()

	<-tcp.ctx.Done()

	closing = true
	socket.Close()

	return nil
}

func (tcp *tcpEventServer) listen(socket net.Listener, router *router.Switch) {
	tcp.Infof("listening on %v", socket.Addr())

	defer socket.Close()

	for {
		client, err := socket.Accept()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			tcp.Errorf("%v %v", err, errors.Is(err, net.ErrClosed))
		}

		if err != nil {
			return
		}

		tcp.Infof("incoming connection (%v)", client.RemoteAddr())

		if socket, ok := client.(*net.TCPConn); !ok {
			tcp.Warnf("invalid TCP socket (%v)", socket)
			client.Close()
		} else {
			tcp.Lock()
			tcp.connections[socket] = struct{}{}
			tcp.Unlock()

			go func(socket *net.TCPConn) {
				for {
					buffer := make([]byte, 2048) // buffer is handed off to router
					if N, err := socket.Read(buffer); err != nil {
						if err == io.EOF {
							tcp.Infof("client connection %v closed ", socket.RemoteAddr())
						} else {
							tcp.Warnf("%v", err)
						}
						break
					} else {
						tcp.received(buffer[:N], router, socket)
					}
				}

				tcp.Lock()
				delete(tcp.connections, socket)
				tcp.Unlock()
			}(socket)
		}
	}
}
