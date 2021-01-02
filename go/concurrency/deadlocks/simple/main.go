package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock1 = sync.Mutex{}
	lock2 = sync.Mutex{}
	wg    = sync.WaitGroup{}
)

func main() {
	describe()
	wg.Add(2)
	go blueRobot()
	go redRobot()
	wg.Wait()
}

func blueRobot() {
	// for {
	fmt.Println("Blue : Acquiring lock 1")
	lock1.Lock()
	fmt.Println("Blue : Acquired lock 1")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Blue : Acquiring lock 2")
	lock2.Lock()
	fmt.Println("Blue : Acquired lock 2")
	fmt.Println("Blue : Both locks acquired")
	lock1.Unlock()
	lock2.Unlock()
	fmt.Println("Blue : Both locks released")
	// }
	wg.Done()
}

func redRobot() {
	// for {
	fmt.Println("Red : Acquiring lock 2")
	lock2.Lock()
	fmt.Println("Red : Acquired lock 2")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Red : Acquiring lock 1")
	lock1.Lock()
	fmt.Println("Red : Acquired lock 1")
	fmt.Println("Red : Both locks acquired")
	lock2.Unlock()
	lock1.Unlock()
	fmt.Println("Red : Both locks released")
	// }
	wg.Done()
}

func describe() {
	str := `
This example depicts a simple deadlock with both goroutines stuck in a circular dependency.

_____________________
	`
	fmt.Println(str)
}
