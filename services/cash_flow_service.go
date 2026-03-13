package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type CashFlowService struct {
	repo *repositories.CashFlowRepository
}

func NewCashFlowService(repo *repositories.CashFlowRepository) *CashFlowService {
	return &CashFlowService{repo: repo}
}

func (s *CashFlowService) GetSummary(startDate, endDate time.Time, loc *time.Location, storeID int, tzName string) (*models.CashFlowSummary, error) {
	if startDate.IsZero() || endDate.IsZero() {
		now := time.Now().In(loc)
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
	}

	return s.repo.GetSummary(startDate, endDate, storeID, tzName)
}

func (s *CashFlowService) GetTrend(startDate, endDate time.Time, loc *time.Location, tzName string, storeID int) (*models.CashFlowTrendResponse, error) {
	if startDate.IsZero() || endDate.IsZero() {
		now := time.Now().In(loc)
		// Default ambil data 30 hari kebelakang
		start := now.AddDate(0, 0, -30)
		startDate = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	}

	diff := endDate.Sub(startDate)
	var format string

	// Jika jangka waktu <= 90 hari = tampil per hari
	// Jika jangka waktu > 90 hari = tampil per bulan
	if diff.Hours() <= 90*24 {
		format = "YYYY-MM-DD"
	} else {
		format = "YYYY-MM"
	}

	return s.repo.GetTrend(startDate, endDate, format, tzName, storeID)
}

func (s *CashFlowService) GetLedger(startDate, endDate time.Time, page, limit int, storeID int, tzName string) (*models.LedgerResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	return s.repo.GetLedger(startDate, endDate, page, limit, storeID, tzName)
}

func (s *CashFlowService) GetFunds(page, limit int, storeID int) (*models.CashFundsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.GetFunds(page, limit, storeID)
}

func (s *CashFlowService) GetInitialBalance(storeID int) (*models.CashInitialBalance, error) {
	return s.repo.GetInitialBalance(storeID)
}

func (s *CashFlowService) CreateFund(req *models.CashFundRequest, createdBy int, storeID int) (*models.CashFund, error) {
	// Validasi type
	if req.Type != "in" && req.Type != "out" {
		return nil, fmt.Errorf("type harus 'in' atau 'out'")
	}
	// Validasi amount
	if req.Amount <= 0 || req.Amount > 999999999999 {
		return nil, fmt.Errorf("amount harus lebih dari 0 dan maksimum 999.999.999.999")
	}
	// Validasi date
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("format date tidak valid, gunakan YYYY-MM-DD")
	}
	if d.After(time.Now()) {
		return nil, fmt.Errorf("date tidak boleh di masa depan")
	}
	// Validasi description
	if req.Description == "" {
		return nil, fmt.Errorf("description wajib diisi")
	}
	if len(req.Description) > 255 {
		return nil, fmt.Errorf("description maksimal 255 karakter")
	}
	return s.repo.CreateFund(req, createdBy, storeID)
}

func (s *CashFlowService) DeleteFund(id int, storeID int) error {
	return s.repo.DeleteFund(id, storeID)
}
