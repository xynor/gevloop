package main

import (
	"fmt"
	"github.com/xinxuwang/gevloop"
	"log"
	"net"
	"syscall"
)

type session struct {
	bytes []byte
	pos   int
}

func main() {
	accept, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal("err:", err)
	}
	defer syscall.Close(accept)

	if err = syscall.SetNonblock(accept, true); err != nil {
		log.Fatal("Set noblock err:", err)
	}
	addr := syscall.SockaddrInet4{Port: 2000}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())

	if err := syscall.Bind(accept, &addr); err != nil {
		log.Fatal("Bind err:", err)
	}
	if err := syscall.Listen(accept, 10); err != nil {
		log.Fatal("Listen err:", err)
	}
	el, err := gevloop.Init()
	if err != nil {
		log.Fatal("err:", err)
	}
	log.Println("Accept fd:", accept)
	acceptIO := gevloop.EvIO{}
	acceptIO.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
		log.Println("AcceptIO Called")
		connFd, _, err := syscall.Accept(event.Fd())
		if err != nil {
			log.Println("accept: ", err)
			return
		}
		syscall.SetNonblock(connFd, true)
		connFdIO := gevloop.EvIO{}
		sess := session{
			bytes: make([]byte, 5),
			pos:   0,
		}
		connFdIO.Init(el, func(evLoop *gevloop.EvLoop, event gevloop.Event, revent uint32) {
			log.Println("connFdIO Called")
			//assume `HELLO`
			buf := make([]byte, 0)
			for {
				nbytes, e := syscall.Read(event.Fd(), buf)
				sess := event.Data().(*session)
				if e != nil {
					log.Println("Read Error:", e)
					return
				}

				if nbytes > 0 {
					fmt.Println("Read n:", nbytes)
					sess.pos = nbytes
					copy(sess.bytes[sess.pos:], buf)
					if 5 == len(sess.bytes) {
						log.Println(string(sess.bytes))
						sess.pos = 0
						return
					}
				}
				if nbytes == 0 {
					log.Println("close")
					syscall.Close(event.Fd())
					event.Stop()
					return
				}
				fmt.Println("Read < 0")
				return
			}
		}, connFd, syscall.EPOLLIN, &sess)
		connFdIO.Start()
	}, accept, syscall.EPOLLIN|syscall.EPOLLET&0xffffffff, nil)

	acceptIO.Start()
	err = el.Run()
	if err != nil {
		log.Println("error:", err)
	}
}
