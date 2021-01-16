package main

import (
	"context"
	"fmt"
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
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &calculatorpb.CalculationRequest{
		Operands: []float64{5.7, 2.3, 9.1, 6.2},
		Operator: "Sum",
	}

	resp, err := c.Calculation(ctx, req)
	if err != nil {
		log.Fatalf("Received an error from the grpc server : %v", err)
	}

	fmt.Printf("Result of calculation : %f\n", resp.GetResult())
}
