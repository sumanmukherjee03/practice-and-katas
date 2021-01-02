package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type Train struct {
	Id     int
	Length int
	Front  int
}

type Intersection struct {
	Id       int
	Mutex    sync.Mutex
	LockedBy int
}

type Crossing struct {
	Position     int
	Intersection *Intersection
}

var (
	screenHeight = 320
	screenWidth  = 320
	colours      = [4]color.RGBA{
		{233, 33, 40, 255},
		{78, 151, 210, 255},
		{251, 170, 26, 255},
		{11, 132, 54, 255},
	}
	white         = color.RGBA{R: 185, G: 185, B: 185, A: 255}
	trains        [4]*Train
	intersections [4]*Intersection
)

const trainLength = 70

/////////////////////// Define Game ///////////////////////////
type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(_, _ int) (int, int) {
	return 320, 320
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawTracks(screen)
	DrawIntersections(screen)
	DrawTrains(screen)
}

///////////////////// Main entrypoint ///////////////////////
func main() {
	describe()

	for i := 0; i < 4; i++ {
		trains[i] = &Train{Id: i, Length: trainLength, Front: 0}
	}

	for i := 0; i < 4; i++ {
		intersections[i] = &Intersection{Id: i, Mutex: sync.Mutex{}, LockedBy: -1}
	}

	go MoveTrain(trains[0], 300, []*Crossing{{Position: 125, Intersection: intersections[0]}, {Position: 175, Intersection: intersections[1]}})
	go MoveTrain(trains[1], 300, []*Crossing{{Position: 125, Intersection: intersections[1]}, {Position: 175, Intersection: intersections[2]}})
	go MoveTrain(trains[2], 300, []*Crossing{{Position: 125, Intersection: intersections[2]}, {Position: 175, Intersection: intersections[3]}})
	go MoveTrain(trains[3], 300, []*Crossing{{Position: 125, Intersection: intersections[3]}, {Position: 175, Intersection: intersections[0]}})

	// Start game
	ebiten.SetWindowSize(320*3, 320*3)
	ebiten.SetWindowTitle("Trains in a box demo")
	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func lockIntersectionsInDistance(id int, reserveStart int, reserveEnd int, crossings []*Crossing) {
	var intersectionsToLock []*Intersection
	for _, crossing := range crossings {
		// If the crossing comes between the end and front of the train and the intersection of the crossing is already locked by the same train
		// then add the intersection from this crossing to the list of crossings to lock
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}

	// Sort the intersections by id to respect the order of locking the intersections
	sort.Slice(intersectionsToLock, func(i, j int) bool {
		return intersectionsToLock[i].Id < intersectionsToLock[j].Id
	})

	for _, intersection := range intersectionsToLock {
		intersection.Mutex.Lock()
		intersection.LockedBy = id
		time.Sleep(10 * time.Millisecond) // Small sleep to slightly increase the chances of avoiding a deadlock
	}
}

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	// Continue the movement of the train until it has travelled the desired distance
	for train.Front < distance {
		train.Front += 1
		// After train moves 1 unit of distance, iterate over all crossings to see which intersections need to be locked/unlocked
		for _, crossing := range crossings {
			intersection := crossing.Intersection
			// When the front of the train reaches a crossing, look ahead and lock both the crossings
			if train.Front == crossing.Position {
				lockIntersectionsInDistance(train.Id, crossing.Position, crossing.Position+train.Length, crossings)
			}
			// As trains end moves past one crossing unlock the mutex of that intersection
			if (train.Front - train.Length) == crossing.Position {
				intersection.LockedBy = -1
				intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}

func describe() {
	str := `
Try and lock the crossings by looking ahead based on the train length, ie acquire 2 locks for 2 crossings.
Also, follow the ordering of crossings, ie try to acquire lock on crossing with lower id first then the
crossing with the higher id. For instance for train 1, crossings 1 first, then 2. For train 2, crossing 2 first, then crossing 3.
For train 3, crossing 3 first, then crossing 4. For train 4, crossing 1 first then crossing 4.
Note the difference here - how we are going by the id of the crossing
as opposed to simply looking ahead and locking the nearest one first and then the fathest crossingi second.

Remember that the ordering solution only works when there are no code branches and you already know which resources you will be locking.

Here's a pictorial representation of the trains and crossings
                        t
                        r
                        a
                        i
                        n

                        2

                        v
               |        |
               |        |
train 1 > -----1 ------ 2------
               |        |
               |        |
               |        |
          -----4--------3------< train 3
               |        |
               |        |
               ^

               t
               r
               a
               i
               n

               4

_______________
	`
	fmt.Println(str)
}

//////////////////// Draw functions /////////////////////

func DrawIntersections(screen *ebiten.Image) {
	drawIntersection(screen, intersections[0], 145, 145)
	drawIntersection(screen, intersections[1], 175, 145)
	drawIntersection(screen, intersections[2], 175, 175)
	drawIntersection(screen, intersections[3], 145, 175)
}

func DrawTracks(screen *ebiten.Image) {
	for i := 0; i < 300; i++ {
		screen.Set(10+i, 135, white)
		screen.Set(185, 10+i, white)
		screen.Set(310-i, 185, white)
		screen.Set(135, 310-i, white)
	}
}

func DrawTrains(screen *ebiten.Image) {
	drawXTrain(screen, 0, 1, 10, 135)
	drawYTrain(screen, 1, 1, 10, 185)
	drawXTrain(screen, 2, -1, 310, 185)
	drawYTrain(screen, 3, -1, 310, 135)
}

func drawIntersection(screen *ebiten.Image, intersection *Intersection, x int, y int) {
	c := white
	if intersection.LockedBy >= 0 {
		c = colours[intersection.LockedBy]
	}
	screen.Set(x-1, y, c)
	screen.Set(x, y-1, c)
	screen.Set(x, y, c)
	screen.Set(x+1, y, c)
	screen.Set(x, y+1, c)
}

func drawXTrain(screen *ebiten.Image, id int, dir int, start int, yPos int) {
	s := start + (dir * (trains[id].Front - trains[id].Length))
	e := start + (dir * trains[id].Front)
	for i := math.Min(float64(s), float64(e)); i <= math.Max(float64(s), float64(e)); i++ {
		screen.Set(int(i)-dir, yPos-1, colours[id])
		screen.Set(int(i), yPos, colours[id])
		screen.Set(int(i)-dir, yPos+1, colours[id])
	}
}

func drawYTrain(screen *ebiten.Image, id int, dir int, start int, xPos int) {
	s := start + (dir * (trains[id].Front - trains[id].Length))
	e := start + (dir * trains[id].Front)
	for i := math.Min(float64(s), float64(e)); i <= math.Max(float64(s), float64(e)); i++ {
		screen.Set(xPos-1, int(i)-dir, colours[id])
		screen.Set(xPos, int(i), colours[id])
		screen.Set(xPos+1, int(i)-dir, colours[id])
	}
}
