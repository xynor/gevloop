package gevloop

import "syscall"

type EvIO struct {
	fd      int
	el      *EvLoop
	active  bool
	event   syscall.EpollEvent
	revent  uint32
	handler HandlerFunc
	data    interface{}
}

func (evIo *EvIO) cb(el *EvLoop, revent uint32) {
	evIo.handler(el, evIo, revent)
}

func (evIo *EvIO) Stop() {
}

func (evIo *EvIO) Start() {
}

func (evIo *EvIO) IsActive() bool {
	return true
}
