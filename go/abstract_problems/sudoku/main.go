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

func printMatrix(m *[9][9]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			v := m[i][j]
			fmt.Print(v)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func findDupRowColForInt(a []int, i int, j int) bool {
	for _, v := range a {
		r := v / 10
		c := v % 10
		if r == i || c == j {
			return true
		}
	}
	return false
}

func validateDupesInRowsCols(m *[9][9]int) error {
	intMap := make(map[int][]int)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			v := m[i][j]
			if v > 0 {
				if found, ok := intMap[v]; ok && findDupRowColForInt(found, i, j) {
					return errors.New("there are dupes in rows or cols")
				}
				intMap[v] = append(intMap[v], i*10+j)
			}
		}
	}
	return nil
}

func validateDupesIn3x3Matrix(m *[9][9]int, startRow int, startCol int) error {
	intMap := make(map[int][]int)
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			v := m[i][j]
			if v > 0 {
				if found, ok := intMap[v]; ok && len(found) > 0 {
					return errors.New("there are dupes in 3x3 matrix")
				}
				intMap[v] = append(intMap[v], i*10+j)
			}
		}
	}
	return nil
}

func validateAll3x3Matrix(m *[9][9]int) error {
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

func boundingFn(m *[9][9]int) bool {
	if err1 := validateDupesInRowsCols(m); err1 != nil {
		return false
	}
	if err2 := validateAll3x3Matrix(m); err2 != nil {
		return false
	}
	return true
}

func hasEmptyCells(m *[9][9]int) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if m[i][j] == 0 {
				return true
			}
		}
	}
	return false
}

func getNextRowColToVisit(i, j int) (int, int) {
	nextJ := (j + 1) % 9
	nextI := i
	if i >= 0 && i < 8 && nextJ == 0 {
		nextI = i + 1
	}
	return nextI, nextJ
}

func backtrack(i int, j int, m *[9][9]int) bool {
	if !hasEmptyCells(m) {
		return true
	}
	if m[i][j] == 0 {
		for k := 1; k <= 9; k++ {
			m[i][j] = k
			if boundingFn(m) {
				nextI, nextJ := getNextRowColToVisit(i, j)
				val := backtrack(nextI, nextJ, m)
				if val {
					return true
				}
			}
			m[i][j] = 0
		}
	} else {
		nextI, nextJ := getNextRowColToVisit(i, j)
		val := backtrack(nextI, nextJ, m)
		return val
	}
	return false
}

func main() {
	m := parseInput(os.Args[1])
	fmt.Println("Initial matrix for sudoku :")
	printMatrix(&m)
	err := validateDupesInRowsCols(&m)
	if err != nil {
		panic(err)
	}
	err = validateAll3x3Matrix(&m)
	if err != nil {
		panic(err)
	}
	fmt.Println("Matrix for sudoku is valid")
	backtrack(0, 0, &m)
	fmt.Println("Final solution :")
	printMatrix(&m)
}
