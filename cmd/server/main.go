package main

import (
	"context"
	"flag"
	"log"
	"net"
	__pb "plairsty/backend/pb"
	"plairsty/backend/service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey            = "secret"
	tokenDuration        = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

var (
	port = flag.String("port", ":8080", "Port to connect to")
)

type serverImpl struct {
	__pb.UnimplementedGreetServiceServer
}

func (*serverImpl) Greet(ctx context.Context,
	in *__pb.GreetRequest,
) (*__pb.GreetResponse, error) {
	firstname := in.GetGreeting().GetFirstName()
	res := &__pb.GreetResponse{
		Result: firstname + in.GetGreeting().GetLastName(),
	}
	return res, nil
}

func main() {
	listenc, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listenc.Close()
	log.Printf("Server is listening on port %s", *port)
	log.Printf("https://localhost%s/", *port)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(service.UniaryInterceptor),
		grpc.StreamInterceptor(service.StreamInterceptor),
	)

	userStore := service.NewInMemoryUserStore()
	err = service.SeedUsers(userStore)
	if err != nil {
		log.Fatalln("Could not seed users", err)
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration, refreshTokenDuration)

	authServer := service.NewAuthServer(userStore, jwtManager)
	__pb.RegisterAuthServiceServer(server, authServer)

	__pb.RegisterGreetServiceServer(server, &serverImpl{})
	reflection.Register(server)
	if err := server.Serve(listenc); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
