package models

import "time"

// Supplier represents a product supplier
type Supplier struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	ContactPerson  *string   `json:"contact_person,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	Email          *string   `json:"email,omitempty"`
	Address        *string   `json:"address,omitempty"`
	Notes          *string   `json:"notes,omitempty"`
	IsActive       bool      `json:"is_active"`
	TotalPurchases int       `json:"total_purchases"`
	TotalSpent     float64   `json:"total_spent"`
	TotalPayable   float64   `json:"total_payable"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateSupplierRequest DTO for POST /suppliers
type CreateSupplierRequest struct {
	Name          string  `json:"name" validate:"required"`
	ContactPerson *string `json:"contact_person"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
	Notes         *string `json:"notes"`
}

// UpdateSupplierRequest DTO for PUT /suppliers/:id
type UpdateSupplierRequest struct {
	Name          *string `json:"name"`
	ContactPerson *string `json:"contact_person"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
	Notes         *string `json:"notes"`
	IsActive      *bool   `json:"is_active"`
}

// SupplierPayable represents a debt/payable entry to a supplier
type SupplierPayable struct {
	ID         int       `json:"id"`
	SupplierID int       `json:"supplier_id"`
	PurchaseID *int      `json:"purchase_id,omitempty"`
	Amount     float64   `json:"amount"`
	PaidAmount float64   `json:"paid_amount"`
	Status     string    `json:"status"` // unpaid, partial, paid
	DueDate    *string   `json:"due_date,omitempty"`
	Notes      *string   `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreatePayableRequest DTO for POST /suppliers/:id/payables
type CreatePayableRequest struct {
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	PaidAmount float64 `json:"paid_amount"`
	Status     string  `json:"status"`
	DueDate    *string `json:"due_date"`
	PurchaseID *int    `json:"purchase_id"`
	Notes      *string `json:"notes"`
}

// UpdatePayableRequest DTO for PUT /suppliers/:id/payables/:id
type UpdatePayableRequest struct {
	Amount  *float64 `json:"amount"`
	DueDate *string  `json:"due_date"`
	Notes   *string  `json:"notes"`
}

// PayablePayment represents a payment made toward a payable
type PayablePayment struct {
	ID          int       `json:"id"`
	PayableID   int       `json:"payable_id"`
	Amount      float64   `json:"amount"`
	PaymentDate string    `json:"payment_date"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreatePaymentRequest DTO for POST /payables/:id/payments
type CreatePaymentRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	PaymentDate string  `json:"payment_date" validate:"required"`
	Notes       *string `json:"notes"`
}
