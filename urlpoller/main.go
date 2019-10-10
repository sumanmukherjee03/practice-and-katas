package main

import (
	"fmt"
	"net/http"
	"time"
)

type State struct {
	url    string
	status int
	err    error
}

// Resource - it is a type that represents the url plus the count of successes and failures
type Resource struct {
	url            string
	errorCounter   int
	successCounter int
}

func (r *Resource) Poll() (string, error) {
	resp, err := http.Head(r.url)
	if err != nil {
		r.errorCounter++
		r.successCounter = 0
		return "", err
	}
	r.successCounter++
	r.errorCounter = 0
	fmt.Println(">>>>", r)
	return resp.Status, nil
}

func (r *Resource) Sleep(pending chan<- *Resource, quit <-chan bool) {
	fmt.Println("Sleeping ...")
	time.Sleep(5000 * time.Millisecond)
	select {
	case pending <- r:
		// Do nothing
	case <-quit:
		fmt.Println("No more sleeping. Killing go routine")
	}
}

func Poller(pending chan *Resource, complete chan<- *Resource, errors chan<- error, quit <-chan bool) {
	for {
		select {
		case r := <-pending:
			_, err := r.Poll()
			if err != nil {
				errors <- err
			}
			complete <- r
		case <-quit:
			fmt.Println("Shutting down poller go routine")
			return
		}
	}
}

func main() {
	pending, complete, errors := make(chan *Resource), make(chan *Resource), make(chan error)
	quit := make(chan bool)
	defer close(quit)

	resources := []*Resource{
		&Resource{url: "http://golang.org/"},
		&Resource{url: "http://www.google.com/"},
	}

	// Launch poller go routines
	for i := 0; i < len(resources); i++ {
		go Poller(pending, complete, errors, quit)
	}

	go func() {
		for _, r := range resources {
			pending <- r
		}
	}()

	for {
		select {
		case r := <-complete:
			go r.Sleep(pending, quit)
		case err := <-errors:
			fmt.Println(err)
		case _, ok := <-quit:
			if !ok {
				fmt.Println("Shutting down main go routine")
			}
		}
	}
}
