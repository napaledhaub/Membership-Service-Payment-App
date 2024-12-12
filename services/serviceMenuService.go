package services

import (
	"paymentapp/models"

	"gorm.io/gorm"
)

type ServiceMenuService struct {
	DB *gorm.DB
}

func (s *ServiceMenuService) GetAllServiceMenus() ([]models.ServiceMenu, error) {
	var serviceMenus []models.ServiceMenu
	if err := s.DB.Preload("Exercises").Find(&serviceMenus).Error; err != nil {
		return nil, err
	}
	return serviceMenus, nil
}

func (s *ServiceMenuService) AddServiceMenu(serviceMenu *models.ServiceMenu) error {
	return s.DB.Create(serviceMenu).Error
}
