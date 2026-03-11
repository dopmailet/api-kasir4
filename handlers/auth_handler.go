package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"log/slog"
	"net/http"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// 1. Parse request body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to parse login request", "error", err)
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 2. Validasi input
	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"Username dan password wajib diisi"}`, http.StatusBadRequest)
		return
	}

	// 3. Proses login
	response, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		slog.Warn("Login failed", "username", req.Username, "error", err)

		// Handle specific errors
		switch err {
		case models.ErrInvalidCredentials:
			http.Error(w, `{"error":"Username atau password salah"}`, http.StatusUnauthorized)
		case models.ErrUserInactive:
			http.Error(w, `{"error":"User tidak aktif"}`, http.StatusForbidden)
		default:
			http.Error(w, `{"error":"Login gagal"}`, http.StatusInternalServerError)
		}
		return
	}

	// 4. Log successful login
	slog.Info("User logged in", "username", req.Username, "role", response.User.Role)

	// 5. Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

// Register handles POST /api/auth/register (public)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Parse request body
	var req models.StoreRegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 2. Validasi input
	if req.StoreName == "" || req.AdminUsername == "" || req.AdminPassword == "" || req.AdminName == "" {
		http.Error(w, `{"error":"Semua field wajib diisi"}`, http.StatusBadRequest)
		return
	}

	// 3. Proses register (Toko + Admin)
	res, err := h.authService.RegisterStore(req)
	if err != nil {
		slog.Error("Registration failed", "error", err)
		http.Error(w, `{"error":"Registrasi gagal: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// 4. Log successful registration
	slog.Info("New store registered", "store", req.StoreName, "admin", req.AdminUsername)

	// 5. Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": res.Message,
		"data":    res,
	})
}

// ChangePassword handles POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement change password
	// Akan diimplementasikan nanti jika diperlukan
	http.Error(w, `{"error":"Not implemented yet"}`, http.StatusNotImplemented)
}
