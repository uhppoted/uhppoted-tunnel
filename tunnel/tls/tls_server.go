package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tlsServer struct {
	conn.Conn
	hwif        string
	addr        *net.TCPAddr
	config      *tls.Config
	retry       conn.Backoff
	connections map[net.Conn]struct{}
	pending     map[uint32]context.CancelFunc
	ctx         context.Context
	closing     bool
	closed      chan struct{}
	sync.RWMutex
}

func NewTLSInServer(hwif string, spec string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool, retry conn.Backoff, ctx context.Context) (*tlsServer, error) {
	server, err := makeTLSServer(hwif, spec, ca, keypair, requireClientCertificate, retry, ctx)

	if err == nil {
		server.Infof("connector::tls-server-in")
	}

	return server, err
}

func NewTLSOutServer(hwif string, spec string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool, retry conn.Backoff, ctx context.Context) (*tlsServer, error) {
	server, err := makeTLSServer(hwif, spec, ca, keypair, requireClientCertificate, retry, ctx)

	if err == nil {
		server.Infof("connector::tls-server-out")
	}

	return server, err
}

func makeTLSServer(hwif string, spec string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool, retry conn.Backoff, ctx context.Context) (*tlsServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)

	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	} else if addr.Port == 0 {
		return nil, fmt.Errorf("TCP host requires a non-zero port")
	}

	config := tls.Config{
		ClientCAs:    ca,
		Certificates: []tls.Certificate{keypair},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		ClientAuth: tls.VerifyClientCertIfGiven,
		MinVersion: tls.VersionTLS12,
	}

	if requireClientCertificate {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	tcp := tlsServer{
		Conn: conn.Conn{
			Tag: "TLS",
		},
		hwif:        hwif,
		addr:        addr,
		config:      &config,
		retry:       retry,
		connections: map[net.Conn]struct{}{},
		pending:     map[uint32]context.CancelFunc{},
		ctx:         ctx,
		closed:      make(chan struct{}),
	}

	return &tcp, nil
}

func (tcp *tlsServer) Close() {
	tcp.Infof("closing")

	for _, f := range tcp.pending {
		f()
	}

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		tcp.Infof("closed")

	case <-timeout.C:
		tcp.Infof("close timeout")
	}
}

func (tcp *tlsServer) Run(router *router.Switch) (err error) {
	tcp.closing = false
	sockets := conn.NewSocketList()

	defer sockets.CloseAll()

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

			if sock, err := listener.Listen(context.Background(), "tcp", fmt.Sprintf("%v", tcp.addr)); err != nil {
				tcp.Warnf("%v", err)
			} else if sock == nil {
				tcp.Warnf("%v", fmt.Errorf("failed to create TCP listen socket (%v)", sock))
			} else {
				socket := tls.NewListener(sock, tcp.config)

				sockets.Add(socket)
				tcp.retry.Reset()
				tcp.listen(socket, router)
				sockets.Closed(socket)
			}

			if tcp.closing || !tcp.retry.Wait(tcp.Tag) {
				break loop
			}
		}

		for k := range tcp.connections {
			k.Close()
		}

		tcp.closed <- struct{}{}
	}()

	<-tcp.ctx.Done()

	tcp.closing = true

	return nil
}

func (tcp *tlsServer) Send(id uint32, message []byte) {
	for c := range tcp.connections {
		go func(conn net.Conn) {
			tcp.send(conn, id, message)
		}(c)
	}
}

func (tcp *tlsServer) listen(socket net.Listener, router *router.Switch) {
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

		if socket, ok := client.(*tls.Conn); !ok {
			tcp.Warnf("invalid TLS socket (%v)", socket)
			client.Close()
		} else if err := tcp.handshake(socket); err != nil {
			tcp.Warnf("%v", err)
			client.Close()
		} else {
			tcp.Lock()
			tcp.connections[socket] = struct{}{}
			tcp.Unlock()

			go func(socket *tls.Conn) {
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

func (tcp *tlsServer) handshake(socket *tls.Conn) error {
	id := atomic.AddUint32(&ID, 1)
	state := socket.ConnectionState()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	tcp.Lock()
	tcp.pending[id] = cancel
	tcp.Unlock()

	defer func() {
		tcp.Lock()
		delete(tcp.pending, id)
		tcp.Unlock()
	}()

	if !state.HandshakeComplete {
		if err := socket.HandshakeContext(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (tcp *tlsServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
	tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, func(message []byte) {
			tcp.send(socket, id, message)
		})
	}
}

func (tcp *tlsServer) send(conn net.Conn, id uint32, message []byte) {
	packet := protocol.Packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
