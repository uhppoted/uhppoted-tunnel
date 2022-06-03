package tunnel

import (
	"sync"
)

type Switch struct {
	relay relay
	sync.RWMutex
}

var handlers = map[uint32]func([]byte){}

func (s *Switch) request(id uint32, message []byte, handler func([]byte)) {
	s.Lock()
	defer s.Unlock()

	handlers[id] = handler

	go func() {
		s.relay(id, message)
	}()
}

func (s *Switch) reply(id uint32, message []byte) {
	s.RLock()
	defer s.RUnlock()

	if h, ok := handlers[id]; ok && h != nil {
		go func(handler func([]byte)) {
			handler(message)
		}(h)
	}
}
