package http

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/protocol"
	"github.com/uhppoted/uhppoted-tunnel/router"
	"github.com/uhppoted/uhppoted-tunnel/tunnel/conn"
)

type httpd struct {
	conn.Conn
	addr    *net.TCPAddr
	retry   conn.Backoff
	timeout time.Duration
	fs      filesystem
	ctx     context.Context
	ch      chan protocol.Message
	closed  chan struct{}
}

type slice []byte

func (s slice) MarshalJSON() ([]byte, error) {
	a := make([]uint16, len(s))
	for i, v := range s {
		a[i] = uint16(v)
	}

	return json.Marshal(a)
}

type duration time.Duration

func (d *duration) UnmarshalJSON(b []byte) (err error) {
	var s string
	var t time.Duration
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}

	if t, err = time.ParseDuration(s); err != nil {
		return
	} else {
		*d = duration(t)
	}

	return
}

const GZIP_MINIMUM = 16384

func NewHTTP(spec string, html string, retry conn.Backoff, ctx context.Context) (*httpd, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve HTTP base address '%v'", spec)
	}

	fs := filesystem{
		FileSystem: http.FS(os.DirFS(html)),
	}

	h := httpd{
		Conn: conn.Conn{
			Tag: "HTTP",
		},
		addr:    addr,
		retry:   retry,
		timeout: 5 * time.Second,
		fs:      fs,
		ctx:     ctx,
		ch:      make(chan protocol.Message, 16),
		closed:  make(chan struct{}),
	}

	return &h, nil
}

func (h *httpd) Close() {
	h.Infof("closing")

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-h.closed:
		h.Infof("closed")

	case <-timeout.C:
		h.Infof("close timeout")
	}
}

func (h *httpd) Run(router *router.Switch) error {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(h.fs))
	mux.HandleFunc("/udp/broadcast", func(w http.ResponseWriter, r *http.Request) { h.dispatch(w, r, router) })
	mux.HandleFunc("/udp/send", func(w http.ResponseWriter, r *http.Request) { h.dispatch(w, r, router) })

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v", h.addr),
		Handler: mux,
	}

	closing := false

	go func() {
	loop:
		for {
			start := time.Now()
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				h.Warnf("%v", err)
			}

			if dt := time.Since(start); dt > 30*time.Second {
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

func (h *httpd) Send(id uint32, msg []byte) {
}

func (h *httpd) dispatch(w http.ResponseWriter, r *http.Request, router *router.Switch) {
	switch {
	case strings.ToUpper(r.Method) == http.MethodPost && r.URL.Path == "/udp/broadcast":
		h.broadcast(w, r, router)

	case strings.ToUpper(r.Method) == http.MethodPost && r.URL.Path == "/udp/send":
		h.send(w, r, router)

	default:
		http.Error(w, "Invalid request", http.StatusMethodNotAllowed)
	}
}

func (h *httpd) broadcast(w http.ResponseWriter, r *http.Request, router *router.Switch) {
	acceptsGzip := false
	contentType := ""

	body := struct {
		ID      int      `json:"ID"`
		Wait    duration `json:"wait,omitempty"`
		Request []byte   `json:"request"`
	}{
		ID:   0,
		Wait: duration(5 * time.Second),
	}

	for k, h := range r.Header {
		if strings.TrimSpace(strings.ToLower(k)) == "content-type" {
			for _, v := range h {
				contentType = strings.TrimSpace(strings.ToLower(v))
			}
		}

		if strings.TrimSpace(strings.ToLower(k)) == "accept-encoding" {
			for _, v := range h {
				if strings.Contains(strings.TrimSpace(strings.ToLower(v)), "gzip") {
					acceptsGzip = true
				}
			}
		}
	}

	switch contentType {
	case "application/json":
		blob, err := io.ReadAll(r.Body)
		if err != nil {
			h.Warnf("%v", err)
			http.Error(w, "Error reading request", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(blob, &body); err != nil {
			h.Warnf("%v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

	default:
		h.Warnf("%v", fmt.Errorf("invalid request content-type (%v)", contentType))
		http.Error(w, fmt.Sprintf("Invalid request content-type (%v)", contentType), http.StatusBadRequest)
		return
	}

	id := protocol.NextID()
	replies := []slice{}
	received := make(chan []byte)
	wait := time.Duration(body.Wait)
	waited := time.After(wait)
	ctx, cancel := context.WithTimeout(h.ctx, wait+5*time.Second)

	defer cancel()

	h.Dumpf(body.Request, "request %v  %v bytes from %v", id, len(body.Request), r.RemoteAddr)

	router.Received(id, body.Request, func(reply []byte) { received <- reply })

	for {
		select {
		case reply := <-received:
			h.Dumpf(reply, "reply %v  %v bytes for %v", id, len(reply), r.RemoteAddr)
			replies = append(replies, reply)

		case <-ctx.Done():
			h.Warnf("%v", ctx.Err())
			http.Error(w, "Request cancelled", http.StatusInternalServerError)
			return

		case <-waited:
			response := struct {
				ID      int     `json:"ID"`
				Replies []slice `json:"replies"`
			}{
				ID:      body.ID,
				Replies: replies,
			}

			h.reply(response, w, acceptsGzip)
			return
		}
	}
}

func (h *httpd) send(w http.ResponseWriter, r *http.Request, router *router.Switch) {
	acceptsGzip := false
	contentType := ""

	body := struct {
		ID      int    `json:"ID"`
		Wait    bool   `json:"wait,omitempty"`
		Request []byte `json:"request"`
	}{
		ID:   0,
		Wait: true,
	}

	for k, h := range r.Header {
		if strings.TrimSpace(strings.ToLower(k)) == "content-type" {
			for _, v := range h {
				contentType = strings.TrimSpace(strings.ToLower(v))
			}
		}

		if strings.TrimSpace(strings.ToLower(k)) == "accept-encoding" {
			for _, v := range h {
				if strings.Contains(strings.TrimSpace(strings.ToLower(v)), "gzip") {
					acceptsGzip = true
				}
			}
		}
	}

	switch contentType {
	case "application/json":
		blob, err := io.ReadAll(r.Body)
		if err != nil {
			h.Warnf("%v", err)
			http.Error(w, "Error reading request", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(blob, &body); err != nil {
			h.Warnf("%v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

	default:
		h.Warnf("%v", fmt.Errorf("invalid request content-type (%v)", contentType))
		http.Error(w, fmt.Sprintf("Invalid request content-type (%v)", contentType), http.StatusBadRequest)
		return
	}

	id := protocol.NextID()
	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)

	defer cancel()

	h.Dumpf(body.Request, "request %v  %v bytes from %v", id, len(body.Request), r.RemoteAddr)

	// ... set-ip request does not expect a response
	if !body.Wait {
		router.Received(id, body.Request, func(reply []byte) {})

		response := struct {
			ID int `json:"ID"`
		}{
			ID: body.ID,
		}

		h.reply(response, w, acceptsGzip)
		return
	}

	// ... normal request/response
	received := make(chan []byte)

	router.Received(id, body.Request, func(reply []byte) { received <- reply })

	for {
		select {
		case reply := <-received:
			h.Dumpf(reply, "reply %v  %v bytes for %v", id, len(reply), r.RemoteAddr)

			response := struct {
				ID    int   `json:"ID"`
				Reply slice `json:"reply,omitempty"`
			}{
				ID:    body.ID,
				Reply: reply,
			}

			h.reply(response, w, acceptsGzip)
			return

		case <-ctx.Done():
			h.Warnf("%v", ctx.Err())
			http.Error(w, "Request cancelled", http.StatusInternalServerError)
			return
		}
	}
}

func (h *httpd) reply(response any, w http.ResponseWriter, acceptsGzip bool) {
	if b, err := json.Marshal(response); err != nil {
		h.Warnf("%v", err)
		http.Error(w, "Internal error generating response", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")

		if acceptsGzip && len(b) > GZIP_MINIMUM {
			w.Header().Set("Content-Encoding", "gzip")

			gz := gzip.NewWriter(w)
			gz.Write(b)
			gz.Close()
		} else {
			w.Write(b)
		}
	}
}
