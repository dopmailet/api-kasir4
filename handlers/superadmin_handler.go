package handlers

import (
	"encoding/json"
	"kasir-api/middleware"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type SuperadminHandler struct {
	service *services.SuperadminService
}

func NewSuperadminHandler(service *services.SuperadminService) *SuperadminHandler {
	return &SuperadminHandler{service: service}
}

// GetPublicPackages handles GET /api/packages (Public - No Auth Required)
func (h *SuperadminHandler) GetPublicPackages(w http.ResponseWriter, r *http.Request) {
	// Menampilkan hanya package yang aktif
	pkgs, err := h.service.GetAllSubscriptionPackages(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": pkgs})
}

// GetAllStores handles GET /api/superadmin/stores
func (h *SuperadminHandler) GetAllStores(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden: Superadmin level access required", http.StatusForbidden)
		return
	}

	stores, err := h.service.GetAllStores()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": stores})
}

// UpdateStoreStatus handles PUT /api/superadmin/stores/{id}/status
func (h *SuperadminHandler) UpdateStoreStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden: Superadmin level access required", http.StatusForbidden)
		return
	}

	// Extract ID: hapus prefix dan hapus suffix
	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/stores/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Store ID", http.StatusBadRequest)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStoreStatus(id, req.IsActive); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Status toko berhasil diperbarui"})
}

// GetAllPackages handles GET /api/superadmin/packages
func (h *SuperadminHandler) GetAllPackages(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden: Superadmin level access required", http.StatusForbidden)
		return
	}

	pkgs, err := h.service.GetAllSubscriptionPackages(false) // false = tampilkan semua, termasuk yg tidak aktif
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": pkgs})
}

// GetPackageByID handles GET /api/superadmin/packages/{id}
func (h *SuperadminHandler) GetPackageByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Package ID", http.StatusBadRequest)
		return
	}

	pkg, err := h.service.GetSubscriptionPackageByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": pkg})
}

// CreatePackage handles POST /api/superadmin/packages
func (h *SuperadminHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req models.SubscriptionPackage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Validasi dasar
	if req.Name == "" {
		http.Error(w, "Nama paket wajib diisi", http.StatusBadRequest)
		return
	}
	if req.MaxKasir < 1 {
		req.MaxKasir = 1
	}
	if req.MaxProducts < 1 {
		req.MaxProducts = 100
	}

	if err := h.service.CreateSubscriptionPackage(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": req})
}

// UpdatePackage handles PUT /api/superadmin/packages/{id}
func (h *SuperadminHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Package ID", http.StatusBadRequest)
		return
	}

	var req models.SubscriptionPackage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	req.ID = id
	if err := h.service.UpdateSubscriptionPackage(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": req})
}

// DeletePackage handles DELETE /api/superadmin/packages/{id}
func (h *SuperadminHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Package ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSubscriptionPackage(id); err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict jika masih dipakai
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Paket berhasil dihapus"})
}

// UpgradePackage handles PUT /api/superadmin/stores/{id}/package
// Mengaktifkan paket langganan tertentu ke sebuah toko secara manual oleh Superadmin
func (h *SuperadminHandler) UpgradePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden: Superadmin level access required", http.StatusForbidden)
		return
	}

	// Extract Store ID dari path: /api/superadmin/stores/{id}/package
	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/stores/")
	idStr = strings.TrimSuffix(idStr, "/package")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Store ID", http.StatusBadRequest)
		return
	}

	var req struct {
		PackageID    int `json:"package_id"`    // ID paket yang akan diaktifkan
		DurationDays int `json:"duration_days"` // Jumlah hari aktif (contoh: 30, 90, 365)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if req.PackageID == 0 {
		http.Error(w, "package_id wajib diisi", http.StatusBadRequest)
		return
	}
	if req.DurationDays <= 0 {
		req.DurationDays = 30 // Default 30 hari jika tidak disertakan
	}

	if err := h.service.UpdateStorePackage(id, req.PackageID, req.DurationDays); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Paket langganan toko berhasil diperbarui"})
}

// DeleteStore handles DELETE /api/superadmin/stores/{id}
func (h *SuperadminHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/stores/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Store ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteStore(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Toko berhasil dihapus"})
}

// DeleteUnverifiedStores handles DELETE /api/superadmin/stores/bulk?type=unverified|never_subscribed
func (h *SuperadminHandler) DeleteBulkStores(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	delType := r.URL.Query().Get("type") // "unverified" atau "never_subscribed"
	var deleted int64
	var err error

	switch delType {
	case "unverified":
		deleted, err = h.service.DeleteUnverifiedStores()
	case "never_subscribed":
		deleted, err = h.service.DeleteNeverSubscribedStores()
	default:
		http.Error(w, "Parameter type wajib: unverified atau never_subscribed", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Toko berhasil dihapus secara massal",
		"deleted": deleted,
	})
}
