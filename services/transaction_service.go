package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/utils"
	"time"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	repo     *repositories.TransactionRepository
	storeSvc *StoreService
}

// NewTransactionService creates a new TransactionService
func NewTransactionService(repo *repositories.TransactionRepository, storeSvc *StoreService) *TransactionService {
	return &TransactionService{
		repo:     repo,
		storeSvc: storeSvc,
	}
}

// Checkout processes a checkout request
func (s *TransactionService) Checkout(req *models.CheckoutRequest) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items tidak boleh kosong")
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return nil, fmt.Errorf("product_id tidak valid")
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity harus lebih dari 0")
		}
	}

	// Cek limit paket — gunakan timezone dari request, fallback ke default
	timezone := utils.GetTimezone(req.Timezone)
	limits, err := s.storeSvc.GetStoreLimits(req.StoreID, timezone)
	if err != nil {
		return nil, err
	}
	if limits.MaxDailySales != nil && limits.TodaySales >= *limits.MaxDailySales {
		return nil, fmt.Errorf("limit paket tercapai: maksimal %d transaksi per hari untuk paket %s", *limits.MaxDailySales, limits.PackageName)
	}

	return s.repo.CreateTransaction(req)
}

// GetAll returns all transactions
func (s *TransactionService) GetAll(userID *int, storeID int) ([]models.Transaction, error) {
	return s.repo.GetAll(userID, storeID)
}

// GetByDateRange returns transactions within a date range
func (s *TransactionService) GetByDateRange(startDate, endDate time.Time, userID *int, storeID int, tzName string) ([]models.Transaction, error) {
	return s.repo.GetByDateRange(startDate, endDate, userID, storeID, tzName)
}

// GetByID returns a transaction with items detail
func (s *TransactionService) GetByID(id int, storeID int) (*models.TransactionWithItems, error) {
	return s.repo.GetByID(id, storeID)
}
