package main

import (
	"context"
	"fmt"
	calculatepb "github.com/PereRohit/go-gRPC/calculator/calculatepb"
	"google.golang.org/grpc"
	"io"
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

	calculatePrimeNumberDecomposition(c)
}

func calculatePrimeNumberDecomposition(c calculatepb.CalculatorServiceClient) {
	fmt.Println("Calling remote PrimeNumberDecomposition function")

	req := &calculatepb.PrimeNumberDecompositionRequest{Number: 120}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't call PrimeNumberDecomposition function: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		fmt.Println(msg.GetFactor())
	}
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
