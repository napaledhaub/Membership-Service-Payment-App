package utils

import (
	"math/rand"
	"paymentapp/models"
	"time"

	"github.com/google/uuid"
)

func GenerateVerificationCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	codeLength := 6
	allowedChars := "0123456789"
	verificationCode := make([]byte, codeLength)

	for i := 0; i < codeLength; i++ {
		randomIndex := r.Intn(len(allowedChars))
		verificationCode[i] = allowedChars[randomIndex]
	}

	return string(verificationCode)
}

func GenerateAuthToken(participant *models.Participant) models.AuthToken {
	token := uuid.NewString()
	expirationDateTime := time.Now().Add(1 * time.Hour)

	return models.AuthToken{
		Token:              token,
		ExpirationDateTime: expirationDateTime,
		ParticipantID:      participant.ID,
	}
}
