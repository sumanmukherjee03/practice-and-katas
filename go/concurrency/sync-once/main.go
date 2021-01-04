package main

import (
	"fmt"
	"sync"
)

var (
	wg   = sync.WaitGroup{}
	once sync.Once
)

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(iter int) {
			defer wg.Done()
			once.Do(func() {
				initFn(fmt.Sprintf("Initializing with val : %d", iter))
			})
		}(i)
	}

	wg.Wait()
}

func initFn(str string) {
	fmt.Println(str)
}
