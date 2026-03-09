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

// Register handles POST /api/auth/register (admin only)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Parse request body
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		NamaLengkap string `json:"nama_lengkap"`
		Role        string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 2. Validasi input
	if req.Username == "" || req.Password == "" || req.NamaLengkap == "" || req.Role == "" {
		http.Error(w, `{"error":"Semua field wajib diisi"}`, http.StatusBadRequest)
		return
	}

	// 3. Proses register
	user, err := h.authService.Register(req.Username, req.Password, req.NamaLengkap, req.Role)
	if err != nil {
		slog.Error("Registration failed", "error", err)

		switch err {
		case models.ErrInvalidRole:
			http.Error(w, `{"error":"Role tidak valid. Gunakan 'admin' atau 'kasir'"}`, http.StatusBadRequest)
		default:
			http.Error(w, `{"error":"Registrasi gagal"}`, http.StatusInternalServerError)
		}
		return
	}

	// 4. Log successful registration
	slog.Info("New user registered", "username", user.Username, "role", user.Role)

	// 5. Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User berhasil didaftarkan",
		"data":    user,
	})
}

// ChangePassword handles POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement change password
	// Akan diimplementasikan nanti jika diperlukan
	http.Error(w, `{"error":"Not implemented yet"}`, http.StatusNotImplemented)
}
