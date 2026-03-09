package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type ExpenseService struct {
	repo *repositories.ExpenseRepository
}

func NewExpenseService(repo *repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// GetAll mengambil semua pengeluaran
func (s *ExpenseService) GetAll(year string, month string) ([]models.Expense, error) {
	return s.repo.GetAll(year, month)
}

// GetByID mengambil satu data pengeluaran
func (s *ExpenseService) GetByID(id int) (*models.Expense, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return s.repo.GetByID(id)
}

// Create menambahkan pengeluaran baru
func (s *ExpenseService) Create(req models.CreateExpenseRequest, createdBy int) (*models.Expense, error) {
	if req.Category == "" || req.Description == "" || req.Amount <= 0 || req.ExpenseDate == "" {
		return nil, errors.New("category, description, expense_date, dan amount (harus > 0) wajib diisi")
	}

	expense := &models.Expense{
		Category:        req.Category,
		Description:     req.Description,
		Amount:          req.Amount,
		ExpenseDate:     req.ExpenseDate,
		IsRecurring:     req.IsRecurring,
		RecurringPeriod: req.RecurringPeriod,
		Notes:           req.Notes,
		CreatedBy:       createdBy,
	}

	return s.repo.Create(expense)
}

// Update memodifikasi data pengeluaran
func (s *ExpenseService) Update(id int, req models.UpdateExpenseRequest) (*models.Expense, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}

	// Ambil data lama
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("pengeluaran tidak ditemukan")
	}

	// Terapkan perubahan parsil
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, errors.New("amount pengeluaran tidak valid")
		}
		existing.Amount = *req.Amount
	}
	if req.ExpenseDate != nil {
		// Asumsi frontend selalu mengirim "YYYY-MM-DD"
		existing.ExpenseDate = *req.ExpenseDate
	}
	if req.IsRecurring != nil {
		existing.IsRecurring = *req.IsRecurring
	}
	if req.RecurringPeriod != nil {
		existing.RecurringPeriod = req.RecurringPeriod
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	return s.repo.Update(id, existing)
}

// Delete menghapus pengeluaran (Hard delete)
func (s *ExpenseService) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.Delete(id)
}
