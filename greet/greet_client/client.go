package main

import (
	"context"
	"fmt"
	"github.com/PereRohit/go-gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello world!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// unary API call
	// doUnary(c)

	// Server streaming API call
	//doServerStreaming(c)

	// Client Streaming
	//doClientStreaming(c)

	// Bi-directional streaming
	//doBiDiStreaming(c)

	// Unary With Deadline call
	doUnaryWithDeadline(c)

	//fmt.Printf("Created client: %f", c)
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary Deadline rpc")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rohit",
			LastName:  "Sadhukhan",
		},
	}

	// Background context is the parent context & timeout is set on it
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)   // 5 sec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) // 1 sec timeout
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statErr, ok := status.FromError(err)
		if ok {
			if statErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout reached! Deadline exceeded:", statErr.Message())
			} else {
				log.Fatalf("Error while calling Unary Deadline function: %v\n", err)
			}
		}
		return
	}
	fmt.Println(res.GetResult())
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting bi-di streaming")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while starting stream: %v\n", err)
	}

	var wg sync.WaitGroup

	// trigger goroutines to send data
	go func() {
		req := []*greetpb.GreetEveryoneRequest{
			&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
				FirstName: "Joe",
				LastName:  "Salman",
			}},
			&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			}},
			&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
				FirstName: "Rohit",
				LastName:  "Sadhukhan",
			}},
			&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
				FirstName: "Jane",
				LastName:  "Doe",
			}},
		}
		for _, req := range req {
			fmt.Println("Sending:", req)
			stream.Send(req)
		}

		// close sending to server
		stream.CloseSend()
	}()

	wg.Add(1)
	// trigger goroutines to receive data
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				wg.Done()
				return
			}
			if err != nil {
				log.Fatalf("Error while receiving stream: %v\n", err)
			}
			fmt.Println("From server:", res.GetResult())
		}
	}()
	wg.Wait()
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Couldn't call LongGreet: %v\n", err)
	}

	req := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "Joe",
			LastName:  "Salman",
		}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "Rohit",
			LastName:  "Sadhukhan",
		}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "Jane",
			LastName:  "Doe",
		}},
	}

	for _, name := range req {
		stream.Send(name)
		fmt.Printf("Sending name: %v\n", name)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving from server: %v\n", err)
	}
	fmt.Println("Combined names", res.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting server streaming")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't call GreetManyTimes: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}
		fmt.Println(msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rohit",
			LastName:  "Sadhukhan",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet rpc: %v", err)
	}
	log.Println(res.GetResult())
}
