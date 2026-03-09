package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
)

type SupplierService struct {
	repo *repositories.SupplierRepository
}

func NewSupplierService(repo *repositories.SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) GetAll(search string, isActive *bool) ([]models.Supplier, error) {
	return s.repo.GetAll(search, isActive)
}

func (s *SupplierService) GetByID(id int) (*models.Supplier, error) {
	return s.repo.GetByID(id)
}

func (s *SupplierService) Create(req *models.CreateSupplierRequest) (*models.Supplier, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("nama supplier wajib diisi")
	}
	return s.repo.Create(req)
}

func (s *SupplierService) Update(id int, req *models.UpdateSupplierRequest) (*models.Supplier, error) {
	return s.repo.Update(id, req)
}

func (s *SupplierService) Delete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *SupplierService) GetPayables(supplierID int) ([]models.SupplierPayable, error) {
	_, err := s.repo.GetByID(supplierID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetPayablesBySupplier(supplierID)
}

func (s *SupplierService) CreatePayable(supplierID int, req *models.CreatePayableRequest) (*models.SupplierPayable, error) {
	_, err := s.repo.GetByID(supplierID)
	if err != nil {
		return nil, err
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount harus lebih dari 0")
	}
	return s.repo.CreatePayable(supplierID, req)
}

func (s *SupplierService) UpdatePayable(supplierID, payableID int, req *models.UpdatePayableRequest) (*models.SupplierPayable, error) {
	payable, err := s.repo.GetPayableByID(payableID)
	if err != nil {
		return nil, err
	}
	if payable.SupplierID != supplierID {
		return nil, fmt.Errorf("payable id %d bukan milik supplier id %d", payableID, supplierID)
	}
	return s.repo.UpdatePayable(payableID, req)
}

func (s *SupplierService) GetPayments(payableID int) ([]models.PayablePayment, error) {
	_, err := s.repo.GetPayableByID(payableID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetPaymentsByPayable(payableID)
}

func (s *SupplierService) CreatePayment(payableID int, req *models.CreatePaymentRequest) (*models.PayablePayment, error) {
	return s.repo.CreatePayment(payableID, req)
}
