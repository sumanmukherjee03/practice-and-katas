package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type ticket struct {
	src  string
	dest string
}

var tickets []ticket

func readlines(ch chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			ch <- line
		} else {
			break
		}
	}
	if scanner.Err() != nil {
		log.Fatalln("Scanner encountered an error")
	}
	close(ch)
}

func store(ch <-chan string) {
	sep := regexp.MustCompile(`->`)
	for l := range ch {
		arr := sep.Split(l, 2)
		tickets = append(tickets, ticket{src: strings.Trim(arr[0], " "), dest: strings.Trim(arr[1], " ")})
	}
}

func findStart() string {
	var srcs []string
	var dests []string
	for _, t := range tickets {
		srcs = append(srcs, t.src)
		dests = append(dests, t.dest)
	}

	var helper func(int) string
	helper = func(i int) string {
		var start string
		if i == len(srcs) {
			return start
		}

		for j := 0; j < len(dests); j++ {
			if srcs[i] == dests[j] {
				return helper(i + 1)
			}
		}

		return srcs[i]
	}

	return helper(0)
}

func print(start string) {
	if len(start) == 0 {
		return
	}

	fmt.Println(start)

	var findNext = func() string {
		var next string
		for _, t := range tickets {
			if t.src == start {
				next = t.dest
				break
			}
		}
		return next
	}

	print(findNext())
}

func main() {
	ch := make(chan string)
	go readlines(ch)
	store(ch)
	start := findStart()
	if len(start) == 0 {
		log.Fatalln("There is a cycle, so cant figure out the start")
	}
	print(start)
}
