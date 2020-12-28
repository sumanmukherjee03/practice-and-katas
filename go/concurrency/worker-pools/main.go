package main

import (
	"fmt"
	"time"
)

func main() {
	describe()

	const numOfJobs = 10
	const numOfWorkers = 3
	jobs := make(chan int, numOfJobs)    // A channel to track the job which is consumed by the workers
	results := make(chan int, numOfJobs) // A results channel where the results of the jobs are published

	// Create a worker pool
	for i := 1; i <= numOfWorkers; i++ {
		go worker(i, jobs, results)
	}

	// Push jobs into the jobs channel which gets picked up by the workers in goroutines
	for j := 1; j <= numOfJobs; j++ {
		jobs <- j
	}
	close(jobs) // Close the jobs channel when you are done creating the jobs

	// Finally, extract the results for each job
	// You cant range over the results because we are not closing the results channel anywhere.
	// If we ranged over this channel we would have ended up getting an error - fatal error: all goroutines are asleep - deadlock!
	// This is because the range operator would keep waiting for more input from the channel
	for k := 1; k <= numOfJobs; k++ {
		r := <-results
		fmt.Println("Processed job result", r)
	}
}

func worker(id int, jobs <-chan int, res chan<- int) {
	// You can range over the jobs channel because the jobs channel eventually gets closed in main
	for j := range jobs {
		fmt.Println("Worker", id, "Started processing job", j)
		time.Sleep(200 * time.Millisecond)
		res <- j * 2
		fmt.Println("Worker", id, "finished processing job", j)
	}
}

func describe() {
	str := `
This example enumerates how you can create a worker pool.

_____________________
	`
	fmt.Println(str)
}
