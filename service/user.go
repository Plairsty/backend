package service

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username       string
	HashedPassword string
	Role           string
}

func NewUser(username, password, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}, nil
}

// If error is nil, then the password is correct.
func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
    return err == nil
}

func (user *User) Clone() *User {
	return &User{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}