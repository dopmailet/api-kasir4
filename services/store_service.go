package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type StoreService struct {
	repo    *repositories.StoreRepository
	pkgRepo *repositories.SubscriptionPackageRepository
}

func NewStoreService(repo *repositories.StoreRepository, pkgRepo *repositories.SubscriptionPackageRepository) *StoreService {
	return &StoreService{
		repo:    repo,
		pkgRepo: pkgRepo,
	}
}

// GetMyStoreInfo mengambil profil detil toko yang sedang aktif, termasuk paket langganannya.
func (s *StoreService) GetMyStoreInfo(storeID int) (*models.Store, error) {
	return s.repo.GetByID(storeID)
}

// GetStoreLimits menghitung limit saat ini vs batas dari paket langganan
func (s *StoreService) GetStoreLimits(storeID int, timezone string) (*models.StoreLimitsInfo, error) {
	// 1. Dapatkan info toko
	store, err := s.repo.GetByID(storeID)
	if err != nil {
		return nil, err
	}

	// 2. Dapatkan detail paket
	pkg, err := s.pkgRepo.GetByID(store.SubscriptionPackageID)
	if err != nil {
		return nil, err
	}

	// 3. Validasi batasan saat ini
	cashiers, err := s.repo.CountActiveCashiers(storeID)
	if err != nil {
		return nil, err
	}

	products, err := s.repo.CountActiveProducts(storeID)
	if err != nil {
		return nil, err
	}

	sales, err := s.repo.CountTodayTransactions(storeID, timezone)
	if err != nil {
		return nil, err
	}

	return &models.StoreLimitsInfo{
		StoreID:         storeID,
		PackageName:     pkg.Name,
		CurrentCashiers: cashiers,
		MaxCashiers:     pkg.MaxKasir,
		CurrentProducts: products,
		MaxProducts:     pkg.MaxProducts,
		TodaySales:      sales,
		MaxDailySales:   pkg.MaxDailySales,
	}, nil
}
