package conn

import (
	"sync"
)

type TSocket interface {
	Close() error
}

type SocketList struct {
	list map[TSocket]struct{}
	sync.Mutex
}

func NewSocketList() SocketList {
	return SocketList{
		list: map[TSocket]struct{}{},
	}
}

func (s *SocketList) Add(socket TSocket) {
	s.Lock()
	defer s.Unlock()

	s.list[socket] = struct{}{}
}

func (s *SocketList) Close(socket TSocket) {
	s.Lock()
	defer s.Unlock()

	socket.Close()
	delete(s.list, socket)
}

func (s *SocketList) CloseAll() {
	s.Lock()
	defer s.Unlock()

	for socket := range s.list {
		socket.Close()
	}
}
