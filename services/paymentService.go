package services

import (
	"errors"
	"paymentapp/models"
	"paymentapp/utils"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
	DB           *gorm.DB
	EmailService *utils.EmailService
}

func (s *PaymentService) VerifyBillAmount(token string, expectedBillAmount float64) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return errors.New("participant not found")
	}

	participant := authToken.Participant

	if participant.BillAmount != 0 && participant.BillAmount == expectedBillAmount {
		participant.IsBillVerified = true
		return s.DB.Save(&participant).Error
	}

	return errors.New("bill verification failed")
}

func (s *PaymentService) SendEmailVerification(email string, participantId string) error {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return errors.New("participant not found")
	}

	if participant.IsBillVerified {
		code := utils.GenerateVerificationCode()

		participant.PaymentOTP = code
		participant.PaymentOTPExpiration = time.Now().Add(15 * time.Minute)
		if err := s.DB.Save(participant).Error; err != nil {
			return err
		}

		/*subject := "OTP Code"
		text := "Your OTP code is: " + code
		if err := s.EmailService.SendEmail(email, subject, text); err != nil {
			return err
		}*/

		return nil
	}

	return errors.New("Bill verification is not valid")
}

func (s *PaymentService) VerifyPaymentOTP(email string, OTP string) error {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return errors.New("participant not found")
	}

	if participant.PaymentOTP == OTP && time.Now().Before(participant.PaymentOTPExpiration) {
		participant.PaymentOTP = ""
		participant.PaymentOTPExpiration = time.Time{}
		participant.BillAmount = 0
		participant.IsBillVerified = false
		return s.DB.Save(&participant).Error
	}

	return errors.New("payment verification failed")
}

func (s *PaymentService) IsPaymentOTPExpired(email string) bool {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return false
	}
	return time.Now().After(participant.PaymentOTPExpiration)
}
