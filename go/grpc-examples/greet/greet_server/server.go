package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This server needs to implement the GreetServiceServer interface defined in the code generated from the protobuf
type server struct{}

// Implements - Greet(context.Context, *GreetRequest) (*GreetResponse, error)
func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with : %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName + "\n"
	res := &greetpb.GreetResponse{Result: result}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	for i := 0; i < 10; i++ {
		resp := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hi %s %s : %d\n", firstName, lastName, i),
		}
		stream.Send(resp)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	var names []string
	for {
		req, err := stream.Recv() // Stream blocks until you receive something from it
		if err == io.EOF {
			resp := &greetpb.LongGreetResponse{Result: "Hello " + strings.Join(names, ",")}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received client streaming request : %v\n", req)
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		names = append(names, fmt.Sprintf("%s %s", firstName, lastName))
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF { // This indicates that the client has stopped streaming
			return nil // Dont return the error because this successfully terminates the loop
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received client streaming request : %v\n", req)
		resp := &greetpb.GreetEveryoneResponse{Result: fmt.Sprintf("Hello %s %s\n", req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName())}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
		time.Sleep(400 * time.Millisecond)
	}
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("Greet with deadline function was invoked with : %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled { // Check if the client has cancelled the request
			fmt.Println("Client cancelled the request")
			return nil, status.Error(codes.DeadlineExceeded, "The client cancelled the request") // return an error from the server with the proper error code
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName + "\n"
	res := &greetpb.GreetWithDeadlineResponse{Result: result}
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
