package utils

import (
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
}

func (e *EmailService) SendEmail(to, subject, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(e.SMTPHost, e.SMTPPort, e.Username, e.Password)

	return d.DialAndSend(m)
}

func (e *EmailService) GenerateVerificationCode() string {
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
