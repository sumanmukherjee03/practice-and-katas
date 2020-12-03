package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches []string
	wg      = sync.WaitGroup{}
	mutex   = &sync.Mutex{}
)

func main() {
	path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/sumanmukherjee03/practice-and-katas")
	filename := "main.go"
	wg.Add(1)
	go filesearch(path, filename)
	wg.Wait()
	for _, f := range matches {
		fmt.Println(f)
	}
}

func filesearch(path, filename string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if strings.Contains(f.Name(), filename) {
			mutex.Lock()
			matches = append(matches, filepath.Join(path, f.Name()))
			mutex.Unlock()
		}
		if f.IsDir() {
			wg.Add(1)
			go filesearch(filepath.Join(path, f.Name()), filename)
		}
	}
	wg.Done()
}
