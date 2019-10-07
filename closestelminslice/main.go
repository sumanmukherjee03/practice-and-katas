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

func genArray(n int) []int {
	arr := make([]int, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000)
	}
	return arr
}

func indexOf(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func findClosestElems(n int, x int, k int) []int {
	arr := genArray(n)
	res := make([]int, k)
	var cutoverPoint int
	fmt.Println("Original array : ", arr)
	sort.Ints(arr)
	fmt.Println("Sorted array : ", arr)

	var binarySearch func([]int) int
	binarySearch = func(a []int) int {
		var cutoverElm int
		fmt.Println(a)
		if x < a[0] {
			cutoverElm = a[0]
			return cutoverElm
		} else if x > a[len(a)-1] {
			cutoverElm = a[len(a)-1]
			return cutoverElm
		}
		if len(a) > 1 {
			m := len(a) / 2
			e1 := a[m-1]
			e2 := a[m]
			if x < e1 {
				cutoverElm = binarySearch(a[0:m])
			} else if x > e2 {
				cutoverElm = binarySearch(a[m:len(a)])
			} else {
				if math.Abs(float64(x-e1)) < math.Abs(float64(x-e2)) {
					cutoverElm = e1
				} else {
					cutoverElm = e2
				}
			}
		} else {
			cutoverElm = a[0]
		}
		return cutoverElm
	}

	elm := binarySearch(arr)
	fmt.Println(elm)
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
	fmt.Println("Closest elements to given value in the array : ", findClosestElems(int(n), int(x), int(k)))
}
