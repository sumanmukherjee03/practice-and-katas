package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
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

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			intersection := crossing.Intersection
			if train.Front == crossing.Position {
				intersection.Mutex.Lock()
				intersection.LockedBy = train.Id
			}
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
For instance for train 1, crossing 1 is locked first and then the train tries to move towards crossing 2
And before it reaches crossing 2, it tries to acquire the lock on crossing 2.
For train 2, crossing 2 is locked first and then the train tries to move towards crossing 3.
And before it reaches crossing 3, it tries to acquire the lock on crossing 3.
For train 3, crossing 3 is locked first and then the train tries to move towards crossing 4.
And before it reaches crossing 4, it tries to acquire the lock on crossing 4.
For train 4, crossing 4 is locked first and then the train tries to move towards crossing 1.
And before it reaches crossing 1, it tries to acquire the lock on crossing 1.
But this leads to a deadlock.
The aim of this program is to depict the deadlock situation.

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
