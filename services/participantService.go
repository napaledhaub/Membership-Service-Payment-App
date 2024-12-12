package services

import (
	"errors"
	"paymentapp/models"

	"gorm.io/gorm"
)

type ParticipantService struct {
	DB *gorm.DB
}

func (s *ParticipantService) FindByEmail(email string) (*models.Participant, error) {
	var participant models.Participant
	if err := s.DB.Where("email = ?", email).First(&participant).Error; err != nil {
		return nil, errors.New("participant not found")
	}
	return &participant, nil
}
