package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type PayrollService struct {
	repo *repositories.PayrollRepository
}

func NewPayrollService(repo *repositories.PayrollRepository) *PayrollService {
	return &PayrollService{repo: repo}
}

func (s *PayrollService) GetAll(employeeID int, startDate, endDate time.Time, page, limit int) ([]models.Payroll, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return s.repo.GetAll(employeeID, startDate, endDate, offset, limit)
}

func (s *PayrollService) GetByID(id int) (*models.Payroll, error) {
	return s.repo.GetByID(id)
}

func (s *PayrollService) Create(req models.CreatePayrollRequest, createdBy int) (*models.Payroll, error) {
	bonus := 0.0
	if req.Bonus != nil {
		bonus = *req.Bonus
	}
	potongan := 0.0
	if req.Potongan != nil {
		potongan = *req.Potongan
	}

	// Auto calculate Total
	total := req.GajiPokok + bonus - potongan

	p := &models.Payroll{
		EmployeeID: req.EmployeeID,
		Periode:    req.Periode,
		GajiPokok:  req.GajiPokok,
		Bonus:      bonus,
		Potongan:   potongan,
		Total:      total,
		Catatan:    req.Catatan,
		CreatedBy:  &createdBy,
	}

	err := s.repo.Create(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PayrollService) Update(id int, req models.UpdatePayrollRequest) (*models.Payroll, error) {
	// Ambil data exist
	p, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("data payroll tidak ditemukan")
	}

	// Pengecekan expired (Edit Limit: 24 Jam)
	// Kita asumsikan acuan waktu adalah paid_at (atau created_at, defaultnya sama)
	elapsed := time.Since(p.PaidAt)
	if elapsed > 24*time.Hour {
		return nil, errors.New("data payroll hanya bisa diedit dalam waktu 24 jam setelah dibuat")
	}

	// Terapkan perubahan dan hitung ulang Total
	if req.Periode != nil {
		p.Periode = req.Periode
	}
	if req.GajiPokok != nil {
		p.GajiPokok = *req.GajiPokok
	}
	if req.Bonus != nil {
		p.Bonus = *req.Bonus
	}
	if req.Potongan != nil {
		p.Potongan = *req.Potongan
	}
	if req.Catatan != nil {
		p.Catatan = req.Catatan
	}

	p.Total = p.GajiPokok + p.Bonus - p.Potongan

	err = s.repo.Update(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PayrollService) Delete(id int) error {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("data payroll tidak ditemukan")
	}

	// Pengecekan expired (Delete Limit: 24 Jam)
	elapsed := time.Since(p.PaidAt)
	if elapsed > 24*time.Hour {
		return errors.New("data payroll hanya bisa dihapus dalam waktu 24 jam setelah dibuat")
	}

	return s.repo.Delete(id)
}

func (s *PayrollService) GetReport(startDate, endDate time.Time, tzName string) (*models.PayrollReport, error) {
	return s.repo.GetReport(startDate, endDate, tzName)
}
