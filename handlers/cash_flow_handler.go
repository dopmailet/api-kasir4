package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"log"
	"net/http"
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

	loc, _ := parseTimezoneCF(r)
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

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

	summary, err := h.service.GetSummary(startDate, endDate, loc)
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

	trend, err := h.service.GetTrend(startDate, endDate, loc, tzName)
	if err != nil {
		log.Printf("Error get cash flow trend: %v", err)
		http.Error(w, "Gagal mengambil trend arus kas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trend)
}
