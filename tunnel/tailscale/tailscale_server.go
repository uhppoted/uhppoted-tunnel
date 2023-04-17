package tailscale

import (
	"context"
	"errors"
	"fmt"
	// "html"
	"io"
	"net"
	"sync"
	// "syscall"
	"net/http"
	"strings"
	"time"

	"tailscale.com/tsnet"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tailscaleServer struct {
	conn.Conn
	hwif        string
	addr        *net.TCPAddr
	retry       conn.Backoff
	connections map[net.Conn]struct{}
	ctx         context.Context
	closed      chan struct{}
	sync.RWMutex
}

func NewTailscaleInServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tailscaleServer, error) {
	server, err := makeTailscaleServer(hwif, spec, retry, ctx)

	if err == nil {
		server.Infof("connector::tailscale-server-in")
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

func makeTailscaleServer(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tailscaleServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)

	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve tailscale address '%v'", spec)
	} else if addr.Port == 0 {
		return nil, fmt.Errorf("tailscale host requires a non-zero port")
	}

	ts := tailscaleServer{
		Conn: conn.Conn{
			Tag: "tailscale",
		},
		hwif:        hwif,
		addr:        addr,
		retry:       retry,
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
	s := new(tsnet.Server)

	s.Hostname = "uhppoted"

	defer s.Close()

	ln, err := s.Listen("tcp", ":80")
	if err != nil {
		ts.Fatalf("%v", err)
	}
	defer ln.Close()

	// lc, err := s.LocalClient()
	// if err != nil {
	// 	ts.Fatalf("%v", err)
	// }

	ts.Fatalf("%v", http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hi there! Welcome to the tailnet!")
		// who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		// if err != nil {
		// 	http.Error(w, err.Error(), 500)
		// 	return
		// }
		// fmt.Fprintf(w, "<html><body><h1>Hello, tailnet!</h1>\n")
		// fmt.Fprintf(w, "<p>You are <b>%s</b> from <b>%s</b> (%s)</p>",
		// 	html.EscapeString(who.UserProfile.LoginName),
		// 	html.EscapeString(firstLabel(who.Node.ComputedName)),
		// 	r.RemoteAddr)
	})))

	// var socket net.Listener
	// var closing = false
	//
	// go func() {
	// loop:
	// 	for {
	//
	// 		listener := net.ListenConfig{
	// 			Control: func(network, address string, connection syscall.RawConn) error {
	// 				if ts.hwif != "" {
	// 					return conn.BindToDevice(connection, ts.hwif, conn.IsIPv4(ts.addr.IP), ts.Conn)
	// 				} else {
	// 					return nil
	// 				}
	// 			},
	// 		}
	//
	// 		socket, err = listener.Listen(context.Background(), "tcp", fmt.Sprintf("%v", ts.addr))
	// 		if err != nil {
	// 			ts.Warnf("%v", err)
	// 		} else if socket == nil {
	// 			ts.Warnf("%v", fmt.Errorf("failed to create tailscale listen socket (%v)", socket))
	// 		} else {
	// 			ts.retry.Reset()
	// 			ts.listen(socket, router)
	// 		}
	//
	// 		if closing || !ts.retry.Wait(ts.Tag) {
	// 			break loop
	// 		}
	// 	}
	//
	// 	for k := range ts.connections {
	// 		k.Close()
	// 	}
	//
	// 	ts.closed <- struct{}{}
	// }()
	//
	// <-ts.ctx.Done()
	//
	// closing = true
	// socket.Close()

	return nil
}

func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
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
		}

		if err != nil {
			return
		}

		ts.Infof("incoming connection (%v)", client.RemoteAddr())

		if socket, ok := client.(*net.TCPConn); !ok {
			ts.Warnf("invalid TCP socket (%v)", socket)
			client.Close()
		} else {
			ts.Lock()
			ts.connections[socket] = struct{}{}
			ts.Unlock()

			go func(socket *net.TCPConn) {
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
