package main

import (
	"fmt"
	"time"
)

func main() {
	describe()

	const numOfJobs = 10
	const numOfWorkers = 3
	jobs := make(chan int, numOfJobs)
	results := make(chan int, numOfJobs)

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
	// You cant range over the results because we are closing the results channel anywhere
	for k := 1; k <= numOfJobs; k++ {
		r := <-results
		fmt.Println("Processed job result", r)
	}
}

func worker(id int, jobs <-chan int, res chan<- int) {
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
