package httpd

import (
	"fmt"
	"net"
	"net/http"
	"os"
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
	ch      chan protocol.Message
	closing chan struct{}
	closed  chan struct{}
}

func NewHTTPD(spec string, retry conn.Backoff) (*httpd, error) {
	addr, err := net.ResolveTCPAddr("tcp", spec)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, fmt.Errorf("unable to resolve HTTP base address '%v'", spec)
	}

	h := httpd{
		Conn: conn.Conn{
			Tag: "HTTP",
		},
		addr:    addr,
		retry:   retry,
		timeout: 5 * time.Second,
		ch:      make(chan protocol.Message, 16),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	return &h, nil
}

func (h *httpd) Close() {
	h.Infof("closing")
	
	close(h.closing)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-h.closed:
	    h.Infof("closed")

	case <-timeout.C:
	    h.Infof("close timeout")
	}
}

func (h *httpd) Run(router *router.Switch) error {
	h.listen(router)
	h.closed <- struct{}{}

	return nil
}

func (h *httpd) Send(id uint32, msg []byte) {
	// select {
	// case tcp.ch <- protocol.Message{ID: id, Message: msg}:
	// default:
	// }
}

func (h *httpd) listen(router *router.Switch) {
	h.Infof("listening on %v", h.addr)

	fs := filesystem{
		FileSystem: http.FS(os.DirFS("html")),
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(fs))
	// mux.HandleFunc("/udp", d.dispatch)

	// ... listen and serve
	srv := &http.Server{
		Addr:    fmt.Sprintf("%v", h.addr),
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		h.Fatalf("%v", err)
	}
}
