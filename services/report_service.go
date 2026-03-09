package services

import (
	"fmt"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

// ReportService handles business logic for reports
type ReportService struct {
	repo *repositories.ReportRepository
}

// NewReportService creates a new ReportService
func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetDailySalesReport retrieves sales report for today (kept for backward compat)
func (s *ReportService) GetDailySalesReport(userID *int) (*models.SalesReport, error) {
	return s.repo.GetDailySalesReport(userID)
}

// GetSalesReportByDateRange retrieves sales report for a date range
// startDate dan endDate sudah mengandung timezone yang benar dari handler
func (s *ReportService) GetSalesReportByDateRange(startDate, endDate time.Time, userID *int) (*models.SalesReport, error) {
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start_date harus sebelum atau sama dengan end_date")
	}

	return s.repo.GetSalesReportByDateRange(startDate, endDate, userID)
}

// GetSalesTrend retrieves sales trend data based on period type
// Jika startDate & endDate diisi, akan digunakan langsung (custom range).
// Jika kosong (zero value), akan fallback ke preset berdasarkan periodType.
// loc = timezone dari user (dikirim dari handler)
// tzName = nama timezone string (contoh: "Asia/Makassar") untuk query SQL AT TIME ZONE
func (s *ReportService) GetSalesTrend(periodType string, loc *time.Location, tzName string, startDate, endDate time.Time) ([]models.SalesTrend, error) {
	now := time.Now().In(loc)
	var interval string

	// Jika start/end date tidak diisi, gunakan preset berdasarkan period
	if startDate.IsZero() || endDate.IsZero() {
		endDate = now
		switch periodType {
		case "monthly":
			startDate = now.AddDate(0, -11, 0)
			startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, loc)
		case "yearly":
			startDate = now.AddDate(-4, 0, 0)
			startDate = time.Date(startDate.Year(), 1, 1, 0, 0, 0, 0, loc)
		default: // "daily" — 7 hari terakhir
			startDate = now.AddDate(0, 0, -6)
			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
		}
	} else {
		// Pastikan startDate dari awal hari dan endDate sampai akhir hari
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, loc)
	}

	// Tentukan interval berdasarkan rentang hari
	// Jika rentang > 90 hari → pakai monthly, > 365 hari → yearly
	diff := endDate.Sub(startDate)
	switch {
	case periodType == "monthly" || diff.Hours() > 90*24:
		interval = "month"
	case periodType == "yearly" || diff.Hours() > 365*24:
		interval = "year"
	default:
		interval = "day"
	}

	return s.repo.GetSalesTrend(startDate, endDate, interval, tzName)
}

// GetDashboardSummary retrieves KPI summary for dashboard
// Mengembalikan ringkasan periode saat ini vs periode sebelumnya untuk perbandingan
func (s *ReportService) GetDashboardSummary(startDate, endDate time.Time, loc *time.Location) (*models.DashboardSummary, error) {
	if startDate.IsZero() || endDate.IsZero() {
		now := time.Now().In(loc)
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	}

	// Hitung durasi periode untuk membandingkan dengan periode sebelumnya
	duration := endDate.Sub(startDate)
	prevEndDate := startDate.Add(-time.Nanosecond)
	prevStartDate := prevEndDate.Add(-duration)

	current, err := s.repo.GetSalesReportByDateRange(startDate, endDate, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data periode saat ini: %w", err)
	}

	prev, err := s.repo.GetSalesReportByDateRange(prevStartDate, prevEndDate, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data periode sebelumnya: %w", err)
	}

	summary := &models.DashboardSummary{
		PeriodStart:       startDate,
		PeriodEnd:         endDate,
		Current:           *current,
		Previous:          *prev,
		RevenueGrowth:     calcGrowth(current.TotalRevenue, prev.TotalRevenue),
		ProfitGrowth:      calcGrowth(current.TotalProfit, prev.TotalProfit),
		TransactionGrowth: calcGrowth(float64(current.TotalTransaksi), float64(prev.TotalTransaksi)),
	}

	return summary, nil
}

// calcGrowth menghitung persentase perubahan dari nilai lama ke nilai baru
// Return: positif = naik, negatif = turun, 0 jika prev = 0
func calcGrowth(current, prev float64) float64 {
	if prev == 0 {
		if current > 0 {
			return 100 // naik 100% dari nol
		}
		return 0
	}
	return ((current - prev) / prev) * 100
}

// GetTopProducts retrieves top selling products
// Jika startDate/endDate diisi → gunakan rentang tersebut.
// Jika kosong (zero value) → fallback ke 30 hari terakhir.
// loc = timezone dari user
func (s *ReportService) GetTopProducts(limit int, loc *time.Location, startDate, endDate time.Time) ([]models.TopProduct, []models.TopProduct, error) {
	if limit <= 0 {
		limit = 5
	}

	if startDate.IsZero() || endDate.IsZero() {
		// Fallback: last 30 hari berdasarkan timezone user
		now := time.Now().In(loc)
		start := now.AddDate(0, 0, -30)
		startDate = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	}

	return s.repo.GetTopProducts(startDate, endDate, limit)
}

// CountLowStockProducts menghitung jumlah produk yang stoknya <= threshold
// threshold = batas minimum stok, default di handler adalah 5
func (s *ReportService) CountLowStockProducts(threshold int) (int, error) {
	if threshold < 0 {
		threshold = 0
	}
	return s.repo.CountLowStockProducts(threshold)
}

// GetDashboardAssets retrieves total asset cost and retail from products
func (s *ReportService) GetDashboardAssets() (*models.AssetReport, error) {
	return s.repo.GetDashboardAssets()
}
