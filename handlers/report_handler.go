package handlers

import (
	"encoding/json"
	"kasir-api/middleware"
	"kasir-api/services"
	"kasir-api/utils"
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

// GetDailySalesReport handles GET /api/report/hari-ini?timezone=Asia/Jakarta
// Fungsi ini handle request untuk laporan penjualan hari ini
func (h *ReportHandler) GetDailySalesReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	// Parse timezone dari query parameter (default: Asia/Makassar via utils)
	tzName := utils.GetTimezone(r.URL.Query().Get("timezone"))

	// Parsing user_id optional
	var userID *int
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = &id
		}
	}

	report, err := h.service.GetDailySalesReport(userID, user.StoreID, tzName)
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

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	// Ambil query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "start_date dan end_date harus diisi (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Parse timezone (default: Asia/Makassar via utils)
	tzName := utils.GetTimezone(r.URL.Query().Get("timezone"))
	loc, _ := time.LoadLocation(tzName)

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
	startDate := time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)
	endDate := time.Date(endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(), 23, 59, 59, 999999999, loc)

	report, err := h.service.GetSalesReportByDateRange(startDate, endDate, userID, user.StoreID, tzName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetSalesTrend handles GET /api/dashboard/sales-trend
func (h *ReportHandler) GetSalesTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "day"
	}

	// Parse timezone
	tzName := utils.GetTimezone(r.URL.Query().Get("timezone"))
	loc, _ := time.LoadLocation(tzName)

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

	trends, err := h.service.GetSalesTrend(period, loc, tzName, startDate, endDate, user.StoreID)
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
func (h *ReportHandler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Parse timezone
	tzName := utils.GetTimezone(r.URL.Query().Get("timezone"))
	loc, _ := time.LoadLocation(tzName)

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

	topQty, topProfit, err := h.service.GetTopProducts(limit, loc, startDate, endDate, user.StoreID, tzName)
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
func (h *ReportHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	tzName := utils.GetTimezone(r.URL.Query().Get("timezone"))
	loc, _ := time.LoadLocation(tzName)

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

	summary, err := h.service.GetDashboardSummary(startDate, endDate, loc, user.StoreID, tzName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hitung jumlah produk stok menipis
	lowStockCount, err := h.service.CountLowStockProducts(lowStockThreshold, user.StoreID)
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

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	report, err := h.service.GetDashboardAssets(user.StoreID)
	if err != nil {
		log.Printf("Error get dashboard assets: %v", err)
		http.Error(w, "Gagal mengambil data aset", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
