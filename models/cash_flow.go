package models

// CashFlowSummary merepresentasikan ringkasan arus kas (cash in & cash out)
type CashFlowSummary struct {
	CashIn           float64 `json:"cash_in"`            // Pemasukan dari penjualan (transactions)
	CashOutPurchases float64 `json:"cash_out_purchases"` // Pengeluaran untuk beli stok
	CashOutPayroll   float64 `json:"cash_out_payroll"`   // Pengeluaran untuk bayar gaji karyawan
	CashOutExpenses  float64 `json:"cash_out_expenses"`  // Pengeluaran operasional tambahan
	CashOutTotal     float64 `json:"cash_out_total"`     // Total semua pengeluaran
	NetCashFlow      float64 `json:"net_cash_flow"`      // Cash In - Cash Out Total
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
