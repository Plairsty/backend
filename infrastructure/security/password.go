package security

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func GeneratePassword() string {
	rand.Seed(time.Now().Unix())
	var password strings.Builder
	for i := 0; i < 8; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteByte(specialCharSet[random])
	}
	for i := 0; i < 8; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteByte(numberSet[random])
	}
	for i := 0; i < 8; i++ {
		random := rand.Intn(len(lowerCharSet))
		password.WriteByte(lowerCharSet[random])
	}

	for i := 0; i < 8; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteByte(upperCharSet[random])
	}

	for i := 0; i < 8; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteByte(allCharSet[random])
	}

	iRune := []rune(password.String())
	rand.Shuffle(len(iRune), func(i, j int) { iRune[i], iRune[j] = iRune[j], iRune[i] })
	return password.String()
}


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
