package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type CashFlowRepository struct {
	db *sql.DB
}

func NewCashFlowRepository(db *sql.DB) *CashFlowRepository {
	return &CashFlowRepository{db: db}
}

// GetSummary menghitung akumulasi seluruh cash in dan cash out pada rentang tanggal yg diberikan
func (r *CashFlowRepository) GetSummary(startDate, endDate time.Time, storeID int) (*models.CashFlowSummary, error) {
	var summary models.CashFlowSummary

	// 1. Initial Balance dari cash_funds (semua waktu, tidak terikat tanggal filter)
	if err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN type = 'out' THEN amount ELSE 0 END), 0)
		FROM cash_funds
		WHERE store_id = $1
	`, storeID).Scan(&summary.InitialBalance); err != nil {
		return nil, fmt.Errorf("GetSummary[1-initial_balance]: %w", err)
	}

	// 2. Cash In (Total Pemasukan dari Transaksi)
	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2 AND store_id = $3
	`, startDate, endDate, storeID).Scan(&summary.CashIn); err != nil {
		return nil, fmt.Errorf("GetSummary[2-cash_in]: %w", err)
	}

	// 3. Cash Out: Purchases
	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(CASE WHEN payment_method = 'cash' THEN total_amount ELSE paid_amount END), 0)
		FROM purchases
		WHERE created_at BETWEEN $1 AND $2 AND (payment_method = 'cash' OR paid_amount > 0) AND store_id = $3
	`, startDate, endDate, storeID).Scan(&summary.CashOutPurchases); err != nil {
		return nil, fmt.Errorf("GetSummary[3-purchases]: %w", err)
	}

	// 4. Cash Out: Payroll (filter via store_id column on payroll table directly)
	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total), 0)
		FROM payroll
		WHERE paid_at BETWEEN $1 AND $2 AND store_id = $3
	`, startDate, endDate, storeID).Scan(&summary.CashOutPayroll); err != nil {
		return nil, fmt.Errorf("GetSummary[4-payroll]: %w", err)
	}

	// 5. Cash Out: Expenses
	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE expense_date BETWEEN $1::date AND $2::date AND store_id = $3
	`, startDate, endDate, storeID).Scan(&summary.CashOutExpenses); err != nil {
		return nil, fmt.Errorf("GetSummary[5-expenses]: %w", err)
	}

	// 6. Cash Out: Debt Payments
	// Catatan: supplier_payables tidak punya store_id — JOIN ke suppliers untuk filter
	if err := r.db.QueryRow(`
		SELECT COALESCE(SUM(pp.amount), 0)
		FROM payable_payments pp
		JOIN supplier_payables sp ON sp.id = pp.payable_id
		JOIN suppliers s ON s.id = sp.supplier_id
		WHERE pp.created_at BETWEEN $1 AND $2
		  AND s.store_id = $3
	`, startDate, endDate, storeID).Scan(&summary.CashOutDebtPayments); err != nil {
		return nil, fmt.Errorf("GetSummary[6-debt_payments]: %w", err)
	}

	// Hitung Aggregasi Akhir
	summary.CashOutTotal = summary.CashOutPurchases + summary.CashOutPayroll + summary.CashOutExpenses + summary.CashOutDebtPayments
	summary.NetCashFlow = summary.CashIn - summary.CashOutTotal
	summary.SaldoKas = summary.InitialBalance + summary.NetCashFlow

	return &summary, nil
}

// GetFunds mengambil semua catatan dana masuk/keluar beserta summary
func (r *CashFlowRepository) GetFunds(page, limit int, storeID int) (*models.CashFundsResponse, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT id, type, amount, to_char(date, 'YYYY-MM-DD'), description, created_by, created_at
		FROM cash_funds
		WHERE store_id = $1
		ORDER BY date DESC, id DESC
		LIMIT $2 OFFSET $3
	`, storeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var funds []models.CashFund
	for rows.Next() {
		var f models.CashFund
		err := rows.Scan(&f.ID, &f.Type, &f.Amount, &f.Date, &f.Description, &f.CreatedBy, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		funds = append(funds, f)
	}
	if funds == nil {
		funds = []models.CashFund{}
	}

	// Hitung summary (total tanpa pagination)
	var summary models.CashFundSummary
	err = r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'out' THEN amount ELSE 0 END), 0)
		FROM cash_funds
		WHERE store_id = $1
	`, storeID).Scan(&summary.TotalIn, &summary.TotalOut)
	if err != nil {
		return nil, err
	}
	summary.Net = summary.TotalIn - summary.TotalOut

	return &models.CashFundsResponse{Data: funds, Summary: summary}, nil
}

// GetInitialBalance menghitung saldo awal dari semua dana masuk - dana keluar
func (r *CashFlowRepository) GetInitialBalance(storeID int) (*models.CashInitialBalance, error) {
	var result models.CashInitialBalance
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'out' THEN amount ELSE 0 END), 0)
		FROM cash_funds
		WHERE store_id = $1
	`, storeID).Scan(&result.Breakdown.TotalIn, &result.Breakdown.TotalOut)
	if err != nil {
		return nil, err
	}
	result.Amount = result.Breakdown.TotalIn - result.Breakdown.TotalOut
	return &result, nil
}

// CreateFund menyimpan catatan dana baru ke tabel cash_funds
func (r *CashFlowRepository) CreateFund(req *models.CashFundRequest, createdBy int, storeID int) (*models.CashFund, error) {
	var f models.CashFund
	err := r.db.QueryRow(`
		INSERT INTO cash_funds (type, amount, date, description, created_by, store_id)
		VALUES ($1, $2, $3::date, $4, $5, $6)
		RETURNING id, type, amount, to_char(date, 'YYYY-MM-DD'), description, created_by, created_at
	`, req.Type, req.Amount, req.Date, req.Description, createdBy, storeID).Scan(
		&f.ID, &f.Type, &f.Amount, &f.Date, &f.Description, &f.CreatedBy, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteFund menghapus satu catatan dana
func (r *CashFlowRepository) DeleteFund(id int, storeID int) error {
	result, err := r.db.Exec("DELETE FROM cash_funds WHERE id = $1 AND store_id = $2", id, storeID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("catatan dana id %d tidak ditemukan", id)
	}
	return nil
}

// GetTrend merangkum pergerakan arus kas seiring waktu (bisa per hari / per bulan)
func (r *CashFlowRepository) GetTrend(startDate, endDate time.Time, format, tzName string, storeID int) (*models.CashFlowTrendResponse, error) {
	query := `
		WITH cash_in AS (
			SELECT 
				TO_CHAR((created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(total_amount) as amount
			FROM transactions
			WHERE created_at BETWEEN $3 AND $4 AND store_id = $5
			GROUP BY period
		),
		cash_out_purchases AS (
			SELECT 
				TO_CHAR((created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(CASE WHEN payment_method = 'cash' THEN total_amount ELSE paid_amount END) as amount
			FROM purchases
			WHERE created_at BETWEEN $3 AND $4 AND (payment_method = 'cash' OR paid_amount > 0) AND store_id = $5
			GROUP BY period
		),
		cash_out_payroll AS (
			SELECT 
				TO_CHAR((paid_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(total) as amount
			FROM payroll
			WHERE paid_at BETWEEN $3 AND $4 AND store_id = $5
			GROUP BY period
		),
		cash_out_expenses AS (
			SELECT 
				TO_CHAR(expense_date, $2) as period,
				SUM(amount) as amount
			FROM expenses
			WHERE expense_date BETWEEN $3 AND $4 AND store_id = $5
			GROUP BY period
		),
		cash_out_debt_payments AS (
			SELECT 
				TO_CHAR((pp.created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(pp.amount) as amount
			FROM payable_payments pp
			JOIN supplier_payables sp ON sp.id = pp.payable_id
			JOIN suppliers s ON s.id = sp.supplier_id
			WHERE pp.created_at BETWEEN $3 AND $4 AND s.store_id = $5
			GROUP BY period
		),
		all_periods AS (
			SELECT period FROM cash_in
			UNION SELECT period FROM cash_out_purchases
			UNION SELECT period FROM cash_out_payroll
			UNION SELECT period FROM cash_out_expenses
			UNION SELECT period FROM cash_out_debt_payments
		)
		SELECT 
			ap.period,
			COALESCE(ci.amount, 0) as cash_in,
			COALESCE(cop.amount, 0) + COALESCE(cpr.amount, 0) + COALESCE(cpe.amount, 0) + COALESCE(cdp.amount, 0) as cash_out
		FROM all_periods ap
		LEFT JOIN cash_in ci ON ap.period = ci.period
		LEFT JOIN cash_out_purchases cop ON ap.period = cop.period
		LEFT JOIN cash_out_payroll cpr ON ap.period = cpr.period
		LEFT JOIN cash_out_expenses cpe ON ap.period = cpe.period
		LEFT JOIN cash_out_debt_payments cdp ON ap.period = cdp.period
		ORDER BY ap.period ASC;
	`

	rows, err := r.db.Query(query, tzName, format, startDate, endDate, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []models.CashFlowTrendData

	for rows.Next() {
		var t models.CashFlowTrendData
		if err := rows.Scan(&t.Period, &t.CashIn, &t.CashOut); err != nil {
			return nil, err
		}
		t.Net = t.CashIn - t.CashOut
		trends = append(trends, t)
	}

	if trends == nil {
		trends = []models.CashFlowTrendData{}
	}

	return &models.CashFlowTrendResponse{Data: trends}, nil
}

// GetLedger returns all cash movements (in/out) in a date range, sorted DESC, with running_balance
func (r *CashFlowRepository) GetLedger(startDate, endDate time.Time, page, limit int, storeID int) (*models.LedgerResponse, error) {
	// 1. Dapatkan initial_balance dari cash_funds (semua dana tanpa filter tanggal)
	var initialBalance float64
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN type = 'out' THEN amount ELSE 0 END), 0)
		FROM cash_funds
		WHERE store_id = $1
	`, storeID).Scan(&initialBalance)
	if err != nil {
		initialBalance = 0
	}

	// Sub-query yang digunakan di semua langkah (reusable CTE)
	allEntriesCTE := `
		all_entries AS (
			-- Cash In: sales
			SELECT created_at, 
			       CONCAT('Penjualan #', id) AS description,
			       'in' AS type,
			       'sale' AS category,
			       total_amount AS amount
			FROM transactions
			WHERE created_at BETWEEN $1 AND $2 AND store_id = $3

			UNION ALL

			-- Cash Out: purchases (cash murni atau ada DP)
			SELECT created_at,
			       CONCAT('Pembelian - ', COALESCE(supplier_name, 'Tanpa Supplier')) AS description,
			       'out' AS type,
			       'purchase' AS category,
			       CASE WHEN payment_method = 'cash' THEN total_amount ELSE paid_amount END AS amount
			FROM purchases
			WHERE created_at BETWEEN $1 AND $2 AND (payment_method = 'cash' OR paid_amount > 0) AND store_id = $3

			UNION ALL

			-- Cash Out: payroll
			SELECT paid_at AS created_at,
			       CONCAT('Gaji - ', COALESCE(e.nama, 'Karyawan')) AS description,
			       'out' AS type,
			       'payroll' AS category,
			       p.total AS amount
			FROM payroll p
			LEFT JOIN employees e ON e.id = p.employee_id
			WHERE paid_at BETWEEN $1 AND $2 AND p.store_id = $3

			UNION ALL

			-- Cash Out: expenses
			SELECT expense_date::timestamptz AS created_at,
			       CONCAT('Pengeluaran - ', description) AS description,
			       'out' AS type,
			       'expense' AS category,
			       amount
			FROM expenses
			WHERE expense_date BETWEEN $1::date AND $2::date AND store_id = $3

			UNION ALL

			-- Cash Out: debt payments (bayar hutang supplier)
			SELECT pp.created_at,
			       CONCAT('Bayar hutang - ', COALESCE(s.name, 'Supplier')) AS description,
			       'out' AS type,
			       'debt_payment' AS category,
			       pp.amount
			FROM payable_payments pp
			JOIN supplier_payables sp ON sp.id = pp.payable_id
			LEFT JOIN suppliers s ON s.id = sp.supplier_id
			WHERE pp.created_at BETWEEN $1 AND $2 AND s.store_id = $3
		)
	`

	// 2. Hitung TOTAL entries dalam rentang tanggal (untuk pagination metadata)
	var totalItems int
	countQuery := fmt.Sprintf(`WITH %s SELECT COUNT(*) FROM all_entries`, allEntriesCTE)
	err = r.db.QueryRow(countQuery, startDate, endDate, storeID).Scan(&totalItems)
	if err != nil {
		return nil, err
	}

	totalPages := (totalItems + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	// 3. Hitung running_balance sampai SEBELUM halaman yang diminta
	offset := (page - 1) * limit
	balanceBeforePage := initialBalance

	if offset > 0 {
		balanceQuery := fmt.Sprintf(`
			WITH %s,
			ordered AS (
				SELECT type, amount, ROW_NUMBER() OVER (ORDER BY created_at ASC) AS rn
				FROM all_entries
			)
			SELECT 
				COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE -amount END), 0)
			FROM ordered
			WHERE rn <= $4
		`, allEntriesCTE)

		var netBefore float64
		err = r.db.QueryRow(balanceQuery, startDate, endDate, storeID, offset).Scan(&netBefore)
		if err != nil {
			return nil, err
		}
		balanceBeforePage += netBefore
	}

	// 4. Ambil HANYA entries untuk halaman ini (ASC, LIMIT + OFFSET)
	dataQuery := fmt.Sprintf(`
		WITH %s
		SELECT 
			created_at::text AS date,
			description,
			type,
			category,
			amount
		FROM all_entries
		ORDER BY created_at ASC
		LIMIT $4 OFFSET $5
	`, allEntriesCTE)

	rows, err := r.db.Query(dataQuery, startDate, endDate, storeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LedgerEntry
	balance := balanceBeforePage

	for rows.Next() {
		var e models.LedgerEntry
		err := rows.Scan(&e.Date, &e.Description, &e.Type, &e.Category, &e.Amount)
		if err != nil {
			return nil, err
		}

		// Update running balance incrementally
		if e.Type == "in" {
			balance += e.Amount
		} else {
			balance -= e.Amount
		}
		e.RunningBalance = balance
		entries = append(entries, e)
	}

	// Reverse the order so the newest entries are at the top (DESC)
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	if entries == nil {
		entries = []models.LedgerEntry{}
	}

	return &models.LedgerResponse{
		Data: entries,
		Pagination: &models.Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}
