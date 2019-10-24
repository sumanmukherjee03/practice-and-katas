package main

import (
	"fmt"
	"os"
)

var board *[][]int
var direction []int
var xPos int
var yPos int

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

func getCurrentDirection() int {
	for p, v := range direction {
		if v == 1 {
			return p
		}
	}
	return -1
}

func changeDirection(isLeft bool) {
	d := getCurrentDirection()
	direction[d] = 0
	if isLeft {
		direction[((d + 1) % len(direction))] = 1
		return
	}
	if d-1 < 0 {
		d = len(direction) + (d - 1)
	} else {
		d = d - 1
	}
	direction[d] = 1
}

func traverse(cmd string) {
	switch cmd {
	case "l":
		changeDirection(true)
		return
	case "r":
		changeDirection(false)
		return
	case "f":
		d := getCurrentDirection()

		if d%2 == 0 {
			if d == 0 {
				yPos = yPos - 1
			} else {
				yPos = yPos + 1
			}
		}

		if d%2 == 1 {
			if d == 1 {
				xPos = xPos - 1
			} else {
				xPos = xPos + 1
			}
		}

		(*board)[yPos][xPos] = 1
		return
	case "b":
		d := getCurrentDirection()

		if d%2 == 0 {
			if d == 0 {
				yPos = yPos + 1
			} else {
				yPos = yPos - 1
			}
		}

		if d%2 == 1 {
			if d == 1 {
				xPos = xPos + 1
			} else {
				xPos = xPos - 1
			}
		}

		(*board)[yPos][xPos] = 1
		return

	}
}

func main() {
	cmds := os.Args[1]
	if len(cmds) == 0 {
		panic("No instructions")
	}
	n := 20
	direction = make([]int, 4)
	direction[0] = 1
	xPos = (n - 1) / 2
	yPos = (n - 1) / 2
	initMatrix(n)
	(*board)[yPos][xPos] = 1
	for _, s := range cmds {
		traverse(string(s))
	}
	printMatrix()
}
