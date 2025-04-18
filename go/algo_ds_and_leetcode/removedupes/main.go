package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	str := os.Args[1]
	if len(str) < 2 {
		fmt.Println(str)
		return
	}
	for i := 0; i < len(str); i++ {
		if strings.Contains(str[0:i], string(str[i])) {
			continue
		}
		fmt.Print(string(str[i]))
	}
}
