package main

import (
	"context"
	"fmt"
	"github.com/PereRohit/go-gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type server struct{}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline called with: %v\n", req)

	// loop to sleep for 3 seconds - suggesting expensive task
	for i := 0; i < 3; i++ {
		// check if the context is cancelled by the client
		if ctx.Err() == context.Canceled {
			fmt.Println("Client cancelled the request")
			return nil, status.Error(codes.Canceled, "Client cancelled, abandoning")
		}
		time.Sleep(1 * time.Second)
	}

	res := &greetpb.GreetWithDeadlineResponse{
		Result: "Hello " + req.GetGreeting().GetFirstName(),
	}
	return res, nil
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receiving client stream: %v\n", err)
		}
		fName := req.GetGreeting().GetFirstName()
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + fName,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
			return err
		}
	}

	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Println("LongGreet invoked")

	var nameSlc []string

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: strings.Join(nameSlc, ", "),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
		nameSlc = append(nameSlc, req.GetGreeting().GetFirstName()+" "+req.GetGreeting().GetLastName())
	}
	return nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Greeting " + firstName + " for " + strconv.Itoa(i+1) + " times"
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1 + time.Second)
	}

	return nil
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello world!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Setting up ssl encryption
	// https://grpc.io/docs/guides/auth/#go

	// set to true to enable ssl
	enableSSL := false
	opts := []grpc.ServerOption{}

	if enableSSL {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed to load certificates: %v\n", sslErr)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
