package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (res *calculatorpb.SumResponse, err error) {
	fmt.Printf("Received Sum RPC: %v\n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res = &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) (err error) {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: int64(divisor),
			}
			stream.Send(res)
			number /= divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %d\n", divisor)
		}
	}
	return
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) (err error) {
	fmt.Println("Received ComputeAverage RPC")
	var sum, count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float64(sum / count),
			})
		} else if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		count++
		sum += req.GetNumber()
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) (err error) {
	fmt.Println("Received FindMaximum RPC.")

	var maximum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		if number := req.GetNumber(); maximum < number {
			maximum = number
			if err := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			}); err != nil {
				log.Fatalf("Error while sending data to client.")
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (res *calculatorpb.SquareRootResponse, err error) {
	fmt.Println("Received SquareRoot RPC.")
	if number := req.GetNumber(); number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("received a negative number: %v", number))
	} else {
		return &calculatorpb.SquareRootResponse{
			NumberRoot: math.Sqrt(float64(number)),
		}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
