package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type PayrollRepository struct {
	db *sql.DB
}

func NewPayrollRepository(db *sql.DB) *PayrollRepository {
	return &PayrollRepository{db: db}
}

func (r *PayrollRepository) GetAll(employeeID int, startDate, endDate time.Time, offset, limit int, storeID int, tzName string) ([]models.Payroll, int, error) {
	// Base Query
	query := `SELECT p.id, p.employee_id, e.nama as employee_nama, p.periode, p.gaji_pokok, 
	                 p.bonus, p.potongan, p.total, p.catatan, p.paid_at, p.created_by, p.store_id, p.created_at, p.updated_at 
	          FROM payroll p
	          JOIN employees e ON p.employee_id = e.id
	          WHERE p.store_id = $1`
	countQuery := `SELECT COUNT(*) FROM payroll p WHERE p.store_id = $1`

	args := []interface{}{storeID}
	argCount := 2

	if employeeID > 0 {
		query += ` AND p.employee_id = $` + itoa(argCount)
		countQuery += ` AND p.employee_id = $` + itoa(argCount)
		args = append(args, employeeID)
		argCount++
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		startStr := startDate.Format("2006-01-02")
		endStr := endDate.Format("2006-01-02")
		query += ` AND p.paid_at >= ($` + itoa(argCount) + `::date AT TIME ZONE $` + itoa(argCount+2) + `) AND p.paid_at < (($` + itoa(argCount+1) + `::date + INTERVAL '1 day') AT TIME ZONE $` + itoa(argCount+2) + `)`
		countQuery += ` AND p.paid_at >= ($` + itoa(argCount) + `::date AT TIME ZONE $` + itoa(argCount+2) + `) AND p.paid_at < (($` + itoa(argCount+1) + `::date + INTERVAL '1 day') AT TIME ZONE $` + itoa(argCount+2) + `)`
		args = append(args, startStr, endStr, tzName)
		argCount += 3
	}

	// Hitung total
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Tambahkan sorting & pagination
	query += ` ORDER BY p.paid_at DESC LIMIT $` + itoa(argCount) + ` OFFSET $` + itoa(argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payrolls []models.Payroll
	for rows.Next() {
		var p models.Payroll
		err := rows.Scan(
			&p.ID, &p.EmployeeID, &p.EmployeeNama, &p.Periode, &p.GajiPokok,
			&p.Bonus, &p.Potongan, &p.Total, &p.Catatan, &p.PaidAt,
			&p.CreatedBy, &p.StoreID, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		payrolls = append(payrolls, p)
	}

	return payrolls, total, nil
}

// GetByID mengambil detail satu record payroll
func (r *PayrollRepository) GetByID(id int, storeID int) (*models.Payroll, error) {
	query := `SELECT p.id, p.employee_id, e.nama, p.periode, p.gaji_pokok, 
	                 p.bonus, p.potongan, p.total, p.catatan, p.paid_at, p.created_by, p.store_id, p.created_at, p.updated_at 
	          FROM payroll p
	          JOIN employees e ON p.employee_id = e.id
	          WHERE p.id = $1 AND p.store_id = $2`

	var p models.Payroll
	err := r.db.QueryRow(query, id, storeID).Scan(
		&p.ID, &p.EmployeeID, &p.EmployeeNama, &p.Periode, &p.GajiPokok,
		&p.Bonus, &p.Potongan, &p.Total, &p.Catatan, &p.PaidAt,
		&p.CreatedBy, &p.StoreID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create menambahkan pembayaran gaji baru
func (r *PayrollRepository) Create(p *models.Payroll) error {
	query := `INSERT INTO payroll (employee_id, periode, gaji_pokok, bonus, potongan, total, catatan, created_by, store_id) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	          RETURNING id, paid_at, created_at, updated_at`

	return r.db.QueryRow(query,
		p.EmployeeID, p.Periode, p.GajiPokok, p.Bonus, p.Potongan, p.Total, p.Catatan, p.CreatedBy, p.StoreID,
	).Scan(&p.ID, &p.PaidAt, &p.CreatedAt, &p.UpdatedAt)
}

// Update memodifikasi record (hanya <= 24 jam)
func (r *PayrollRepository) Update(p *models.Payroll) error {
	query := `UPDATE payroll 
	          SET periode = $1, gaji_pokok = $2, bonus = $3, potongan = $4, total = $5, catatan = $6
	          WHERE id = $7 AND store_id = $8 RETURNING updated_at`

	return r.db.QueryRow(query,
		p.Periode, p.GajiPokok, p.Bonus, p.Potongan, p.Total, p.Catatan, p.ID, p.StoreID,
	).Scan(&p.UpdatedAt)
}

// Delete menghapus record (hanya <= 24 jam)
func (r *PayrollRepository) Delete(id int, storeID int) error {
	query := `DELETE FROM payroll WHERE id = $1 AND store_id = $2`
	_, err := r.db.Exec(query, id, storeID)
	return err
}

// GetReport melakukan agregasi gaji (Timezone-Aware)
func (r *PayrollRepository) GetReport(startDate, endDate time.Time, tzName string, storeID int) (*models.PayrollReport, error) {
	if tzName == "" {
		tzName = "Asia/Jakarta" // Default
	}

	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	// 1. Ambil Summary Global dengan Timezone Aware
	// Menggunakan pola yang persis dengan dashboard (AT TIME ZONE 'UTC' AT TIME ZONE $1) diganti dengan index-friendly
	summaryQuery := `
		SELECT 
			COALESCE(SUM(gaji_pokok), 0) as total_gaji,
			COALESCE(SUM(bonus), 0) as total_bonus,
			COALESCE(SUM(potongan), 0) as total_potongan,
			COALESCE(SUM(total), 0) as total_dibayar,
			COUNT(id) as jumlah_pembayaran
		FROM payroll
		WHERE store_id = $1
		  AND paid_at >= ($3::date AT TIME ZONE $2)
		  AND paid_at < (($4::date + INTERVAL '1 day') AT TIME ZONE $2)
	`

	report := &models.PayrollReport{}
	err := r.db.QueryRow(summaryQuery, storeID, tzName, startStr, endStr).Scan(
		&report.TotalGaji, &report.TotalBonus, &report.TotalPotongan, &report.TotalDibayar, &report.JumlahPembayaran,
	)
	if err != nil {
		return nil, err
	}

	// 2. Ambil rincian per karyawan
	perEmpQuery := `
		SELECT 
			e.id as employee_id,
			e.nama,
			e.posisi,
			COALESCE(SUM(p.total), 0) as total_dibayar,
			COUNT(p.id) as jumlah_pembayaran
		FROM employees e
		JOIN payroll p ON e.id = p.employee_id
		WHERE p.store_id = $1
		  AND p.paid_at >= ($3::date AT TIME ZONE $2)
		  AND p.paid_at < (($4::date + INTERVAL '1 day') AT TIME ZONE $2)
		GROUP BY e.id, e.nama, e.posisi
		ORDER BY total_dibayar DESC
	`

	rows, err := r.db.Query(perEmpQuery, storeID, tzName, startStr, endStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pe models.PayrollSummaryPerEmp
		err := rows.Scan(&pe.EmployeeID, &pe.Nama, &pe.Posisi, &pe.TotalDibayar, &pe.JumlahPembayaran)
		if err != nil {
			return nil, err
		}
		report.PerKaryawan = append(report.PerKaryawan, pe)
	}

	// Jika tidak ada data
	if report.PerKaryawan == nil {
		report.PerKaryawan = []models.PayrollSummaryPerEmp{}
	}

	return report, nil
}

// Utility: integer to string formatter untuk query args
func itoa(i int) string {
	importStrConv := func(x int) string {
		var s string
		if x == 0 {
			return "0"
		}
		for x > 0 {
			s = string(rune('0'+(x%10))) + s
			x /= 10
		}
		return s
	}
	return importStrConv(i)
}
