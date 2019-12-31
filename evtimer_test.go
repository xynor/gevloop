package gevloop

import (
	"container/heap"
	"fmt"
	"testing"
)

func Test_EvTimerHeap(t *testing.T) {
	timer1 := EvTimer{
		at: 2,
	}
	timer2 := EvTimer{
		at: 1,
	}
	timer3 := EvTimer{
		at: 5,
	}
	timer4 := EvTimer{
		at: 3,
	}
	timer5 := EvTimer{
		at: 3,
	}
	h := &EvTimerHeap{&timer1, &timer2, &timer3, &timer5}
	//h := &EvTimerHeap{}
	heap.Init(h)
	heap.Push(h, &timer4)
	fmt.Println("minimum: ", (*h)[0].at)
	fmt.Println(h)
	toRemove := &timer3
	for i := 0; i < h.Len(); i++ {
		n := (*h)[i]
		if n == toRemove {
			(*h)[i], (*h)[h.Len()-1] = (*h)[h.Len()-1], (*h)[i]
			*h = (*h)[:h.Len()-1]
			i--
		}
	}

	heap.Init(h)
	for h.Len() > 0 {
		fmt.Println(" ", *(heap.Pop(h).(*EvTimer)))
	}
}
