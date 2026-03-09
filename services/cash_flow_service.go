package services

import (
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

func (s *CashFlowService) GetSummary(startDate, endDate time.Time, loc *time.Location) (*models.CashFlowSummary, error) {
	if startDate.IsZero() || endDate.IsZero() {
		now := time.Now().In(loc)
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
	}

	return s.repo.GetSummary(startDate, endDate)
}

func (s *CashFlowService) GetTrend(startDate, endDate time.Time, loc *time.Location, tzName string) (*models.CashFlowTrendResponse, error) {
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

	return s.repo.GetTrend(startDate, endDate, format, tzName)
}
