package main

import (
	"fmt"
	"sync"
)

func main() {
	describe()
	variant1()
}

func variant1() {
	jobs := make(chan string)
	done := make(chan bool)
	go func() {
		for {
			// In this special 2-value form of receive,
			// the more value will be false if jobs has been closed and
			// all values in the channel have already been received.
			if j, more := <-jobs; more {
				fmt.Println("Processed job : ", j)
			} else {
				fmt.Println("No more jobs to process")
				done <- true
				return
			}
		}
	}()

	for i := 0; i < 3; i++ {
		jobs <- fmt.Sprintf("job_%d", i)
	}
	close(jobs) // Remember to close channel after the job is done
	<-done      // Finally read from the done channel to indicate end of job
}

// This is what you can avoid doing with the better variant above
//   - by having the done complete inside the goroutine by publishing to the done channel inside the goroutine
//     and reading from it in main proc
// because otherwise you will require wait groups to make the goroutines wait before the main function completes
func bad_variant2() {
	jobs := make(chan string)
	done := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for {
			select {
			case j := <-jobs:
				fmt.Println("Processed job : ", j)
			case <-done:
				fmt.Println("No more jobs to process - closing jobs")
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 3; i++ {
			jobs <- fmt.Sprintf("job_%d", i)
		}
		close(jobs)
		done <- true
		wg.Done()
	}()

	wg.Wait()
}

func describe() {
	str := `
Use a jobs channel to do some work sent from the main() goroutine to a worker goroutine.
When there are no more jobs to process close the jobs channel.

_____________________
	`
	fmt.Println(str)
}
