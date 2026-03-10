package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type SupplierHandler struct {
	service *services.SupplierService
}

func NewSupplierHandler(service *services.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

// jsonOK writes a JSON success response with { "data": ... }
func jsonData(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}

func jsonMsg(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func jsonErr(w http.ResponseWriter, status int, err error) {
	jsonMsg(w, status, err.Error())
}

// ─── SUPPLIER HANDLERS ───

// GetAll handles GET /api/suppliers
func (h *SupplierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var isActive *bool
	if v := r.URL.Query().Get("is_active"); v != "" {
		b := v == "true"
		isActive = &b
	}

	suppliers, err := h.service.GetAll(search, isActive)
	if err != nil {
		jsonErr(w, http.StatusInternalServerError, err)
		return
	}
	if suppliers == nil {
		suppliers = []models.Supplier{}
	}
	jsonData(w, http.StatusOK, suppliers)
}

// GetDebtSummary handles GET /api/suppliers/debt-summary
func (h *SupplierHandler) GetDebtSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetDebtSummary()
	if err != nil {
		jsonErr(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetByID handles GET /api/suppliers/:id
func (h *SupplierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/suppliers/")
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID supplier tidak valid")
		return
	}
	supplier, err := h.service.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonData(w, http.StatusOK, supplier)
}

// Create handles POST /api/suppliers
func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMsg(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	supplier, err := h.service.Create(&req)
	if err != nil {
		if strings.Contains(err.Error(), "wajib") {
			jsonErr(w, http.StatusBadRequest, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonData(w, http.StatusCreated, supplier)
}

// Update handles PUT /api/suppliers/:id
func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/suppliers/")
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID supplier tidak valid")
		return
	}
	var req models.UpdateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMsg(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	supplier, err := h.service.Update(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonData(w, http.StatusOK, supplier)
}

// Delete handles DELETE /api/suppliers/:id
func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/api/suppliers/")
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID supplier tidak valid")
		return
	}
	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonMsg(w, http.StatusOK, "Supplier deleted")
}

// ─── PAYABLE HANDLERS ───

// GetPayables handles GET /api/suppliers/:id/payables
func (h *SupplierHandler) GetPayables(w http.ResponseWriter, r *http.Request) {
	supplierID, err := extractSupplierIDFromPayablePath(r.URL.Path)
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID supplier tidak valid")
		return
	}
	payables, err := h.service.GetPayables(supplierID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	if payables == nil {
		payables = []models.SupplierPayable{}
	}
	jsonData(w, http.StatusOK, payables)
}

// CreatePayable handles POST /api/suppliers/:id/payables
func (h *SupplierHandler) CreatePayable(w http.ResponseWriter, r *http.Request) {
	supplierID, err := extractSupplierIDFromPayablePath(r.URL.Path)
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID supplier tidak valid")
		return
	}
	var req models.CreatePayableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMsg(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	payable, err := h.service.CreatePayable(supplierID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") || strings.Contains(err.Error(), "harus") {
			jsonErr(w, http.StatusBadRequest, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonData(w, http.StatusCreated, payable)
}

// UpdatePayable handles PUT /api/suppliers/:supplierId/payables/:payableId
func (h *SupplierHandler) UpdatePayable(w http.ResponseWriter, r *http.Request) {
	supplierID, payableID, err := extractSupplierAndPayableID(r.URL.Path)
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID tidak valid")
		return
	}
	var req models.UpdatePayableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMsg(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	payable, err := h.service.UpdatePayable(supplierID, payableID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") || strings.Contains(err.Error(), "bukan milik") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonData(w, http.StatusOK, payable)
}

// ─── PAYMENT HANDLERS ───

// GetPayments handles GET /api/payables/:payableId/payments
func (h *SupplierHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	payableID, err := extractPayableIDFromPaymentPath(r.URL.Path)
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID payable tidak valid")
		return
	}
	payments, err := h.service.GetPayments(payableID)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	if payments == nil {
		payments = []models.PayablePayment{}
	}
	jsonData(w, http.StatusOK, payments)
}

// CreatePayment handles POST /api/payables/:payableId/payments
func (h *SupplierHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	payableID, err := extractPayableIDFromPaymentPath(r.URL.Path)
	if err != nil {
		jsonMsg(w, http.StatusBadRequest, "ID payable tidak valid")
		return
	}
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonMsg(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	payment, err := h.service.CreatePayment(payableID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "melebihi") || strings.Contains(err.Error(), "harus") || strings.Contains(err.Error(), "format") {
			jsonErr(w, http.StatusBadRequest, err)
		} else if strings.Contains(err.Error(), "tidak ditemukan") {
			jsonErr(w, http.StatusNotFound, err)
		} else {
			jsonErr(w, http.StatusInternalServerError, err)
		}
		return
	}
	jsonData(w, http.StatusCreated, payment)
}

// ─── URL HELPERS ───

func extractID(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	// hilangkan trailing slash atau sub-path
	idStr = strings.SplitN(idStr, "/", 2)[0]
	return strconv.Atoi(idStr)
}

// /api/suppliers/{supplierID}/payables  → supplierID
func extractSupplierIDFromPayablePath(path string) (int, error) {
	// path = /api/suppliers/5/payables
	trimmed := strings.TrimPrefix(path, "/api/suppliers/")
	parts := strings.Split(trimmed, "/")
	return strconv.Atoi(parts[0])
}

// /api/suppliers/{supplierID}/payables/{payableID} → supplierID, payableID
func extractSupplierAndPayableID(path string) (int, int, error) {
	trimmed := strings.TrimPrefix(path, "/api/suppliers/")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 3 {
		return 0, 0, strconv.ErrSyntax
	}
	sID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	pID, err := strconv.Atoi(parts[2])
	return sID, pID, err
}

// /api/payables/{payableID}/payments → payableID
func extractPayableIDFromPaymentPath(path string) (int, error) {
	trimmed := strings.TrimPrefix(path, "/api/payables/")
	parts := strings.Split(trimmed, "/")
	return strconv.Atoi(parts[0])
}
