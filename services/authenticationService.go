package services

import (
	"errors"
	"paymentapp/models"
	"paymentapp/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	DB            *gorm.DB
	EmailService  *utils.EmailService
	EncryptionKey []byte
}

func (s *AuthenticationService) Register(request models.Participant) (string, error) {
	var existingParticipant models.Participant
	if err := s.DB.Where("email = ?", request.Email).First(&existingParticipant).Error; err == nil {
		return "", errors.New("email is already in use")
	}

	verificationCode := s.EmailService.GenerateVerificationCode()
	request.VerificationCode = verificationCode
	request.IsVerified = false

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	request.Password = string(hashedPassword)

	creditCardInfo, err := utils.EncryptCreditCardInfo(request.CreditCardInfo, s.EncryptionKey)
	if err != nil {
		return "", err
	}
	request.CreditCardInfo = creditCardInfo

	if err := s.DB.Create(&request).Error; err != nil {
		return "", err
	}

	/*subject := "Registration Confirmation"
	text := "Thank you for registering at our fitness center. Here is your OTP to activate your account: " + verificationCode
	if err := s.EmailService.SendEmail(request.Email, subject, text); err != nil {
		return "", err
	}*/

	return verificationCode, nil
}

func (s *AuthenticationService) Login(email, password string) (string, error) {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return "", errors.New("incorrect email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(participant.Password), []byte(password)); err != nil {
		return "", errors.New("incorrect email or password")
	}

	authToken := s.generateAuthToken(&participant)
	if err := s.DB.Save(&authToken).Error; err != nil {
		return "", err
	}

	return authToken.Token, nil
}

func (s *AuthenticationService) Logout(token string) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Delete(&authToken).Error; err != nil {
		return errors.New("logout failed: " + err.Error())
	}
	return nil
}

func (service *AuthenticationService) ValidateOtp(token string, otp string) (string, error) {
	var authToken models.AuthToken
	if err := service.DB.Preload("Participant").Where("token = ?", token).First(&authToken).Error; err != nil {
		return "", errors.New("invalid token")
	}

	participant := authToken.Participant
	if otp == participant.VerificationCode {
		participant.IsVerified = true
		participant.VerificationCode = ""
		if err := service.DB.Save(&participant).Error; err != nil {
			return "", err
		}
		return "Your account has been successfully validated. Membership status: REGISTERED", nil
	}
	return "The OTP you entered is incorrect", nil
}

func (service *AuthenticationService) RefreshToken(token string) (*models.AuthToken, error) {
	var existingToken models.AuthToken
	if err := service.DB.Where("token = ?", token).First(&existingToken).Error; err != nil {
		return nil, errors.New("invalid or expired token")
	}

	existingToken.UpdateExpiration()
	if err := service.DB.Save(&existingToken).Error; err != nil {
		return nil, err
	}
	return &existingToken, nil
}

func (s *AuthenticationService) generateAuthToken(participant *models.Participant) models.AuthToken {
	token := uuid.NewString()
	expirationDateTime := time.Now().Add(1 * time.Hour)

	return models.AuthToken{
		Token:              token,
		ExpirationDateTime: expirationDateTime,
		ParticipantID:      participant.ID,
	}
}
