package tailscale

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"sync"
	"time"

	"tailscale.com/tsnet"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tailscaleServer struct {
	conn.Conn
	dir         string
	hostname    string
	addr        string
	port        uint16
	retry       conn.Backoff
	logging     string
	connections map[net.Conn]struct{}
	ctx         context.Context
	closed      chan struct{}
	sync.RWMutex
}

func NewTailscaleInServer(workdir string, hostname string, spec string, retry conn.Backoff, logging string, ctx context.Context) (*tailscaleServer, error) {
	server, err := makeTailscaleServer(workdir, hostname, spec, retry, logging, ctx)

	if err == nil {
		server.Infof("connector::tailscale-server-in  %v", server.hostname)
		server.Infof("connector::tailscale-server-in  %v", server.dir)
	}

	return server, err
}

// func NewTCPOutServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpServer, error) {
//     server, err := makeTCPServer(hwif, spec, retry, ctx)
//
//     if err == nil {
//         server.Infof("connector::tcp-server-out")
//     }
//
//     return server, err
// }

func makeTailscaleServer(workdir string, hostname string, spec string, retry conn.Backoff, logging string, ctx context.Context) (*tailscaleServer, error) {
	addr, port, err := resolveTailscaleAddr(spec)
	if err != nil {
		return nil, err
	} else if port == 0 {
		return nil, fmt.Errorf("tailscale server requires a non-zero port")
	}

	var dir string
	var name string

	if hostname == "" {
		name = fmt.Sprintf("%v", addr)
	} else {
		name = hostname
	}

	if hostname == "" {
		dir = filepath.Join(workdir, "tailscale", addr, "server")
	} else {
		dir = filepath.Join(workdir, "tailscale", hostname)
	}

	ts := tailscaleServer{
		Conn: conn.Conn{
			Tag: "tailscale",
		},
		dir:         dir,
		hostname:    name,
		addr:        addr,
		port:        port,
		retry:       retry,
		logging:     logging,
		connections: map[net.Conn]struct{}{},
		ctx:         ctx,
		closed:      make(chan struct{}),
	}

	return &ts, nil
}

func (ts *tailscaleServer) Close() {
	ts.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-ts.closed:
		ts.Infof("closed")

	case <-timeout.C:
		ts.Infof("close timeout")
	}
}

func (ts *tailscaleServer) Run(router *router.Switch) (err error) {
	var socket net.Listener
	var closing = false

	logf := func(f string, args ...any) {
		if ts.logging != "" && ts.logging != "no-log" {
			ts.Debugf(f, args...)
		}
	}

	server := &tsnet.Server{
		Logf:      logf,
		Hostname:  ts.hostname,
		Dir:       ts.dir,
		Ephemeral: false,
	}

	defer server.Close()

	go func() {
	loop:
		for {
			socket, err := server.Listen("tcp", fmt.Sprintf(":%v", ts.port))
			if err != nil {
				ts.Warnf("%v", err)
			} else if socket == nil {
				ts.Warnf("%v", fmt.Errorf("failed to create tailscale listen socket (%v)", socket))
			} else {
				ts.retry.Reset()
				ts.listen(socket, router)
			}

			if closing || !ts.retry.Wait(ts.Tag) {
				break loop
			}
		}

		for k := range ts.connections {
			k.Close()
		}

		ts.closed <- struct{}{}
	}()

	<-ts.ctx.Done()

	closing = true
	socket.Close()

	return nil
}

func (ts *tailscaleServer) Send(id uint32, message []byte) {
	for c := range ts.connections {
		go func(conn net.Conn) {
			ts.send(conn, id, message)
		}(c)
	}
}

func (ts *tailscaleServer) listen(socket net.Listener, router *router.Switch) {
	ts.Infof("listening on %v", socket.Addr())

	defer socket.Close()

	for {
		client, err := socket.Accept()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			ts.Errorf("%v %v", err, errors.Is(err, net.ErrClosed))
		} else if err != nil {
			return
		}

		ts.Infof("incoming connection (%v)", client.RemoteAddr())

		if socket, ok := client.(net.Conn); !ok {
			ts.Warnf("invalid tailscale socket (%v)", socket)
			client.Close()
		} else {
			ts.Lock()
			ts.connections[socket] = struct{}{}
			ts.Unlock()

			go func(socket net.Conn) {
				for {
					buffer := make([]byte, 2048) // buffer is handed off to router
					if N, err := socket.Read(buffer); err != nil {
						if err == io.EOF {
							ts.Infof("client connection %v closed ", socket.RemoteAddr())
						} else {
							ts.Warnf("%v", err)
						}
						break
					} else {
						ts.received(buffer[:N], router, socket)
					}

					time.Sleep(5000)
				}

				ts.Lock()
				delete(ts.connections, socket)
				ts.Unlock()
			}(socket)
		}
	}
}

func (ts *tailscaleServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
	ts.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		fmt.Printf("%v %v %v\n", id, msg, remaining)

		router.Received(id, msg, func(message []byte) {
			ts.send(socket, id, message)
		})
	}
}

func (ts *tailscaleServer) send(conn net.Conn, id uint32, message []byte) {
	packet := protocol.Packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		ts.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		ts.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		ts.Infof("msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
