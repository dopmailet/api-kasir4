package models

import "time"

// Employee merepresentasikan tabel employees
type Employee struct {
	ID           int        `json:"id"`
	Nama         string     `json:"nama"`
	Posisi       string     `json:"posisi"`
	GajiPokok    float64    `json:"gaji_pokok"`
	NoHp         *string    `json:"no_hp"`
	Alamat       *string    `json:"alamat"`
	TanggalMasuk *time.Time `json:"tanggal_masuk"`
	Aktif        bool       `json:"aktif"`
	UserID       *int       `json:"user_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relasi ke tabel payroll (opsional, untuk endpoint detail Employee)
	RecentPayrolls []Payroll `json:"recent_payrolls,omitempty"`
}

// CreateEmployeeRequest DTO
type CreateEmployeeRequest struct {
	Nama         string  `json:"nama" validate:"required"`
	Posisi       string  `json:"posisi" validate:"required"`
	GajiPokok    float64 `json:"gaji_pokok" validate:"required,min=0"`
	NoHp         *string `json:"no_hp,omitempty"`
	Alamat       *string `json:"alamat,omitempty"`
	TanggalMasuk *string `json:"tanggal_masuk,omitempty"` // format YYYY-MM-DD
	UserID       *int    `json:"user_id,omitempty"`
}

// UpdateEmployeeRequest DTO
type UpdateEmployeeRequest struct {
	Nama         string   `json:"nama"`
	Posisi       string   `json:"posisi"`
	GajiPokok    *float64 `json:"gaji_pokok"`
	NoHp         *string  `json:"no_hp"`
	Alamat       *string  `json:"alamat"`
	TanggalMasuk *string  `json:"tanggal_masuk"`
	UserID       *int     `json:"user_id"`
	Aktif        *bool    `json:"aktif"`
}
