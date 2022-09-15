package main

import (
	"context"
	"flag"
	"log"
	"plairsty/backend/cmd/client/service"
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
	addr      = flag.String("addr", "localhost:8080", "Address to connect to")
	firstName = flag.String("first_name", DefaultFirstName, "FirstName of user")
	lastName  = flag.String("last_name", DefaulLastName, "LastName of user")
)

func main() {
	conn1, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(conn1)

	// AuthClient
	authClient := service.NewAuthClient(conn1, "admin1", "secret")

	// Interceptor
	interceptor, err := service.NewAuthInterceptor(authClient, authMethods(), 30*time.Second)

	// Second connection
	conn2, err := grpc.Dial(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}(conn2)

	c := __pb.NewGreetServiceClient(conn2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	Greet(c, ctx)
}

func Greet(c __pb.GreetServiceClient, ctx context.Context) {
	r, err := c.Greet(ctx, &__pb.GreetRequest{
		Greeting: &__pb.Greeting{
			FirstName: *firstName,
			LastName:  *lastName,
		},
	})

	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetResult())
}

func authMethods() map[string]bool {
	const greetServicePath = "/greet.GreetService/"
	return map[string]bool{
		greetServicePath + "Greet": true,
	}
}
