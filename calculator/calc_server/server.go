package main

import (
	"context"
	"fmt"
	calculatepb "github.com/PereRohit/go-gRPC/calculator/calculatepb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type calServer struct{}

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
