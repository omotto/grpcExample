package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/testdata"
	"grpcExample/pb"
	"log"
	"net"
	"strings"
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

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

func serverUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("errMissingMetadata")
	}
	if !valid(md["authorization"]) {
		return nil, errors.New("errInvalidToken")
	}
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return m, err
}

// main start a gRPC server and waits for connection
func main() {

	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	// create a gRPC server object
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(serverUnaryInterceptor))

	pb.RegisterExmapleServiceServer(grpcServer, &exampleService{})

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}




