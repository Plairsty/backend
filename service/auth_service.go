package service

import (
	"context"

	__pb "plairsty/backend/pb"
	"plairsty/backend/util"

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
	_ context.Context,
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

// Logout Invalidate access token and refresh token
// Possible ways:
// We can blacklist the token in the database
// For this using redis would be a good idea as the data has a short TTL
// Access Token expires in 15 minutes
// Refresh Token expires in 7 days
//
// We can decode the token and check the expiry time
// Set the expiry time in redis to the expiry time of the token
func (server *AuthService) Logout(
	_ context.Context,
	_ *__pb.LogoutRequest,
) (*__pb.LogoutResponse, error) {
	// TODO: implement
	return nil, nil
}

func (server *AuthService) Refresh(
	_ context.Context,
	req *__pb.AccessTokenRequest,
) (*__pb.AccessTokenResponse, error) {
	// verify the refresh token
	claims, err := server.jwtManager.Verify(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}
	// Generate new access token
	token, err := server.jwtManager.Generate(&User{
		Username: claims.Username,
		Role:     claims.Role,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate token: %v", err)
	}
	// Invalidate the old refresh token
	err = server.jwtManager.Invalidate(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot invalidate token: %v", err)
	}
	res := &__pb.AccessTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return res, nil
}

func (server *AuthService) Register(
	_ context.Context,
	req *__pb.RegisterRequest,
) (*__pb.RegisterResponse, error) {
	userStore := NewInMemoryUserStore()
	registerUser := RequiredUserFields{
		First_name: req.FirstName,
		Last_name:  req.LastName,
		Username:   req.Username,
		Password:   util.GeneratePassword(),
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		Email:      req.Email,
		CreatedBy:  "admin", // Since only admin can access this endpoint
		// But this should be the user who is logged in
		// TODO: implement
		Role: req.Role,
	}
	// Create a new user
	err := createUser(userStore, registerUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}
	res := &__pb.RegisterResponse{
		Success: true,
	}
	return res, nil
}
