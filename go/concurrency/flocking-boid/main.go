package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth    = 640
	screenHeight   = 360
	boidCount      = 500   // total number of boids
	viewRadius     = 14    // radius of circle around boid to adjust for clustering
	adjustmentRate = 0.015 // factor for smooth transition of velocity for boids to get alignment with the neighbours
)

var (
	green   = color.RGBA{10, 255, 50, 255}
	boids   [boidCount]*Boid                       // an array of boid pointers
	boidMap [screenWidth + 1][screenHeight + 1]int // stores where on the grid which boid is
	mutex   = &sync.RWMutex{}                      // Get a reader writer lock to allow multiple reads and exclusive writes
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
	boidMap[int(b.position.x)][int(b.position.y)] = b.id // Mark initial positon in the grid with the boid id
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
	acceleration := b.calculateAcceleration()
	mutex.Lock()                                                 // Acquire a write lock here because we are updating the state variable boidMap which is used to draw the graphics on the screen
	b.velocity = b.velocity.AddVector(acceleration).Limit(-1, 1) // Calculate new velocity
	boidMap[int(b.position.x)][int(b.position.y)] = -1           // Mark old position as empty in grid
	b.position = b.position.AddVector(b.velocity)                // Update to new position based on velocity
	boidMap[int(b.position.x)][int(b.position.y)] = b.id         // Update grid at new position with boid id
	mutex.Unlock()                                               // Release writers lock
}

func (b *Boid) calculateAcceleration() Vector2D {
	xLeft := math.Max(0, b.position.x-viewRadius)
	yBottom := math.Max(0, b.position.y-viewRadius)
	xRight := math.Min(screenWidth, b.position.x+viewRadius)
	yTop := math.Min(screenHeight, b.position.y+viewRadius)
	count := 0.0
	accelerationVector := Vector2D{x: b.borderBounce(b.position.x, screenWidth), y: b.borderBounce(b.position.y, screenHeight)}
	targetPositionVector, targetVelocityVector, separationVector := Vector2D{x: 0, y: 0}, Vector2D{x: 0, y: 0}, Vector2D{x: 0, y: 0}
	mutex.RLock() // Acquire a read lock here so that you can safely access the state variable boidMap to get the boid ids located in the view radius
	for i := xLeft; i <= xRight; i++ {
		for j := yBottom; j <= yTop; j++ {
			bId := boidMap[int(i)][int(j)] // Iterate over view square area (because circles are hard to operate on in grid systems) to find out the neighbouring boid ids
			if bId >= 0 && bId != b.id {   // Only if not the same as the current boid, enter the conditional section
				if d := boids[bId].position.Distance(b.position); d < viewRadius { // Calculate distance and then compare with view radius to handle the circular case
					count += 1                                                                 // Keep count of neighbours because you need to find averages
					targetVelocityVector = targetVelocityVector.AddVector(boids[bId].velocity) // Needed to find average target velocity for flocking
					targetPositionVector = targetPositionVector.AddVector(boids[bId].position) // Needed to find average target position for flocking
					temp := b.position.SubtractVector(boids[bId].position)
					separationVector = separationVector.AddVector(temp.DivideScalar(d)) // Look up separation vector in vector calculus here - https://www.youtube.com/watch?v=KU1VNb5ls6g
				}
			}
		}
	}
	mutex.RUnlock() // Release read lock because we do not need to access the state variable boidMap any more
	if count > 0 {
		targetVelocityVector = targetVelocityVector.DivideScalar(count) // Get average target velocity
		targetPositionVector = targetPositionVector.DivideScalar(count) // Get average position vector
	}
	accelerationAlignmentVector := targetVelocityVector.SubtractVector(b.velocity)           // Get acceleration required to align speed
	accelerationAlignmentVector = accelerationAlignmentVector.MultiplyScalar(adjustmentRate) // Adjust with the adjustment rate to make it easy for viewing
	accelerationCohesionVector := targetPositionVector.SubtractVector(b.position)            // Get acceleration required to align position, ie to bring cohesion in the flock
	accelerationCohesionVector = accelerationCohesionVector.MultiplyScalar(adjustmentRate)   // Adjust with the adjustment rate to make it easy for viewing
	accelerationSeparationVector := separationVector.MultiplyScalar(adjustmentRate)          // Adjust the separation vector with the adjustment rate to make it easy for viewing
	// Finally add all the vectors - alignment vector, cohesion vector and the separation vector to get the final acceleration vector
	accelerationVector = accelerationVector.AddVector(accelerationAlignmentVector)
	accelerationVector = accelerationVector.AddVector(accelerationCohesionVector)
	accelerationVector = accelerationVector.AddVector(accelerationSeparationVector)
	return accelerationVector
}

func (b *Boid) borderBounce(pos, max float64) float64 {
	if pos < viewRadius { // If current position dimension is within the 0 - view radius range, then accelerate slower in that direction
		return 1 / pos
	} else if pos > max-viewRadius { // Similarly if current position dimension is within the view radius - screen limit range, then accelerate slower in that direction
		return 1 / (pos - max)
	}
	return 0
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

func (v *Vector2D) Distance(v1 Vector2D) float64 {
	return math.Sqrt(math.Pow((v.x-v1.x), 2) + math.Pow((v.y-v1.y), 2))
}

func (v Vector2D) Limit(lower, upper float64) Vector2D {
	return Vector2D{math.Min(math.Max(v.x, lower), upper),
		math.Min(math.Max(v.y, lower), upper)}
}
