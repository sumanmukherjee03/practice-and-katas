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

	// Maintain 1 goroutine that runs in an infinite loop and performs select operation
	// to read from a channel of read ops or write ops
	// and publish a response in the response channel of the read ops or white ops struct.
	// Also, access the states variable only in this goroutine.
	// Any values that you get from the state should be passede back as a message in the response channel of the read/write ops.
	// For thread safety never access the states variable outside this goroutine.
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

	// Start 100 goroutines
	// Each goroutine runs an infinite loop to continuously read state, not directly but by message passing.
	// On each iteration create a read op and publish to the reads channel.
	// Then wait for a response to appear in the response channel of the read op
	for i := 0; i < 100; i++ {
		go func() {
			for {
				r := readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- r                           // Publish to reads channel which gets consumed by the goroutine above
				<-r.resp                             // Block and wait for the goroutine above to publish a response to the response channel of this read op
				atomic.AddUint64(&readOpsCounter, 1) // Increment atomic counter when you receive a successful response
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Start 10 goroutines
	// Each goroutine runs an infinite loop to continuously write state, not directly but by message passing.
	// On each iteration create a read op and publish to the reads channel.
	// Then wait for a response to appear in the response channel of the read op
	for j := 0; j < 10; j++ {
		go func() {
			for {
				w := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				writes <- w                           // Publish to writes channel which gets consumed by the goroutine above
				<-w.resp                              // Block and wait for the goroutine above to publish a response to the response channel of this write op
				atomic.AddUint64(&writeOpsCounter, 1) // Increment atomic counter when you receive a successful response
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	// Use atomic.LoadUint64 to safely read the final values of the read and write ops counters
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

_____________________
	`
	fmt.Println(str)
}
