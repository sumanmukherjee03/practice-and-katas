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
	"time"
)

var (
	windRegex      = regexp.MustCompile(`\d* METAR.*EGLL \d*Z [A-Z ]*(\d{5}KT|VRB\d{2}KT).*=`)
	tafValidation  = regexp.MustCompile(`.*TAF.*`)
	comment        = regexp.MustCompile(`\w*#.*`)
	metarClose     = regexp.MustCompile(`.*=`)
	variableWind   = regexp.MustCompile(`.*VRB\d{2}KT`)
	validWind      = regexp.MustCompile(`\d{5}KT`)
	windDirOnly    = regexp.MustCompile(`(\d{3})\d{2}KT`)
	windDirections = make([]string, 0, 8)
	windDist       = make(map[string]int)
)

func main() {
	textChan := make(chan string)
	metarChan := make(chan []string)
	windChan := make(chan []string)
	degreesChan := make(chan float64)

	// Start the goroutines that handle the pipeline
	// by passing output channel from one as input to the next goroutine
	go parseToArray(textChan, metarChan)
	go extractWindDirection(metarChan, windChan)
	go mineWindDistribution(windChan, degreesChan)

	path, err := filepath.Abs("./pipeline/metarfiles")
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	initDist()

	startTime := time.Now()
	for _, f := range files {
		data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		text := string(data)
		textChan <- text
	}
	close(textChan) // You can close a non-empty channel and still have the remaining values received
	elapsedTime := time.Since(startTime)

	fmt.Printf("%v\n", windDist)
	fmt.Println("Time taken : ", elapsedTime)
}

func initDist() {
	windDist["N"] = 0
	windDist["NE"] = 0
	windDist["E"] = 0
	windDist["SE"] = 0
	windDist["S"] = 0
	windDist["SW"] = 0
	windDist["W"] = 0
	windDist["NW"] = 0
	windDirections = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
}

func parseToArray(inChan chan string, outChan chan []string) {
	// Iterate over each element in the channel as it is received
	for text := range inChan {
		lines := strings.Split(text, "\n")
		metarSlice := make([]string, 0, len(lines))
		metarStr := ""
		for _, line := range lines {
			if tafValidation.MatchString(line) {
				break
			}
			if !comment.MatchString(line) {
				metarStr += strings.Trim(line, " ")
			}
			if metarClose.MatchString(line) {
				metarSlice = append(metarSlice, metarStr)
				metarStr = ""
			}
		}
		outChan <- metarSlice
	}
	close(outChan)
}

func extractWindDirection(inChan chan []string, outChan chan []string) {
	for metars := range inChan {
		winds := make([]string, 0, len(metars))
		for _, metar := range metars {
			if windRegex.MatchString(metar) {
				winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
			}
		}
		outChan <- winds
	}
	close(outChan)
}

func mineWindDistribution(inChan chan []string, outChan chan float64) {
	for winds := range inChan {
		for _, wind := range winds {
			if variableWind.MatchString(wind) {
				for k, _ := range windDist {
					windDist[k] = windDist[k] + 1
				}
			} else if validWind.MatchString(wind) {
				windStr := windDirOnly.FindAllStringSubmatch(wind, -1)[0][1]
				if degrees, err := strconv.ParseFloat(windStr, 64); err == nil {
					dirIndex := int(math.Round(degrees/45.0)) % 8
					windDist[windDirections[dirIndex]] = windDist[windDirections[dirIndex]] + 1
				}
			}
		}
	}
	close(outChan)
}
