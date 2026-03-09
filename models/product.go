package models

// Product adalah struct untuk data produk
type Product struct {
	ID                   int       `json:"id" db:"id"`
	Nama                 string    `json:"nama" db:"nama"`
	Harga                float64   `json:"harga" db:"harga"`                     // Harga jual
	HargaBeli            *float64  `json:"harga_beli,omitempty" db:"harga_beli"` // Harga beli/modal (nullable)
	Stok                 int       `json:"stok" db:"stok"`
	Barcode              *string   `json:"barcode,omitempty" db:"barcode"`                               // Barcode produk (nullable, unique)
	CategoryID           *int      `json:"category_id,omitempty" db:"category_id"`                       // Foreign key ke categories (nullable)
	DefaultDiscountType  *string   `json:"default_discount_type,omitempty" db:"default_discount_type"`   // "percentage" atau "fixed" (nullable)
	DefaultDiscountValue *float64  `json:"default_discount_value,omitempty" db:"default_discount_value"` // Nilai diskon default (nullable)
	IsFeatured           bool      `json:"is_featured" db:"is_featured"`                                 // Flag fitur unggulan
	CreatedBy            *int      `json:"created_by,omitempty" db:"created_by"`                         // User ID yang menambahkan produk
	Category             *Category `json:"category,omitempty" db:"-"`                                    // Untuk hasil JOIN (tidak disimpan di DB)
	Margin               *float64  `json:"margin,omitempty" db:"-"`                                      // Margin keuntungan % (calculated field)
}

// CalculateMargin menghitung margin keuntungan dalam persen
// Formula: ((harga_jual - harga_beli) / harga_jual) * 100
func (p *Product) CalculateMargin() *float64 {
	if p.HargaBeli == nil || p.Harga == 0 {
		return nil
	}

	margin := ((p.Harga - *p.HargaBeli) / p.Harga) * 100
	return &margin
}

// GetProfit menghitung profit per unit
func (p *Product) GetProfit() *float64 {
	if p.HargaBeli == nil {
		return nil
	}

	profit := p.Harga - *p.HargaBeli
	return &profit
}

// ValidatePrice memvalidasi harga jual dan harga beli
func (p *Product) ValidatePrice() error {
	if p.Harga < 0 {
		return ErrInvalidPrice
	}

	if p.HargaBeli != nil && *p.HargaBeli < 0 {
		return ErrInvalidPurchasePrice
	}

	// Warning jika harga beli >= harga jual (dijual rugi)
	if p.HargaBeli != nil && *p.HargaBeli >= p.Harga {
		return ErrNegativeMargin
	}

	return nil
}

// TableName untuk override nama table di database
func (Product) TableName() string {
	return "products"
}
