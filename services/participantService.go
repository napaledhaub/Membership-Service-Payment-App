package services

import (
	"encoding/json"
	"errors"
	"paymentapp/models"
	"paymentapp/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ParticipantService struct {
	DB            *gorm.DB
	EncryptionKey []byte
}

func (s *ParticipantService) FindByEmail(email string) (*models.Participant, error) {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return nil, errors.New("participant not found")
	}
	return &participant, nil
}

func (s *ParticipantService) UpdateFullName(token string, newFullName string) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return errors.New("participant not found")
	}

	authToken.Participant.Name = newFullName
	if err := s.DB.Save(&authToken.Participant).Error; err != nil {
		return err
	}

	return nil
}

func (s *ParticipantService) UpdateCreditCardInfo(token string, newCreditCardInfo models.CreditCard) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return errors.New("participant not found")
	}

	creditCardJSON, err := json.Marshal(newCreditCardInfo)
	if err != nil {
		return errors.New("failed to marshal credit card info to JSON")
	}
	encryptedCreditCardInfo, err := utils.EncryptCreditCardInfo(string(creditCardJSON), s.EncryptionKey)
	if err != nil {
		return err
	}

	authToken.Participant.CreditCardInfo = encryptedCreditCardInfo

	if err := s.DB.Save(&authToken.Participant).Error; err != nil {
		return err
	}

	return nil
}

func (s *ParticipantService) UpdatePassword(token string, oldPassword string, newPassword string) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return errors.New("participant not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authToken.Participant.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password does not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	authToken.Participant.Password = string(hashedPassword)

	if err := s.DB.Save(&authToken.Participant).Error; err != nil {
		return err
	}

	return nil
}
