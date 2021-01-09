/*
Given 2 jugs with max cap of m and n, with no markings on the sides,
the goal is to reach d litres using these 2 jugs.
So, the state is to reach from (m,n) -> (0, d) or (d, 0)
*/

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type step struct {
	jug1 int
	jug2 int
}

var jug1Cap int
var jug2Cap int
var goal int
var solution = []step{}

func isStepVisited(s step) bool {
	for _, v := range solution {
		if v.jug1 == s.jug1 && v.jug2 == s.jug2 {
			return true
		}
	}
	return false
}

func waterJug(amt1 int, amt2 int) bool {
	s := step{amt1, amt2}
	if (amt1 == goal && amt2 == 0) || (amt1 == 0 && amt2 == goal) {
		solution = append(solution, s)
		return true
	}

	if isStepVisited(s) {
		return false
	}

	solution = append(solution, s)

	return waterJug(amt1, 0) ||
		waterJug(0, amt2) ||
		waterJug(jug1Cap, amt2) ||
		waterJug(amt1, jug2Cap) ||
		waterJug(amt1+int(math.Min(float64(amt2), float64(jug1Cap-amt1))), amt2-int(math.Min(float64(amt2), float64(jug1Cap-amt1)))) ||
		waterJug(amt1-int(math.Min(float64(amt1), float64(jug2Cap-amt2))), amt2+int(math.Min(float64(amt1), float64(jug2Cap-amt2))))
}

func main() {
	m, err1 := strconv.ParseInt(os.Args[1], 0, 0)
	if err1 != nil {
		panic(err1)
	}
	jug1Cap = int(m)

	n, err2 := strconv.ParseInt(os.Args[2], 0, 0)
	if err2 != nil {
		panic(err2)
	}
	jug2Cap = int(n)

	d, err3 := strconv.ParseInt(os.Args[3], 0, 0)
	if err3 != nil {
		panic(err3)
	}
	goal = int(d)

	waterJug(0, 0)

	fmt.Println(solution)
}
