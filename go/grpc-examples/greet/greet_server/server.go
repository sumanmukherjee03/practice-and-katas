package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/greet/greetpb"
	"google.golang.org/grpc"
)

// This server needs to implement the GreetServiceServer interface defined in the code generated from the protobuf
type server struct{}

// Implements - Greet(context.Context, *GreetRequest) (*GreetResponse, error)
func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetResponse{Result: result}
	return res, nil
}

func main() {
	fmt.Println("Starting new grpc server")
	lis, err := net.Listen("tcp", "localhost:50051") // 50051 is the default port for grpc
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	// boilerplate code to create a new grpc server listening on a host:port
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start server : %s", err)
	}
}
