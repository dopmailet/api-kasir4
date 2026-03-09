package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
)

type CustomerService struct {
	repo        *repositories.CustomerRepository
	loyaltyRepo *repositories.LoyaltyRepository
}

func NewCustomerService(repo *repositories.CustomerRepository, loyaltyRepo *repositories.LoyaltyRepository) *CustomerService {
	return &CustomerService{
		repo:        repo,
		loyaltyRepo: loyaltyRepo,
	}
}

func (s *CustomerService) Create(req *models.CreateCustomerRequest) (*models.Customer, error) {
	code, err := s.repo.GenerateCustomerCode()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat customer code: %w", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	customer := &models.Customer{
		CustomerCode: code,
		Name:         req.Name,
		Phone:        req.Phone,
		Address:      req.Address,
		Notes:        req.Notes,
		IsActive:     isActive,
	}

	err = s.repo.Create(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetByID(id int) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) GetAll(search string, status string, page, limit int, sortBy, sortOrder string) ([]models.Customer, int, error) {
	// Defaults
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetAll(search, status, page, limit, sortBy, sortOrder)
}

func (s *CustomerService) Update(id int, req *models.UpdateCustomerRequest) (*models.Customer, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.Address != nil {
		existing.Address = req.Address
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	err = s.repo.Update(existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *CustomerService) GetTransactions(id int) ([]models.TransactionWithItems, error) {
	// Validasi eksistensi user
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTransactions(id)
}

func (s *CustomerService) GetLoyaltyHistory(id int) ([]models.LoyaltyTransaction, error) {
	// Validasi eksistensi user
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.loyaltyRepo.GetHistoryByCustomerID(id)
}
