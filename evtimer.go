package gevloop

type EvTimer struct {
	at      int
	active  bool
	repeat  int
	handler HandlerFunc
	data    interface{}
}

func (evTimer *EvTimer) cb(el *evLoop, revent uint32) {
	evTimer.handler(el, evTimer, revent)
}

func (evTimer *EvTimer) Stop() error {
	return nil
}

func (evTimer *EvTimer) Start() error {
	return nil
}

func (evTimer *EvTimer) IsActive() bool {
	return true
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
