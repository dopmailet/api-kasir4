package models

import "time"

// LoyaltyTransaction represents a loyalty point transaction
type LoyaltyTransaction struct {
	ID            int       `json:"id" db:"id"`
	CustomerID    int       `json:"customer_id" db:"customer_id"`
	TransactionID *int      `json:"transaction_id,omitempty" db:"transaction_id"`
	Type          string    `json:"type" db:"type"` // "earn" or "adjust"
	Points        int       `json:"points" db:"points"`
	Description   *string   `json:"description,omitempty" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	StoreID       int       `json:"store_id" db:"store_id"`
}

// Default settings response
type AppSettings struct {
	ShowCustomerInPOS       bool `json:"showCustomerInPOS"`
	EnableLoyaltyPoints     bool `json:"enableLoyaltyPoints"`
	ShowTransactionDiscount bool `json:"showTransactionDiscount"`
}

// PlatformSettings mewakili pengaturan global aplikasi (SaaS level)
type PlatformSettings struct {
	AdminWhatsApp *string `json:"admin_whatsapp"`
}
