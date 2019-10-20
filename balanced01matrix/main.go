package main

import (
	"fmt"
	"os"
	"strconv"
)

var board *[][]int
var permutations [][]int

func initMatrix(n int) {
	t := make([][]int, n)
	for i := range t {
		t[i] = make([]int, n)
	}
	board = &t
}

func printMatrix() {
	for _, x := range *board {
		for _, y := range x {
			fmt.Printf("%d ", y)
		}
		fmt.Println()
	}
}

func boundingFn() bool {
	n := len((*board))
	rowSum := make([]int, n)
	colSum := make([]int, n)
	for i := range *board {
		for j := range (*board)[i] {
			rowSum[i] = rowSum[i] + (*board)[i][j]
			colSum[j] = colSum[j] + (*board)[i][j]
		}
	}

	for _, x := range rowSum {
		if x != n/2 {
			return false
		}
	}

	for _, y := range colSum {
		if y != n/2 {
			return false
		}
	}

	return true
}

// Permutation implementation with Heaps algorithm
func permutation(a []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, k int) {
		if k == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < k; i++ {
				helper(arr, k-1)
				if k%2 == 1 {
					arr[i], arr[k-1] = arr[k-1], a[i]
				} else {
					arr[0], arr[k-1] = arr[k-1], arr[0]
				}
			}
		}
	}

	helper(a, len(a))
	return res
}

func genInitPermutationArr() []int {
	n := len(*board)
	initArr := make([]int, n)
	for x := 0; x < n/2; x++ {
		initArr[x] = 1
	}
	return initArr
}

func backtrack(solutionSpace []int, k int) bool {
	if k == len(solutionSpace) {
		return true
	}

	temp := solutionSpace[k]
	found := false

	for j := solutionSpace[k] + 1; j < len(permutations); j++ {
		solutionSpace[k] = j
		copy((*board)[k], permutations[solutionSpace[k]])
		if boundingFn() {
			found = true
			break
		}
	}

	if !found {
		solutionSpace[k] = temp
		copy((*board)[k], permutations[solutionSpace[k]])
		backtrack(solutionSpace, k+1)
	}

	return found
}

func main() {
	n, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	if n%2 == 1 {
		panic("The argument can only be even")
	}
	initMatrix(int(n))
	permutations = permutation(genInitPermutationArr())

	for i := 0; i < len((*board)); i++ {
		copy((*board)[i], permutations[0])
	}

	backtrack([]int{1, 2, 3, 4}, 0)
	printMatrix()
}
