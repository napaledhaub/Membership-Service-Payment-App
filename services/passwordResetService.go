package services

import (
	"errors"
	"paymentapp/models"
	"paymentapp/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PasswordResetService struct {
	DB           *gorm.DB
	EmailService *utils.EmailService
}

func (s *PasswordResetService) RequestPasswordReset(email string) error {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	resetToken := uuid.New().String()
	participant.ResetToken = resetToken

	if err := s.DB.Save(&participant).Error; err != nil {
		return err
	}

	subject := "Reset Password"
	text := "You have requested a password reset. Here is the token to reset your password: " + resetToken
	if err := s.EmailService.SendEmail(email, subject, text); err != nil {
		return err
	}

	return nil
}

func (s *PasswordResetService) ConfirmPasswordReset(resetToken, password string) error {
	var participant models.Participant
	if err := s.DB.Where("reset_token = ?", resetToken).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	participant.ResetToken = ""
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	participant.Password = string(hashedPassword)

	if err := s.DB.Save(&participant).Error; err != nil {
		return err
	}

	return nil
}
