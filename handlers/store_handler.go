package handlers

import (
	"encoding/json"
	"kasir-api/middleware"
	"kasir-api/services"
	"net/http"
)

type StoreHandler struct {
	service *services.StoreService
}

func NewStoreHandler(service *services.StoreService) *StoreHandler {
	return &StoreHandler{service: service}
}

// GetMyStoreInfo handles GET /api/store/info
// Endpoint krusial bagi frontend untuk merender badge langganan ("Free" / "Pro") di header
func (h *StoreHandler) GetMyStoreInfo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: Token tidak tersedia", http.StatusUnauthorized)
		return
	}

	store, err := h.service.GetMyStoreInfo(user.StoreID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": store})
}

// GetLimits handles GET /api/store/limits
// Endpoint untuk menampilkan sisa kuota (kasir, produk, sales hari ini)
func (h *StoreHandler) GetLimits(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized: Token tidak tersedia", http.StatusUnauthorized)
		return
	}

	limits, err := h.service.GetStoreLimits(user.StoreID)
	if err != nil {
		http.Error(w, "Gagal mendapatkan info limit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": limits})
}
