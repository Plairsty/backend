package service

import (
	__auth "plairsty/backend/pb"
)

type AuthService struct {
	userStore  UserStore
	jwtManager *JWTManager
	__auth.UnimplementedAuthServiceServer
}