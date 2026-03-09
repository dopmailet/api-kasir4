package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"

	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomerHandler struct {
	service  *services.CustomerService
	validate *validator.Validate
}

func NewCustomerHandler(service *services.CustomerService, validate *validator.Validate) *CustomerHandler {
	return &CustomerHandler{
		service:  service,
		validate: validate,
	}
}

// Create handles creating a new customer
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "Validasi gagal: "+err.Error(), http.StatusBadRequest)
		return
	}

	customer, err := h.service.Create(&req)
	if err != nil {
		http.Error(w, "Gagal membuat customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetAll handles reading customers with pagination and search
func (h *CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	status := r.URL.Query().Get("status") // all, active, inactive
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}

	customers, total, err := h.service.GetAll(search, status, page, limit, sortBy, sortOrder)
	if err != nil {
		http.Error(w, "Gagal mengambil data customer", http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	response := map[string]interface{}{
		"data": customers,
		"pagination": map[string]interface{}{
			"page":         page,
			"limit":        limit,
			"total_items":  total,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Search handles searching active customers quickly for POS
func (h *CustomerHandler) Search(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	if search == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Customer{})
		return
	}

	// Always status: active, page: 1, limit: 10
	customers, _, err := h.service.GetAll(search, "active", 1, 10, "id", "asc")
	if err != nil {
		http.Error(w, "Gagal mencari data customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// GetByID handles getting single customer details
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/customers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	customer, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// Update handles editing existing customer
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/customers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var req models.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	customer, err := h.service.Update(id, &req)
	if err != nil {
		http.Error(w, "Gagal update customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// GetTransactions retrieves the order history for a customer
func (h *CustomerHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/customers/")
	idStr = strings.TrimSuffix(idStr, "/transactions")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	txs, err := h.service.GetTransactions(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if txs == nil {
		txs = []models.TransactionWithItems{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

// GetLoyaltyHistory retrieves the point earning history
func (h *CustomerHandler) GetLoyaltyHistory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/customers/")
	idStr = strings.TrimSuffix(idStr, "/loyalty-transactions")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	history, err := h.service.GetLoyaltyHistory(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if history == nil {
		history = []models.LoyaltyTransaction{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
