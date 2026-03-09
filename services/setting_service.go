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

func (s *SettingService) GetCustomerSettings() (*models.AppSettings, error) {
	return s.repo.GetCustomerSettings()
}

func (s *SettingService) UpdateCustomerSettings(stg *models.AppSettings) error {
	return s.repo.UpdateCustomerSettings(stg)
}
