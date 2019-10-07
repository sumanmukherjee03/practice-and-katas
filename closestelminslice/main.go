package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

// generate an array of random elements
func genArray(n int) []int {
	arr := make([]int, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000)
	}
	return arr
}

// slice indexOf
func indexOf(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

// binary search for closest element in slice - returns the element in the original array
func binarySearch(val int, arr []int) int {
	if val < arr[0] || len(arr) == 1 {
		return arr[0]
	}
	if val > arr[len(arr)-1] {
		return arr[len(arr)-1]
	}

	var cutoverElm int
	middle := len(arr) / 2
	leftLast := arr[middle-1]
	rightFirst := arr[middle]
	if val == leftLast {
		cutoverElm = leftLast
	} else if val == rightFirst {
		cutoverElm = rightFirst
	} else if val < leftLast {
		cutoverElm = binarySearch(val, arr[0:middle])
	} else if val > rightFirst {
		cutoverElm = binarySearch(val, arr[middle:len(arr)])
	} else {
		if math.Abs(float64(val-leftLast)) < math.Abs(float64(val-rightFirst)) {
			cutoverElm = leftLast
		} else {
			cutoverElm = rightFirst
		}
	}
	return cutoverElm
}

func findClosestElems(x int, k int, arr []int) []int {
	var cutoverPoint int
	res := make([]int, k)
	elm := binarySearch(x, arr)
	cutoverPoint = indexOf(elm, arr)
	if cutoverPoint >= 1 && cutoverPoint < len(arr)-1 {
		res[0] = arr[cutoverPoint-1]
		res[1] = arr[cutoverPoint]
		res[2] = arr[cutoverPoint+1]
	} else if cutoverPoint == 0 {
		copy(res, arr[0:3])
	} else if cutoverPoint == len(arr)-1 {
		copy(res, arr[len(arr)-4:])
	}
	return res
}

func main() {
	n, err1 := strconv.ParseInt(os.Args[1], 0, 0)
	if err1 != nil {
		panic(err1)
	}
	x, err2 := strconv.ParseInt(os.Args[2], 0, 0)
	if err2 != nil {
		panic(err2)
	}
	k, err3 := strconv.ParseInt(os.Args[3], 0, 0)
	if err3 != nil {
		panic(err3)
	}

	arr := genArray(int(n))
	fmt.Println("Original array : ", arr)
	sort.Ints(arr)
	fmt.Println("Sorted array : ", arr)
	fmt.Println("Closest elements to given value in the array : ", findClosestElems(int(x), int(k), arr))
}
