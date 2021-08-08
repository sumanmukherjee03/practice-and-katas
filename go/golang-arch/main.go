package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

const (
	slurpOutputType  = "slurp"
	streamOutputType = "stream"
)

var (
	allowedOutTypes = map[string]bool{
		slurpOutputType:  true,
		streamOutputType: true,
	}
)

type person struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	signals := make(chan os.Signal)

	// The Notify function will pass the incoming signals that you provided, in this case os.Interrupt
	// to the signals channel, which you can then read from to customize how you handle OS signals
	// This comes from CTRL+c or kill -2 <pid>
	signal.Notify(signals, os.Interrupt)

	// Process the OS interrupt signal in a goroutine
	go func() {
		s := <-signals
		errorf("Received OS signal - %v", s)
	}()

	http.HandleFunc("/encode", handleEncode)
	http.HandleFunc("/decode", handleDecode)
	http.ListenAndServe(":8080", nil)
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Error("ERROR - This path only handles a GET request")
		return
	}
	p1 := &person{
		FirstName: "John",
		LastName:  "Doe",
	}
	p2 := &person{
		FirstName: "Jane",
		LastName:  "Doe",
	}
	ps := []*person{p1, p2}
	err := json.NewEncoder(w).Encode(ps)
	if err != nil {
		log.Error(fmt.Sprintf("ERROR - Could not encode into json - %v", err))
	}
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	var ps []person
	if r.Method != "POST" {
		log.Error("ERROR - This path only handles a POST request")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		log.Error(fmt.Sprintf("ERROR - Could not decode json - %v", err))
	}
	log.Info(ps)
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}
