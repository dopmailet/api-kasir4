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

type PayrollHandler struct {
	service *services.PayrollService
}

func NewPayrollHandler(service *services.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

// GetAll handles GET /api/payroll
func (h *PayrollHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	employeeID, _ := strconv.Atoi(r.URL.Query().Get("employee_id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	loc, _ := parseTimezone(r) // function dari report_handler.go

	var startDate, endDate time.Time
	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		startDate, _ = time.ParseInLocation("2006-01-02", startStr, loc)
	}
	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		endDate, _ = time.ParseInLocation("2006-01-02", endStr, loc)
		// Geser ke akhir hari
		if !endDate.IsZero() {
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, loc)
		}
	}

	payrolls, total, err := h.service.GetAll(employeeID, startDate, endDate, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       payrolls,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	})
}

// GetByID handles GET /api/payroll/{id}
func (h *PayrollHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/payroll/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Create handles POST /api/payroll
func (h *PayrollHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Ambil user dari JWT Context Token
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	createdBy := user.ID

	var req models.CreatePayrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation tambahan bila lib validate belum disetup di repo ini
	if req.EmployeeID <= 0 || req.GajiPokok < 0 {
		http.Error(w, "EmployeeID dan GajiPokok valid wajib diisi", http.StatusBadRequest)
		return
	}

	p, err := h.service.Create(req, createdBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// Update handles PUT /api/payroll/{id}
func (h *PayrollHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/payroll/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePayrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	p, err := h.service.Update(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Delete handles DELETE /api/payroll/{id}
func (h *PayrollHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/payroll/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data gaji berhasil dihapus"})
}

// GetReport handles GET /api/payroll/report
func (h *PayrollHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	loc, tzName := parseTimezone(r) // function internal helper

	// Penentuan rentang default (last 30 days fallback)
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr == "" || endDateStr == "" {
		// Default 30 hari lokal
		now := time.Now().In(loc)
		start := now.AddDate(0, 0, -30)
		startDate = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	} else {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, loc)
		if err != nil {
			http.Error(w, "Format start_date tidak valid", http.StatusBadRequest)
			return
		}

		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, loc)
		if err != nil {
			http.Error(w, "Format end_date tidak valid", http.StatusBadRequest)
			return
		}
	}

	// Panggil service
	report, err := h.service.GetReport(startDate, endDate, tzName)
	if err != nil {
		log.Printf("Error get payroll report: %v", err)
		http.Error(w, "Gagal mengambil laporan penggajian", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
