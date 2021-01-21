package main

import (
	"context"
	"fmt"
	calculatepb "github.com/PereRohit/go-gRPC/calculator/calculatepb"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
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

	calculateAverage(c)

	calculateMaximum(c)
}

func calculateMaximum(c calculatepb.CalculatorServiceClient) {
	fmt.Println("Calling remote FindMaximum bi-di streaming function")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling FindMaximum: %v\n", err)
	}

	go func() {
		numbers := [...]int32{43, 56, 7, 2, 45, 67, 9}
		for _, num := range numbers {
			stream.Send(&calculatepb.FindMaximumRequest{Number: num})
		}
		stream.CloseSend()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				wg.Done()
				return
			}
			if err != nil {
				log.Fatalf("Error  while receiving stream from server: %v\n", err)
			}
			fmt.Println("New maximum", res.GetMaximum())
		}
	}()
	wg.Wait()
}

func calculateAverage(c calculatepb.CalculatorServiceClient) {
	fmt.Println("Calling remote ComputeAverage streaming function")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v\n", err)
	}
	numbers := [...]int32{4, 90, 3, 10, 56, 78}
	for _, num := range numbers {
		fmt.Println("Sending:", num)
		stream.Send(&calculatepb.CalculateAverageRequest{
			Number: num,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reading response: %v\n", err)
	}
	fmt.Println(res.GetResponse())
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
