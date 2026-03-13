package services

import (
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

func (s *CustomerService) Create(req *models.CreateCustomerRequest, storeID int) (*models.Customer, error) {
	// Active default true jika tidak dikirim
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	customer := &models.Customer{
		Name:       req.Name,
		Phone:      req.Phone,
		CardNumber: req.CardNumber,
		Address:    req.Address,
		Notes:      req.Notes,
		IsActive:   isActive,
		StoreID:    storeID,
	}

	err := s.repo.Create(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetByID(id int, storeID int) (*models.Customer, error) {
	return s.repo.GetByID(id, storeID)
}

func (s *CustomerService) GetAll(search string, status string, page, limit int, sortBy, sortOrder string, storeID int) ([]models.Customer, int, error) {
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

	return s.repo.GetAll(search, status, page, limit, sortBy, sortOrder, storeID)
}

func (s *CustomerService) Update(id int, req *models.UpdateCustomerRequest, storeID int) (*models.Customer, error) {
	existing, err := s.repo.GetByID(id, storeID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.CardNumber != nil {
		existing.CardNumber = req.CardNumber
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

func (s *CustomerService) GetTransactions(id int, storeID int) ([]models.TransactionWithItems, error) {
	// Validasi eksistensi user
	_, err := s.repo.GetByID(id, storeID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTransactions(id, storeID)
}

func (s *CustomerService) GetLoyaltyHistory(id int, storeID int) ([]models.LoyaltyTransaction, error) {
	// Validasi eksistensi user
	_, err := s.repo.GetByID(id, storeID)
	if err != nil {
		return nil, err
	}
	return s.loyaltyRepo.GetHistoryByCustomerID(id, storeID)
}
