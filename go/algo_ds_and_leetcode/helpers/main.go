package helpers

func GenGrid(m int, n int) [][]int {
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

func Dupe2DIntsArray(grid [][]int) [][]int {
	duplicate := make([][]int, len(grid))
	for i := range grid {
		duplicate[i] = make([]int, len(grid[i]))
		copy(duplicate[i], grid[i])
	}
	return duplicate
}

func RemoveElementFromIntsArray(arr []int, index int) []int {
	copy(arr[index:], arr[index+1:]) // Shift arr[index+1:] left one index.
	arr[len(arr)-1] = 0              // Erase last element (write zero value).
	arr = arr[:len(arr)-1]
	return arr
}

func UniqueInts(input []int) []int {
	u := make([]int, 0, len(input)) // length is 0 but allocated capacity is same as input slice
	m := make(map[int]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func ValsFromMap(h map[string]int) []int {
	vals := make([]int, 0, len(h))
	for _, v := range h {
		vals = append(vals, v)
	}
	return vals
}
