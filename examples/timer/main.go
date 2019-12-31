package main

import (
	"fmt"
	"github.com/xinxuwang/gevloop"
	"os"
)

func main() {
	el, err := gevloop.Init()
	if err != nil {
		os.Exit(0)
	}

	timer1 := gevloop.EvTimer{}
	timer1.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
		fmt.Println("timer1 Called", timer1, revent)
	}, 2000, 2000, nil)
	timer1.Start()

	timer2 := gevloop.EvTimer{}
	timer1.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
		fmt.Println("timer2 Called", timer2, revent)
	}, 2000, 4000, nil)
	timer2.Start()

	err = el.Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}
