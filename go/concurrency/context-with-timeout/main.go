package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// const timeout = 3000
const timeout = 300

func main() {
	describe()
	req, err := http.NewRequest("GET", "https://nike.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Use the request context from above to generate the context with timeout
	ctx, cancel := context.WithTimeout(req.Context(), timeout*time.Millisecond)
	defer cancel() // Defer cancel so that resources associated with the context get released when the function returns

	req = req.WithContext(ctx) // Bind this new context with timeout to the request

	resp, err := http.DefaultClient.Do(req) // This sends the http request and returns the response
	if err != nil {
		log.Println("ERROR : ", err)
		return
	}

	defer resp.Body.Close() // close the request body on return

	io.Copy(os.Stdout, resp.Body)
}

func describe() {
	str := `
This is an example of a cancellable context with timeout. The cancel function is automatically called when the timeout expires.
When cancel function is called, it sends a done signal on the done channel of the context.
The child goroutine must handle the done signal and bail out to prevent goroutine leaks.

_____________________
	`
	fmt.Println(str)
}
