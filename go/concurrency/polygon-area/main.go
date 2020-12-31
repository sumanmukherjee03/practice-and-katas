package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point2D struct {
	x int
	y int
}

const numOfThreads int = 8

var (
	rex = regexp.MustCompile(`\((\d*),(\d*)\)`)
	wg  = sync.WaitGroup{}
)

func main() {
	describe()

	inputChan := make(chan string, 1000)

	path, err := filepath.Abs("./polygon-area")
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.ReadFile(filepath.Join(path, "polygons.txt"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numOfThreads; i++ {
		wg.Add(1)
		go findArea(inputChan)
	}

	data := string(file)
	lines := strings.Split(data, "\n")

	start := time.Now()
	for _, line := range lines {
		inputChan <- line
	}
	close(inputChan)
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Processing time : ", elapsed)
}

func findArea(ch chan string) {
	for pointsStr := range ch {
		var points []Point2D
		for _, p := range rex.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x: x, y: y})
		}

		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
	wg.Done()
}

func describe() {
	str := `
This example demonstrates a worker pool implementation.
Imagine building developers sending in cartesian co-ordinates of vertices of a polygon.
Our job is to find the area and return the value as a result.
To implement that we are gonna use a buffered channel and a worker pool.

_____________________
	`
	fmt.Println(str)
}
