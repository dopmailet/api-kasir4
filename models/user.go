package models

import (
	"time"
)

// User adalah struct untuk data user (Admin & Kasir)
type User struct {
	ID          int       `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Password    string    `json:"-" db:"password"` // "-" = tidak ditampilkan di JSON response
	NamaLengkap string    `json:"nama_lengkap" db:"nama_lengkap"`
	Role        string    `json:"role" db:"role"` // "admin" atau "kasir"
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole constants untuk validasi
const (
	RoleAdmin = "admin"
	RoleKasir = "kasir"
)

// IsAdmin mengecek apakah user adalah admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsKasir mengecek apakah user adalah kasir
func (u *User) IsKasir() bool {
	return u.Role == RoleKasir
}

// LoginRequest adalah struct untuk request login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse adalah struct untuk response login
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
