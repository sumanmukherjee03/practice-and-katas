package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"
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

	if len(os.Args) != 2 {
		errorf("The provided args is not valid. You need to provide at least one output type - slurp|stream")
	}

	outType := os.Args[1]
	if _, ok := allowedOutTypes[outType]; !ok {
		errorf("The provided output type is not valid - %s", outType)
	}

	p1 := &person{
		FirstName: "John",
		LastName:  "Doe",
	}

	p2 := &person{
		FirstName: "Jane",
		LastName:  "Doe",
	}

	persons := []*person{
		p1,
		p2,
	}

	switch outType {
	case slurpOutputType:
		b, err := json.Marshal(persons)
		if err != nil {
			errorf("The persons slice is not marshalable : %v", err)
		}
		fmt.Fprintf(os.Stdout, string(b)+"\n")
	case streamOutputType:
		ticker := time.NewTicker(3 * time.Second)
		for _, p := range persons {
			b, err := json.Marshal(p)
			if err != nil {
				errorf("The person is not marshalable : %v", err)
			}
			<-ticker.C
			fmt.Fprintf(os.Stdout, string(b)+"\n")
		}
		ticker.Stop()
	}
	fmt.Println("-----------------")

	sampleStr := `[{"first_name":"Jack","last_name":"Doyle","initials": "Mr."},{"first_name":"Jenny","last_name":"Doyle","initials":"Mrs."}]`
	people := []person{}
	if err := json.Unmarshal([]byte(sampleStr), &people); err != nil {
		errorf("The given string could not be unmarshaled - %v", err)
	}
	for _, np := range people {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("FirstName: %s, LastName: %s\n", np.FirstName, np.LastName))
	}
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}
