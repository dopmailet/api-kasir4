package models

import "time"

// Payroll merepresentasikan tabel payroll
type Payroll struct {
	ID         int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	Periode    *string   `json:"periode,omitempty"`
	GajiPokok  float64   `json:"gaji_pokok"`
	Bonus      float64   `json:"bonus"`
	Potongan   float64   `json:"potongan"`
	Total      float64   `json:"total"`
	Catatan    *string   `json:"catatan,omitempty"`
	PaidAt     time.Time `json:"paid_at"`
	CreatedBy  *int      `json:"created_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Join detail
	EmployeeNama *string `json:"employee_nama,omitempty"`
}

// CreatePayrollRequest DTO untuk POST /api/payroll
type CreatePayrollRequest struct {
	EmployeeID int      `json:"employee_id" validate:"required"`
	Periode    *string  `json:"periode,omitempty"`
	GajiPokok  float64  `json:"gaji_pokok" validate:"required,min=0"`
	Bonus      *float64 `json:"bonus,omitempty"`
	Potongan   *float64 `json:"potongan,omitempty"`
	Catatan    *string  `json:"catatan,omitempty"`
	// Total akan dihitung oleh backend
	// CreatedBy akan diambil dari user JWT Auth Middleware
}

// UpdatePayrollRequest DTO untuk PUT /api/payroll/:id
type UpdatePayrollRequest struct {
	Periode   *string  `json:"periode"`
	GajiPokok *float64 `json:"gaji_pokok"`
	Bonus     *float64 `json:"bonus"`
	Potongan  *float64 `json:"potongan"`
	Catatan   *string  `json:"catatan"`
}

// -------- Payroll Report structs --------

// PayrollReport struct utama untuk /api/payroll/report
type PayrollReport struct {
	TotalGaji        float64                `json:"total_gaji"`
	TotalBonus       float64                `json:"total_bonus"`
	TotalPotongan    float64                `json:"total_potongan"`
	TotalDibayar     float64                `json:"total_dibayar"`
	JumlahPembayaran int                    `json:"jumlah_pembayaran"`
	PerKaryawan      []PayrollSummaryPerEmp `json:"per_karyawan"`
}

// PayrollSummaryPerEmp struct rincian per karyawan di report
type PayrollSummaryPerEmp struct {
	EmployeeID       int     `json:"employee_id"`
	Nama             string  `json:"nama"`
	Posisi           string  `json:"posisi"`
	TotalDibayar     float64 `json:"total_dibayar"`
	JumlahPembayaran int     `json:"jumlah_pembayaran"`
}
