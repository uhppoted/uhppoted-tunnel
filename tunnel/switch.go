package tunnel

import ()

type Switch struct {
	relay relay
}

var handlers = map[uint32]func([]byte){}

func (s *Switch) request(id uint32, message []byte, handler func([]byte)) {
	go func() {
		handlers[id] = handler
		if reply := s.relay(id, message); reply != nil {
			handler(reply)
		}
	}()
}

func (s *Switch) reply(id uint32, message []byte) {
	go func() {
		if handler, ok := handlers[id]; ok && handler != nil {
			handler(message)
		}
	}()
}
