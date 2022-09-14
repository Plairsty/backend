package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func UniaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println("--> Unary Interceptor", info.FullMethod)
	return handler(ctx, req)
}

func StreamInterceptor(srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("--> Stream Interceptor", info.FullMethod)
	return handler(srv, ss)
}

func main() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UniaryInterceptor),
		grpc.StreamInterceptor(StreamInterceptor),
	)
	reflection.Register(grpcServer)
	address := fmt.Sprintf(":%d", 8080)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	log.Println("Server is running on", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
