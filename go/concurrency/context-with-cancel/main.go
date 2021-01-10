package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	describe()

	generator := func(c context.Context) <-chan int {
		out := make(chan int)

		go func() {
			defer close(out) // Defer close the output channel here so that the channel closes no matter the return path from this goroutine
			for {
				// Allow pre-emptive termination of the infinite loop with a combination of select and done channel
				select {
				case out <- rand.Intn(1000):
				case <-c.Done():
					return // bail out if program termination is indicated using the done channel
				}
			}
		}()

		return out
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := generator(ctx)
	count := 0
	for x := range ch {
		if count >= 5 {
			cancel() // indicate program termination by calling cancel when condition is met
			return   // bail out of the for loop
		}
		fmt.Println(x)
		count++
	}
}

func describe() {
	str := `
This is an example of a cancellable context. The cancel function is generally called with a defer where the cancellable context is defined.
When cancel function is called, it sends a done signal on the done channel of the context.
The child goroutine must handle the done signal and bail out to prevent goroutine leaks.

In this program the generator function returns a channel where it pushes random integers.
After the first 5 integers the main function which consumes from this channel should terminate the program.

_____________________
	`
	fmt.Println(str)
}
