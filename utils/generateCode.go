package utils

import (
	"math/rand"
	"time"
)

func GenerateVerificationCode() string {
	codeLength := 6
	allowedChars := "0123456789"
	rand.Seed(time.Now().UnixNano())

	verificationCode := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		randomIndex := rand.Intn(len(allowedChars))
		verificationCode[i] = allowedChars[randomIndex]
	}

	return string(verificationCode)
}
