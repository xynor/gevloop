package gevloop

import (
	"container/list"
	"fmt"
	"syscall"
	"time"
)

type HandlerFunc func(evLoop *EvLoop, event Event, revent uint32)
type Event interface {
	Stop() error
	Start() error
	IsActive() bool
	Data() interface{}
	cb(el *EvLoop)
}

type EvLoop struct {
	fd           int
	active       bool
	timeOut      int
	eventIO      []*EvIO
	timerList    *list.List
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
	el.timerList = list.New()
	//TODO maybe should use minHeap
	el.pendingQueue = make([]Event, 0)
	el.eventIO = make([]*EvIO, 0)
	el.numbEvIO()
	return &el, nil
}

func (el *EvLoop) Run() error {
	el.active = true
	for el.active {
		timeNow := int(time.Now().UnixNano() / 1e6)
		if el.timerList.Len() > 0 {
			triggerTime := el.timerList.Front().Value.(*EvTimer).triggerTime
			if timeNow >= triggerTime {
				el.timeOut = 0
			} else {
				el.timeOut = triggerTime - timeNow
			}
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
			for e := el.timerList.Front(); e != nil; e = e.Next() {
				if e.Value.(*EvTimer).triggerTime <= timeNow {
					el.timerList.Remove(e)
					el.add2PendingQueue([]Event{e.Value.(*EvTimer)})
				} else {
					break
				}
			}
		} else if nevents > 0 {
			for _, v := range events {
				for _, j := range el.eventIO {
					if v.Fd == int32(j.fd) {
						j.revents = v.Events
						el.add2PendingQueue([]Event{j})
					}
				}
			}
		}
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
		fmt.Println("Numb EVIO Called,Loop running")
	}, fd, syscall.EPOLLIN|syscall.EPOLLET&0xffffffff, nil)
	numb.Start()
	return
}
