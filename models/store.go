package models

import "time"

// SubscriptionPackage - Paket langganan yang tersedia di platform
type SubscriptionPackage struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`         // contoh: Gratis, Basic, Pro
	MaxKasir    int         `json:"max_kasir"`    // batas jumlah kasir per toko
	MaxProducts int         `json:"max_products"` // batas jumlah produk per toko
	Price       float64     `json:"price"`        // harga per bulan (0 = gratis)
	Features    interface{} `json:"features"`     // fitur tambahan dalam format JSON
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
}

// Store - Data toko yang mendaftar ke platform
type Store struct {
	ID                    int        `json:"id"`
	Name                  string     `json:"name"`
	Address               *string    `json:"address,omitempty"`
	Phone                 *string    `json:"phone,omitempty"`
	Email                 *string    `json:"email,omitempty"`
	SubscriptionPackageID int        `json:"subscription_package_id"`
	SubscriptionEndDate   *time.Time `json:"subscription_end_date,omitempty"`
	IsActive              bool       `json:"is_active"`
	IsVerified            bool       `json:"is_verified"` // true jika pernah berlangganan paket berbayar
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`

	// Relasi (diisi saat dibutuhkan)
	SubscriptionPackage *SubscriptionPackage `json:"subscription_package,omitempty"`
}

// StoreRegisterRequest - Data yang dikirim user saat daftar baru
type StoreRegisterRequest struct {
	// Data Toko
	StoreName string `json:"store_name" validate:"required,min=3"`

	// Data Admin Toko (pemilik)
	AdminUsername string `json:"username" validate:"required,min=3"`
	AdminPassword string `json:"password" validate:"required,min=6"`
	AdminName     string `json:"full_name" validate:"required"`
	AdminEmail    string `json:"email" validate:"required,email"`
}

// StoreRegisterResponse - Respons setelah berhasil daftar
type StoreRegisterResponse struct {
	Message string `json:"message"`
	Store   *Store `json:"store"`
	Token   string `json:"token"` // langsung login otomatis setelah daftar
}
