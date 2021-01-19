package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting grpc client for calculator service")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not setup connection to grpc server from client : %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &calculatorpb.SumRequest{
		Operands: []float64{5.7, 2.3, 9.1, 6.2},
		Operator: "Sum",
	}

	resp, err := c.Sum(ctx, req)
	if err != nil {
		log.Fatalf("Received an error from the grpc server : %v", err)
	}

	fmt.Printf("Result of calculation : %f\n", resp.GetResult())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &calculatorpb.PrimeNumberDecompositionRequest{Number: 120}
	stream, err := c.PrimeNumberDecomposition(ctx, req)
	if err != nil {
		log.Fatalf("Received an error from the grpc server : %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Received an error during streaming of the response : %v", err)
		}
		fmt.Printf("Received factor : %d\n", msg.GetFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := c.ComputedAverage(ctx)
	if err != nil {
		log.Fatalf("Encountered an error setting up the connection to the server for client streaming : %v", err)
	}
	for i := 0; i < 10; i++ {
		req := &calculatorpb.ComputedAverageRequest{Number: rand.Float64() * 100}
		stream.Send(req)
		time.Sleep(300 * time.Millisecond)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Encountered an error closing stream and fetching final result : %v", err)
	}
	fmt.Printf("Final comnputed average : %f\n", resp.GetResult())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	rand.Seed(time.Now().UnixNano())
	waitCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := c.FindMaximum(ctx)
	if err != nil {
		log.Fatalf("Encountered an error setting up the connection to the server for bidirectional streaming : %v", err)
	}

	go func() {
		defer stream.CloseSend()
		fmt.Println("Starting sender goroutine")
		for i := 0; i < 10; i++ {
			err := stream.Send(&calculatorpb.FindMaximumRequest{Number: rand.Int31()})
			if err != nil {
				log.Printf("Encountered an error sending requests to server with bi-directional streaming : %v", err)
				return
			}
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		defer close(waitCh)
		fmt.Println("Starting receiver goroutine")
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server closed the stream")
				return
			}
			if err != nil {
				log.Fatalf("Encountered an error fetching result from server : %v", err)
				return
			}
			fmt.Printf("Current max is : %d\n", msg.GetResult())
		}
	}()

	<-waitCh
}
