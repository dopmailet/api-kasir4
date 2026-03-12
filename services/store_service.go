package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type StoreService struct {
	repo *repositories.StoreRepository
}

func NewStoreService(repo *repositories.StoreRepository) *StoreService {
	return &StoreService{repo: repo}
}

// GetMyStoreInfo mengambil profil detil toko yang sedang aktif, termasuk paket langganannya.
func (s *StoreService) GetMyStoreInfo(storeID int) (*models.Store, error) {
	return s.repo.GetByID(storeID)
}
