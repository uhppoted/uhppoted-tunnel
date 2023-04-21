package tailscale

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"tailscale.com/tsnet"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tailscaleClient struct {
	conn.Conn
	dir      string
	hostname string
	addr     string
	port     uint16
	auth     string
	retry    conn.Backoff
	logging  string
	timeout  time.Duration
	ch       chan protocol.Message
	ctx      context.Context
	closed   chan struct{}
}

// func NewTCPInClient(hwif string, spec string, retry conn.Backoff, ctx context.Context) (*tcpClient, error) {
//     client, err := makeTCPClient(hwif, spec, retry, ctx)
//
//     if err == nil {
//         client.Infof("connector::tcp-client-in")
//     }
//
//     return client, err
// }

func NewTailscaleOutClient(workdir string, hostname string, spec string, auth string, retry conn.Backoff, logging string, ctx context.Context) (*tailscaleClient, error) {
	client, err := makeTailscaleClient(workdir, hostname, spec, auth, retry, logging, ctx)

	if err == nil {
		client.Infof("connector::tailscale-client-out  %v", client.hostname)
		client.Infof("connector::tailscale-client-out  %v", client.dir)
	}

	return client, err
}

func makeTailscaleClient(workdir string, hostname, spec string, auth string, retry conn.Backoff, logging string, ctx context.Context) (*tailscaleClient, error) {
	addr, port, err := resolveTailscaleAddr(spec)
	if err != nil {
		return nil, err
	} else if addr == "" {
		return nil, fmt.Errorf("unable to resolve tailscale address '%v'", spec)
	}

	var dir string
	var name string

	if hostname == "" {
		name = fmt.Sprintf("%v.client", addr)
	} else {
		name = hostname
	}

	if hostname == "" {
		dir = filepath.Join(workdir, "tailscale", addr, "client")
	} else {
		dir = filepath.Join(workdir, "tailscale", hostname)
	}

	in := tailscaleClient{
		Conn: conn.Conn{
			Tag: "tailscale",
		},
		dir:      dir,
		hostname: name,
		addr:     addr,
		port:     port,
		auth:     auth,
		retry:    retry,
		logging:  logging,
		timeout:  5 * time.Second,
		ch:       make(chan protocol.Message, 16),
		ctx:      ctx,
		closed:   make(chan struct{}),
	}

	return &in, nil
}

func (ts *tailscaleClient) Close() {
	ts.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-ts.closed:
		ts.Infof("closed")

	case <-timeout.C:
		ts.Infof("close timeout")
	}
}

func (ts *tailscaleClient) Run(router *router.Switch) error {
	if err := ts.connect(router); err != nil {
		return err
	} else {
		ts.closed <- struct{}{}

		return nil
	}
}

func (ts *tailscaleClient) Send(id uint32, msg []byte) {
	select {
	case ts.ch <- protocol.Message{ID: id, Message: msg}:
	default:
	}
}

func (ts *tailscaleClient) connect(router *router.Switch) error {
	// ... get authkey
	var authKey string

	if key, err := getAuthKey(ts.auth); err != nil {
		return err
	} else {
		authKey = key
	}

	// ... initialise server
	logf := func(f string, args ...any) {
		if ts.logging == "debug" {
			ts.Debugf(f, args...)
		}
	}

	server := &tsnet.Server{
		Logf:      logf,
		Hostname:  ts.hostname,
		AuthKey:   authKey,
		Dir:       ts.dir,
		Ephemeral: false,
	}

	defer server.Close()

	if err := server.Start(); err != nil {
		return err
	}

	for {
		// ... bring server up
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if status, err := server.Up(ctx); err != nil {
			ts.Warnf("%v", err)
		} else {
			cancel()
			ts.Debugf("state  %v", status.BackendState)
		}

		// ... manual authorisation if required
		if lc, err := server.LocalClient(); err != nil {
			ts.Warnf("%v", err)
		} else if status, err := lc.Status(ts.ctx); err != nil {
			ts.Warnf("%v", err)
		} else {
			ts.Infof("status %v", status.BackendState)
			if status.BackendState == "NeedsLogin" && status.AuthURL != "" {
				ts.Errorf("node authorisation required - please authorise node at URL %v", status.AuthURL)
			}
		}

		// ... 'k, we're good to go
		ts.Infof("connecting to %v:%v", ts.addr, ts.port)

		if socket, err := server.Dial(context.Background(), "tcp", fmt.Sprintf("%v:%v", ts.addr, ts.port)); err != nil {
			ts.Warnf("%v", err)
		} else if socket == nil {
			ts.Warnf("connect %v failed (%v)", ts.addr, socket)
		} else {
			ts.retry.Reset()
			eof := make(chan struct{})

			go func() {
				for {
					select {
					case msg := <-ts.ch:
						ts.Infof("msg %v  relaying to %v", msg.ID, socket.RemoteAddr())
						ts.send(socket, msg.ID, msg.Message)

					case <-eof:
						return

					case <-ts.ctx.Done():
						socket.Close()
						return
					}

					println("woot/sleep")
					time.Sleep(1000)
				}
			}()

			if err := ts.listen(socket, router); err != nil && !errors.Is(err, net.ErrClosed) {
				ts.Warnf("%v", err)
			}

			close(eof)

			time.Sleep(5000)
		}

		if !ts.retry.Wait(ts.Tag) {
			return nil
		}
	}
}

func (ts *tailscaleClient) listen(socket net.Conn, router *router.Switch) error {
	ts.Infof("connected  to %v", socket.RemoteAddr())

	defer socket.Close()

	for {
		buffer := make([]byte, 2048)
		N, err := socket.Read(buffer)
		if err != nil {
			return err
		}

		go func() {
			ts.received(buffer[:N], router, socket)
		}()
	}
}

func (ts *tailscaleClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	ts.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			ts.send(socket, id, message)
		})
	}
}

func (ts *tailscaleClient) send(conn net.Conn, id uint32, msg []byte) []byte {
	packet := protocol.Packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		ts.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		ts.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		ts.Infof("msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}

	return nil
}
