package handlers

import (
	"encoding/json"
	"kasir-api/middleware"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TransactionHandler handles HTTP requests for transactions
type TransactionHandler struct {
	service *services.TransactionService
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Checkout handles POST /api/checkout
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CheckoutRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Get current user
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser != nil {
		req.CreatedBy = currentUser.ID
	}

	transaction, err := h.service.Checkout(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// HandleTransactions handles GET /api/transactions?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&timezone=Asia/Makassar
// Mendukung filter tanggal opsional. Tanpa filter → semua transaksi.
func (h *TransactionHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil optional query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Parsing user_id optional
	var userID *int
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = &id
		}
	}

	var transactions []models.Transaction
	var err error

	if startDateStr != "" && endDateStr != "" {
		// Parse timezone (default: Asia/Jakarta)
		loc, _ := parseTimezone(r)

		// Parse tanggal
		startDateParsed, parseErr := time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			http.Error(w, "Format start_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		endDateParsed, parseErr := time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			http.Error(w, "Format end_date tidak valid (gunakan: YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		// Buat boundary waktu berdasarkan timezone user
		startDate := time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)
		endDate := time.Date(endDateParsed.Year(), endDateParsed.Month(), endDateParsed.Day(), 23, 59, 59, 999999999, loc)

		transactions, err = h.service.GetByDateRange(startDate, endDate, userID)
	} else {
		// Tanpa filter tanggal → ambil semua
		transactions, err = h.service.GetAll(userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// HandleTransactionByID handles /api/transactions/{id} (GET by ID)
func (h *TransactionHandler) HandleTransactionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.GetByID(w, r)
}

// GetByID handles GET /api/transactions/{id}
// Returns full transaction detail with items
func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID transaksi tidak valid", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": result,
	})
}
