package main

import (
	"fmt"
	"os"
)

func main() {
	var found []byte
	str := os.Args[1]
	if len(str) == 1 {
		fmt.Println(str)
		return
	}
	found = append(found, str[0])
	for i := 1; i < len(str); i++ {
		if str[i] == found[len(found)-1] {
			continue
		} else {
			found = append(found, str[i])
		}
	}
	fmt.Println(string(found))
}
