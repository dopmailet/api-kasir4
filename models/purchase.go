package models

import "time"

// Purchase represents a purchase header (pembelian dari supplier)
// Struct ini menyimpan informasi header setiap pembelian
type Purchase struct {
	ID              int            `json:"id"`
	SupplierID      *int           `json:"supplier_id"`
	SupplierName    *string        `json:"supplier_name"`
	TotalAmount     float64        `json:"total_amount"`
	PaymentMethod   string         `json:"payment_method"`
	PaymentStatus   string         `json:"payment_status"`
	PaidAmount      float64        `json:"paid_amount"`
	RemainingAmount float64        `json:"remaining_amount"`
	DueDate         *string        `json:"due_date"`
	PaymentNotes    *string        `json:"payment_notes"`
	Notes           *string        `json:"notes"`
	CreatedBy       *int           `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	StoreID         int            `json:"store_id" db:"store_id"` // Multi-tenant isolation
	Items           []PurchaseItem `json:"items,omitempty"`
}

// PurchaseItem represents a purchase detail item
// Struct ini menyimpan detail setiap item dalam pembelian
type PurchaseItem struct {
	ID          int       `json:"id" db:"id"`
	PurchaseID  int       `json:"purchase_id" db:"purchase_id"`
	ProductID   *int      `json:"product_id,omitempty" db:"product_id"`   // NULL jika produk baru
	ProductName string    `json:"product_name" db:"product_name"`         // Nama produk (snapshot)
	Quantity    int       `json:"quantity" db:"quantity"`                 // Jumlah beli
	BuyPrice    float64   `json:"buy_price" db:"buy_price"`               // Harga beli per unit
	SellPrice   *float64  `json:"sell_price,omitempty" db:"sell_price"`   // Harga jual (hanya produk baru)
	CategoryID  *int      `json:"category_id,omitempty" db:"category_id"` // Kategori (hanya produk baru)
	Subtotal    float64   `json:"subtotal" db:"subtotal"`                 // quantity × buy_price
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// PurchaseRequest represents the request body for creating a purchase
// Struct ini untuk menerima request pembelian baru dari frontend
type PurchaseRequest struct {
	SupplierID    *int                  `json:"supplier_id"`
	SupplierName  *string               `json:"supplier_name"`
	PaymentMethod string                `json:"payment_method"` // "cash" | "credit" | "partial"
	PaidAmount    *float64              `json:"paid_amount"`    // opsional, default = total jika cash
	DueDate       *string               `json:"due_date"`       // wajib jika credit/partial
	PaymentNotes  *string               `json:"payment_notes"`
	Notes         *string               `json:"notes"`
	StoreID       int                   `json:"-"` // User StoreID (diisi dari context auth)
	Items         []PurchaseItemRequest `json:"items"`
}

// PurchaseItemRequest represents an item in the purchase request
// Struct ini untuk setiap item dalam request pembelian
type PurchaseItemRequest struct {
	ProductID   *int     `json:"product_id"`   // NULL = produk baru, ada ID = restok
	ProductName *string  `json:"product_name"` // Wajib jika produk baru
	Quantity    int      `json:"quantity"`     // Jumlah beli (harus > 0)
	BuyPrice    float64  `json:"buy_price"`    // Harga beli per unit (harus >= 0)
	SellPrice   *float64 `json:"sell_price"`   // Harga jual (wajib jika produk baru)
	CategoryID  *int     `json:"category_id"`  // Kategori (optional, untuk produk baru)
}
