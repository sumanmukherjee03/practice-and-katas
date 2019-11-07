package main

import (
	"fmt"
	"os"
)

func unique(s []int) []int {
	keys := make(map[int]bool)
	l := []int{}
	for _, i := range s {
		if _, v := keys[i]; !v {
			keys[i] = true
			l = append(l, i)
		}
	}
	return l
}

func findAggregatedPath(path1 string, path2 string) string {
	var solution = ""

	var lcs func(int, int) string
	lcs = func(i int, j int) string {
		if i == len(path1) || j == len(path2) {
			return ""
		}
		if string(path1[i]) == string(path2[j]) {
			val := string(path1[i]) + lcs(i+1, j+1)
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
	lcs(0, 0)

	var str string
	if len(path1) > len(path2) {
		str = path1
	} else {
		str = path2
	}

	var mismatches []int
	var pointsOfMismatch func(int, int)
	pointsOfMismatch = func(i int, j int) {
		if i == len(str) || j == len(solution) {
			return
		}
		if string(solution[j]) == "*" {
			pointsOfMismatch(i, j+1)
			return
		}
		if string(str[i]) != string(solution[j]) {
			mismatches = append(mismatches, j)
			pointsOfMismatch(i+1, j)
			return
		}
		pointsOfMismatch(i+1, j+1)
	}
	pointsOfMismatch(0, 0)
	mismatches = unique(mismatches)

	starcount := 0
	for i := 0; i < len(mismatches); i++ {
		x := mismatches[i] + starcount
		val := solution[:x] + "*" + solution[x:]
		starcount = starcount + 1
		solution = val
	}

	return solution
}

func main() {
	str1 := os.Args[1]
	str2 := os.Args[2]
	res := findAggregatedPath(str1, str2)
	fmt.Println(res)
}
