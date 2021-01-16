package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type calculatorServer struct{}

func (c *calculatorServer) Calculation(ctx context.Context, req *calculatorpb.CalculationRequest) (*calculatorpb.CalculationResponse, error) {
	fmt.Println("Calculation method called in calculator service")
	var sum float64
	for _, num := range req.GetOperands() {
		sum += num
	}
	resp := &calculatorpb.CalculationResponse{Result: sum}
	return resp, nil
}

func main() {
	fmt.Println("Starting grpc server for calculator service")
	lis, err := net.Listen("tcp", "localhost:50051") // default port for grpc is 50051
	if err != nil {
		log.Fatalf("Encountered an error in creating tcp listener for grpc server : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &calculatorServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Encountered error starting grpc server : %v", err)
	}
}
