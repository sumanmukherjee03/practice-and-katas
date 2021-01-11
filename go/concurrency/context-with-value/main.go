package main

import (
	"context"
	"fmt"
)

type database map[string]bool
type userIdKeyType string // We need a user defined data type for the context key

var db = database{
	"foo": true,
}

func main() {
	describe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	processRequest(ctx, "foo") // It is a go idiom to pass ctx as the first argument
}

func processRequest(ctx context.Context, userid string) {
	// When creating a new context with value from the parent context, we should keep in mind
	// that contexts are passed by value, ie this new context is a copy of the parent context.
	// But when a parent context is cancelled, all contexts derived from it are also cancelled.
	// That is why when creating a context with value we do not need the cancel function separately to be returned.
	// Handle the cancel function in the parent context.
	newCtx := context.WithValue(ctx, userIdKeyType("userid"), userid)
	ch := checkMembership(newCtx)
	status := <-ch
	fmt.Printf("Membership of userid %s is : %v\n", userid, status)
}

func checkMembership(ctx context.Context) <-chan bool {
	ch := make(chan bool)
	go func() {
		defer close(ch)
		k := ctx.Value(userIdKeyType("userid")).(string)
		status := db[k]
		ch <- status
	}()
	return ch
}

func describe() {
	str := `
This is an example of a cancellable context with value.
Context with values are used as data bags to carry data from the context of a request to the goroutines downstream.


_____________________
	`
	fmt.Println(str)
}
