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

func (s *SupplierService) GetAll(search string, isActive *bool, storeID int) ([]models.Supplier, error) {
	return s.repo.GetAll(search, isActive, storeID)
}

func (s *SupplierService) GetDebtSummary(storeID int) (*models.SupplierDebtSummary, error) {
	return s.repo.GetDebtSummary(storeID)
}

func (s *SupplierService) GetByID(id int, storeID int) (*models.Supplier, error) {
	return s.repo.GetByID(id, storeID)
}

func (s *SupplierService) Create(req *models.CreateSupplierRequest, storeID int) (*models.Supplier, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("nama supplier wajib diisi")
	}
	return s.repo.Create(req, storeID)
}

func (s *SupplierService) Update(id int, req *models.UpdateSupplierRequest, storeID int) (*models.Supplier, error) {
	return s.repo.Update(id, req, storeID)
}

func (s *SupplierService) Delete(id int, storeID int) error {
	_, err := s.repo.GetByID(id, storeID)
	if err != nil {
		return err
	}
	return s.repo.Delete(id, storeID)
}

func (s *SupplierService) GetPayables(supplierID int, storeID int) ([]models.SupplierPayable, error) {
	_, err := s.repo.GetByID(supplierID, storeID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetPayablesBySupplier(supplierID, storeID)
}

func (s *SupplierService) CreatePayable(supplierID int, req *models.CreatePayableRequest, storeID int) (*models.SupplierPayable, error) {
	_, err := s.repo.GetByID(supplierID, storeID)
	if err != nil {
		return nil, err
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount harus lebih dari 0")
	}
	return s.repo.CreatePayable(supplierID, req)
}

func (s *SupplierService) UpdatePayable(supplierID, payableID int, req *models.UpdatePayableRequest, storeID int) (*models.SupplierPayable, error) {
	payable, err := s.repo.GetPayableByID(payableID, storeID)
	if err != nil {
		return nil, err
	}
	if payable.SupplierID != supplierID {
		return nil, fmt.Errorf("payable id %d bukan milik supplier id %d", payableID, supplierID)
	}
	return s.repo.UpdatePayable(payableID, req, storeID)
}

func (s *SupplierService) GetPayments(payableID int, storeID int) ([]models.PayablePayment, error) {
	_, err := s.repo.GetPayableByID(payableID, storeID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetPaymentsByPayable(payableID, storeID)
}

func (s *SupplierService) CreatePayment(payableID int, req *models.CreatePaymentRequest, storeID int) (*models.PayablePayment, error) {
	return s.repo.CreatePayment(payableID, req, storeID)
}

func (s *SupplierService) GetAllPayablePayments(startDate, endDate string, storeID int) ([]models.PayablePaymentWithSupplier, error) {
	return s.repo.GetAllPayablePayments(startDate, endDate, storeID)
}
