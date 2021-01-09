package main

import "fmt"

func main() {
	describe()
	for n := range square(square(generator(2, 3))) {
		fmt.Println(n)
	}
}

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

func describe() {
	str := `
This is a simple example of a pipeline to generate some random numbers and then square them.
A big feature of pipelines is composability, so we should attempt to keep that in mind when generating a solution.

_____________________
	`
	fmt.Println(str)
}
