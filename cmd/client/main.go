package main

import (
	"context"
	"flag"
	"log"
	__pb "plairsty/backend/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultFirstName = "Gulshan"
	DefaulLastName   = "Yadav"
)

var (
	addr       = flag.String("addr", "localhost:8080", "Address to connect to")
	first_name = flag.String("first_name", DefaultFirstName, "FirstName of user")
	last_name  = flag.String("last_name", DefaulLastName, "LastName of user")
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := __pb.NewGreetServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	Greet(c, ctx)
	// unary(c, ctx)
	// serverStreaming(c, ctx)
	// clientStreaming(c, ctx)
	// bidiClientStreaming(c, ctx)
}

func Greet(c __pb.GreetServiceClient, ctx context.Context) {
	r, err := c.Greet(ctx, &__pb.GreetRequest{
		Greeting: &__pb.Greeting{
			FirstName: *first_name,
			LastName:  *last_name,
		},
	})

	if err != nil {
		log.Panicf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetResult())
}
