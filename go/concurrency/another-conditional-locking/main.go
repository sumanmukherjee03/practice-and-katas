package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg             = sync.WaitGroup{}
	mutex          = sync.Mutex{}
	cond           = sync.NewCond(&mutex)
	sharedResource = make(map[string]interface{})
)

func main() {
	ch := make(chan interface{}, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cond.L.Unlock()
		cond.L.Lock()
		for len(sharedResource) == 0 {
			cond.Wait()
		}
		ch <- sharedResource["foo"]
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer mutex.Unlock()
		time.Sleep(2 * time.Second)
		mutex.Lock()
		sharedResource["foo"] = "bar"
		cond.Signal()
	}()

	select {
	case v := <-ch:
		fmt.Println("Found value : ", v)
	case <-time.After(3 * time.Second):
		fmt.Println("timed out")
	}
	wg.Wait()
}
