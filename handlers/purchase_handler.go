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
)

// PurchaseHandler handles HTTP requests for purchases
// Handler untuk pembelian/pengadaan barang
type PurchaseHandler struct {
	service *services.PurchaseService
}

// NewPurchaseHandler creates a new PurchaseHandler
func NewPurchaseHandler(service *services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

// HandlePurchases handles /api/purchases (GET all, POST new)
func (h *PurchaseHandler) HandlePurchases(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetAll(w, r)
	case "POST":
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandlePurchaseByID handles /api/purchases/{id} (GET by ID)
func (h *PurchaseHandler) HandlePurchaseByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.GetByID(w, r)
}

// Create handles POST /api/purchases
// Fungsi ini mencatat pembelian baru
func (h *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Ambil user dari context (harus Admin)
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden: Hanya Admin yang bisa mencatat pembelian", http.StatusForbidden)
		return
	}

	// Decode request body
	var req models.PurchaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("⚠️ Handler: Invalid request body for purchase: %v", err)
		http.Error(w, "Format request tidak valid", http.StatusBadRequest)
		return
	}

	// Panggil service untuk proses pembelian
	purchase, err := h.service.Create(&req, user.ID)
	if err != nil {
		log.Printf("❌ Handler: Error creating purchase: %v", err)
		// Cek apakah error validasi
		errMsg := err.Error()
		if strings.Contains(errMsg, "wajib") || strings.Contains(errMsg, "harus") ||
			strings.Contains(errMsg, "minimal") || strings.Contains(errMsg, "tidak ditemukan") ||
			strings.Contains(errMsg, "tidak boleh") {
			http.Error(w, errMsg, http.StatusBadRequest)
		} else {
			http.Error(w, errMsg, http.StatusInternalServerError)
		}
		return
	}

	// Response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(purchase)
}

// GetAll handles GET /api/purchases
// Fungsi ini mengambil riwayat semua pembelian
func (h *PurchaseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	purchases, err := h.service.GetAll()
	if err != nil {
		log.Printf("❌ Handler: Error getting purchases: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchases)
}

// GetByID handles GET /api/purchases/{id}
// Fungsi ini mengambil detail 1 pembelian
func (h *PurchaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/purchases/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("⚠️ Handler: Invalid purchase ID: %s", idStr)
		http.Error(w, "ID pembelian tidak valid", http.StatusBadRequest)
		return
	}

	purchase, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("❌ Handler: Error getting purchase ID %d: %v", id, err)
		if strings.Contains(err.Error(), "tidak ditemukan") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchase)
}
