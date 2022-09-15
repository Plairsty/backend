package main

import (
	"flag"
	"log"
	"net"
	__auth "plairsty/backend/pb"
	"plairsty/backend/service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

var (
	port = flag.String("port", ":8080", "Port to connect to")
)

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
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)

	authServer := service.NewAuthServer(userStore, jwtManager)
	__auth.RegisterAuthServiceServer(server, authServer)

	reflection.Register(server)
	// __auth.RegisterGreetServiceServer(server, &serverImpl{})
	if err := server.Serve(listenc); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
