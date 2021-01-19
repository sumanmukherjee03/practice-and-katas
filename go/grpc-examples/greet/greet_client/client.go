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
	// doClientStreaming(c)
	doBiDirectionalStreaming(c)
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

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	waitChan := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := c.GreetEveryone(ctx)
	if err != nil {
		log.Fatalf("Encountered an error making call to server with bi-directional streaming : %v", err)
	}

	go func() {
		defer stream.CloseSend() // Make sure to defer close the stream in the goroutine where you are sending data
		for i := 0; i < 10; i++ {
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: fmt.Sprintf("John %d", i),
					LastName:  "Doe",
				},
			}
			err := stream.Send(req)
			if err != nil {
				log.Printf("Encountered an error sending requests to server with bi-directional streaming : %v", err)
				return
			}
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		defer close(waitChan)
		// Since this is an infinite loop dependent upon the fact that the server will
		// close the stream, you cant send a value to the wait channel here.
		// That will cause the wait channel to block forever.
		// One way to unblock the wait channel is by closing the channel.
		// So, we defer close the channel as soon as the goroutine finishes which can handle the error path of the goroutine.
		// You can also send a value to the wait channel when you encounter an error and are about to exit the goroutine,
		// but when combined with the close(waitChan) that logic becomes redundant.
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Encountered an error receiving responses from server with bi-directional streaming : %v", err)
				return
			}
			fmt.Println(resp.GetResult())
		}
	}()

	<-waitChan // This blocks the function until the goroutines are done
}
