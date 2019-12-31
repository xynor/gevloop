package gevloop

import (
	"container/heap"
	"fmt"
	"syscall"
)

type HandlerFunc func(evLoop *evLoop, event Event, revent uint32)
type Event interface {
	Stop() error
	Start() error
	IsActive() bool
	cb(el *evLoop, revent uint32)
}

type evLoop struct {
	fd           int
	active       bool
	timeOut      int
	eventIO      []*EvIO
	timerHeap    *EvTimerHeap
	pendingQueue []Event
}

func Init() (*evLoop, error) {
	el := &evLoop{}
	fd, err := syscall.EpollCreate(1)
	if err != nil {
		return nil, err
	}
	el.fd = fd
	el.active = false
	el.timeOut = -1
	el.timerHeap = &EvTimerHeap{}
	heap.Init(el.timerHeap)
	el.pendingQueue = make([]Event, 0)
	el.eventIO = make([]*EvIO, 0)
	return el, nil
}

func (el *evLoop) Run() error {
	el.active = true
	for el.active {
		if el.timerHeap.Len() > 0 {
			el.timeOut = (*el.timerHeap)[0].at
		}
		var events []syscall.EpollEvent
		for _, v := range el.eventIO {
			events = append(events, v.event)
		}
		nevents, err := syscall.EpollWait(el.fd, events, el.timeOut)
		if err != nil {
			return err
		}
		if nevents < 0 { //timeout
			fmt.Println("evloop timeout....")
			//add first timeout timer to pendingQueue
			el.add2PendingQueue([]Event{(*el.timerHeap)[0]})
			timeOut := heap.Pop(el.timerHeap).(*EvTimer)
			if timeOut.repeat > 0 {
				timeOut.at = timeOut.repeat
				heap.Push(el.timerHeap, timeOut)
			}
		} else {
			fmt.Println("io event...")
			for _, v := range events {
				for _, j := range el.eventIO {
					if v.Fd == int32(j.fd) {
						el.add2PendingQueue([]Event{j})
					}
				}
			}
		}
		fmt.Println("CALL...")
		el.pendingCB()
		if !el.active {
			break
		}
	}
	return nil
}

func (el *evLoop) Stop() {
	el.active = false
}

func (el *evLoop) pendingCB() {
	for _, v := range el.pendingQueue {
		if v.IsActive() {
			v.cb(el, 1)
		}
	}
	el.pendingQueue = make([]Event, 0)
}

func (el *evLoop) add2PendingQueue(events []Event) {
	el.pendingQueue = append(el.pendingQueue, events...)
}
