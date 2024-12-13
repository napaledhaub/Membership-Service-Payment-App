package services

import (
	"errors"
	"paymentapp/models"
	"time"

	"gorm.io/gorm"
)

type SubscriptionService struct {
	DB *gorm.DB
}

func (s *SubscriptionService) GetSubscriptionList(token string) ([]models.Subscription, error) {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return nil, errors.New("participant not found")
	}

	var subscriptions []models.Subscription
	if err := s.DB.Where("participant_id = ?", authToken.Participant.ID).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *SubscriptionService) SubscribeToService(token string, serviceMenuID uint) (*models.Subscription, error) {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return nil, errors.New("participant not found")
	}

	var serviceMenu models.ServiceMenu
	if err := s.DB.First(&serviceMenu, serviceMenuID).Error; err != nil {
		return nil, errors.New("service menu not found")
	}

	var subscription models.Subscription
	if err := s.DB.Where("participant_id = ? AND service_menu_id = ?", authToken.Participant.ID, serviceMenu.ID).First(&subscription).Error; err != nil {
		subscription = models.Subscription{}
	}

	subscription.ParticipantID = authToken.Participant.ID
	subscription.ServiceMenuID = serviceMenu.ID
	subscription.StartDate = time.Now()
	subscription.EndDate = time.Now().Add(30 * 24 * time.Hour)
	subscription.RemainingSessions = serviceMenu.TotalSessions

	authToken.Participant.BillAmount = serviceMenu.PricePerSession * float64(serviceMenu.TotalSessions)
	authToken.Participant.IsBillVerified = false

	if err := s.DB.Save(&authToken.Participant).Error; err != nil {
		return nil, err
	}

	if err := s.DB.Save(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *SubscriptionService) CancelSubscription(token string, subscriptionId uint) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return errors.New("participant not found")
	}

	var subscription models.Subscription
	if err := s.DB.Preload("ServiceMenu").First(&subscription, subscriptionId).Error; err != nil {
		return errors.New("subscription not found")
	}

	participant := authToken.Participant
	serviceMenu := subscription.ServiceMenu

	currentBillAmount := participant.BillAmount
	pricePerSession := serviceMenu.PricePerSession
	remainingSessions := subscription.RemainingSessions
	deductionAmount := pricePerSession * float64(remainingSessions)
	newBillAmount := currentBillAmount - deductionAmount
	if newBillAmount < 0 {
		newBillAmount = 0
	}

	participant.BillAmount = newBillAmount
	if err := s.DB.Save(&participant).Error; err != nil {
		return err
	}

	return s.DB.Delete(&subscription).Error
}

func (s *SubscriptionService) ExtendSubscription(token string, subscriptionId uint) error {
	var authToken models.AuthToken
	if err := s.DB.Where("token = ?", token).Preload("Participant").First(&authToken).Error; err != nil {
		return errors.New("participant not found")
	}

	var subscription models.Subscription
	if err := s.DB.Preload("ServiceMenu").First(&subscription, subscriptionId).Error; err != nil {
		return errors.New("subscription not found")
	}

	serviceMenu := subscription.ServiceMenu
	participant := authToken.Participant

	subscription.EndDate = subscription.EndDate.AddDate(0, 0, serviceMenu.TotalSessions)
	subscription.RemainingSessions += serviceMenu.TotalSessions

	if err := s.DB.Save(&subscription).Error; err != nil {
		return err
	}

	pricePerSession := serviceMenu.PricePerSession
	totalSessions := serviceMenu.TotalSessions
	currentBillAmount := participant.BillAmount
	newBillAmount := (pricePerSession * float64(totalSessions)) + currentBillAmount

	participant.BillAmount = newBillAmount
	participant.IsBillVerified = false

	return s.DB.Save(&participant).Error
}
