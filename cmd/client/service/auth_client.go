package service

import (
	"context"
	"google.golang.org/grpc"
	__pb "plairsty/backend/pb"
	"time"
)

type AuthClient struct {
	service  __pb.AuthServiceClient
	username string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	authServiceClient := __pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service:  authServiceClient,
		username: username,
		password: password,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	res, err := client.service.Login(ctx, &__pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	})
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
