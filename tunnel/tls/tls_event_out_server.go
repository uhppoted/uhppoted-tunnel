package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tlsEventOutServer struct {
	tlsEventServer
}

func NewTLSEventOutServer(hwif string, spec string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool, retry conn.Backoff, ctx context.Context) (*tlsEventOutServer, error) {
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

	tcp := tlsEventOutServer{
		tlsEventServer{
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
		},
	}

	tcp.tlsEventServer.received = tcp.received
	tcp.tlsEventServer.send = tcp.send

	tcp.Infof("connector::tls-event-out-server")

	return &tcp, nil
}

func (tcp *tlsEventOutServer) received(buffer []byte, router *router.Switch, socket net.Conn) {
}

func (tcp *tlsEventOutServer) send(conn net.Conn, id uint32, message []byte) {
	packet := protocol.Packetize(id, message)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(message), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v sent %v bytes to %v", id, len(message), conn.RemoteAddr())
	}
}
