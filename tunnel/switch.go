package tunnel

type Switch struct {
	relay relay
}

func (s *Switch) received(id uint32, message []byte, handler func([]byte)) {
	go func() {
		if reply := s.relay(id, message); reply != nil {
			handler(reply)
		}
	}()
}
