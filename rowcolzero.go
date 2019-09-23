package main

import (
	"fmt"
)

func genSampleMatrix() [][]int {
	return [][]int{
		{1, 2, 3, 4, 5},
		{6, 0, 7, 8, 9},
		{10, 11, 12, 0, 14},
		{15, 0, 16, 17, 18},
	}
}

func printMatrix(matrix [][]int) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Print(matrix[i][j])
			if j < len(matrix[i])-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func main() {
	matrix := genSampleMatrix()
	printMatrix(matrix)
}
