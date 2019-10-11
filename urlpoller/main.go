package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	fmt.Println(">>>>", *r)
	return resp.Status, nil
}

func (r *Resource) Sleep(pending chan<- *Resource) {
	fmt.Println("Sleeping ...")
	time.Sleep(5000 * time.Millisecond)
	pending <- r
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
		case q := <-quit:
			if q {
				fmt.Println("Shutting down poller. Closing pending and complete channels")
				close(complete)
				close(pending)
				close(errors)
				return
			}
		}
	}
}

func setupCloseHandler(q chan bool) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			q <- true
			close(q)
			fmt.Println("\r\n Killing program and cleaning up")
			os.Exit(0)
		}
	}()
}

func main() {
	pending, complete, errors := make(chan *Resource), make(chan *Resource), make(chan error)
	quit := make(chan bool)
	setupCloseHandler(quit)

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
			go r.Sleep(pending)
		case err := <-errors:
			if err != nil {
				fmt.Println("ERROR :>> ", err)
			}
		}
	}
}
