package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	repo *repositories.TransactionRepository
}

// NewTransactionService creates a new TransactionService
func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
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

	return s.repo.CreateTransaction(req)
}

// GetAll returns all transactions
func (s *TransactionService) GetAll(userID *int) ([]models.Transaction, error) {
	return s.repo.GetAll(userID)
}

// GetByDateRange returns transactions within a date range
func (s *TransactionService) GetByDateRange(startDate, endDate time.Time, userID *int) ([]models.Transaction, error) {
	return s.repo.GetByDateRange(startDate, endDate, userID)
}

// GetByID returns a transaction with items detail
func (s *TransactionService) GetByID(id int) (*models.TransactionWithItems, error) {
	return s.repo.GetByID(id)
}
