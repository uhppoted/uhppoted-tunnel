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

type tlsEventInServer struct {
	conn.Conn
	hwif        string
	addr        *net.TCPAddr
	config      *tls.Config
	key         string
	retry       conn.Backoff
	connections map[net.Conn]struct{}
	pending     map[uint32]context.CancelFunc
	ctx         context.Context
	closed      chan struct{}
	sync.RWMutex
}

func NewTLSEventInServer(hwif string, spec string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool, retry conn.Backoff, ctx context.Context) (*tlsEventInServer, error) {
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

	tcp := tlsEventInServer{
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

	tcp.Infof("connector::tls-event-in-server")

	return &tcp, nil
}

func (tcp *tlsEventInServer) Close() {
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

func (tcp *tlsEventInServer) Run(router *router.Switch) (err error) {
	var socket net.Listener

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
				tcp.Warnf("%v", fmt.Errorf("Failed to create TCP listen socket (%v)", sock))
			} else {
				socket = tls.NewListener(sock, tcp.config)

				tcp.retry.Reset()
				tcp.listen(socket, router)
			}

			if !tcp.retry.Wait(tcp.Tag) {
				break loop
			}
		}

		for k, _ := range tcp.connections {
			k.Close()
		}

		tcp.closed <- struct{}{}
	}()

	<-tcp.ctx.Done()

	socket.Close()

	return nil
}

func (tcp *tlsEventInServer) Send(id uint32, message []byte) {
}

func (tcp *tlsEventInServer) listen(socket net.Listener, router *router.Switch) {
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

func (tcp *tlsEventInServer) handshake(socket *tls.Conn) error {
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

func (tcp *tlsEventInServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
	tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, nil)
	}
}
