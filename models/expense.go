package models

import "time"

// Expense merepresentasikan tabel pengeluaran operasional
type Expense struct {
	ID              int       `json:"id"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	ExpenseDate     string    `json:"expense_date"` // YYYY-MM-DD format
	IsRecurring     bool      `json:"is_recurring"`
	RecurringPeriod *string   `json:"recurring_period,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	CreatedBy       int       `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relasi untuk view/summary jika diperlukan
	CreatorName string `json:"creator_name,omitempty"`
}

// CreateExpenseRequest DTO untuk payload POST
type CreateExpenseRequest struct {
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	ExpenseDate     string  `json:"expense_date"`
	IsRecurring     bool    `json:"is_recurring"`
	RecurringPeriod *string `json:"recurring_period,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

// UpdateExpenseRequest DTO untuk payload PUT
type UpdateExpenseRequest struct {
	Category        *string  `json:"category"`
	Description     *string  `json:"description"`
	Amount          *float64 `json:"amount"`
	ExpenseDate     *string  `json:"expense_date"`
	IsRecurring     *bool    `json:"is_recurring"`
	RecurringPeriod *string  `json:"recurring_period"` // disengaja tidak di-omit agar bisa diset null jika false
	Notes           *string  `json:"notes"`
}
