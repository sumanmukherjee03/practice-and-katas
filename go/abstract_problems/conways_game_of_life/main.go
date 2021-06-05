/*
Given a grid or a matrix and each cell of the matrix can have values either 0 or 1,
and a tick or a timer where the snapshot of the grid changes from generation 0 to generation 1 and so on.
The next snapshot of the grid depends on the new values of the cells.
But the cells' new value depends on these following rules :
  - A cell changes state from 0 -> 1 if the state of 3 of it's neighbours is 1
  - A cell changes state from 1 -> 0 if the state of < 2 neighbours is 1 or > 3 neighbours is 1
    - Meaning cell dies either due to underpopulation (tough conditions) or overpopulation
  - In other cases the state of a cell stays the same
*/
package main

import "fmt"

func main() {
	fmt.Println("Starting conways game of life")
}
