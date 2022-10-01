package main

import (
	"context"
	"flag"
	"log"
	"net"
	"plairsty/backend/infrastructure/service"
	__pb "plairsty/backend/pb"
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

func (*serverImpl) Greet(_ context.Context,
	in *__pb.GreetRequest,
) (*__pb.GreetResponse, error) {
	firstname := in.GetGreeting().GetFirstName()
	res := &__pb.GreetResponse{
		Result: firstname + in.GetGreeting().GetLastName(),
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer func(conn net.Listener) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v\n", err)
		}
	}(lis)
	log.Printf("Server is listening on port %s", *port)
	log.Printf("https://localhost%s/", *port)

	jwtManager := service.NewJWTManager(secretKey, tokenDuration, refreshTokenDuration)
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	userStore := service.NewInMemoryUserStore()
	err = service.SeedUsers(userStore)
	if err != nil {
		log.Fatalln("Could not seed users", err)
	}
	authServer := service.NewAuthServer(userStore, jwtManager)

	__pb.RegisterAuthServiceServer(server, authServer)
	__pb.RegisterGreetServiceServer(server, &serverImpl{})

	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// If target RPC method is not in accessibleRoles map, then it is publicly accessible
func accessibleRoles() map[string][]string {
	const greetServicePath = "/greet.GreetService/"
	const authServicePath = "/auth.AuthService/"
	return map[string][]string{
		greetServicePath + "Greet": {"user"},
		// authServicePath + "Register": {"admin"},
	}
}
