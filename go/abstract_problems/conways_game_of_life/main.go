/*
Given a grid or a matrix where each cell of the matrix can have values either 0 or 1,
and a tick or a timer where the snapshot of the grid changes from generation 0 to generations 1, 2, 3 ... and so on.
The next snapshot of the grid depends on the new values of the cells.
But the cells' new value depends on these following rules :
  - A cell changes state from 0 -> 1 if the state of 3 of it's neighbours is 1
  - A cell changes state from 1 -> 0 if the state of < 2 neighbours is 1 or > 3 neighbours is 1
    - Meaning cell dies either due to underpopulation (tough conditions) or overpopulation
  - In other cases the state of a cell stays the same
  - A cell can have 8 neighbours unless it is at the edge of the grid
*/
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	gridSize = 41
)

var (
	grid *[][]int
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Starting conways game of life")
	initGrid()
	printMatrix(grid)
}

func initGrid() {
	var m [][]int
	for i := 0; i < gridSize; i++ {
		r := make([]int, gridSize)
		for j, _ := range r {
			r[j] = rand.Intn(100000) % 2
		}
		m = append(m, r)
	}
	grid = &m
}

func printMatrix(m *[][]int) {
	for i := 0; i < len((*m)); i++ {
		str := ""
		for j := 0; j < len((*m)[i]); j++ {
			str += " " + strconv.Itoa((*grid)[i][j])
		}
		fmt.Println(str)
	}
}
