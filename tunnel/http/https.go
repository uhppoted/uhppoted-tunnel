package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type https struct {
	httpd
	TLS *tls.Config
}

func NewHTTPS(spec string, html string, ca *x509.CertPool, keypair tls.Certificate, requireClientCertificate bool, retry conn.Backoff, ctx context.Context) (*https, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve HTTPS base address '%v'", spec)
	}

	fs := filesystem{
		FileSystem: http.FS(os.DirFS(html)),
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

	h := https{
		httpd: httpd{
			Conn: conn.Conn{
				Tag: "HTTPS",
			},
			addr:    addr,
			retry:   retry,
			timeout: 5 * time.Second,
			fs:      fs,
			ctx:     ctx,
			ch:      make(chan protocol.Message, 16),
			closed:  make(chan struct{}),
		},
		TLS: &config,
	}

	return &h, nil
}

func (h *https) Run(router *router.Switch) error {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(h.fs))
	mux.HandleFunc("/udp/broadcast", func(w http.ResponseWriter, r *http.Request) { h.dispatch(w, r, router) })
	mux.HandleFunc("/udp/send", func(w http.ResponseWriter, r *http.Request) { h.dispatch(w, r, router) })

	srv := &http.Server{
		Addr:      fmt.Sprintf("%v", h.addr),
		TLSConfig: h.TLS,
		Handler:   mux,
	}

	closing := false

	go func() {
	loop:
		for {
			start := time.Now()
			if err := srv.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				h.Warnf("%v", err)
			}

			if dt := time.Now().Sub(start); dt > 30*time.Second {
				h.retry.Reset()
			}

			if closing || !h.retry.Wait(h.Tag) {
				break loop
			}
		}
	}()

	<-h.ctx.Done()

	closing = true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		h.Warnf("%v", err)
	}

	h.closed <- struct{}{}

	return nil
}
