package gevloop

import (
	"container/heap"
	"fmt"
	"syscall"
)

type HandlerFunc func(evLoop *EvLoop, event Event, revent uint32)
type Event interface {
	Stop() error
	Start() error
	IsActive() bool
	cb(el *EvLoop)
}

type EvLoop struct {
	fd           int
	active       bool
	timeOut      int
	eventIO      []*EvIO
	timerHeap    *EvTimerHeap
	pendingQueue []Event
}

func Init() (*EvLoop, error) {
	el := EvLoop{}
	fd, err := syscall.EpollCreate1(0)
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
	el.numbEvIO()
	return &el, nil
}

func (el *EvLoop) Run() error {
	el.active = true
	for el.active {
		if el.timerHeap.Len() > 0 {
			el.timeOut = (*el.timerHeap)[0].at
		}
		var events []syscall.EpollEvent
		for _, v := range el.eventIO {
			events = append(events, v.events)
		}
		nevents, err := syscall.EpollWait(el.fd, events, el.timeOut)
		if err != nil {
			return err
		}
		if nevents == 0 { //timeout
			fmt.Println("evloop timeout....")
			//add first timeout timer to pendingQueue
			el.add2PendingQueue([]Event{(*el.timerHeap)[0]})
			timeOut := heap.Pop(el.timerHeap).(*EvTimer)
			if timeOut.repeat > 0 {
				timeOut.at = timeOut.repeat
				heap.Push(el.timerHeap, timeOut)
			}
		} else if nevents > 0 {
			fmt.Println("io event...")
			for _, v := range events {
				for _, j := range el.eventIO {
					if v.Fd == int32(j.fd) {
						j.revents = v.Events
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

func (el *EvLoop) Stop() {
	el.active = false
}

func (el *EvLoop) pendingCB() {
	for _, v := range el.pendingQueue {
		if v.IsActive() {
			v.cb(el)
		}
	}
	el.pendingQueue = make([]Event, 0)
}

func (el *EvLoop) add2PendingQueue(events []Event) {
	el.pendingQueue = append(el.pendingQueue, events...)
}

func (el *EvLoop) numbEvIO() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}
	if err = syscall.SetNonblock(fd, true); err != nil {
		return
	}
	numb := EvIO{}
	numb.Init(el, func(evLoop *EvLoop, event Event, revent uint32) {
		fmt.Println("Numb EVIO Called,must be something wrong")
	}, fd, syscall.EPOLLIN|syscall.EPOLLET&0xffffffff, nil)
	numb.Start()
	return
}
