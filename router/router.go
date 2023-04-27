package router

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

type Switch struct {
	relay func(uint32, []byte)
}

type Router struct {
	handlers ihandlers
	idletime time.Duration
	closing  chan struct{}
	closed   chan struct{}
	sync.RWMutex
}

type ihandlers interface {
	get(uint32) *handler
	put(uint32, *handler)
	apply(f func(map[uint32]*handler))
}

type handler struct {
	f       func([]byte)
	touched time.Time
}

var router = Router{
	handlers: hmake(),
	idletime: 15 * time.Second,
	closing:  make(chan struct{}),
	closed:   make(chan struct{}),
}

var ticker = time.NewTicker(15 * time.Second)
var limiter = rate.NewLimiter(1, 120)

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

func SetRateLimiter(l *rate.Limiter) {
	if l != nil {
		limiter = l
	}
}

func NewSwitch(f func(uint32, []byte)) Switch {
	return Switch{
		relay: f,
	}
}

func (s *Switch) Received(id uint32, message []byte, h func([]byte)) {
	if !limiter.Allow() {
		warnf("ROUTER", "rate limit exceeded")
		return
	}

	if message != nil {
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
}

func (r *Router) add(id uint32, h func([]byte)) {
	r.handlers.put(id,
		&handler{
			f:       h,
			touched: time.Now(),
		})
}

func (r *Router) get(id uint32) func([]byte) {
	if h := r.handlers.get(id); h != nil && h.f != nil {
		h.touched = time.Now()
		return h.f
	}

	return nil
}

func (r *Router) Sweep() {
	f := func(handlers map[uint32]*handler) {
		cutoff := time.Now().Add(-r.idletime)
		idle := []uint32{}

		for k, v := range handlers {
			if v.touched.Before(cutoff) {
				idle = append(idle, k)
			}
		}

		for _, k := range idle {
			debugf("ROUTER", "removing idle handler function (%v)", k)
			delete(handlers, k)
		}
	}

	r.handlers.apply(f)
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
	f := fmt.Sprintf("%-10v %v", tag, format)

	log.Debugf(f, args...)
}

func infof(tag, format string, args ...any) {
	f := fmt.Sprintf("%-10v %v", tag, format)

	log.Infof(f, args...)
}

func warnf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-10v %v", tag, format)

	log.Warnf(f, args...)
}
