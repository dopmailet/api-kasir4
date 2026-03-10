package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
	"time"
)

type SupplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

// ─── SUPPLIER CRUD ───

// GetDebtSummary returns total payable and count of suppliers with outstanding debt
func (r *SupplierRepository) GetDebtSummary() (*models.SupplierDebtSummary, error) {
	var summary models.SupplierDebtSummary
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(amount - paid_amount), 0) AS total_payable,
			COUNT(DISTINCT supplier_id) AS total_suppliers_with_debt
		FROM supplier_payables
		WHERE status != 'paid'
	`).Scan(&summary.TotalPayable, &summary.TotalSuppliersWithDebt)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *SupplierRepository) GetAll(search string, isActive *bool) ([]models.Supplier, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(
			"(s.name ILIKE $%d OR s.contact_person ILIKE $%d)", argIdx, argIdx,
		))
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if isActive != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("s.is_active = $%d", argIdx))
		args = append(args, *isActive)
		argIdx++
	}

	whereStmt := ""
	if len(whereClauses) > 0 {
		whereStmt = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			s.id, s.name, s.contact_person, s.phone, s.email, s.address, s.notes,
			s.is_active, s.total_purchases, s.total_spent,
			COALESCE(SUM(sp.amount - sp.paid_amount), 0) AS total_payable,
			s.created_at, s.updated_at
		FROM suppliers s
		LEFT JOIN supplier_payables sp ON sp.supplier_id = s.id AND sp.status != 'paid'
		%s
		GROUP BY s.id
		ORDER BY s.name ASC
	`, whereStmt)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(
			&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Email, &s.Address, &s.Notes,
			&s.IsActive, &s.TotalPurchases, &s.TotalSpent, &s.TotalPayable,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, nil
}

func (r *SupplierRepository) GetByID(id int) (*models.Supplier, error) {
	query := `
		SELECT 
			s.id, s.name, s.contact_person, s.phone, s.email, s.address, s.notes,
			s.is_active, s.total_purchases, s.total_spent,
			COALESCE(SUM(sp.amount - sp.paid_amount), 0) AS total_payable,
			s.created_at, s.updated_at
		FROM suppliers s
		LEFT JOIN supplier_payables sp ON sp.supplier_id = s.id AND sp.status != 'paid'
		WHERE s.id = $1
		GROUP BY s.id
	`
	var s models.Supplier
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Email, &s.Address, &s.Notes,
		&s.IsActive, &s.TotalPurchases, &s.TotalSpent, &s.TotalPayable,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("supplier id %d tidak ditemukan", id)
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepository) Create(req *models.CreateSupplierRequest) (*models.Supplier, error) {
	query := `
		INSERT INTO suppliers (name, contact_person, phone, email, address, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, contact_person, phone, email, address, notes,
		          is_active, total_purchases, total_spent, 0::numeric, created_at, updated_at
	`
	var s models.Supplier
	err := r.db.QueryRow(query,
		req.Name, req.ContactPerson, req.Phone, req.Email, req.Address, req.Notes,
	).Scan(
		&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Email, &s.Address, &s.Notes,
		&s.IsActive, &s.TotalPurchases, &s.TotalSpent, &s.TotalPayable,
		&s.CreatedAt, &s.UpdatedAt,
	)
	return &s, err
}

func (r *SupplierRepository) Update(id int, req *models.UpdateSupplierRequest) (*models.Supplier, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.ContactPerson != nil {
		existing.ContactPerson = req.ContactPerson
	}
	if req.Phone != nil {
		existing.Phone = req.Phone
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.Address != nil {
		existing.Address = req.Address
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	query := `
		UPDATE suppliers
		SET name = $1, contact_person = $2, phone = $3, email = $4,
		    address = $5, notes = $6, is_active = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`
	err = r.db.QueryRow(query,
		existing.Name, existing.ContactPerson, existing.Phone, existing.Email,
		existing.Address, existing.Notes, existing.IsActive, id,
	).Scan(&existing.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *SupplierRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM suppliers WHERE id = $1", id)
	return err
}

// ─── PAYABLES ───

func (r *SupplierRepository) GetPayablesBySupplier(supplierID int) ([]models.SupplierPayable, error) {
	query := `
		SELECT id, supplier_id, purchase_id, amount, paid_amount, status, 
		       to_char(due_date, 'YYYY-MM-DD'), notes, created_at, updated_at
		FROM supplier_payables
		WHERE supplier_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payables []models.SupplierPayable
	for rows.Next() {
		var p models.SupplierPayable
		var dueDate sql.NullString
		err := rows.Scan(
			&p.ID, &p.SupplierID, &p.PurchaseID, &p.Amount, &p.PaidAmount,
			&p.Status, &dueDate, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			p.DueDate = &dueDate.String
		}
		payables = append(payables, p)
	}
	return payables, nil
}

func (r *SupplierRepository) GetPayableByID(payableID int) (*models.SupplierPayable, error) {
	query := `
		SELECT id, supplier_id, purchase_id, amount, paid_amount, status, 
		       to_char(due_date, 'YYYY-MM-DD'), notes, created_at, updated_at
		FROM supplier_payables WHERE id = $1
	`
	var p models.SupplierPayable
	var dueDate sql.NullString
	err := r.db.QueryRow(query, payableID).Scan(
		&p.ID, &p.SupplierID, &p.PurchaseID, &p.Amount, &p.PaidAmount,
		&p.Status, &dueDate, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payable id %d tidak ditemukan", payableID)
	}
	if err != nil {
		return nil, err
	}
	if dueDate.Valid {
		p.DueDate = &dueDate.String
	}
	return &p, nil
}

func (r *SupplierRepository) CreatePayable(supplierID int, req *models.CreatePayableRequest) (*models.SupplierPayable, error) {
	// Hitung status awal
	status := req.Status
	if status == "" {
		if req.PaidAmount >= req.Amount {
			status = "paid"
		} else if req.PaidAmount > 0 {
			status = "partial"
		} else {
			status = "unpaid"
		}
	}

	query := `
		INSERT INTO supplier_payables (supplier_id, purchase_id, amount, paid_amount, status, due_date, notes)
		VALUES ($1, $2, $3, $4, $5, $6::date, $7)
		RETURNING id, supplier_id, purchase_id, amount, paid_amount, status, 
		          to_char(due_date, 'YYYY-MM-DD'), notes, created_at, updated_at
	`
	var p models.SupplierPayable
	var dueDate sql.NullString
	err := r.db.QueryRow(query,
		supplierID, req.PurchaseID, req.Amount, req.PaidAmount, status, req.DueDate, req.Notes,
	).Scan(
		&p.ID, &p.SupplierID, &p.PurchaseID, &p.Amount, &p.PaidAmount,
		&p.Status, &dueDate, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if dueDate.Valid {
		p.DueDate = &dueDate.String
	}

	// Update total_payable di supplier
	r.recalcSupplierPayable(supplierID)
	return &p, nil
}

func (r *SupplierRepository) UpdatePayable(payableID int, req *models.UpdatePayableRequest) (*models.SupplierPayable, error) {
	existing, err := r.GetPayableByID(payableID)
	if err != nil {
		return nil, err
	}

	if req.Amount != nil {
		existing.Amount = *req.Amount
	}
	if req.DueDate != nil {
		existing.DueDate = req.DueDate
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	// Recalculate status
	status := "unpaid"
	if existing.PaidAmount >= existing.Amount {
		status = "paid"
	} else if existing.PaidAmount > 0 {
		status = "partial"
	}

	query := `
		UPDATE supplier_payables
		SET amount = $1, due_date = $2::date, notes = $3, status = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err = r.db.QueryRow(query,
		existing.Amount, existing.DueDate, existing.Notes, status, payableID,
	).Scan(&existing.UpdatedAt)
	if err != nil {
		return nil, err
	}
	existing.Status = status
	r.recalcSupplierPayable(existing.SupplierID)
	return existing, nil
}

// ─── PAYMENTS ───

func (r *SupplierRepository) GetPaymentsByPayable(payableID int) ([]models.PayablePayment, error) {
	query := `
		SELECT id, payable_id, amount, to_char(payment_date, 'YYYY-MM-DD'), notes, created_at
		FROM payable_payments
		WHERE payable_id = $1
		ORDER BY payment_date DESC
	`
	rows, err := r.db.Query(query, payableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PayablePayment
	for rows.Next() {
		var p models.PayablePayment
		err := rows.Scan(&p.ID, &p.PayableID, &p.Amount, &p.PaymentDate, &p.Notes, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *SupplierRepository) CreatePayment(payableID int, req *models.CreatePaymentRequest) (*models.PayablePayment, error) {
	// Ambil info payable
	payable, err := r.GetPayableByID(payableID)
	if err != nil {
		return nil, err
	}

	remaining := payable.Amount - payable.PaidAmount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("jumlah pembayaran harus lebih dari 0")
	}
	if req.Amount > remaining {
		return nil, fmt.Errorf("jumlah pembayaran (%.2f) melebihi sisa hutang (%.2f)", req.Amount, remaining)
	}

	// Validasi format tanggal
	_, errDate := time.Parse("2006-01-02", req.PaymentDate)
	if errDate != nil {
		return nil, fmt.Errorf("format payment_date tidak valid, gunakan YYYY-MM-DD")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert payment
	var p models.PayablePayment
	err = tx.QueryRow(`
		INSERT INTO payable_payments (payable_id, amount, payment_date, notes)
		VALUES ($1, $2, $3::date, $4)
		RETURNING id, payable_id, amount, to_char(payment_date, 'YYYY-MM-DD'), notes, created_at
	`, payableID, req.Amount, req.PaymentDate, req.Notes).Scan(
		&p.ID, &p.PayableID, &p.Amount, &p.PaymentDate, &p.Notes, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Update paid_amount dan status di payable
	newPaid := payable.PaidAmount + req.Amount
	newStatus := "partial"
	if newPaid >= payable.Amount {
		newStatus = "paid"
	}

	_, err = tx.Exec(`
		UPDATE supplier_payables
		SET paid_amount = $1, status = $2, updated_at = NOW()
		WHERE id = $3
	`, newPaid, newStatus, payableID)
	if err != nil {
		return nil, err
	}

	// Recalculate total_payable supplier
	_, err = tx.Exec(`
		UPDATE suppliers
		SET total_payable = (
			SELECT COALESCE(SUM(amount - paid_amount), 0)
			FROM supplier_payables
			WHERE supplier_id = (SELECT supplier_id FROM supplier_payables WHERE id = $1)
			AND status != 'paid'
		),
		updated_at = NOW()
		WHERE id = (SELECT supplier_id FROM supplier_payables WHERE id = $1)
	`, payableID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// recalcSupplierPayable recalculates and updates total_payable for a supplier
func (r *SupplierRepository) recalcSupplierPayable(supplierID int) {
	r.db.Exec(`
		UPDATE suppliers
		SET total_payable = (
			SELECT COALESCE(SUM(amount - paid_amount), 0)
			FROM supplier_payables
			WHERE supplier_id = $1 AND status != 'paid'
		), updated_at = NOW()
		WHERE id = $1
	`, supplierID)
}

// GetAllPayablePayments gets all payments with optional date filters and supplier names.
func (r *SupplierRepository) GetAllPayablePayments(startDate, endDate string) ([]models.PayablePaymentWithSupplier, error) {
	query := `
		SELECT 
			pp.id, pp.payable_id, pp.amount, to_char(pp.payment_date, 'YYYY-MM-DD'), pp.notes, pp.created_at,
			s.name as supplier_name,
			sp.amount as payable_amount
		FROM payable_payments pp
		JOIN supplier_payables sp ON pp.payable_id = sp.id
		JOIN suppliers s ON sp.supplier_id = s.id
	`
	whereClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if startDate != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("pp.payment_date >= $%d", argIdx))
		args = append(args, startDate)
		argIdx++
	}
	if endDate != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("pp.payment_date <= $%d", argIdx))
		args = append(args, endDate)
		argIdx++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY pp.payment_date DESC, pp.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PayablePaymentWithSupplier
	for rows.Next() {
		var p models.PayablePaymentWithSupplier
		err := rows.Scan(&p.ID, &p.PayableID, &p.Amount, &p.PaymentDate, &p.Notes, &p.CreatedAt, &p.SupplierName, &p.PayableAmount)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}
