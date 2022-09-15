package service

import (
	"log"

	"time"
)

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethod  map[string]bool
	accessToken string
}

// NewAuthInterceptor In this function,
// first we will create a new interceptor object.
// Then we will call an internal function:
// scheduleRefreshToken() to schedule refreshing access token and pass in the refresh duration.
//
//If an error occurs, just return it. Or else, return the interceptor.
func NewAuthInterceptor(
	authClient *AuthClient,
	authMethod map[string]bool,
	refreshDuration time.Duration,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient: authClient,
		authMethod: authMethod,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}
	return interceptor, nil
}

func (interceptor *AuthInterceptor) refreshToken() error {
	accessToken, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}

	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)

	return nil
}

func (interceptor *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				log.Printf("failed to refresh token: %v", err)
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}
