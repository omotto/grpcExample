package main

import (
	"context"
	"google.golang.org/grpc"
	"grpcExample/pb"
	"log"
	"time"
)

const address = "localhost:7777"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExmapleServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Generate(ctx, &pb.ExampleRequest{Email: "name", Size: 10})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Returned: %s", r.Url)
}