package main

import (
	"github.com/xinxuwang/gevloop"
	"log"
	"os"
)

func main() {
	el, err := gevloop.Init()
	if err != nil {
		os.Exit(0)
	}

	timer1 := gevloop.EvTimer{}
	timer1Data := int32(0)
	timer1.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
		log.Println("timer1 Called", timer1, revent)
		data := event.Data().(*int32)
		*data = *data + 1
		if *data == 3 {
			log.Println("timer1 Called 3 times,stop")
			timer1.Stop()
		}
	}, 1000, 2000, &timer1Data)
	timer1.Start()

	timer2 := gevloop.EvTimer{}
	timer2Data := int32(0)
	timer2.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
		log.Println("timer2 Called", timer2, revent)
		data := event.Data().(*int32)
		*data = *data + 1
		if *data == 3 {
			log.Println("timer2 Called 3 times,start timer3")
			timer3 := gevloop.EvTimer{}
			timer3.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
				log.Println("timer3 Called", timer3, revent)
			}, 1000, 1000, nil)
			timer3.Start()
		}
	}, 2000, 4000, &timer2Data)
	timer2.Start()

	timer4 := gevloop.EvTimer{}
	timer4.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
		log.Println("timer4 Called", timer4, revent)
	}, 3000, 0, &timer1Data)
	timer4.Start()
	log.Println("Loop Run:")
	err = el.Run()
	if err != nil {
		log.Println("error:", err)
	}
}
