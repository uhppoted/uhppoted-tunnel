package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type tlsEventInClient struct {
	tlsEventClient
}

func NewTLSEventInClient(hwif string, spec string, ca *x509.CertPool, keypair *tls.Certificate, retry conn.Backoff, ctx context.Context) (*tlsEventInClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	} else if addr == nil {
		return nil, fmt.Errorf("unable to resolve TCP address '%v'", spec)
	}

	config := tls.Config{
		RootCAs: ca,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	if keypair != nil {
		config.Certificates = []tls.Certificate{*keypair}
	}

	tcp := tlsEventInClient{
		tlsEventClient{
			Conn: conn.Conn{
				Tag: "TLS",
			},
			hwif:    hwif,
			addr:    addr,
			config:  &config,
			retry:   retry,
			timeout: 5 * time.Second,
			ch:      make(chan protocol.Message, 16),
			ctx:     ctx,
			closed:  make(chan struct{}),
		},
	}

	tcp.tlsEventClient.received = tcp.received
	tcp.tlsEventClient.send = tcp.send

	tcp.Infof("connector::tls-event-in-client")

	return &tcp, nil
}

func (tcp *tlsEventInClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
	tcp.Dumpf(buffer, "received %v bytes from %v", len(buffer), socket.RemoteAddr())

	for len(buffer) > 0 {
		id, msg, remaining := protocol.Depacketize(buffer)
		buffer = remaining

		router.Received(id, msg, nil)
	}
}

func (tcp *tlsEventInClient) send(conn net.Conn, id uint32, msg []byte) {
}
