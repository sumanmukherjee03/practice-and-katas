package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var maxItrs, err = strconv.ParseInt(os.Args[len(os.Args)-1], 0, 0)
	if err != nil {
		panic(err)
	}
	for i := 1; i <= int(maxItrs); i++ {
		switch {
		case i%3 == 0 && i%5 == 0:
			fmt.Print("FizzBuzz")
		case i%3 == 0:
			fmt.Print("Fizz")
		case i%5 == 0:
			fmt.Print("Buzz")
		default:
			fmt.Print(i)
		}
		fmt.Printf("\n")
	}
}
