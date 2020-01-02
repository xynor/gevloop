package gevloop

import (
	"errors"
	"syscall"
)

type EvIO struct {
	fd      int
	el      *EvLoop
	active  bool
	events  syscall.EpollEvent
	revents uint32
	handler HandlerFunc
	data    interface{}
}

func (evIo *EvIO) cb(el *EvLoop) {
	revent := evIo.revents
	evIo.handler(el, evIo, revent)
}

func (evIo *EvIO) Init(el *EvLoop, handler HandlerFunc, fd int, Events uint32, data interface{}) error {
	if el == nil {
		return errors.New("evLoop is nil")
	}

	evIo.fd = fd
	evIo.events.Fd = int32(fd)
	evIo.events.Events = Events
	evIo.handler = handler
	evIo.data = data
	evIo.el = el
	return nil
}

func (evIo *EvIO) Stop() error {
	evIo.active = false
	evIo.el.eventIO = append(evIo.el.eventIO, evIo)
	if err := syscall.EpollCtl(evIo.el.fd, syscall.EPOLL_CTL_DEL, evIo.fd, &evIo.events); err != nil {
		return err
	}

	return nil
}

func (evIo *EvIO) Start() error {
	evIo.active = true
	evIo.el.eventIO = append(evIo.el.eventIO, evIo)
	if err := syscall.EpollCtl(evIo.el.fd, syscall.EPOLL_CTL_ADD, evIo.fd, &evIo.events); err != nil {
		return err
	}

	return nil
}

func (evIo *EvIO) Data() interface{} {
	return evIo.data
}

func (evIo *EvIO) Fd()int {
	return evIo.fd
}

func (evIo *EvIO) IsActive() bool {
	return evIo.active
}
