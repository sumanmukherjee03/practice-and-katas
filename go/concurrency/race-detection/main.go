package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	describe()
	variant2()
}

// func variant1() {
// start := time.Now()
// var t *time.Timer
// t = time.AfterFunc(randomDuration(), func() {
// fmt.Println(time.Now().Sub(start))
// t.Reset(randomDuration())
// })
// time.Sleep(5 * time.Second)
// }

func variant2() {
	start := time.Now()
	var t *time.Timer
	ch := make(chan bool)

	// time.AfterFunc calls the callback function after a specified duration in it's own goroutine
	// and returns a timer object. You can call Stop on the timer to cancel
	t = time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		// Resetting the timer here causes race condition between the main proc and the goroutine in which this callback runs.
		// Instead of resetting the timer here use channel for threadsafe processing
		// Due to the message passing, the scope of the variable "t" gets confined to one goroutine only - the main proc.
		ch <- true
	})

	for time.Since(start) < 5*time.Second {
		if ok := <-ch; ok {
			t.Reset(randomDuration())
		}
	}
}

// returns a random duration between 0 - 1 second
func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}

func describe() {
	str := `
You can detect race conditions with the -race flag.
It can be used in run, build, test etc.
This example demonstrates 2 variants - one with a race condition and the other without a race condition.
The one with the race condition is commented out.

_____________________
	`
	fmt.Println(str)
}
