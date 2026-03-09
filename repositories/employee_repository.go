package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// GetAll mengambil semua karyawan, opsional filter berdasarkan status aktif
func (r *EmployeeRepository) GetAll(aktif *bool) ([]models.Employee, error) {
	query := `SELECT id, nama, posisi, gaji_pokok, no_hp, alamat, tanggal_masuk, aktif, user_id, created_at, updated_at 
	          FROM employees WHERE 1=1`
	var args []interface{}

	if aktif != nil {
		query += ` AND aktif = $1`
		args = append(args, *aktif)
	}

	query += ` ORDER BY nama ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		err := rows.Scan(
			&e.ID, &e.Nama, &e.Posisi, &e.GajiPokok, &e.NoHp, &e.Alamat,
			&e.TanggalMasuk, &e.Aktif, &e.UserID, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	return employees, nil
}

// GetByID mengambil detail karyawan berdasarkan ID, termasuk up to 5 payroll history
func (r *EmployeeRepository) GetByID(id int) (*models.Employee, error) {
	query := `SELECT id, nama, posisi, gaji_pokok, no_hp, alamat, tanggal_masuk, aktif, user_id, created_at, updated_at 
	          FROM employees WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var e models.Employee
	err := row.Scan(
		&e.ID, &e.Nama, &e.Posisi, &e.GajiPokok, &e.NoHp, &e.Alamat,
		&e.TanggalMasuk, &e.Aktif, &e.UserID, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch up to 5 recent payrolls
	queryPayrolls := `SELECT id, employee_id, periode, gaji_pokok, bonus, potongan, total, catatan, paid_at, created_by, created_at, updated_at 
	                  FROM payroll WHERE employee_id = $1 ORDER BY paid_at DESC LIMIT 5`

	rows, err := r.db.Query(queryPayrolls, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p models.Payroll
			rows.Scan(
				&p.ID, &p.EmployeeID, &p.Periode, &p.GajiPokok, &p.Bonus, &p.Potongan,
				&p.Total, &p.Catatan, &p.PaidAt, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
			)
			e.RecentPayrolls = append(e.RecentPayrolls, p)
		}
	}

	return &e, nil
}

// Create menambahkan karyawan baru
func (r *EmployeeRepository) Create(e *models.Employee) error {
	query := `INSERT INTO employees (nama, posisi, gaji_pokok, no_hp, alamat, tanggal_masuk, user_id) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, aktif, created_at, updated_at`

	return r.db.QueryRow(
		query, e.Nama, e.Posisi, e.GajiPokok, e.NoHp, e.Alamat, e.TanggalMasuk, e.UserID,
	).Scan(&e.ID, &e.Aktif, &e.CreatedAt, &e.UpdatedAt)
}

// Update memperbarui data karyawan
func (r *EmployeeRepository) Update(e *models.Employee) error {
	query := `UPDATE employees 
	          SET nama = $1, posisi = $2, gaji_pokok = $3, no_hp = $4, alamat = $5, tanggal_masuk = $6, user_id = $7, aktif = $8
	          WHERE id = $9
	          RETURNING id, nama, posisi, gaji_pokok, no_hp, alamat, tanggal_masuk, aktif, user_id, created_at, updated_at`

	return r.db.QueryRow(
		query, e.Nama, e.Posisi, e.GajiPokok, e.NoHp, e.Alamat, e.TanggalMasuk, e.UserID, e.Aktif, e.ID,
	).Scan(
		&e.ID, &e.Nama, &e.Posisi, &e.GajiPokok, &e.NoHp, &e.Alamat,
		&e.TanggalMasuk, &e.Aktif, &e.UserID, &e.CreatedAt, &e.UpdatedAt,
	)
}

// SoftDelete melakukan nonaktif pada karyawan (set aktif = false)
func (r *EmployeeRepository) SoftDelete(id int) error {
	query := `UPDATE employees SET aktif = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
