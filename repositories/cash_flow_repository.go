package repositories

import (
	"database/sql"
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

	// 1. Cash In (Total Pemasukan dari Transaksi)
	queryCashIn := `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions
		WHERE created_at BETWEEN $1 AND $2
	`
	err := r.db.QueryRow(queryCashIn, startDate, endDate).Scan(&summary.CashIn)
	if err != nil {
		return nil, err
	}

	// 2. Cash Out: Purchases
	queryPurchases := `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM purchases
		WHERE created_at BETWEEN $1 AND $2
	`
	err = r.db.QueryRow(queryPurchases, startDate, endDate).Scan(&summary.CashOutPurchases)
	if err != nil {
		return nil, err
	}

	// 3. Cash Out: Payroll (Pastikan status dibayar / memiliki "paid_at" time log)
	queryPayroll := `
		SELECT COALESCE(SUM(total), 0)
		FROM payroll
		WHERE paid_at BETWEEN $1 AND $2
	`
	err = r.db.QueryRow(queryPayroll, startDate, endDate).Scan(&summary.CashOutPayroll)
	if err != nil {
		return nil, err
	}

	// 4. Cash Out: Expenses
	queryExpenses := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE expense_date BETWEEN $1 AND $2
	`
	err = r.db.QueryRow(queryExpenses, startDate, endDate).Scan(&summary.CashOutExpenses)
	if err != nil {
		return nil, err
	}

	// Hitung Aggregasi Akhir
	summary.CashOutTotal = summary.CashOutPurchases + summary.CashOutPayroll + summary.CashOutExpenses
	summary.NetCashFlow = summary.CashIn - summary.CashOutTotal

	return &summary, nil
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
