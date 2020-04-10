package main

func numIslands2(m int, n int, positions [][]int) []int {
	grid := genGrid(m, n)
	var countArr []int

	for i, x := range positions {
		if i == 0 && grid[x[0]][x[1]] == 0 {
			countArr = append(countArr, 1)
		} else if grid[x[0]][x[1]] == 0 {
			r := numberOfIslandsGettingAddedOrReduced(x[0], x[1], grid)
			countArr = append(countArr, countArr[i-1]+r)
		} else {
			countArr = append(countArr, countArr[i-1])
		}
		grid[x[0]][x[1]] = 1
	}

	return countArr
}

func genGrid(m int, n int) [][]int {
	var g [][]int
	for i := 0; i < m; i++ {
		row := make([]int, n)
		for j, _ := range row {
			row[j] = 0
		}
		g = append(g, row)
	}
	return g
}

func numberOfIslandsGettingAddedOrReduced(row int, col int, grid [][]int) int {
	res := 0
	neighbours := findNeighbouringLands(row, col, grid)
	if len(neighbours) == 0 {
		res = 1
	} else if len(neighbours) == 1 {
		res = 0
	} else {
		res = -(findNumDistinctIslandsContainingNeighbours(grid, neighbours) - 1)
	}
	return res
}

func findNeighbouringLands(row int, col int, grid [][]int) [][2]int {
	var found [][2]int
	if (row-1 >= 0) && grid[row-1][col] == 1 {
		found = append(found, [2]int{row - 1, col})
	}
	if (col-1 >= 0) && grid[row][col-1] == 1 {
		found = append(found, [2]int{row, col - 1})
	}
	if (row+1 < len(grid)) && grid[row+1][col] == 1 {
		found = append(found, [2]int{row + 1, col})
	}
	if (col+1 < len(grid[row])) && grid[row][col+1] == 1 {
		found = append(found, [2]int{row, col + 1})
	}
	return found
}

func dupeGrid(grid [][]int) [][]int {
	duplicate := make([][]int, len(grid))
	for i := range grid {
		duplicate[i] = make([]int, len(grid[i]))
		copy(duplicate[i], grid[i])
	}
	return duplicate
}

func uniqueInts(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func findNumDistinctIslandsContainingNeighbours(grid [][]int, found [][2]int) int {
	h := make(map[[2]int]int) // maintain hash of neighbours found and the island they belong to

	// start initially as if they dont belong to any island
	for m := 0; m < len(found); m++ {
		h[found[m]] = m + 1
	}

	var discoverNeighbours func([][]int, int, int, int)
	discoverNeighbours = func(newgrid [][]int, i int, j int, islandId int) {
		// guard clause to check boundaries and if new element visited isnt 0
		if i < 0 || i >= len(newgrid) || j < 0 || j >= len(newgrid[i]) || newgrid[i][j] == 0 {
			return
		}

		// find if the current position in grid matches any of the other neighbour positions
		for n := 0; n < len(found); n++ {
			if h[found[n]] != islandId {
				if found[n][0] == i && found[n][1] == j {
					h[found[n]] = islandId // if the current position is one of the neighbours update island id to match that of the staring position
				}
			}
		}

		newgrid[i][j] = 0                             // convert the visited char to "0"
		discoverNeighbours(newgrid, i, j+1, islandId) // traverse right
		discoverNeighbours(newgrid, i, j-1, islandId) // traverse left
		discoverNeighbours(newgrid, i+1, j, islandId) // traverse up
		discoverNeighbours(newgrid, i-1, j, islandId) // traverse down
	}

	// copy the grid so that you can mark visited nodes
	duped := dupeGrid(grid)

	for k, v := range h {
		discoverNeighbours(duped, k[0], k[1], v)
	}

	vals := make([]int, 0, len(h))

	for _, v := range h {
		vals = append(vals, v)
	}
	uniqueVals := uniqueInts(vals)

	return len(uniqueVals)
}
