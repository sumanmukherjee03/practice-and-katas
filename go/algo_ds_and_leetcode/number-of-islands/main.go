package main

import "fmt"

func numIslands(grid [][]byte) int {
	count := 0
	for i, x := range grid {
		for j, y := range x {
			// convert byte to string and check for "1"
			// if at least one "1" is found, an island is found, and explore it's neighbours
			if string(y) == "1" {
				count += 1
				bfs(&grid, i, j)
			}
		}
	}
	return count
}

func bfs(grid *[][]byte, i int, j int) {
	// guard clause
	if i < 0 || i >= len((*grid)) || j < 0 || j >= len((*grid)[i]) || string((*grid)[i][j]) == "0" {
		return
	}

	// convert string to byte array and then get the first char as byte
	(*grid)[i][j] = []byte("0")[0] // convert the visited char to "0"

	bfs(grid, i, j+1) // traverse right
	bfs(grid, i, j-1) // traverse left
	bfs(grid, i+1, j) // traverse up
	bfs(grid, i-1, j) // traverse down
}

func main() {
	grid := [][]byte{
		[]byte("11110"),
		[]byte("11010"),
		[]byte("11000"),
		[]byte("00000"),
	}

	res := numIslands(grid)
	fmt.Println(res)
}
