package gevloop

import (
	"container/heap"
	"errors"
	"syscall"
)

type EvTimer struct {
	el      *EvLoop
	at      int
	active  bool
	repeat  int
	handler HandlerFunc
	data    interface{}
}

func (evTimer *EvTimer) cb(el *EvLoop, revent uint32) {
	revent = syscall.SYS_TIMES
	evTimer.handler(el, evTimer, revent)
}

func (evTimer *EvTimer) Init(el *EvLoop, handler HandlerFunc, at, repeat int, data interface{}) error {
	if el == nil {
		return errors.New("evLoop is nil")
	}
	evTimer.el = el
	evTimer.at = at
	evTimer.handler = handler
	evTimer.repeat = repeat
	evTimer.data = data
	evTimer.active = false
	return nil
}

func (evTimer *EvTimer) Stop() {
	evTimer.active = false
	for i := 0; i < evTimer.el.timerHeap.Len(); i++ {
		n := (*evTimer.el.timerHeap)[i]
		if n == evTimer {
			(*evTimer.el.timerHeap)[i], (*evTimer.el.timerHeap)[evTimer.el.timerHeap.Len()-1] =
				(*evTimer.el.timerHeap)[evTimer.el.timerHeap.Len()-1], (*evTimer.el.timerHeap)[i]
			*evTimer.el.timerHeap = (*evTimer.el.timerHeap)[:evTimer.el.timerHeap.Len()-1]
			i--
		}
	}

	heap.Init(evTimer.el.timerHeap)
}

func (evTimer *EvTimer) Start() {
	evTimer.active = true
	heap.Push(evTimer.el.timerHeap, evTimer)
}

func (evTimer *EvTimer) IsActive() bool {
	return evTimer.active
}

type EvTimerHeap []*EvTimer

func (h EvTimerHeap) Len() int           { return len(h) }
func (h EvTimerHeap) Less(i, j int) bool { return h[i].at < h[j].at }
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
