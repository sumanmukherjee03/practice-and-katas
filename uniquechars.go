package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	foundDupePos := -1
	str := os.Args[1]
	for pos, r := range str {
		strPart := str[0:pos]
		if strings.Contains(strPart, string(r)) {
			foundDupePos = pos
			break
		}
	}
	if foundDupePos > 0 {
		fmt.Println(string(str[foundDupePos]))
	} else {
		fmt.Println("No dupes found")
	}
}
