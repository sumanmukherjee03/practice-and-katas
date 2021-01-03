package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// Server code to allow multiple concurrent connections to the server
	listener, err := net.Listen("tcp", "localhost:23000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue // If there's an error, continue with getting the next client connection
		}
		go handleConn(conn) // For each client connection start a goroutine to handle the responses
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, "response from server\n")
		if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}
