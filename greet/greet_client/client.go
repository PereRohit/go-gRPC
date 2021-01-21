package main

import (
	"context"
	"fmt"
	"github.com/PereRohit/go-gRPC/greet/greetpb"
	"google.golang.org/grpc"
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
	doUnary(c)
	//fmt.Printf("Created client: %f", c)
}

func doUnary (c greetpb.GreetServiceClient){

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rohit",
			LastName: "Sadhukhan",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil{
		log.Fatalf("Error while calling greet rpc: %v", err)
	}
	log.Println(res.GetResult())
}
