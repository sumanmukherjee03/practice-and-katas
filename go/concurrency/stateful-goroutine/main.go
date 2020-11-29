package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type readOp struct {
	key  int
	resp chan int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func main() {
	describe()
	rand.Seed(time.Now().UnixNano())
	var readOpsCounter uint64
	var writeOpsCounter uint64
	reads := make(chan readOp)
	writes := make(chan writeOp)

	go func() {
		states := make(map[int]int)
		for {
			select {
			case r := <-reads:
				r.resp <- states[r.key]
			case w := <-writes:
				states[w.key] += w.val
				w.resp <- true
			}
		}
	}()

	for i := 0; i < 100; i++ {
		go func() {
			for {
				r := readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- r
				<-r.resp
				atomic.AddUint64(&readOpsCounter, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for j := 0; j < 10; j++ {
		go func() {
			for {
				w := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				writes <- w
				<-w.resp
				atomic.AddUint64(&writeOpsCounter, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)
	readOpsFinal := atomic.LoadUint64(&readOpsCounter)
	writeOpsFinal := atomic.LoadUint64(&writeOpsCounter)
	fmt.Println("readOps : ", readOpsFinal)
	fmt.Println("writeOps : ", writeOpsFinal)
}

func describe() {
	str := `
The state will be owned by a single goroutine.
This will guarantee that the data is never corrupted with concurrent access.
In order to read or write that state, other goroutines will send messages to the owning goroutine and receive corresponding replies.
These readOp and writeOp structs encapsulate those requests and a resp channel for the owning goroutine to respond.
	`
	fmt.Println(str)
}
