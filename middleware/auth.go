package middleware

import (
	"context"
	"kasir-api/models"
	"kasir-api/utils"
	"net/http"
	"strings"
)

// ContextKey adalah type untuk context key
type ContextKey string

const (
	// UserContextKey adalah key untuk menyimpan user di context
	UserContextKey ContextKey = "user"
)

// AuthMiddleware adalah middleware untuk validasi JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Invalid authorization header format. Use: Bearer <token>"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Simpan user info ke context
		user := &models.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext mengambil user dari context
func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
