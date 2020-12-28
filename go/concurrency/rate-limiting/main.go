package main

import (
	"fmt"
	"time"
)

func main() {
	describe()

	fmt.Println("Simple rate limit example : ")
	// Make this channel buffered otherwise you cant push requests into it
	// because a receiver doesnt exist yet. The receiver is defined in the func simpleRateLimit.
	// If the channel is not buffered you will get an error - fatal error: all goroutines are asleep - deadlock!
	requests := make(chan string, 10)
	for i := 0; i < 10; i++ {
		requests <- fmt.Sprintf("Request_%d", i)
	}
	close(requests)
	simpleRateLimit(requests)

	fmt.Println("\n\nBurstable rate limit example : ")
	// Make this channel buffered otherwise you cant push requests into it
	// because a receiver doesnt exist yet. The receiver is defined in the func simpleRateLimit.
	// If the channel is not buffered you will get an error - fatal error: all goroutines are asleep - deadlock!
	burstableRequests := make(chan string, 10)
	for j := 0; j < 10; j++ {
		burstableRequests <- fmt.Sprintf("BurstableRequest_%d", j)
	}
	close(burstableRequests)
	burstRateLimit(burstableRequests)
}

func simpleRateLimit(reqs chan string) {
	// Create a ticker
	limiter := time.Tick(100 * time.Millisecond)

	// Consume requests
	for req := range reqs {
		<-limiter // Block for receive from the limiter
		fmt.Println(time.Now(), "Request processed : ", req)
	}
}

func burstRateLimit(reqs chan string) {
	// Create a channel of time that serves as an enhanced ticker
	// with a buffer of 3 to handle bursts of requests
	limiter := make(chan time.Time, 3)

	// Initiate the limiter with 3 immediate time stamps so that the first 3 requests can be immediately processed
	for k := 0; k < 3; k++ {
		limiter <- time.Now()
	}

	// Start a goroutine that reads from a ticker and pushes a timestamp into our rate limiter every 100ms
	// Since we are not directly ranging over the limiter, this goroutine can die when main exits.
	// Nothing should cause a deadlock in that case. Also, no cleanup is necessary.
	go func() {
		ticker := time.Tick(100 * time.Millisecond)
		for t := range ticker {
			limiter <- t
		}
	}()

	// Consume requests
	for req := range reqs {
		<-limiter // Block for receive from the limiter
		fmt.Println(time.Now(), "Request processed : ", req)
	}
}

func describe() {
	str := `
This example demonstrates simple rate limiting.
Say, the rate limit is 10 requests per second, or 1 every 100ms.
If 2 requests arrive within 100ms time span, they are not served immediately,
but rather at the interval of 1 request every 100 ms.

This example also demonstrates rate limiting with burst.
This is useful when we may want to allow short bursts of requests in our rate limiting scheme,
while preserving the overall rate limit.
If 5 requests arrive, and your rate limit is 1 request every 100ms and burst limit is 3 requests,
then the first 3 requests are served immediately. The next 2 are served as per the rate limit,
ie 1 every 100 ms.

_____________________
	`
	fmt.Println(str)
}
