package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ReportHandler handles HTTP requests for reports
// Handler untuk report/laporan
type ReportHandler struct {
	service *services.ReportService
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// parseTimezone parses timezone from query parameter, defaults to Asia/Jakarta.
// Return dua nilai: *time.Location (untuk Go time ops) dan string nama timezone
// (untuk diteruskan ke query SQL PostgreSQL via AT TIME ZONE).
func parseTimezone(r *http.Request) (*time.Location, string) {
	tzStr := r.URL.Query().Get("timezone")
	if tzStr == "" {
		tzStr = "Asia/Jakarta" // Default WIB
	}

	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		log.Printf("⚠️ Timezone '%s' tidak valid, menggunakan Asia/Jakarta: %v", tzStr, err)
		tzStr = "Asia/Jakarta"
		loc, _ = time.LoadLocation(tzStr)
		if loc == nil {
			loc = time.FixedZone("WIB", 7*60*60)
			tzStr = "Asia/Jakarta"
		}
	}
	return loc, tzStr
}

// GetDailySalesReport handles GET /api/report/hari-ini?timezone=Asia/Jakarta
// Fungsi ini handle request untuk laporan penjualan hari ini
func (h *ReportHandler) GetDailySalesReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse timezone dari query parameter (default: Asia/Jakarta)
	loc, _ := parseTimezone(r)

	// Parsing user_id optional
	var userID *int
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = &id
		}
	}

	// Hitung "hari ini" berdasarkan timezone user
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)

	// Panggil service dengan date range yang sudah di-timezone
	report, err := h.service.GetSalesReportByDateRange(startOfDay, endOfDay, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetSalesReportByDateRange handles GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&timezone=Asia/Makassar
// Fungsi ini handle request untuk laporan penjualan berdasarkan rentang tanggal
func (h *ReportHandler) GetSalesReportByDateRange(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date dan end_date harus diisi (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Parse timezone (default: Asia/Jakarta)
	loc, _ := parseTimezone(r)

	// Parsing user_id optional
	var userID *int
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = &id
		}
	}

	// Parse tanggal
	startDateParsed, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Format start_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	endDateParsed, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Format end_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Buat boundary waktu berdasarkan timezone user
	// Contoh: 2026-02-17 di Asia/Makassar (UTC+8)
	// → startOfDay = 2026-02-17 00:00:00 WITA = 2026-02-16 16:00:00 UTC
	// → endOfDay   = 2026-02-17 23:59:59 WITA = 2026-02-17 15:59:59 UTC
	startDate := time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)
	endDate := time.Date(endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(), 23, 59, 59, 999999999, loc)

	report, err := h.service.GetSalesReportByDateRange(startDate, endDate, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetSalesTrend handles GET /api/dashboard/sales-trend
// Query params:
//
//	period=day|month|year  (default: day)
//	start_date=YYYY-MM-DD  (opsional, jika diisi maka end_date wajib diisi juga)
//	end_date=YYYY-MM-DD    (opsional)
//	timezone=Asia/Jakarta  (default: Asia/Jakarta)
func (h *ReportHandler) GetSalesTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "day"
	}

	// Parse timezone
	loc, tzName := parseTimezone(r)

	// Parse start_date & end_date (opsional)
	var startDate, endDate time.Time
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr != "" && endDateStr != "" {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, loc)
		if err != nil {
			http.Error(w, "Format start_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, loc)
		if err != nil {
			http.Error(w, "Format end_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		if startDate.After(endDate) {
			http.Error(w, "start_date harus sebelum atau sama dengan end_date", http.StatusBadRequest)
			return
		}
	}

	trends, err := h.service.GetSalesTrend(period, loc, tzName, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"period":     period,
		"start_date": startDateStr,
		"end_date":   endDateStr,
		"data":       trends,
	})
}

// GetTopProducts handles GET /api/dashboard/top-products
// Query params:
//
//	limit=N                    (default: 5)
//	start_date=YYYY-MM-DD      (opsional, default: 30 hari terakhir)
//	end_date=YYYY-MM-DD        (opsional)
//	timezone=Asia/Jakarta      (default: Asia/Jakarta)
func (h *ReportHandler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Parse timezone
	loc, _ := parseTimezone(r)

	// Parse start_date & end_date (opsional)
	var startDate, endDate time.Time
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr != "" && endDateStr != "" {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, loc)
		if err != nil {
			http.Error(w, "Format start_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, loc)
		if err != nil {
			http.Error(w, "Format end_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}

	topQty, topProfit, err := h.service.GetTopProducts(limit, loc, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"by_quantity": topQty,
		"by_profit":   topProfit,
	})
}

// GetDashboardSummary handles GET /api/dashboard/summary
// Mengembalikan KPI cards: omzet, profit, transaksi, items, pengeluaran, laba bersih
// beserta % pertumbuhan vs periode sebelumnya dan jumlah produk stok menipis.
// Query params:
//
//	start_date=YYYY-MM-DD  (opsional, default: hari ini)
//	end_date=YYYY-MM-DD    (opsional, default: hari ini)
//	low_stock_threshold=N  (opsional, default: 5)
//	timezone=Asia/Jakarta  (default: Asia/Jakarta)
func (h *ReportHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	loc, _ := parseTimezone(r)

	// Parse start_date & end_date (opsional, default: hari ini)
	var startDate, endDate time.Time
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr != "" && endDateStr != "" {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, loc)
		if err != nil {
			http.Error(w, "Format start_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, loc)
		if err != nil {
			http.Error(w, "Format end_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		if startDate.After(endDate) {
			http.Error(w, "start_date harus sebelum atau sama dengan end_date", http.StatusBadRequest)
			return
		}
	}

	// Parse low_stock_threshold (default: 5)
	lowStockThreshold := 5
	if thStr := r.URL.Query().Get("low_stock_threshold"); thStr != "" {
		if th, err := strconv.Atoi(thStr); err == nil && th >= 0 {
			lowStockThreshold = th
		}
	}

	summary, err := h.service.GetDashboardSummary(startDate, endDate, loc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hitung jumlah produk stok menipis
	lowStockCount, err := h.service.CountLowStockProducts(lowStockThreshold)
	if err != nil {
		log.Printf("⚠️ Gagal hitung low stock: %v", err)
		// Tidak fatal, tetap return data lainnya
	}
	summary.LowStockCount = lowStockCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetDashboardAssets handles GET /api/dashboard/assets
// Endpoint tidak membutuhkan query parameter, menampilkan total cost HPP dan potensi retail
func (h *ReportHandler) GetDashboardAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.service.GetDashboardAssets()
	if err != nil {
		log.Printf("Error get dashboard assets: %v", err)
		http.Error(w, "Gagal mengambil data aset", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
