package main

import (
	"fmt"
	"os"
)

func main() {
	str := os.Args[1]
	for i, _ := range str {
		fmt.Print(string(str[len(str)-1-i]))
	}
}
