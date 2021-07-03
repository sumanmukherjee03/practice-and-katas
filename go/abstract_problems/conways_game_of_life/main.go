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
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Game struct {
	GridSize int
	Grid     *[][]int
}

var (
	game  *Game
	clear map[string]func()
)

func init() {
	clear = make(map[string]func())
	commonClearFn := func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = commonClearFn
	clear["darwin"] = commonClearFn
}

func main() {
	rand.Seed(time.Now().UnixNano())
	x, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Could not parse input for grid size - %v", err)
	}
	gridSize := x + 2
	grid := genMatrix(gridSize, func(p, q int) int {
		return rand.Intn(100000) % 2
	})
	game := &Game{
		GridSize: gridSize,
		Grid:     grid,
	}

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				Clear()
				fmt.Print(game)
				game.Step()
			case <-done:
				return
			}
		}
	}()

	time.Sleep(60 * time.Second)
	ticker.Stop()
	done <- true
}

func Clear() {
	fn, ok := clear[runtime.GOOS]
	if !ok {
		log.Fatal("Clear function not supported for this OS")
	}
	fn()
}

func genMatrix(size int, fn func(int, int) int) *[][]int {
	var m [][]int
	for i := 0; i < size; i++ {
		r := make([]int, size)
		for j, _ := range r {
			if i == 0 || i == size-1 || j == 0 || j == size-1 {
				r[j] = 0
			} else {
				r[j] = fn(i, j)
			}
		}
		m = append(m, r)
	}
	return &m
}

func (g *Game) String() string {
	str := ""
	for i := 0; i < g.GridSize; i++ {
		s := ""
		for j := 0; j < g.GridSize; j++ {
			s += " " + strconv.Itoa((*g.Grid)[i][j])
		}
		str += fmt.Sprintf("%s\n", s)
	}
	return str
}

func (g *Game) Step() {
	nextValFn := func(i, j int) int {
		sum := 0
		for r := i - 1; r <= i+1; r++ {
			for c := j - 1; c <= j+1; c++ {
				if r != i && c != j {
					sum += (*g.Grid)[r][c]
				}
			}
		}
		switch {
		case sum == 3: // If 3 neighbours are alive the node becomes 1
			return 1
		case sum == 2: // If 2 neighbours are alive the node stays as is
			return (*g.Grid)[i][j]
		default: // In all other cases the node transitions to 0
			return 0
		}
	}
	g.Grid = genMatrix(g.GridSize, nextValFn)
}
