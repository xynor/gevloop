package gevloop

import (
	"errors"
	"syscall"
	"time"
)

type EvTimer struct {
	el          *EvLoop
	active      bool
	repeat      int
	triggerTime int
	handler     HandlerFunc
	data        interface{}
}

func (evTimer *EvTimer) cb(el *EvLoop) {
	revent := uint32(syscall.SYS_TIMES)
	evTimer.active = false
	if evTimer.repeat != 0 {
		evTimer.triggerTime = int(time.Now().UnixNano()/1e6) + evTimer.repeat
		evTimer.Start()
	}
	evTimer.handler(el, evTimer, revent)
}

func (evTimer *EvTimer) Init(el *EvLoop, handler HandlerFunc, at, repeat int, data interface{}) error {
	if el == nil {
		return errors.New("evLoop is nil")
	}
	evTimer.el = el
	evTimer.handler = handler
	evTimer.repeat = repeat
	evTimer.data = data
	evTimer.active = false
	evTimer.triggerTime = int(time.Now().UnixNano()/1e6) + at
	return nil
}

func (evTimer *EvTimer) Stop() error {
	evTimer.active = false
	for e := evTimer.el.timerList.Front(); e != nil; e = e.Next() {
		if e.Value.(*EvTimer) == evTimer {
			evTimer.el.timerList.Remove(e)
			break
		}
	}
	return nil
}

func (evTimer *EvTimer) Fd()int {
	return 0
}

func (evTimer *EvTimer) Start() error {
	evTimer.active = true
	for e := evTimer.el.timerList.Front(); e != nil; e = e.Next() {
		if e.Value.(*EvTimer).triggerTime >= evTimer.triggerTime {
			evTimer.el.timerList.InsertBefore(evTimer, e)
			return nil
		}
	}
	evTimer.el.timerList.PushBack(evTimer)
	return nil
}

func (evTimer *EvTimer) Data() interface{} {
	return evTimer.data
}

func (evTimer *EvTimer) IsActive() bool {
	return evTimer.active
}

type EvTimerHeap []*EvTimer

func (h EvTimerHeap) Len() int           { return len(h) }
func (h EvTimerHeap) Less(i, j int) bool { return h[i].triggerTime < h[j].triggerTime }
func (h EvTimerHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *EvTimerHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*EvTimer))
}
func (h *EvTimerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
