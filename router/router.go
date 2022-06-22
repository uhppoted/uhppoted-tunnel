package router

import (
	"fmt"
	"sync"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

type Switch struct {
	relay func(uint32, []byte)
}

type Router struct {
	handlers map[uint32]handler
	idletime time.Duration
	closing  chan struct{}
	closed   chan struct{}
	sync.RWMutex
}

type handler struct {
	f       func([]byte)
	touched time.Time
}

var router = Router{
	handlers: map[uint32]handler{},
	idletime: 15 * time.Second,
	closing:  make(chan struct{}),
	closed:   make(chan struct{}),
}

var ticker = time.NewTicker(15 * time.Second)

func init() {
	go func() {
		defer func() {
			router.closed <- struct{}{}
		}()

		for {
			select {
			case <-router.closing:
				return

			case <-ticker.C:
				router.Sweep()
			}
		}
	}()
}

func NewSwitch(f func(uint32, []byte)) Switch {
	return Switch{
		relay: f,
	}
}

func (s *Switch) Received(id uint32, message []byte, h func([]byte)) {
	hf := router.get(id)

	switch {
	case hf != nil:
		go func() {
			hf(message)
		}()

	default:
		if h != nil {
			router.add(id, h)
		}

		go func() {
			s.relay(id, message)
		}()
	}
}

func (r *Router) add(id uint32, h func([]byte)) {
	router.Lock()
	defer router.Unlock()

	router.handlers[id] = handler{
		f:       h,
		touched: time.Now(),
	}
}

func (r *Router) get(id uint32) func([]byte) {
	if h, ok := r.handlers[id]; ok && h.f != nil {
		h.touched = time.Now()
		return h.f
	}

	return nil
}

func (r *Router) Sweep() {
	r.Lock()
	defer r.Unlock()

	cutoff := time.Now().Add(-r.idletime)
	idle := []uint32{}

	for k, v := range r.handlers {
		if v.touched.Before(cutoff) {
			idle = append(idle, k)
		}
	}

	for _, k := range idle {
		debugf("ROUTER", "removing idle handler function (%v)", k)
		delete(r.handlers, k)
	}
}

func Close() {
	infof("ROUTER", "closing")
	router.closing <- struct{}{}

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-router.closed:
		infof("ROUTER", "closed")

	case <-timeout.C:
		infof("ROUTER", "close timeout")
	}
}

func debugf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Debugf(f, args...)
}

func infof(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Infof(f, args...)
}
