package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// SubscriptionPackage - Paket langganan yang tersedia di platform
// SubscriptionPackage - Paket langganan yang tersedia di platform
type SubscriptionPackage struct {
	ID              int              `json:"id" db:"id"`
	Name            string           `json:"name" db:"name"`
	MaxKasir        int              `json:"max_kasir" db:"max_kasir"`
	MaxProducts     int              `json:"max_products" db:"max_products"`
	Price           float64          `json:"price" db:"price"`
	IsActive        bool             `json:"is_active" db:"is_active"`
	Description     *string          `json:"description" db:"description"`
	Features        JSONBStringArray `json:"features" db:"features"` // Use custom scanner for JSONB array
	Period          string           `json:"period" db:"period"`
	DiscountPercent float64          `json:"discount_percent" db:"discount_percent"`
	DiscountLabel   *string          `json:"discount_label" db:"discount_label"`
	IsPopular       bool             `json:"is_popular" db:"is_popular"`
	SortOrder       int              `json:"sort_order" db:"sort_order"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
}

// JSONBStringArray - Tipe khusus untuk mem-parsing JSONB array of string dari database
type JSONBStringArray []string

// Scan implements the sql.Scanner interface
func (a *JSONBStringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("JSONBStringArray: expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, a)
}

// Value implements the driver.Valuer interface
func (a JSONBStringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
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
