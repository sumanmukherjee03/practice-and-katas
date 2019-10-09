package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type SafeCounter struct {
	vals map[string]int
	mux  sync.Mutex
}

func (s *SafeCounter) Increment(k string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.vals[k]++
}

func (s *SafeCounter) Read(k string) int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.vals[k]
}

func main() {
	maxCount, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	s := SafeCounter{vals: make(map[string]int)}
	for i := 0; i < int(maxCount); i++ {
		go s.Increment("words")
		go s.Increment("words")
	}
	time.Sleep(300 * time.Millisecond)
	fmt.Println(s.Read("words"))
}
