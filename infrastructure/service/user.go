package service

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string
	first_name     string
	last_name      string
	phone          string
	mobile         string
	email          string
	HashedPassword string
	Role           string
	createdBy      string
}

type RequiredUserFields struct {
	Username   string
	Password   string
	First_name string
	Last_name  string
	Phone      string
	Mobile     string
	Email      string
	CreatedBy  string
	Role       string
}

func NewUser(
	user RequiredUserFields,
) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:       user.Username,
		HashedPassword: string(hashedPassword),
		first_name:     user.First_name,
		last_name:      user.Last_name,
		phone:          user.Phone,
		mobile:         user.Mobile,
		email:          user.Email,
		createdBy:      user.CreatedBy,
		Role:           user.Role,
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

func createUser(userStore UserStore, newUser RequiredUserFields) error {

	user, err := NewUser(newUser)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}

func SeedUsers(userStore UserStore) error {
	adminUser := RequiredUserFields{
		Username:   "admin",
		Password:   "secret",
		First_name: "admin",
		Last_name:  "admin",
		Phone:      "admin",
		Mobile:     "admin",
		Email:      "admin",
		CreatedBy:  "system",
		Role:       "admin",
	}
	hrUser := RequiredUserFields{
		Username:   "hr",
		Password:   "secret",
		First_name: "hr",
		Last_name:  "hr",
		Phone:      "hr",
		Mobile:     "hr",
		Email:      "hr",
		CreatedBy:  "system",
		Role:       "hr",
	}
	studentUser := RequiredUserFields{
		Username:   "student",
		Password:   "secret",
		First_name: "student",
		Last_name:  "student",
		Phone:      "student",
		Mobile:     "student",
		Email:      "student",
		CreatedBy:  "system",
		Role:       "user",
	}

	err := createUser(userStore, adminUser)
	if err != nil {
		return err
	}
	log.Println("Admin user created")
	err = createUser(userStore, hrUser)
	if err != nil {
		return err
	}
	err = createUser(userStore, studentUser)
	if err != nil {
		return err
	}
	log.Println("User created")
	return nil
}
