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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Participant not found")
		}
		return nil, err
	}
	return &participant, nil
}

func (s *ParticipantService) UpdateFullName(token string, newFullName string) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	creditCardJSON, err := json.Marshal(newCreditCardInfo)
	if err != nil {
		return err
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

func (service *ParticipantService) UpdatePassword(token string, oldPassword string, newPassword string) error {
	var authToken models.AuthToken
	if err := service.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Participant not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authToken.Participant.Password), []byte(oldPassword)); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	authToken.Participant.Password = string(hashedPassword)

	if err := service.DB.Save(&authToken.Participant).Error; err != nil {
		return err
	}

	return nil
}
