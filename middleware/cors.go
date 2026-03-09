package middleware

import (
	"net/http"
)

// CORSMiddleware menambahkan CORS headers untuk mengizinkan akses dari frontend
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers untuk development (Allow All)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key, Accept, Origin, X-Requested-With")
		// Note: Access-Control-Allow-Credentials tidak boleh true jika Origin adalah *
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
