package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	describe()
	rand.Seed(time.Now().UnixNano())

	c1 := make(chan string, 1) // Use buffered channel to make sure that the send operation is non-blocking
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "hello_Channel_1"
	}()

	// Single select because we are publishing only 1 message in the channel above
	select {
	case m1 := <-c1:
		fmt.Println("Msg from channel 1 : ", m1)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout waiting for channel 1 to send anything")
	}

	c2 := make(chan string, 2) // Use buffered channel to make sure that the send operation is non-blocking
	for i := 0; i < 2; i++ {
		go func() {
			time.Sleep(2 * time.Second)
			c2 <- fmt.Sprintf("hello_Channel_2_%d", rand.Intn(100))
		}()
	}

	// 2 selects because we are publishing 2 messages in the channel above
	for j := 0; j < 2; j++ {
		select {
		case m2 := <-c2:
			fmt.Println("Msg from channel 2 : ", m2)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout waiting for channel 2 to send anything")
		}
	}
}

func describe() {
	str := `
This program depicts the use of timeouts for time bound execution of programs.
Note that the channel is buffered, so the send in the goroutine is nonblocking.
This is a common pattern to prevent goroutine leaks in case the channel is never read.

_____________________
	`
	fmt.Println(str)
}
