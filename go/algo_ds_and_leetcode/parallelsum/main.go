package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func genRandomSlice(n int) []int {
	rand.Seed(time.Now().UnixNano())
	arr := []int{}
	for i := 0; i < n; i++ {
		arr = append(arr, rand.Intn(1000))
	}
	return arr
}

func sum(s []int, c chan int) {
	sum := 0
	for _, val := range s {
		sum += val
	}
	c <- sum
}

func main() {
	n, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	arr := genRandomSlice(int(n))
	c := make(chan int)
	fmt.Println("Original array : ", arr)
	go sum(arr[:len(arr)/2], c)
	go sum(arr[len(arr)/2:], c)
	x, y := <-c, <-c
	fmt.Println("Total sum : ", x+y)
}
