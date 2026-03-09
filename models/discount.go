package models

import "time"

// DiscountType represents type of discount: PERCENTAGE (0) or FIXED (1)
// Tipe diskon: Persentase (misal 10%) atau Nominal (misal Rp 5.000)
type DiscountType string

const (
	DiscountPercentage DiscountType = "PERCENTAGE"
	DiscountFixed      DiscountType = "FIXED"
)

// Discount represents a promotional code or automatic discount
// Struct ini untuk diskon yang tersedia
type Discount struct {
	ID             int          `json:"id" db:"id"`
	Name           string       `json:"name" db:"name"`
	Type           DiscountType `json:"type" db:"type"`                         // Const: PERCENTAGE / FIXED
	Value          float64      `json:"value" db:"value"`                       // 10.0 (10%) or 5000 (Rp 5,000)
	MinOrderAmount float64      `json:"min_order_amount" db:"min_order_amount"` // Minimal belanja Rp 50,000 baru aktif
	ProductID      *int         `json:"product_id" db:"product_id"`             // Nullable: Jika set, hanya apply ke produk ini
	CategoryID     *int         `json:"category_id" db:"category_id"`           // Nullable: Jika set, hanya apply ke kategori ini
	StartDate      time.Time    `json:"start_date" db:"start_date"`
	EndDate        time.Time    `json:"end_date" db:"end_date"`
	IsActive       bool         `json:"is_active" db:"is_active"`
}

// CalculateDiscount menghitung jumlah potongan berdasarkan total belanja
// Return: jumlah potongan (amount), bukan harga akhir
func (d *Discount) CalculateDiscount(totalAmount float64) float64 {
	// Cek minimal belanja
	if totalAmount < d.MinOrderAmount {
		return 0
	}

	// Cek periode aktif
	now := time.Now()
	if now.Before(d.StartDate) || now.After(d.EndDate) || !d.IsActive {
		return 0
	}

	var discountAmount float64
	if d.Type == DiscountPercentage {
		discountAmount = totalAmount * (d.Value / 100)
	} else {
		discountAmount = d.Value
	}

	// Safety: Diskon tidak boleh melebihi total belanja (gratis oke, minus jangan)
	if discountAmount > totalAmount {
		return totalAmount
	}
	return discountAmount
}
