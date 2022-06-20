package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
)

type tlsServer struct {
	tag         string
	addr        *net.TCPAddr
	config      *tls.Config
	key         string
	retryDelay  time.Duration
	connections map[net.Conn]struct{}
	pending     map[uint32]context.CancelFunc
	closing     chan struct{}
	closed      chan struct{}
	sync.RWMutex
}

var ID uint32 = 0

func NewTLSServer(spec string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool) (*tlsServer, error) {
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
		tag:         "TLS",
		addr:        addr,
		config:      &config,
		retryDelay:  15 * time.Second,
		connections: map[net.Conn]struct{}{},
		pending:     map[uint32]context.CancelFunc{},
		closing:     make(chan struct{}),
		closed:      make(chan struct{}),
	}

	return &tcp, nil
}

func (tcp *tlsServer) Close() {
	infof(tcp.tag, "closing")
	close(tcp.closing)

	for _, f := range tcp.pending {
		f()
	}

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-tcp.closed:
		infof(tcp.tag, "closed")

	case <-timeout.C:
		infof(tcp.tag, "close timeout")
	}
}

func (tcp *tlsServer) Run(router *router.Switch) (err error) {
	var socket net.Listener
	var closing = false
	var delay = 0 * time.Second

	go func() {
		for !closing {
			time.Sleep(delay)

			socket, err = tls.Listen("tcp", fmt.Sprintf("%v", tcp.addr), tcp.config)
			if err != nil {
				warnf(tcp.tag, "%v", err)

			} else if socket == nil {
				warnf(tcp.tag, "%v", fmt.Errorf("Failed to create TCP listen socket (%v)", socket))
			}

			delay = tcp.retryDelay

			tcp.listen(socket, router)
		}

		for k, _ := range tcp.connections {
			k.Close()
		}

		tcp.closed <- struct{}{}
	}()

	<-tcp.closing

	closing = true
	socket.Close()

	return nil
}

func (tcp *tlsServer) Send(id uint32, message []byte) {
	for c, _ := range tcp.connections {
		go func(conn net.Conn) {
			tcp.send(conn, id, message)
		}(c)
	}
}

func (tcp *tlsServer) listen(socket net.Listener, router *router.Switch) {
	infof(tcp.tag, "listening on %v", socket.Addr())

	defer socket.Close()

	for {
		client, err := socket.Accept()
		if err != nil {
			errorf(tcp.tag, "%v", err)
			return
		}

		infof(tcp.tag, "incoming connection (%v)", client.RemoteAddr())

		if socket, ok := client.(*tls.Conn); !ok {
			warnf(tcp.tag, "invalid TLS socket (%v)", socket)
			client.Close()
		} else if err := tcp.handshake(socket); err != nil {
			warnf(tcp.tag, "%v", err)
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
							infof(tcp.tag, "client connection %v closed ", socket.RemoteAddr())
						} else {
							warnf(tcp.tag, "%v", err)
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
	dumpf(tcp.tag, buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

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
		warnf(tcp.tag, "msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		warnf(tcp.tag, "msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		infof(tcp.tag, "msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
