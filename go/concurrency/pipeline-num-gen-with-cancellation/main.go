package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	describe()

	done := make(chan struct{}) // Create a done channel of empty struct as we just want to send the close signal and not any data
	in := generator(done, 2, 3, 4, 5)

	// This is the fan out section.
	// We assume that square is a computationally expensive function, so are creating multiple square goroutines.
	ch1 := square(done, in)
	ch2 := square(done, in)

	// Make sure the channel that merge returns is closed in merge, otherwise this will block
	m := merge(done, ch1, ch2)
	fmt.Println(<-m)
	close(done) // Close the done channel to trigger cancellation after receiving 1 value

	time.Sleep(10 * time.Millisecond)                // allow some time to terminate the goroutines
	g := runtime.NumGoroutine()                      // Returns the number of active goroutines
	fmt.Println("Number of active goroutines : ", g) // We will see only 1 in this output - the main goroutine
}

// Remember to close the output channel here to make sure the next stage of the pipeline is unblocked
func generator(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // We make the closing of the channel deferred to ensure that the out channel closes on all return paths, including the pre-emptive cancel
		for _, num := range nums {
			select {
			case out <- num:
			case <-done:
				return
			}
		}
	}()
	return out
}

// Remember to close the output channel here to make sure the next stage of the pipeline is unblocked
func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // We make the closing of the channel deferred so that we close the out channel from all return paths
		for num := range in {
			select {
			case out <- num * num:
			case <-done:
				return
			}
		}
	}()
	return out
}

// This function performs the fan in section
// We need a wait group to synchronize that all the fan in goroutines complete before finishing this function
// Also, make sure to close the fan in channel to unblock the main function
func merge(done <-chan struct{}, chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		defer wg.Done() // Make this a deffered call, so that it gets called on all return paths, even if we cancel and pre-emptively return
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	for _, ch := range chans {
		wg.Add(1)
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func describe() {
	str := `
This is a simple example of a pipeline to generate some random numbers and then square them.
A big feature of pipelines is composability, so we should attempt to keep that in mind when generating a solution.

             ----square----
             |             |
generators---|             |---merge
             |             |
             ---square-----

_____________________
	`
	fmt.Println(str)
}
