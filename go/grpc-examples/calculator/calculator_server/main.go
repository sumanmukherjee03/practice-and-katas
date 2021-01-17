package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Sum method called in calculator service")
	var sum float64
	for _, num := range req.GetOperands() {
		sum += num
	}
	resp := &calculatorpb.SumResponse{Result: sum}
	return resp, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	var num int32
	num = req.GetNumber()
	factor := int32(2)
	for num > int32(1) {
		// If this factor divides num then enter else try to find the next factor
		if (num % factor == 0) {
			// Keep iterating with the same factor as long as it is divisible
			for (num % factor == 0) {
				resp := &calculatorpb.PrimeNumberDecompositionResponse{Factor: factor}
				num = num/factor
				stream.Send(resp)
				time.Sleep(1 * time.Second) // Sleeping for a bit to simulate a real life working example
			}
			factor = factor + 1
		} else {
			factor = factor + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Starting grpc server for calculator service")
	lis, err := net.Listen("tcp", "localhost:50051") // default port for grpc is 50051
	if err != nil {
		log.Fatalf("Encountered an error in creating tcp listener for grpc server : %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Encountered error starting grpc server : %v", err)
	}
}
