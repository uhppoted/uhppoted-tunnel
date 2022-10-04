package router

import (
	"sync"
)

type handlers struct {
	handlers map[uint32]*handler
	sync.RWMutex
}

func hmake() *handlers {
	return &handlers{
		handlers: map[uint32]*handler{},
	}
}

func (h *handlers) get(id uint32) *handler {
	h.RLock()
	defer h.RUnlock()

	if h, ok := h.handlers[id]; ok {
		return h
	}

	return nil
}

func (h *handlers) put(id uint32, handler *handler) {
	h.Lock()
	defer h.Unlock()

	h.handlers[id] = handler
}

func (h *handlers) apply(f func(map[uint32]*handler)) {
	h.Lock()
	defer h.Unlock()

	f(h.handlers)
}
