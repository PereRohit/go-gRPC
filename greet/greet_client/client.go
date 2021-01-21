package main

import (
	"context"
	"fmt"
	"github.com/PereRohit/go-gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
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
	doClientStreaming(c)
	//fmt.Printf("Created client: %f", c)
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
			LastName:  "Saldana",
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
