package main

import (
	"fmt"
	"time"
)

func main() {
	describe()

	timer1 := time.NewTimer(2 * time.Second)

	<-timer1.C // This blocks program execution until after 2 seconds a timer event gets fired when things continue further
	fmt.Println("Timer 1 fired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C // The timer event in timer 2 gets cancelled, so this keeps waiting in the goroutine until the main proc exists
		fmt.Println("Timer 2 fired")
	}()
	stop := timer2.Stop() // This cancels the timer event
	if stop {
		fmt.Println("Timer 2 stopped")
	}

	time.Sleep(2 * time.Second) // Give enough time for timer 2 to fire
}

func describe() {
	str := `
For executing a task in the future as a one time event use a timer.

_____________________
	`
	fmt.Println(str)
}
