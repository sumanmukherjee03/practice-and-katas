package main

import (
	"fmt"
	"os"
	"strconv"
)

func genMatrix(n int) *[][]int {
	m := make([][]int, n)
	for i := 0; i < n; i++ {
		m[i] = append(m[i], make([]int, n)...)
	}
	return &m
}

func genInitialSolution(n int) *[]int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		m[i] = -1
	}
	return &m
}

func printMatrix(m *[][]int) {
	for _, x := range *m {
		for _, y := range x {
			fmt.Printf("%d ", y)
		}
		fmt.Println()
	}
}

func boundingFn(m *[][]int, s *[]int) bool {
	n := len(*m)
	for k := 0; k < n; k++ {
		col := (*s)[k]
		if col < 0 {
			continue
		}
		for i := 0; i < n; i++ {
			if i != k {
				diagLeft := -1
				diagRight := -1
				if i < k {
					diagLeft = col - (k - i)
					diagRight = col + (k - i)
				} else {
					diagLeft = col - (i - k)
					diagRight = col + (i + k)
				}
				if (*m)[i][col] == 1 || (diagLeft > 0 && diagLeft < n && (*m)[i][diagLeft] == 1) || (diagRight > 0 && diagRight < n && (*m)[i][diagRight] == 1) {
					return false
				}
			}
		}
	}
	return true
}

func isSolutionSpaceFull(s *[]int) bool {
	for _, v := range *s {
		if v < 0 {
			return false
		}
	}
	return true
}

func backtrack(q int, m *[][]int, s *[]int) bool {
	n := len(*s)
	if isSolutionSpaceFull(s) {
		return true
	}
	for col := 0; col < n; col++ {
		(*s)[q] = col
		(*m)[q][col] = 1
		if boundingFn(m, s) {
			nextQ := (q + 1) % n
			val := backtrack(nextQ, m, s)
			if val {
				return true
			}
		}
		(*s)[q] = -1
		(*m)[q][col] = 0
	}
	return false
}

func main() {
	n, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(n); i++ {
		s := genInitialSolution(int(n))
		m := genMatrix(int(n))
		backtrack(i, m, s)
		printMatrix(m)
		fmt.Println(">>>>>>>>>>>>>")
	}
}
