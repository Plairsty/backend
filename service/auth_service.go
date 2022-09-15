package service

import (
	"context"
	__auth "plairsty/backend/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	userStore  UserStore
	jwtManager *JWTManager
	__auth.UnimplementedAuthServiceServer
}

func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthService {
	return &AuthService{
		userStore,
		jwtManager,
		__auth.UnimplementedAuthServiceServer{},
	}
}

func (server *AuthService) Login(
	ctx context.Context,
	req *__auth.LoginRequest,
) (*__auth.LoginResponse, error) {
	user, err := server.userStore.Find(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "Invalid username or password")
	}
	ok := user.IsCorrectPassword(req.Password)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Invalid username or password")
	}
	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate token: %v", err)
	}

	res := &__auth.LoginResponse{
		AccessToken: token,
	}
	return res, nil
}
