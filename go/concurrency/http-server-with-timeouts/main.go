package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error // creating a var because this is being re-assigned multiple times

	describe()

	connStr := "host=localhost port=5432 user=bar password=pqrs6789 dbname=foo sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// We use a context with timeout here to ping the database during the initial start
	// This context is not used anywhere else
	// Also we defer cancel the context so that we release the resources upon termination of main
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	// WriteTimeout - This essentially covers the lifetime of the ServeHTTP method called by the server to respond to the request
	// http.TimeoutHandler - This returns a handler that runs the input handler with the given time limit
	// The timeout on http.TimeoutHandler should be less than the WriteTimeout to allow the timeout handler to write a message to the client caused by timeout expiry
	// If handler times out the client receives a response of 503 - service unavailable
	srv := http.Server{
		Addr:         "localhost:23000",
		WriteTimeout: 2 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(slowHandler), 1*time.Second, "Timeout!!!"),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// http.HandlerFunc is a func which accepts a function with a signature func(http.ResponseWriter, *http.Request)
// and turns that into a http handler. The function f passed to HandlerFunc gets called during ServeHTTP method execution.
// Here slowHandler is that function f passed to HandlerFunc
func slowHandler(w http.ResponseWriter, req *http.Request) {
	// The *http.Request already has a context and the context is cancelled when the client closes the connection,
	// or the request times out or ServeHTTP method returns
	start := time.Now()
	err := slowQuery(req.Context())
	elapsed := time.Since(start)
	if err != nil {
		log.Printf("ERROR : %s\n", err.Error())
		return
	}
	fmt.Fprintf(w, "OK")
	fmt.Printf("slowHandler took :%v\n", elapsed)
}

// Without any context being passed in, that handles timeouts, the slowQuery does not cancel and completes it's operation.
// Even though from the servers perspective, it is not going to write that response to the response body.
// Which means the client gets a response back after 5 seconds, but it is an empty response.
// This is what needs to change in the implementation :
func slowQuery(ctx context.Context) error {
	// _, err := db.Exec("SELECT pg_sleep(5)")
	_, err := db.ExecContext(ctx, "SELECT pg_sleep(5)") // This is the context aware database call which terminates the query when the context is cancelled
	return err
}

func describe() {
	str := `
The different timeouts in go for http servers.

<---Wait---><---TLS Handshake---><---Req Header---><---Req body---><---Response---><---Idle--->
                                 |----------------|                                             : ReadHeader Timeout
|-----------------------------------------------------------------|                             : Read Timeout
                                                                  |----------------|            : Write Timeout with no TLS
|----------------------------------------------------------------------------------|            : Write Timeout with TLS
                                                                                   |----------| : Idle Timeout - waiting for next request if keep-alive is enabled

These are different http timeouts in go http server.
We should set these timeouts when we are dealing with untrusted clients, networks, or want us to be somewhat safe from attacks.
This also protects the servers from clients holding up connections that are slow to write or slow to read.
A connection consumes a file descriptor. So the servers have to be dilligent in how to handle these client requests efficiently.
Http handler funcs are usually unaware of these timeouts set by the server. So they will run to completion even if the timers expire.
To prevent this from happening we must pass the contexts down to the handler funcs which are usually run in separate child goroutines.

_____________________
	`
	fmt.Println(str)
}
