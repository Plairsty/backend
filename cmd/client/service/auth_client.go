package service

import (
	"context"
	"log"
	__pb "plairsty/backend/pb"
	"time"

	"google.golang.org/grpc"
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

func (client *AuthClient) Register() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	res, err := client.service.Register(ctx, &__pb.RegisterRequest{
		Username:   client.username,
		FirstName:  "Gulshan",
		LastName:   "Yadav",
		MiddleName: "Mohan",
		Email:      "gulshan@duck.com",
		Phone:      "1234567890",
		Mobile:     "1234567890",
		Role:       "admin",
	})
	if err != nil {
		return false , err
	}
	log.Println("Register response: ", res)

	return res.GetSuccess(), nil
}
