package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func UnaryInterceptor(ctx context.Context,
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

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Println("--> Unary Auth Interceptor", info.FullMethod)
		return handler(ctx, req)
	}
}
