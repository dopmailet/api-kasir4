package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/repositories"
	"net/http"
	"strconv"
)

type DiscountHandler struct {
	repo *repositories.DiscountRepository
}

func NewDiscountHandler(repo *repositories.DiscountRepository) *DiscountHandler {
	return &DiscountHandler{repo: repo}
}

// GetAll handles GET /api/discounts (Admin)
func (h *DiscountHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	discounts, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discounts)
}

// GetActive handles GET /api/discounts/active (Kasir/Checkout UI)
func (h *DiscountHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	discounts, err := h.repo.GetActive()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discounts)
}

// Create handles POST /api/discounts (Admin)
func (h *DiscountHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var d models.Discount
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple validation
	if d.Name == "" || d.Value <= 0 {
		http.Error(w, "Name and Value are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

// Update handles PUT /api/discounts/{id} (Admin)
func (h *DiscountHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ID from URL path (assuming /api/discounts/{id})
	// In standard ServeMux, this is tricky.
	// We will rely on r.URL.Path logic or simpler query param for now?
	// But in main.go we use specific path matching.
	// Since we use http.StripPrefix or similar in main.go, let's parse from path.
	// Example path: /api/discounts/123
	// Last element is ID.

	idStr := r.URL.Path[len("/api/discounts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var d models.Discount
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(id, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Discount updated"})
}

// Delete handles DELETE /api/discounts/{id} (Admin)
func (h *DiscountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/discounts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Discount deleted"})
}
