package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type data struct {
	value int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	describe()

	compute := func(c context.Context) <-chan data {
		out := make(chan data)
		go func() {
			defer close(out)

			d1 := data{241}
			out <- d1

			d2 := data{243}

			deadline, ok := c.Deadline()
			if ok {
				fmt.Println("Checking if there is a deadline in the context")
				if deadline.Sub(time.Now().Add(200*time.Millisecond)) < 0 {
					fmt.Println("Not enough time given for operation 2 in the goroutine")
				}
			}

			time.Sleep(200 * time.Millisecond) // delay sending the second value so that the deadline expires before that

			select {
			case out <- d2:
			case <-c.Done():
				fmt.Println("Context has expired, exiting goroutine")
				return
			}
		}()
		return out
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
	defer cancel() // We defer cancel irrespective of whether the deadline expires or not so that the resources associated with the context are released on program termination
	ch := compute(ctx)

	v1, ok := <-ch // only read 1 value from the channel because the deadline will trigger before the next value is sent to the channel
	if ok {        // check if channel has already closed or not. if open, then print the value obtained
		fmt.Println("Printing the first value : ", v1.value)
	}

	time.Sleep(300 * time.Millisecond) // Wait for some time to see if the deadline expires
	_, ok = <-ch                       // try to read the second value
	if !ok {                           // check if channel has already closed or not. if closed, then print a message.
		fmt.Println("The second value could not be obtained - context expired")
	}
}

func describe() {
	str := `
This is an example of a cancellable context with deadline. The cancel function is automatically called when the deadline expires.
When cancel function is called, it sends a done signal on the done channel of the context.
The child goroutine must handle the done signal and bail out to prevent goroutine leaks.

In this program the generator function returns a channel where it pushes random integers.
After the first 5 integers the main function which consumes from this channel should terminate the program.

_____________________
	`
	fmt.Println(str)
}
