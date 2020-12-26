package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup // wait group to synchronize

func main() {
	describe()
	var counter uint64

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				atomic.AddUint64(&counter, 1) // Use an atomic counter to make sure you can access a thread safe variable
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("counter : ", counter)
}

func describe() {
	str := `
Start 50 go routines and make each go routine increment a counter exactly 1000 times.
At the end we are supposed to get back exactly 50000 such operations.

The program shouldnt show any errors when race detection is enabled.
  For example using : counter += 1 // go run -race main.go - This will give a race condition error when executed

_____________________
	`
	fmt.Println(str)
}
