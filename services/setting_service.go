package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type SettingService struct {
	repo *repositories.SettingRepository
}

func NewSettingService(repo *repositories.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetCustomerSettings(storeID int) (*models.AppSettings, error) {
	return s.repo.GetCustomerSettings(storeID)
}

func (s *SettingService) UpdateCustomerSettings(storeID int, stg *models.AppSettings) error {
	return s.repo.UpdateCustomerSettings(storeID, stg)
}

func (s *SettingService) GetPlatformSettings() (*models.PlatformSettings, error) {
	return s.repo.GetPlatformSettings()
}

func (s *SettingService) UpdatePlatformSettings(stg *models.PlatformSettings) error {
	return s.repo.UpdatePlatformSettings(stg)
}
