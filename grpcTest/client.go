package main

import (
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/testdata"
	"grpcExample/pb"
	"log"
	"time"
)

const (
	address = "localhost:7777"
	fallbackToken = "some-secret-token"
)

// unaryInterceptor in order to add auth for each sent message
func clientUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: fallbackToken,
		})))
	}
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	log.Printf("RPC: %s, start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)
	return err
}


func main() {
	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(
		address, // Address where gRPC server is located
 	    grpc.WithTransportCredentials(creds)/*, grpc.WithInsecure()*/, // Credentials?
		grpc.WithUnaryInterceptor(clientUnaryInterceptor), // Unary Interceptor example to add token
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExmapleServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()

	r, err := c.Generate(ctx, // Context
		&pb.ExampleRequest{Email: "name", Size: 10}, // Method to call through gRPC
		grpc.UseCompressor(gzip.Name), // Add compression on message
		grpc.WaitForReady(true), // Add Wait for server ready on gRPC communication
	)

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Returned: %s", r.Url)
}