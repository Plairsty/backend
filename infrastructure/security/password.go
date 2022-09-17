package security

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func SaltedHash(password string, salt []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+string(salt)), bcrypt.DefaultCost)
}

// CompareHashAndPassword If error is nil, then the password is correct.
func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func CompareSaltedHashAndPassword(hashedPassword, password, salt []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(string(password)+string(salt)))
}

func GenerateRandomString(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(salt)
}
