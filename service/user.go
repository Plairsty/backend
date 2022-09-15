package service

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

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

func createUser(userStore UserStore, username, password, role string) error {

	user, err := NewUser(username, password, role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}

func SeedUsers(userStore UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	log.Println("Admin user created")
	err = createUser(userStore, "user1", "secret", "user")
	if err != nil {
		return err
	}
	log.Println("User created")
	return nil
}
