package main

import (
	"fmt"
	"log"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting new grpc client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldnt connect to the server : %v", err)
	}
	defer conn.Close() // Close the connection when the client returns, no matter the return path

	// boilerplate code to create a new grpc client connecting to a grpc server at host:port
	c := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client : %f", c)
}
