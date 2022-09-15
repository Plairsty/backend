package service

import (
	"context"

	__pb "plairsty/backend/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	userStore  UserStore
	jwtManager *JWTManager
	__pb.UnimplementedAuthServiceServer
}

func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthService {
	return &AuthService{
		userStore,
		jwtManager,
		__pb.UnimplementedAuthServiceServer{},
	}
}

func (server *AuthService) Login(
	ctx context.Context,
	req *__pb.LoginRequest,
) (*__pb.LoginResponse, error) {
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

	res := &__pb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return res, nil
}

// Invalidate access token and refresh token
// Possible ways:
// We can blacklist the token in the database
// For this using redis would be a good idea as the data has a short TTL
// Access Token expires in 15 minutes
// Refresh Token expires in 7 days
//
// We can decode the token and check the expiry time
// Set the expiry time in redis to the expiry time of the token
func (server *AuthService) Logout(
	ctx context.Context,
	req *__pb.LogoutRequest,
) (*__pb.LogoutResponse, error) {
	// TODO: implement
	return nil, nil
}
