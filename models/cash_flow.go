package models

// CashFlowSummary merepresentasikan ringkasan arus kas (cash in & cash out)
type CashFlowSummary struct {
	CashIn              float64 `json:"cash_in"`
	CashOutPurchases    float64 `json:"cash_out_purchases"`
	CashOutPayroll      float64 `json:"cash_out_payroll"`
	CashOutExpenses     float64 `json:"cash_out_expenses"`
	CashOutDebtPayments float64 `json:"cash_out_debt_payments"`
	CashOutTotal        float64 `json:"cash_out_total"`
	NetCashFlow         float64 `json:"net_cash_flow"`
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
