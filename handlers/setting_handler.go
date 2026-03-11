package handlers

import (
	"encoding/json"
	"kasir-api/middleware"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
)

type SettingHandler struct {
	service *services.SettingService
}

func NewSettingHandler(service *services.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

// GetCustomerSettings returns the configurations for POS Customer UI and loyalty
func (h *SettingHandler) GetCustomerSettings(w http.ResponseWriter, r *http.Request) {
	// Ambil user dari context (diset oleh AuthMiddleware)
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	settings, err := h.service.GetCustomerSettings(user.StoreID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// UpdateCustomerSettings updates the configurations
func (h *SettingHandler) UpdateCustomerSettings(w http.ResponseWriter, r *http.Request) {
	// Ambil user dari context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Sesuai permintaan: Hanya Admin yang boleh update settings (bukan kasir)
	if user.Role != "admin" && user.Role != "superadmin" {
		http.Error(w, "Hanya admin yang dapat mengubah pengaturan", http.StatusForbidden)
		return
	}

	var req models.AppSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateCustomerSettings(user.StoreID, &req)
	if err != nil {
		http.Error(w, "Gagal mengupdate pengaturan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}
