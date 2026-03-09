package models

import "time"

// Customer mewakili data pelanggan
type Customer struct {
	ID                int        `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	Phone             string     `json:"phone" db:"phone"`
	Address           *string    `json:"address,omitempty" db:"address"`
	Notes             *string    `json:"notes,omitempty" db:"notes"`
	LoyaltyPoints     int        `json:"loyalty_points" db:"loyalty_points"`
	TotalSpent        float64    `json:"total_spent" db:"total_spent"`
	TotalTransactions int        `json:"total_transactions" db:"total_transactions"`
	LastTransactionAt *time.Time `json:"last_transaction_at,omitempty" db:"last_transaction_at"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateCustomerRequest DTO for POST
type CreateCustomerRequest struct {
	Name     string  `json:"name" validate:"required"`
	Phone    string  `json:"phone" validate:"required"`
	Address  *string `json:"address"`
	Notes    *string `json:"notes"`
	IsActive *bool   `json:"is_active"` // Default true
}

// UpdateCustomerRequest DTO for PUT
type UpdateCustomerRequest struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	Notes    *string `json:"notes"`
	IsActive *bool   `json:"is_active"`
}

// CustomerSummary DTO for partial response in transaction/checkout
type CustomerSummary struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	LoyaltyPoints int    `json:"loyalty_points"`
}
