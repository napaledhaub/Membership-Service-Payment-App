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

func (s *PasswordResetService) RequestPasswordReset(email string) (bool, error) {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return false, errors.New("participant not found")
	}

	resetToken := uuid.New().String()
	participant.ResetToken = resetToken

	if err := s.DB.Save(&participant).Error; err != nil {
		return false, err
	}

	/*subject := "Reset Password"
	text := "You have requested a password reset. Here is the token to reset your password: " + resetToken
	if err := s.EmailService.SendEmail(request.Email, subject, text); err != nil {
		return "", err
	}*/

	return true, nil
}

func (s *PasswordResetService) ConfirmPasswordReset(resetToken, password string) (bool, error) {
	var participant models.Participant
	if err := s.DB.Where("reset_token = ?", resetToken).First(&participant).Error; err != nil {
		return false, errors.New("participant not found")
	}

	participant.ResetToken = ""
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	participant.Password = string(hashedPassword)

	if err := s.DB.Save(&participant).Error; err != nil {
		return false, err
	}

	return true, nil
}
