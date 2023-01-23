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

type tlsEventOutClient struct {
	tlsEventClient
}

func NewTLSEventOutClient(hwif string, spec string, ca *x509.CertPool, keypair *tls.Certificate, retry conn.Backoff, ctx context.Context) (*tlsEventOutClient, error) {
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

	tcp := tlsEventOutClient{
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

	tcp.Infof("connector::tls-event-out-client")

	return &tcp, nil
}

func (tcp *tlsEventOutClient) received(buffer []byte, router *router.Switch, socket net.Conn) {
}

func (tcp *tlsEventOutClient) send(conn net.Conn, id uint32, msg []byte) {
	packet := protocol.Packetize(id, msg)

	if N, err := conn.Write(packet); err != nil {
		tcp.Warnf("msg %v  error sending message to %v (%v)", id, conn.RemoteAddr(), err)
	} else if N != len(packet) {
		tcp.Warnf("msg %v  sent %v of %v bytes to %v", id, N, len(msg), conn.RemoteAddr())
	} else {
		tcp.Infof("msg %v  sent %v bytes to %v", id, len(msg), conn.RemoteAddr())
	}
}
