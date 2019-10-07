package main

import (
	"fmt"
	"os"
	"strconv"
)

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func main() {
	var maxItrs, err = strconv.ParseInt(os.Args[len(os.Args)-1], 0, 0)
	if err != nil {
		panic(err)
	}
	numbers := make(chan int)
	go fibonacci(int(maxItrs), numbers)
	for x := range numbers {
		fmt.Println(x)
	}
}
