package middleware

import (
	"kasir-api/models"
	"net/http"
)

// RequireRole adalah middleware untuk membatasi akses berdasarkan role
// Contoh: RequireRole("admin") atau RequireRole("admin", "kasir")
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil user dari context (sudah diset oleh AuthMiddleware)
			user := GetUserFromContext(r.Context())
			if user == nil {
				http.Error(w, `{"error":"User not found in context"}`, http.StatusUnauthorized)
				return
			}

			// Cek apakah role user ada di allowedRoles
			hasAccess := false
			for _, role := range allowedRoles {
				if user.Role == role {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				http.Error(w, `{"error":"Forbidden: Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			// User punya akses, lanjutkan ke handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin adalah shortcut untuk RequireRole("admin")
func RequireAdmin(next http.Handler) http.Handler {
	return RequireRole(models.RoleAdmin)(next)
}

// RequireKasir adalah shortcut untuk RequireRole("kasir")
func RequireKasir(next http.Handler) http.Handler {
	return RequireRole(models.RoleKasir)(next)
}

// RequireAdminOrKasir adalah shortcut untuk RequireRole("admin", "kasir")
func RequireAdminOrKasir(next http.Handler) http.Handler {
	return RequireRole(models.RoleAdmin, models.RoleKasir)(next)
}
