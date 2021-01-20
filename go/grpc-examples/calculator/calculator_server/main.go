package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
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
		if num%factor == 0 {
			// Keep iterating with the same factor as long as it is divisible
			for num%factor == 0 {
				resp := &calculatorpb.PrimeNumberDecompositionResponse{Factor: factor}
				num = num / factor
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

func (s *server) ComputedAverage(stream calculatorpb.CalculatorService_ComputedAverageServer) error {
	fmt.Println("Starting client streaming -")
	var sum float64
	count := 0
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received message from client : [%v]\n", msg)
		sum += msg.GetNumber()
		count++
	}
	resp := &calculatorpb.ComputedAverageResponse{Result: sum / float64(count)}
	return stream.SendAndClose(resp)
}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	maxNum := int32(0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Client closed the stream")
			return nil // return nil because this successfully terminates the loop
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received message from client : [%v]\n", msg)
		if msg.GetNumber() > maxNum {
			maxNum = msg.GetNumber()
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{Result: maxNum})
			if sendErr != nil {
				return sendErr
			}
		}
	}
}

// Simple example with error handling
func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("SquareRoot invoked")
	num := req.GetNumber()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number : %d", num))
	}
	return &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(num))}, nil
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
