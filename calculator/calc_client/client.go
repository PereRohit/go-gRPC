package main

import (
	"context"
	"fmt"
	calculatepb "github.com/PereRohit/go-gRPC/calculator/calculatepb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer cc.Close()

	c := calculatepb.NewCalculatorServiceClient(cc)

	calculateSum(c)
}

func calculateSum(c calculatepb.CalculatorServiceClient) {
	fmt.Println("Calling remote Sum function")

	req := &calculatepb.CalcRequest{
		Num1: 20,
		Num2: 45,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Unable to call remote Sum function: %v", err)
	}
	fmt.Println("Sum:", res.GetResult())
}
