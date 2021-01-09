package main

import (
	"fmt"
	"os"
)

var str1 string
var str2 string
var solution string

func lcs(i int, j int) string {
	if i == len(str1) || j == len(str2) {
		return ""
	}

	if string(str1[i]) == string(str2[j]) {
		val := string(str1[i]) + lcs(i+1, j+1)
		if len(val) > len(solution) {
			solution = val
		}
		return val
	}

	v1 := lcs(i, j+1)
	v2 := lcs(i+1, j)
	if len(v1) > len(v2) {
		return v1
	}
	return v2
}

func main() {
	str1 = os.Args[1]
	str2 = os.Args[2]
	lcs(0, 0)
	fmt.Println(solution)
}
