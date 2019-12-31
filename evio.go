package gevloop

import "syscall"

type EvIO struct {
	fd      int
	active  bool
	event   syscall.EpollEvent
	revent  uint32
	handler HandlerFunc
	data    interface{}
}

func (evIo *EvIO) cb(el *evLoop, revent uint32) {
	evIo.handler(el, evIo, revent)
}

func (evIo *EvIO) Stop() error {
	return nil
}

func (evIo *EvIO) Start() error {
	return nil
}

func (evIo *EvIO) IsActive() bool {
	return true
}
