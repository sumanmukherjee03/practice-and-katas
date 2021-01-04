package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(4) // set the max number of processors to use
	var balance = 0
	var wg = sync.WaitGroup{}
	var mutex = sync.Mutex{}

	for i := 0; i < 100; i++ {
		// Do not add to a waitgroup inside a goroutine because the main function can exit before a goroutine is scheduled to be run by the scheduler
		wg.Add(1)
		go func() {
			// In deferred execution to make sure that the mutex unlock happens first and then the waitgroup done
			// we put things in the reverse order in the code - ie wg.Done() first and then mutex.Unlock()
			defer wg.Done()      // Ensure that waitgroup done gets called when goroutine exits
			defer mutex.Unlock() // Ensure mutex is unlocked when goroutine exits
			mutex.Lock()         // Everything below this line in this function is the critical section
			balance--            // balance is the state with shared memory
		}()
	}

	for j := 0; j < 100; j++ {
		// Do not add to a waitgroup inside a goroutine because the main function can exit before a goroutine is scheduled to be run by the scheduler
		wg.Add(1)
		go func() {
			// In deferred execution to make sure that the mutex unlock happens first and then the waitgroup done
			// we put things in the reverse order in the code - ie wg.Done() first and then mutex.Unlock()
			defer wg.Done()      // Ensure that waitgroup done gets called when goroutine exits
			defer mutex.Unlock() // Ensure mutex is unlocked when goroutine exits
			mutex.Lock()         // Everything below this line in this function is the critical section
			balance++            // balance is the state with shared memory
		}()
	}

	wg.Wait()
	fmt.Println("Final balance : ", balance)
}

func describe() {
	str := `
This simple example demonstrates mutexes with a bank account example where we deposit and withdraw money.

_____________________
	`
	fmt.Println(str)
}
