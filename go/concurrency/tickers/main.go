package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Ticker fired")
			case <-done:
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

func describe() {
	str := `
To keep executing a task repeatedly at an interval use a ticker.

_____________________
	`
	fmt.Println(str)
}
