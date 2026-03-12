package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kasir-api/middleware"
	"kasir-api/models"
	"kasir-api/services"
	"log"
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

// normalizeFeatures menerima json.RawMessage dan mengembalikan []string yang dinormalisasi.
// Mendukung format:
//   - []string             → dipakai langsung
//   - []object             → ambil field "label", "name", atau "value" (string pertama yang non-empty)
//   - null / absent / ""   → return []string{}
//
// Jika format tidak dikenali → return error (→ 400 di handler)
func normalizeFeatures(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return []string{}, nil
	}

	// Coba []string terlebih dahulu
	var strSlice []string
	if err := json.Unmarshal(raw, &strSlice); err == nil {
		// Filter string kosong
		result := make([]string, 0, len(strSlice))
		for _, s := range strSlice {
			if strings.TrimSpace(s) != "" {
				result = append(result, s)
			}
		}
		return result, nil
	}

	// Coba []object (mis. [{label:"Fitur A"}, ...])
	var objSlice []map[string]interface{}
	if err := json.Unmarshal(raw, &objSlice); err == nil {
		result := make([]string, 0, len(objSlice))
		for _, obj := range objSlice {
			// Ambil field label / name / value (prioritas label)
			val := ""
			for _, key := range []string{"label", "name", "value", "text"} {
				if v, ok := obj[key]; ok {
					if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
						val = strings.TrimSpace(s)
						break
					}
				}
			}
			if val != "" {
				result = append(result, val)
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("format features tidak valid: harus array string atau array object")
}

// marshalFeatures mengubah []string menjadi JSON string untuk SQL ::jsonb cast
func marshalFeatures(features []string) (string, error) {
	if len(features) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(features)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// normalizeStringPtr — pastikan *string yang berisi "" atau whitespace saja diubah ke nil (SQL NULL)
func normalizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	if strings.TrimSpace(*s) == "" {
		return nil
	}
	return s
}

// respondJSON helper — write JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondError helper — write JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// GetPublicPackages handles GET /api/packages (Public - No Auth Required)
func (h *SuperadminHandler) GetPublicPackages(w http.ResponseWriter, r *http.Request) {
	pkgs, err := h.service.GetAllSubscriptionPackages(true)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": pkgs})
}

// GetAllStores handles GET /api/superadmin/stores
func (h *SuperadminHandler) GetAllStores(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden: Superadmin level access required")
		return
	}

	stores, err := h.service.GetAllStores()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": stores})
}

// UpdateStoreStatus handles PUT /api/superadmin/stores/{id}/status
func (h *SuperadminHandler) UpdateStoreStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden: Superadmin level access required")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/stores/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := h.service.UpdateStoreStatus(id, req.IsActive); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Status toko berhasil diperbarui"})
}

// GetAllPackages handles GET /api/superadmin/packages
func (h *SuperadminHandler) GetAllPackages(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden: Superadmin level access required")
		return
	}

	pkgs, err := h.service.GetAllSubscriptionPackages(false)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": pkgs})
}

// GetPackageByID handles GET /api/superadmin/packages/{id}
func (h *SuperadminHandler) GetPackageByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Package ID")
		return
	}

	pkg, err := h.service.GetSubscriptionPackageByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Paket tidak ditemukan")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"data": pkg})
}

// CreatePackage handles POST /api/superadmin/packages
func (h *SuperadminHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	var req models.PackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validasi dasar
	if strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "name wajib diisi")
		return
	}
	if req.MaxKasir < 1 {
		req.MaxKasir = 1
	}
	if req.MaxProducts < 1 {
		req.MaxProducts = 100
	}

	// Normalisasi features → []string
	features, err := normalizeFeatures(req.Features)
	if err != nil {
		log.Printf("[POST /packages] Normalisasi features gagal: %v | raw=%s", err, string(req.Features))
		respondError(w, http.StatusBadRequest, "Format features tidak valid: "+err.Error())
		return
	}

	// Marshal features ke JSON string untuk SQL ::jsonb
	featuresJSON, err := marshalFeatures(features)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Gagal marshal features")
		return
	}

	log.Printf("[POST /packages] name=%s features=%v featuresJSON=%s", req.Name, features, featuresJSON)

	// Build model dari request DTO
	// normalizeStringPtr: pastikan "" → nil agar respons konsisten dengan DB (NULL)
	pkg := &models.SubscriptionPackage{
		Name:            strings.TrimSpace(req.Name),
		MaxKasir:        req.MaxKasir,
		MaxProducts:     req.MaxProducts,
		Price:           req.Price,
		IsActive:        req.IsActive,
		Description:     normalizeStringPtr(req.Description),
		Features:        features,
		Period:          req.Period,
		DiscountPercent: req.DiscountPercent,
		DiscountLabel:   normalizeStringPtr(req.DiscountLabel),
		IsPopular:       req.IsPopular,
		MaxDailySales:   req.MaxDailySales,
	}

	if err := h.service.CreateSubscriptionPackage(pkg, featuresJSON); err != nil {
		log.Printf("[POST /packages] Error create: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal membuat paket")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": pkg})
}

// UpdatePackage handles PUT /api/superadmin/packages/{id}
func (h *SuperadminHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	// Extract ID — handle trailing slash dan sub-path seperti "/popular"
	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Package ID")
		return
	}

	var req models.PackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validasi dasar
	if strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "name wajib diisi")
		return
	}

	// Normalisasi features → []string
	features, err := normalizeFeatures(req.Features)
	if err != nil {
		log.Printf("[PUT /packages/%d] Normalisasi features gagal: %v | raw=%s", id, err, string(req.Features))
		respondError(w, http.StatusBadRequest, "Format features tidak valid: "+err.Error())
		return
	}

	// Marshal features ke JSON string untuk SQL ::jsonb
	featuresJSON, err := marshalFeatures(features)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Gagal marshal features")
		return
	}

	log.Printf("[PUT /packages/%d] name=%s isPopular=%v isActive=%v features=%v featuresJSON=%s",
		id, req.Name, req.IsPopular, req.IsActive, features, featuresJSON)

	// Build model
	// normalizeStringPtr: pastikan "" → nil agar respons konsisten dengan DB (NULL)
	pkg := &models.SubscriptionPackage{
		ID:              id,
		Name:            strings.TrimSpace(req.Name),
		MaxKasir:        req.MaxKasir,
		MaxProducts:     req.MaxProducts,
		Price:           req.Price,
		IsActive:        req.IsActive,
		Description:     normalizeStringPtr(req.Description),
		Features:        features,
		Period:          req.Period,
		DiscountPercent: req.DiscountPercent,
		DiscountLabel:   normalizeStringPtr(req.DiscountLabel),
		IsPopular:       req.IsPopular,
		MaxDailySales:   req.MaxDailySales,
	}

	if err := h.service.UpdateSubscriptionPackage(pkg, featuresJSON); err != nil {
		if err.Error() == "package not found" {
			respondError(w, http.StatusNotFound, "Paket tidak ditemukan")
			return
		}
		log.Printf("[PUT /packages/%d] Error update: %v", id, err)
		respondError(w, http.StatusInternalServerError, "Gagal memperbarui paket")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": pkg})
}

// DeletePackage handles DELETE /api/superadmin/packages/{id}
func (h *SuperadminHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Package ID")
		return
	}

	if err := h.service.DeleteSubscriptionPackage(id); err != nil {
		respondError(w, http.StatusConflict, err.Error()) // 409 jika masih dipakai
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Paket berhasil dihapus"})
}

// TogglePackagePopular handles PATCH /api/superadmin/packages/{id}/popular
func (h *SuperadminHandler) TogglePackagePopular(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	// Extract ID: /api/superadmin/packages/{id}/popular
	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/packages/")
	idStr = strings.TrimSuffix(idStr, "/popular")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Package ID")
		return
	}

	var req struct {
		IsPopular bool `json:"is_popular"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	log.Printf("[PATCH /packages/%d/popular] is_popular=%v", id, req.IsPopular)

	if err := h.service.TogglePackagePopular(id, req.IsPopular); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Status populer paket berhasil diperbarui",
		"id":         id,
		"is_popular": req.IsPopular,
	})
}

// UpgradePackage handles PUT /api/superadmin/stores/{id}/package
func (h *SuperadminHandler) UpgradePackage(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden: Superadmin level access required")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/stores/")
	idStr = strings.TrimSuffix(idStr, "/package")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	var req struct {
		PackageID    int `json:"package_id"`
		DurationDays int `json:"duration_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if req.PackageID == 0 {
		respondError(w, http.StatusBadRequest, "package_id wajib diisi")
		return
	}
	if req.DurationDays <= 0 {
		req.DurationDays = 30
	}

	if err := h.service.UpdateStorePackage(id, req.PackageID, req.DurationDays); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Paket langganan toko berhasil diperbarui"})
}

// DeleteStore handles DELETE /api/superadmin/stores/{id}
func (h *SuperadminHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/superadmin/stores/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Store ID")
		return
	}

	if err := h.service.DeleteStore(id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Toko berhasil dihapus"})
}

// DeleteBulkStores handles DELETE /api/superadmin/stores/bulk?type=unverified|never_subscribed
func (h *SuperadminHandler) DeleteBulkStores(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil || !user.IsSuperadmin {
		respondError(w, http.StatusForbidden, "Forbidden")
		return
	}

	delType := r.URL.Query().Get("type")
	var deleted int64
	var err error

	switch delType {
	case "unverified":
		deleted, err = h.service.DeleteUnverifiedStores()
	case "never_subscribed":
		deleted, err = h.service.DeleteNeverSubscribedStores()
	default:
		respondError(w, http.StatusBadRequest, "Parameter type wajib: unverified atau never_subscribed")
		return
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Toko berhasil dihapus secara massal",
		"deleted": deleted,
	})
}
