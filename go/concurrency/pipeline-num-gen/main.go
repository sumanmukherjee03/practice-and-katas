package main

import (
	"fmt"
	"sync"
)

func main() {
	describe()

	in := generator(2, 3, 4, 5)

	// This is the fan out section.
	// We assume that square is a computationally expensive function, so are creating multiple square goroutines.
	// Take note of the composability of the square goroutine that we have used here.
	ch1 := square(square(in))
	ch2 := square(square(in))

	// Make sure the channel that merge returns is closed in merge, otherwise this will block
	for n := range merge(ch1, ch2) {
		fmt.Println(n)
	}
}

// Remember to close the output channel here to make sure the next stage of the pipeline is unblocked
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, num := range nums {
			out <- num
		}
		close(out)
	}()
	return out
}

// Remember to close the output channel here to make sure the next stage of the pipeline is unblocked
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for num := range in {
			out <- num * num
		}
		close(out)
	}()
	return out
}

// This function performs the fan in section
// We need a wait group to synchronize that all the fan in goroutines complete before finishing this function
// Also, make sure to close the fan in channel to unblock the main function
func merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
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
