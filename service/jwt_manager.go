package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey            string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWTManager(secretKey string, tokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		tokenDuration:        tokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func userClaimProvider(username, role string, expiresAt time.Duration) UserClaims {
	return UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			Issuer:    "plairsty",
		},
		Username: username,
		Role:     role,
	}
}

func (manager *JWTManager) Generate(user *User) (JWTToken, error) {
	claims := userClaimProvider(user.Username, user.Role, manager.tokenDuration)
	refreshClaims := userClaimProvider(user.Username, user.Role, manager.refreshTokenDuration)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	accessTokenVal, err := accessToken.SignedString([]byte(manager.secretKey))
	if err != nil {
		return JWTToken{}, err
	}
	refreshTokenVal, err := refreshToken.SignedString([]byte(manager.secretKey))
	if err != nil {
		return JWTToken{}, err
	}
	return JWTToken{
		AccessToken:  accessTokenVal,
		RefreshToken: refreshTokenVal,
	}, nil
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	/*
		Pass in the access token, an empty user claims, and a custom key function.
		In this function,
		It’s very important to check the signing method of the token to make sure that
		it matches with the algorithm our server uses, which in our case is RSA. (Now HMAC 🥹)
		If it matches, then we just return the secret key that is used to sign the token.
	*/

	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
