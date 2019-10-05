package main

import (
	"fmt"
	"os"
	"strconv"
)

func fibonacci(num1 int64, num2 int64, limit int64, c chan int64) {
	if num1 == 0 && num2 == 1 {
		c <- num1
		c <- num2
	}
	if len(c) == int(limit) {
		close(c)
		return
	}
	next := num1 + num2
	c <- next
	go fibonacci(num2, next, limit, c)
}

func main() {
	var maxItrs, err = strconv.ParseInt(os.Args[len(os.Args)-1], 0, 0)
	if err != nil {
		panic(err)
	}
	numbers := make(chan int64, maxItrs)
	go fibonacci(0, 1, maxItrs, numbers)
	for x := range numbers {
		fmt.Println(x)
	}
}
