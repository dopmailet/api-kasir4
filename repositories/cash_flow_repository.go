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
func (r *CashFlowRepository) GetSummary(startDate, endDate time.Time) (*models.CashFlowSummary, error) {
	var summary models.CashFlowSummary

	// 1. Initial Balance dari cash_funds (semua waktu, tidak terikat tanggal filter)
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN type = 'out' THEN amount ELSE 0 END), 0)
		FROM cash_funds
	`).Scan(&summary.InitialBalance)
	if err != nil {
		return nil, err
	}

	// 2. Cash In (Total Pemasukan dari Transaksi)
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&summary.CashIn)
	if err != nil {
		return nil, err
	}

	// 3. Cash Out: Purchases
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0)
		FROM purchases
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&summary.CashOutPurchases)
	if err != nil {
		return nil, err
	}

	// 4. Cash Out: Payroll
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(total), 0)
		FROM payroll
		WHERE paid_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&summary.CashOutPayroll)
	if err != nil {
		return nil, err
	}

	// 5. Cash Out: Expenses
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE expense_date BETWEEN $1::date AND $2::date
	`, startDate, endDate).Scan(&summary.CashOutExpenses)
	if err != nil {
		return nil, err
	}

	// 6. Cash Out: Debt Payments
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM payable_payments
		WHERE created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&summary.CashOutDebtPayments)
	if err != nil {
		return nil, err
	}

	// Hitung Aggregasi Akhir
	summary.CashOutTotal = summary.CashOutPurchases + summary.CashOutPayroll + summary.CashOutExpenses + summary.CashOutDebtPayments
	summary.NetCashFlow = summary.CashIn - summary.CashOutTotal
	summary.SaldoKas = summary.InitialBalance + summary.NetCashFlow

	return &summary, nil
}

// GetFunds mengambil semua catatan dana masuk/keluar beserta summary
func (r *CashFlowRepository) GetFunds(page, limit int) (*models.CashFundsResponse, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT id, type, amount, to_char(date, 'YYYY-MM-DD'), description, created_by, created_at
		FROM cash_funds
		ORDER BY date DESC, id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
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
	`).Scan(&summary.TotalIn, &summary.TotalOut)
	if err != nil {
		return nil, err
	}
	summary.Net = summary.TotalIn - summary.TotalOut

	return &models.CashFundsResponse{Data: funds, Summary: summary}, nil
}

// GetInitialBalance menghitung saldo awal dari semua dana masuk - dana keluar
func (r *CashFlowRepository) GetInitialBalance() (*models.CashInitialBalance, error) {
	var result models.CashInitialBalance
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'in' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'out' THEN amount ELSE 0 END), 0)
		FROM cash_funds
	`).Scan(&result.Breakdown.TotalIn, &result.Breakdown.TotalOut)
	if err != nil {
		return nil, err
	}
	result.Amount = result.Breakdown.TotalIn - result.Breakdown.TotalOut
	return &result, nil
}

// CreateFund menyimpan catatan dana baru ke tabel cash_funds
func (r *CashFlowRepository) CreateFund(req *models.CashFundRequest, createdBy int) (*models.CashFund, error) {
	var f models.CashFund
	err := r.db.QueryRow(`
		INSERT INTO cash_funds (type, amount, date, description, created_by)
		VALUES ($1, $2, $3::date, $4, $5)
		RETURNING id, type, amount, to_char(date, 'YYYY-MM-DD'), description, created_by, created_at
	`, req.Type, req.Amount, req.Date, req.Description, createdBy).Scan(
		&f.ID, &f.Type, &f.Amount, &f.Date, &f.Description, &f.CreatedBy, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteFund menghapus satu catatan dana
func (r *CashFlowRepository) DeleteFund(id int) error {
	result, err := r.db.Exec("DELETE FROM cash_funds WHERE id = $1", id)
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
// tzName: nama timezone untuk mapping timestamp UTC ke regional user.
// format: "YYYY-MM-DD" untuk daily atau "YYYY-MM" untuk monthly
func (r *CashFlowRepository) GetTrend(startDate, endDate time.Time, format, tzName string) (*models.CashFlowTrendResponse, error) {
	// CTE (Common Table Expression) untuk menggabungkan Cash In (transactions)
	// dan Cash Out (purchases + payroll + expenses) pada timezone specifik lalu group by Period format.
	query := `
		WITH cash_in AS (
			SELECT 
				TO_CHAR((created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(total_amount) as amount
			FROM transactions
			WHERE created_at BETWEEN $3 AND $4
			GROUP BY period
		),
		cash_out_purchases AS (
			SELECT 
				TO_CHAR((created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(total_amount) as amount
			FROM purchases
			WHERE created_at BETWEEN $3 AND $4
			GROUP BY period
		),
		cash_out_payroll AS (
			SELECT 
				TO_CHAR((paid_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
				SUM(total) as amount
			FROM payroll
			WHERE paid_at BETWEEN $3 AND $4
			GROUP BY period
		),
		cash_out_expenses AS (
			SELECT 
				TO_CHAR(expense_date, $2) as period,
				SUM(amount) as amount
			FROM expenses
			WHERE expense_date BETWEEN $3 AND $4
			GROUP BY period
		),
		all_periods AS (
			SELECT period FROM cash_in
			UNION SELECT period FROM cash_out_purchases
			UNION SELECT period FROM cash_out_payroll
			UNION SELECT period FROM cash_out_expenses
		)
		SELECT 
			ap.period,
			COALESCE(ci.amount, 0) as cash_in,
			COALESCE(cop.amount, 0) + COALESCE(cpr.amount, 0) + COALESCE(cpe.amount, 0) as cash_out
		FROM all_periods ap
		LEFT JOIN cash_in ci ON ap.period = ci.period
		LEFT JOIN cash_out_purchases cop ON ap.period = cop.period
		LEFT JOIN cash_out_payroll cpr ON ap.period = cpr.period
		LEFT JOIN cash_out_expenses cpe ON ap.period = cpe.period
		ORDER BY ap.period ASC;
	`

	rows, err := r.db.Query(query, tzName, format, startDate, endDate)
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
func (r *CashFlowRepository) GetLedger(startDate, endDate time.Time) (*models.LedgerResponse, error) {
	query := `
		SELECT 
			created_at::text AS date,
			description,
			type,
			category,
			amount
		FROM (
			-- Cash In: sales
			SELECT created_at, 
			       CONCAT('Penjualan #', id) AS description,
			       'in' AS type,
			       'sale' AS category,
			       total_amount AS amount
			FROM transactions
			WHERE created_at BETWEEN $1 AND $2

			UNION ALL

			-- Cash Out: purchases (non-credit)
			SELECT created_at,
			       CONCAT('Pembelian - ', COALESCE(supplier_name, 'Tanpa Supplier')) AS description,
			       'out' AS type,
			       'purchase' AS category,
			       total_amount AS amount
			FROM purchases
			WHERE created_at BETWEEN $1 AND $2

			UNION ALL

			-- Cash Out: payroll
			SELECT paid_at AS created_at,
			       CONCAT('Gaji - ', COALESCE(e.name, 'Karyawan')) AS description,
			       'out' AS type,
			       'payroll' AS category,
			       p.total AS amount
			FROM payroll p
			LEFT JOIN employees e ON e.id = p.employee_id
			WHERE paid_at BETWEEN $1 AND $2

			UNION ALL

			-- Cash Out: expenses
			SELECT expense_date::timestamptz AS created_at,
			       CONCAT('Pengeluaran - ', description) AS description,
			       'out' AS type,
			       'expense' AS category,
			       amount
			FROM expenses
			WHERE expense_date BETWEEN $1::date AND $2::date

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
			WHERE pp.created_at BETWEEN $1 AND $2
		) all_entries
		ORDER BY date DESC
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LedgerEntry
	for rows.Next() {
		var e models.LedgerEntry
		err := rows.Scan(&e.Date, &e.Description, &e.Type, &e.Category, &e.Amount)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	// Hitung running_balance dari awal (reverse order, karena sorted DESC)
	// Hitung total dulu
	var runningBalance float64
	for _, e := range entries {
		if e.Type == "in" {
			runningBalance += e.Amount
		} else {
			runningBalance -= e.Amount
		}
	}
	// Set running balance dari terbesar ke terkecil (karena DESC)
	balance := runningBalance
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Type == "out" {
			balance += entries[i].Amount
		}
		entries[i].RunningBalance = balance
		if entries[i].Type == "in" {
			balance -= entries[i].Amount
		}
	}

	if entries == nil {
		entries = []models.LedgerEntry{}
	}
	return &models.LedgerResponse{Data: entries}, nil
}
