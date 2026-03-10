package models

import "time"

// CashFlowSummary merepresentasikan ringkasan arus kas (cash in & cash out)
type CashFlowSummary struct {
	InitialBalance      float64 `json:"initial_balance"`
	CashIn              float64 `json:"cash_in"`
	CashOutPurchases    float64 `json:"cash_out_purchases"`
	CashOutPayroll      float64 `json:"cash_out_payroll"`
	CashOutExpenses     float64 `json:"cash_out_expenses"`
	CashOutDebtPayments float64 `json:"cash_out_debt_payments"`
	CashOutTotal        float64 `json:"cash_out_total"`
	NetCashFlow         float64 `json:"net_cash_flow"`
	SaldoKas            float64 `json:"saldo_kas"`
}

// CashFlowTrendData merepresentasikan data per periode (harian/bulanan)
type CashFlowTrendData struct {
	Period  string  `json:"period"`   // YYYY-MM-DD atau YYYY-MM
	CashIn  float64 `json:"cash_in"`  // Pemasukan periode tersebut
	CashOut float64 `json:"cash_out"` // Total pengeluaran (pembelian + gaji + expenses) periode tsb
	Net     float64 `json:"net"`      // Pemasukan - Pengeluaran
}

// CashFlowTrendResponse menampung response list trend arus kas
type CashFlowTrendResponse struct {
	Data []CashFlowTrendData `json:"data"`
}

// LedgerEntry adalah satu baris di buku kas (cash ledger)
type LedgerEntry struct {
	Date           string  `json:"date"`
	Description    string  `json:"description"`
	Type           string  `json:"type"`     // "in" atau "out"
	Category       string  `json:"category"` // "sale", "purchase", "payroll", "expense", "debt_payment"
	Amount         float64 `json:"amount"`
	RunningBalance float64 `json:"running_balance"`
}

// LedgerResponse menampung response ledger
type LedgerResponse struct {
	Data []LedgerEntry `json:"data"`
}

// SupplierDebtSummary menampung ringkasan hutang supplier
type SupplierDebtSummary struct {
	TotalPayable           float64 `json:"total_payable"`
	TotalSuppliersWithDebt int     `json:"total_suppliers_with_debt"`
}

// CashFund merepresentasikan satu catatan dana masuk/keluar
type CashFund struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"` // "in" atau "out"
	Amount      float64   `json:"amount"`
	Date        string    `json:"date"` // YYYY-MM-DD
	Description string    `json:"description"`
	CreatedBy   *int      `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CashFundRequest DTO untuk POST /api/cash-flow/funds
type CashFundRequest struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
}

// CashFundSummary ringkasan total dana
type CashFundSummary struct {
	TotalIn  float64 `json:"total_in"`
	TotalOut float64 `json:"total_out"`
	Net      float64 `json:"net"`
}

// CashFundsResponse response GET /api/cash-flow/funds
type CashFundsResponse struct {
	Data    []CashFund      `json:"data"`
	Summary CashFundSummary `json:"summary"`
}

// CashInitialBalance response GET /api/cash-flow/initial-balance
type CashInitialBalance struct {
	Amount    float64 `json:"amount"`
	Breakdown struct {
		TotalIn  float64 `json:"total_in"`
		TotalOut float64 `json:"total_out"`
	} `json:"breakdown"`
}
