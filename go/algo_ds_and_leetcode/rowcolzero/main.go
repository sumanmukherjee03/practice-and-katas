// Run : go run rowcolzero.go 5 4 | column -t -s ' '
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
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

func findRowsColsWithZero(matrix [][]int) (map[int]bool, map[int]bool) {
	var rowsWithZero = make(map[int]bool)
	var colsWithZero = make(map[int]bool)
	for i := range matrix {
		if _, foundRow := rowsWithZero[i]; !foundRow {
			for j := range matrix[i] {
				if _, foundCol := colsWithZero[j]; !foundCol {
					if matrix[i][j] == 0 {
						rowsWithZero[i] = true
						colsWithZero[j] = true
					}
				}
			}
		}
	}
	return rowsWithZero, colsWithZero
}

func main() {
	rows, err1 := strconv.ParseInt(os.Args[1], 0, 0)
	if err1 != nil {
		panic(err1)
	}
	cols, err2 := strconv.ParseInt(os.Args[2], 0, 0)
	if err2 != nil {
		panic(err2)
	}
	matrix := genSampleMatrix(int(rows), int(cols))
	printMatrix(matrix)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>")
	rowsWithZero, colsWithZero := findRowsColsWithZero(matrix)
	for r := range rowsWithZero {
		fmt.Print(r)
		fmt.Print(",")
	}
	fmt.Println()
	for c := range colsWithZero {
		fmt.Print(c)
		fmt.Print(",")
	}
	fmt.Println()
}
