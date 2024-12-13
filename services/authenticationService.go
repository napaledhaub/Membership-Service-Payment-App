package services

import (
	"errors"
	"paymentapp/models"
	"paymentapp/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	DB *gorm.DB
}

func (s *AuthenticationService) Login(email, password string) (string, error) {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("Incorrect email or password")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(participant.Password), []byte(password)); err != nil {
		return "", errors.New("Incorrect email or password")
	}

	authToken := utils.GenerateAuthToken(&participant)
	if err := s.DB.Save(&authToken).Error; err != nil {
		return "", err
	}

	return authToken.Token, nil
}

func (s *AuthenticationService) Logout(token string) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Delete(&authToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Logout failed")
		}
		return err
	}
	return nil
}

func (service *AuthenticationService) ValidateOtp(token string, otp string) error {
	var authToken models.AuthToken
	if err := service.DB.Preload("Participant").Where("token = ?", token).First(&authToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Invalid or expired token")
		}
		return err
	}

	participant := authToken.Participant
	if otp == participant.VerificationCode {
		participant.IsVerified = true
		participant.VerificationCode = ""
		if err := service.DB.Save(&participant).Error; err != nil {
			return err
		}
		return nil
	}
	return errors.New("The OTP you entered is incorrect")
}

func (service *AuthenticationService) RefreshToken(token string) (*models.AuthToken, error) {
	var existingToken models.AuthToken
	if err := service.DB.Where("token = ?", token).First(&existingToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid or expired token")
		}
		return nil, err
	}

	existingToken.UpdateExpiration()
	if err := service.DB.Save(&existingToken).Error; err != nil {
		return nil, err
	}
	return &existingToken, nil
}
