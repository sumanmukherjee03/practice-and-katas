package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	describe()
	rand.Seed(time.Now().UnixNano())

	var state = make(map[int]int)
	var mutex = &sync.Mutex{}
	var readOps uint64
	var writeOps uint64

	// Run 100 goroutines to read state from variable
	for i := 0; i < 100; i++ {
		go func() {
			total := 0
			for { // Continue infinitely
				key := rand.Intn(5)
				mutex.Lock() // Acquire lock on state var
				val := state[key]
				total += val                  // Just something to use the read value
				mutex.Unlock()                // Release lock on state var
				atomic.AddUint64(&readOps, 1) // Use atomic counter to maintain count
				time.Sleep(time.Duration(1000))
			}
		}()
	}

	// Run 100 goroutines to read state from variable
	for j := 0; j < 10; j++ {
		go func() {
			for { // Continue infinitely
				key := rand.Intn(5)
				mutex.Lock() // Acquire lock on state var
				state[key] = rand.Intn(100)
				mutex.Unlock()                 // Release lock on state var
				atomic.AddUint64(&writeOps, 1) // Use atomic counter to maintain count
				time.Sleep(time.Duration(1000))
			}
		}()
	}

	time.Sleep(10 * time.Second)
	readOpsFinal := atomic.LoadUint64(&readOps)   // Read final value of counter using atomic.LoadUint64
	writeOpsFinal := atomic.LoadUint64(&writeOps) // Read final value of counter using atomic.LoadUint64
	fmt.Println("readOps : ", readOpsFinal)
	fmt.Println("writeOps : ", writeOpsFinal)
	mutex.Lock()
	fmt.Println("state : ", state)
	mutex.Unlock()
}

func describe() {
	str := `
Start 10 go routines to write data to a map which is a shared state between multiple go routines.
Start 100 go routines to read state from the map which is the shared state variable.
Increment a counter for the read operations.
Also, increment the counters for write operations.
Remember, we arent using waitgroups in this example.
Finally, report the read and write operations.
Also, report the final state by taking a lock on the state.

_____________________
	`
	fmt.Println(str)
}
