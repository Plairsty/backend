package util

import (
	"math/rand"
	"strings"
	"time"
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
