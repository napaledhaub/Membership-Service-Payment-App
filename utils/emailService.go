package utils

import (
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
