package main

import (
	"context"
	"fmt"
	calculatepb "github.com/PereRohit/go-gRPC/calculator/calculatepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type calServer struct{}

func (cs *calServer) SquareRoot(ctx context.Context, req *calculatepb.SquareRootRequest) (*calculatepb.SquareRootResponse, error) {
	fmt.Println("SquareRoot() requested")

	num := req.GetNumber()
	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v\n", num))
	}
	return &calculatepb.SquareRootResponse{
		Root: math.Sqrt(num),
	}, nil
}

func (cs *calServer) FindMaximum(stream calculatepb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum() requested")

	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		num := req.GetNumber()
		if max < num {
			max = num
			err := stream.Send(&calculatepb.FindMaximumResponse{Maximum: max})
			if err != nil {
				log.Fatalf("Error while streaming data to client: %v\n", err)
			}
		}
	}

	return nil
}

func (cs *calServer) ComputeAverage(stream calculatepb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage() requested")

	count := 0
	sum := 0.0

	for {
		num, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatepb.CalculateAverageResponse{
				Response: sum / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
		count++
		sum += float64(num.GetNumber())
	}

	return nil
}

func (cs *calServer) PrimeNumberDecomposition(req *calculatepb.PrimeNumberDecompositionRequest, stream calculatepb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition() requested with request: %v\n", req)

	num := req.GetNumber()
	div := int64(2)
	for num > 1 {
		if num%div == 0 {
			stream.Send(&calculatepb.PrimeNumberDecompositionResponse{
				Factor: div,
			})
			num /= div
		} else {
			div++
		}
	}
	return nil
}

func (cs *calServer) Sum(ctx context.Context, req *calculatepb.CalcRequest) (*calculatepb.CalcResponse, error) {
	fmt.Printf("Sum() requested with request: %v\n", req)
	res := &calculatepb.CalcResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return res, nil
}

func main() {
	fmt.Println("Calculator server started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Unable to bind port: %v\n", err)
	}

	s := grpc.NewServer()
	calculatepb.RegisterCalculatorServiceServer(s, &calServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
