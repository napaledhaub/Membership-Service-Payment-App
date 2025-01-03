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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	participant := authToken.Participant

	if participant.BillAmount != 0 && participant.BillAmount == expectedBillAmount {
		participant.IsBillVerified = true
		return s.DB.Save(&participant).Error
	}

	return errors.New("Bill verification failed")
}

func (s *PaymentService) SendEmailVerification(email string, participantId string) error {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	if participant.IsBillVerified {
		code := utils.GenerateVerificationCode()

		participant.PaymentOTP = code
		participant.PaymentOTPExpiration = time.Now().Add(15 * time.Minute)
		if err := s.DB.Save(participant).Error; err != nil {
			return err
		}

		subject := "OTP Code"
		text := "Your OTP code is: " + code
		if err := s.EmailService.SendEmail(email, subject, text); err != nil {
			return err
		}

		return nil
	}

	return errors.New("Bill verification is not valid")
}

func (s *PaymentService) VerifyPaymentOTP(email string, OTP string) error {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	if participant.PaymentOTP == OTP && time.Now().Before(participant.PaymentOTPExpiration) {
		participant.PaymentOTP = ""
		participant.PaymentOTPExpiration = time.Time{}
		participant.BillAmount = 0
		participant.IsBillVerified = false
		return s.DB.Save(&participant).Error
	}

	return errors.New("Payment verification failed")
}

func (s *PaymentService) IsPaymentOTPExpired(email string) (bool, error) {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("Participant not found")
		}
		return false, err
	}

	return time.Now().After(participant.PaymentOTPExpiration), nil
}
