package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"google.golang.org/grpc"
	"grpcExample/pb"
	"log"
	"net"
)

const port = ":7777"

type exampleService struct{}


func giveHash(email string) [16]byte {
	val := md5.Sum([]byte(email))
	return val
}

func giveURL(hash[16]byte, size int32) string {
	return fmt.Sprintf("https://www.omotto.com/%x?s=%d", hash, size)
}

func MyMethod(email string, size int32) string {
	hash := giveHash(email)
	return giveURL(hash, size)
}

func (s *exampleService) Generate(ctx context.Context, in *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	log.Printf("Received email %v with size %v", in.Email, in.Size)
	return &pb.ExampleResponse{Url: MyMethod(in.Email, in.Size)}, nil
}

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	// create a gRPC server object
	grpcServer := grpc.NewServer()

	pb.RegisterExmapleServiceServer(grpcServer, &exampleService{})

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}




