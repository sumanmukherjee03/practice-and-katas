package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func genArray(n int) []int {
	arr := make([]int, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000)
	}
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func maxOfSlidingWindow(n, k int) []int {
	maxs := []int{}
	arr := genArray(n)
	fmt.Println("Original array : ", arr)
	for i := 0; i <= len(arr)-k; i++ {
		temp := make([]int, k)
		copy(temp, arr[i:i+k])
		sort.Ints(temp)
		maxs = append(maxs, temp[len(temp)-1])
	}
	return maxs
}

func main() {
	n, err1 := strconv.ParseInt(os.Args[1], 0, 0)
	if err1 != nil {
		panic(err1)
	}
	k, err2 := strconv.ParseInt(os.Args[2], 0, 0)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("Maxs' of sliding windows : ", maxOfSlidingWindow(int(n), int(k)))
}
