package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	resp, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("Encountered an error making a request : %v", err)
	}
	fmt.Printf("Got response from server : %s", resp.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	resStream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("Encountered an error in server streaming : %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			fmt.Println("We have reached the end of the stream")
			break
		}
		if err != nil {
			log.Fatalf("Encountered an error in receiving streaming result : %v", err)
		}
		fmt.Printf("Got response from server : %s", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientStream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("Encountered an error making call to server with client streaming : %v", err)
	}
	for i := 0; i < 10; i++ {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: fmt.Sprintf("John %d", i),
				LastName:  "Doe",
			},
		}
		clientStream.Send(req)
		time.Sleep(1 * time.Second)
	}
	resp, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Encountered an error receiving response from server with client streaming : %v", err)
	}
	fmt.Printf("Client streaming response : %s\n", resp.GetResult())
}
