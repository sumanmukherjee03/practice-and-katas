package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var philosophers = []string{"Joe", "John", "Dave", "Mike", "Dan"}

const (
	hunger = 3
	eat    = time.Second
	think  = time.Second
)

var dining sync.WaitGroup

func main() {
	describe()

	rand.Seed(time.Now().UnixNano())

	// Go around the table and invoke go routines while synchronizing access to the forks with mutexes
	fmt.Println("Table empty")
	dining.Add(len(philosophers))
	fork0 := &sync.Mutex{}
	forkLeft := fork0
	for i := 1; i < len(philosophers); i++ {
		forkRight := &sync.Mutex{}
		go solve(philosophers[i], forkLeft, forkRight)
		forkLeft = forkRight
	}
	go solve(philosophers[0], forkLeft, fork0)
	dining.Wait()
	fmt.Println("Table empty")
}

func solve(philosopher string, leftHand, rightHand *sync.Mutex) {
	rSleep := func(t time.Duration) {
		time.Sleep(t + time.Duration(rand.Intn(500)))
	}
	fmt.Println(philosopher, "seated")
	for h := hunger; h > 0; h-- {
		fmt.Println(philosopher, "hungry")
		leftHand.Lock()
		rightHand.Lock()
		fmt.Println(philosopher, "eating")
		rSleep(eat)
		leftHand.Unlock()
		rightHand.Unlock()
		fmt.Println(philosopher, "thinking")
		rSleep(think)
	}
	fmt.Println(philosopher, "done")
	dining.Done()
	fmt.Println(philosopher, "left")
}

func describe() {
	str := `
Five silent philosophers sit at a round table with bowls of spaghetti.
Forks are placed between each pair of adjacent philosophers.
Each philosopher must alternately think and eat.
However, a philosopher can only eat spaghetti when they have both left and right forks.
Each fork can be held by only one philosopher.
So a philosopher can use the fork only if it is not being used by another philosopher.
After an individual philosopher finishes eating,
they need to put down both forks so that the forks become available to others.
A philosopher can take the fork on their right or the one on their left as they become available,
but cannot start eating before getting both forks.
Eating is not limited by the remaining amounts of spaghetti or stomach space,
an infinite supply and an infinite demand are assumed.
The problem is how to design a discipline of behavior such that no philosopher will starve,
i.e., each can forever continue to alternate between eating and thinking,
assuming that no philosopher can know when others may want to eat or think.

Example output :
________________
Table empty
Mark seated
Mark Hungry
Mark Eating
..................
..................
Haris Thinking
Haris Done
Haris Left the table
Table empty

_______________
	`
	fmt.Println(str)
}
