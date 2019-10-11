// Valid input string : ...1.5.68......7.19.1....3...7.26...5.......3...87.4...3....8.51.5......79.4.1...
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput(str string) [9][9]int {
	matrix := [9][9]int{}
	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanRunes)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			scanner.Scan()
			x, _ := strconv.Atoi(scanner.Text())
			matrix[i][j] = x
		}
	}
	return matrix
}

func printMatrix(matrix [9][9]int) {
	for _, m := range matrix {
		for _, v := range m {
			fmt.Print(v)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func validateDupesInRowsCols(m [9][9]int) error {
	intMap := make(map[int][]int)
	for i, r := range m {
		for j, v := range r {
			if found, ok := intMap[v]; ok && v > 0 && (found[0] == i || found[1] == j) {
				return errors.New("There are dupes in rows or cols")
			}
			intMap[v] = []int{i, j}
		}
	}
	return nil
}

func validateDupesIn3x3Matrix(m [9][9]int, startRow int, startCol int) error {
	intMap := make(map[int][]int)
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			v := m[i][j]
			if _, ok := intMap[v]; ok && v > 0 {
				return errors.New("There are dupes in 3x3 matrix")
			}
			intMap[v] = []int{i, j}
		}
	}
	return nil
}

func validateAll3x3Matrix(m [9][9]int) error {
	for i := 0; i < 9; i = i + 3 {
		for j := 0; j < 9; j = j + 3 {
			err := validateDupesIn3x3Matrix(m, i, j)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	m := parseInput(os.Args[1])
	fmt.Println("Initial matrix for sudoku :")
	printMatrix(m)
	err := validateDupesInRowsCols(m)
	if err != nil {
		panic(err)
	}
	err = validateAll3x3Matrix(m)
	if err != nil {
		panic(err)
	}
	fmt.Println("Matrix for sudoku is valid")
	fmt.Println("Final solution :")
}
