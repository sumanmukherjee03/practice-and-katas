/*
Dynamic programming solution for LCS - lowest common subsequence
*/

package main

import (
	"fmt"
	"math"
	"os"
)

var str1 string
var str2 string
var ssm *[][]int // solution space matrix

func initSpaceMatrix() {
	m := make([][]int, len(str1)+1)
	for i := range m {
		m[i] = make([]int, len(str2)+1)
	}
	ssm = &m
}

func nonRecursiveLCS() {
	for i := 1; i <= len(str1); i++ {
		for j := 1; j <= len(str2); j++ {
			if str1[i-1] == str2[j-1] {
				(*ssm)[i][j] = 1 + (*ssm)[i-1][j-1]
			} else {
				(*ssm)[i][j] = int(math.Max(float64((*ssm)[i-1][j]), float64((*ssm)[i][j-1])))
			}
		}
	}
}

func printSolution(i int, j int) {
	v := (*ssm)[i][j]

	if i-1 > 0 && j-1 > 0 {
		l := (*ssm)[i-1][j]
		t := (*ssm)[i][j-1]
		d := (*ssm)[i-1][j-1]
		switch v {
		case l:
			printSolution(i-1, j)
		case t:
			printSolution(i, j-1)
		case (d + 1):
			printSolution(i-1, j-1)
			fmt.Printf(string(str1[i-1]))
		}
	}
}

func main() {
	str1 = os.Args[1]
	str2 = os.Args[2]
	initSpaceMatrix()
	nonRecursiveLCS()
	printSolution(len(str1), len(str2))
}
