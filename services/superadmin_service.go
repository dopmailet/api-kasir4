package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type SuperadminService struct {
	storeRepo *repositories.StoreRepository
	pkgRepo   *repositories.SubscriptionPackageRepository
}

func NewSuperadminService(storeRepo *repositories.StoreRepository, pkgRepo *repositories.SubscriptionPackageRepository) *SuperadminService {
	return &SuperadminService{
		storeRepo: storeRepo,
		pkgRepo:   pkgRepo,
	}
}

// GetAllStores mengambil daftar pendaftar/toko SaaS
func (s *SuperadminService) GetAllStores() ([]models.Store, error) {
	return s.storeRepo.GetAll()
}

// UpdateStoreStatus (Banned/Unbanned) toko dari platform
func (s *SuperadminService) UpdateStoreStatus(id int, isActive bool) error {
	store, err := s.storeRepo.GetByID(id)
	if err != nil {
		return err
	}
	store.IsActive = isActive
	return s.storeRepo.Update(store)
}

// GetAllSubscriptionPackages melihat daftar paket SaaS
func (s *SuperadminService) GetAllSubscriptionPackages(publicOnly bool) ([]models.SubscriptionPackage, error) {
	return s.pkgRepo.GetAll(publicOnly)
}

// GetSubscriptionPackageByID
func (s *SuperadminService) GetSubscriptionPackageByID(id int) (*models.SubscriptionPackage, error) {
	return s.pkgRepo.GetByID(id)
}

// CreateSubscriptionPackage
// featuresJSON: string JSON valid untuk kolom JSONB, mis: "[]" atau '["Fitur A"]'
func (s *SuperadminService) CreateSubscriptionPackage(p *models.SubscriptionPackage, featuresJSON string) error {
	return s.pkgRepo.Create(p, featuresJSON)
}

// UpdateSubscriptionPackage
// featuresJSON: string JSON valid untuk kolom JSONB, mis: "[]" atau '["Fitur A"]'
func (s *SuperadminService) UpdateSubscriptionPackage(p *models.SubscriptionPackage, featuresJSON string) error {
	return s.pkgRepo.Update(p, featuresJSON)
}

// DeleteSubscriptionPackage
func (s *SuperadminService) DeleteSubscriptionPackage(id int) error {
	return s.pkgRepo.Delete(id)
}

// TogglePackagePopular mengaktifkan/mematikan flag paket populer secara atomik.
func (s *SuperadminService) TogglePackagePopular(id int, isPopular bool) error {
	return s.pkgRepo.TogglePopular(id, isPopular)
}

// UpdateStorePackage mengaktifkan/mengganti paket langganan sebuah toko secara manual
func (s *SuperadminService) UpdateStorePackage(storeID int, packageID int, durationDays int) error {
	endDate := time.Now().AddDate(0, 0, durationDays)
	return s.storeRepo.UpdatePackage(storeID, packageID, &endDate)
}

// DeleteStore menghapus satu toko berdasarkan ID
func (s *SuperadminService) DeleteStore(id int) error {
	if id == 1 {
		return fmt.Errorf("toko default (ID=1) tidak dapat dihapus")
	}
	return s.storeRepo.Delete(id)
}

// DeleteUnverifiedStores menghapus semua toko yang belum pernah berlangganan berbayar
func (s *SuperadminService) DeleteUnverifiedStores() (int64, error) {
	return s.storeRepo.DeleteUnverified()
}

// DeleteNeverSubscribedStores menghapus toko yang hanya pakai paket gratis
func (s *SuperadminService) DeleteNeverSubscribedStores() (int64, error) {
	return s.storeRepo.DeleteNeverSubscribed()
}
