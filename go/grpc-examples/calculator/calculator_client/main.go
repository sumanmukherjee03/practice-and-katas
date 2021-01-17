package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	doUnary(c)
	doServerStreaming(c)
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
