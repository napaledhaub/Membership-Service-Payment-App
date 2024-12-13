package services

import (
	"errors"
	"paymentapp/models"
	"paymentapp/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegistrationService struct {
	DB            *gorm.DB
	EmailService  *utils.EmailService
	EncryptionKey []byte
}

func (s *RegistrationService) Register(request models.Participant) error {
	var existingParticipant models.Participant
	if err := s.DB.Where("email = ?", request.Email).First(&existingParticipant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			verificationCode := utils.GenerateVerificationCode()
			request.VerificationCode = verificationCode
			request.IsVerified = false

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			request.Password = string(hashedPassword)

			creditCardInfo, err := utils.EncryptCreditCardInfo(request.CreditCardInfo, s.EncryptionKey)
			if err != nil {
				return err
			}
			request.CreditCardInfo = creditCardInfo

			if err := s.DB.Create(&request).Error; err != nil {
				return err
			}

			subject := "Registration Confirmation"
			text := "Thank you for registering at our fitness center. Here is your OTP to activate your account: " + verificationCode
			if err := s.EmailService.SendEmail(request.Email, subject, text); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	return errors.New("Email is already in use")
}
