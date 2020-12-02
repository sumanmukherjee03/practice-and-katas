package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 360
	boidCount    = 500 // total number of boids
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidCount]*Boid                       // an array of boid pointers
	boidMap [screenWidth + 1][screenHeight + 1]int // stores where on the grid which boid is
)

/////////////////////// Define Game ///////////////////////////
type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		screen.Set(int(boid.position.x+1), int(boid.position.y), green)
		screen.Set(int(boid.position.x-1), int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y-1), green)
		screen.Set(int(boid.position.x), int(boid.position.y+1), green)
	}
}

///////////////////// Main entrypoint ///////////////////////
func main() {
	rand.Seed(time.Now().UnixNano())
	initBoidMap()

	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}

	// Start game
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids in a box demo")
	g := &Game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

//////////////////////// Utility funcs ///////////////

func initBoidMap() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
}

func createBoid(bId int) {
	// Keep velocity to 1 pixel or below for visibility sake
	b := Boid{
		position: Vector2D{
			x: rand.Float64() * screenWidth,
			y: rand.Float64() * screenHeight,
		},
		velocity: Vector2D{
			x: (rand.Float64()*2 - 1.0),
			y: (rand.Float64()*2 - 1.0),
		},
		id: bId,
	}
	boids[bId] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}

/////////////////////// Boid /////////////////////

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) moveOne() {
	b.position = b.position.AddVector(b.velocity)
	nextPosition := b.position.AddVector(b.velocity)
	if nextPosition.x >= screenWidth || nextPosition.x < 0 {
		b.velocity.x = -b.velocity.x
	}
	if nextPosition.y >= screenHeight || nextPosition.y < 0 {
		b.velocity.y = -b.velocity.y
	}
}

/////////////////////// Vector2D /////////////////////

type Vector2D struct {
	x float64
	y float64
}

func (v *Vector2D) AddVector(v1 Vector2D) Vector2D {
	return Vector2D{x: v.x + v1.x, y: v.y + v1.y}
}

func (v *Vector2D) SubtractVector(v1 Vector2D) Vector2D {
	return Vector2D{x: v.x - v1.x, y: v.y - v1.y}
}

func (v *Vector2D) MultiplyVector(v1 Vector2D) Vector2D {
	return Vector2D{x: v.x * v1.x, y: v.y * v1.y}
}

func (v *Vector2D) AddScalar(d float64) Vector2D {
	return Vector2D{x: v.x + d, y: v.y + d}
}

func (v *Vector2D) SubtractScalar(d float64) Vector2D {
	return Vector2D{x: v.x - d, y: v.y - d}
}

func (v *Vector2D) MultiplyScalar(d float64) Vector2D {
	return Vector2D{x: v.x * d, y: v.y * d}
}

func (v *Vector2D) DivideScalar(d float64) Vector2D {
	return Vector2D{x: v.x / d, y: v.y / d}
}

func (v *Vector2D) DistanceVector(v1 Vector2D) float64 {
	return math.Sqrt(math.Pow((v.x-v1.x), 2) + math.Pow((v.y-v1.y), 2))
}

func (v Vector2D) Limit(lower, upper float64) Vector2D {
	return Vector2D{math.Min(math.Max(v.x, lower), upper),
		math.Min(math.Max(v.y, lower), upper)}
}
