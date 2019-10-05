// Run : go run rowcolzero.go | column -t -s ' '
package main

import (
	"fmt"
	"math/rand"
)

func genSampleMatrix(rows int, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = rand.Int() % 11
		}
	}
	return matrix
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
	matrix := genSampleMatrix(5, 4)
	printMatrix(matrix)
}
