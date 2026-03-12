package handlers

import (
	"encoding/json"
	"kasir-api/middleware"
	"kasir-api/models"
	"kasir-api/services"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CashFlowHandler struct {
	service *services.CashFlowService
}

func NewCashFlowHandler(service *services.CashFlowService) *CashFlowHandler {
	return &CashFlowHandler{service: service}
}

// Reuse parseTimezone dari report_handler jika memungkinkan,
// namun untuk standalone module agar decouple kita buat instance utility sendiri disini.
func parseTimezoneCF(r *http.Request) (*time.Location, string) {
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

// GetSummary handles GET /api/cash-flow/summary
func (h *CashFlowHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	loc, _ := parseTimezoneCF(r)
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startStr != "" && endStr != "" {
		startDateParsed, errParse1 := time.Parse("2006-01-02", startStr)
		endDateParsed, errParse2 := time.Parse("2006-01-02", endStr)
		if errParse1 != nil || errParse2 != nil {
			http.Error(w, "Format tanggal tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		// Set batas awal jam 00:00:00 dan batas akhir 23:59:59 (lokal)
		startDate = time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(), 23, 59, 59, 999999999, loc)
	}

	summary, err := h.service.GetSummary(startDate, endDate, loc, user.StoreID)
	if err != nil {
		log.Printf("Error get cash flow summary: %v", err)
		http.Error(w, "Gagal mengambil data arus kas summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetTrend handles GET /api/cash-flow/trend
func (h *CashFlowHandler) GetTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	loc, tzName := parseTimezoneCF(r)
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startStr != "" && endStr != "" {
		startDateParsed, errParse1 := time.Parse("2006-01-02", startStr)
		endDateParsed, errParse2 := time.Parse("2006-01-02", endStr)
		if errParse1 != nil || errParse2 != nil {
			http.Error(w, "Format tanggal tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		startDate = time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(), 23, 59, 59, 999999999, loc)
	}

	trend, err := h.service.GetTrend(startDate, endDate, loc, tzName, user.StoreID)
	if err != nil {
		log.Printf("Error get cash flow trend: %v", err)
		http.Error(w, "Gagal mengambil trend arus kas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trend)
}

// GetLedger handles GET /api/cash-flow/ledger?start_date=X&end_date=Y&timezone=Asia/Makassar
func (h *CashFlowHandler) GetLedger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	loc, _ := parseTimezoneCF(r)
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startStr != "" && endStr != "" {
		startDateParsed, err1 := time.Parse("2006-01-02", startStr)
		endDateParsed, err2 := time.Parse("2006-01-02", endStr)
		if err1 != nil || err2 != nil {
			http.Error(w, "Format tanggal tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		startDate = time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(), 23, 59, 59, 999999999, loc)
	} else {
		// Default: 30 hari terakhir
		now := time.Now().In(loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
		startDate = endDate.AddDate(0, -1, 0)
	}

	// Parse pagination params (default: page=1, limit=100)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	ledger, err := h.service.GetLedger(startDate, endDate, page, limit, user.StoreID)
	if err != nil {
		log.Printf("Error get cash flow ledger: %v", err)
		http.Error(w, "Gagal mengambil data ledger", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ledger)
}

// GetFunds handles GET /api/cash-flow/funds
func (h *CashFlowHandler) GetFunds(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	result, err := h.service.GetFunds(page, limit, user.StoreID)
	if err != nil {
		log.Printf("Error get funds: %v", err)
		http.Error(w, "Gagal mengambil data dana", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// CreateFund handles POST /api/cash-flow/funds
func (h *CashFlowHandler) CreateFund(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, `{"message":"Forbidden: hanya admin"}`, http.StatusForbidden)
		return
	}
	var req models.CashFundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Format JSON tidak valid"}`, http.StatusBadRequest)
		return
	}
	fund, err := h.service.CreateFund(&req, user.ID, user.StoreID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dana berhasil dicatat",
		"data":    fund,
	})
}

// DeleteFund handles DELETE /api/cash-flow/funds/:id
func (h *CashFlowHandler) DeleteFund(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, `{"message":"Forbidden: hanya admin"}`, http.StatusForbidden)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/cash-flow/funds/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"message":"ID tidak valid"}`, http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteFund(id, user.StoreID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "tidak ditemukan") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data dana berhasil dihapus"})
}

// GetInitialBalance handles GET /api/cash-flow/initial-balance
func (h *CashFlowHandler) GetInitialBalance(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: User context missing", http.StatusUnauthorized)
		return
	}

	result, err := h.service.GetInitialBalance(user.StoreID)
	if err != nil {
		log.Printf("Error get initial balance: %v", err)
		http.Error(w, "Gagal menghitung saldo awal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
