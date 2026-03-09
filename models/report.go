package models

import "time"

// SalesReport represents daily sales summary
// Struct untuk laporan penjualan harian
type SalesReport struct {
	TotalRevenue     float64      `json:"total_revenue"`
	TotalTransaksi   int          `json:"total_transaksi"`
	TotalItemsSold   int          `json:"total_items_sold"`  // Total items terjual
	TotalProfit      float64      `json:"total_profit"`      // Total keuntungan kotor (revenue - modal barang terjual)
	TotalPengeluaran float64      `json:"total_pengeluaran"` // Total pembelian/pengadaan barang
	TotalPembelian   int          `json:"total_pembelian"`   // Jumlah transaksi pembelian
	TotalPayroll     float64      `json:"total_payroll"`     // Total gaji karyawan dibayarkan
	TotalExpenses    float64      `json:"total_expenses"`    // Total pengeluaran operasional
	LabaBersih       float64      `json:"laba_bersih"`       // Profit - Payroll - Expenses
	ProdukTerlaris   []TopProduct `json:"produk_terlaris"`   // Array semua produk terjual
}

// TopProduct represents the best selling product
// Struct untuk produk terlaris
type TopProduct struct {
	NamaProduk  string  `json:"nama_produk"`
	Jumlah      int     `json:"jumlah"`       // Quantity terjual
	TotalSales  float64 `json:"total_sales"`  // Total omzet
	TotalProfit float64 `json:"total_profit"` // Total keuntungan (profit)
}

// SalesTrend represents sales data over a period
// Struct untuk grafik trend penjualan
type SalesTrend struct {
	Date             string  `json:"date"`              // Tanggal/Bulan/Tahun
	TotalSales       float64 `json:"total_sales"`       // Total penjualan (omzet)
	TotalProfit      float64 `json:"total_profit"`      // Total keuntungan (profit)
	TransactionCount int     `json:"transaction_count"` // Jumlah transaksi
}

// DashboardSummary represents KPI summary for dashboard with growth comparison
// Berisi data periode saat ini, periode sebelumnya, dan perbandingan persentase pertumbuhan
type DashboardSummary struct {
	PeriodStart       time.Time   `json:"period_start"`       // Awal periode yang diminta
	PeriodEnd         time.Time   `json:"period_end"`         // Akhir periode yang diminta
	Current           SalesReport `json:"current"`            // Data periode saat ini
	Previous          SalesReport `json:"previous"`           // Data periode sebelumnya (untuk perbandingan)
	RevenueGrowth     float64     `json:"revenue_growth"`     // % perubahan omzet vs periode sebelumnya
	ProfitGrowth      float64     `json:"profit_growth"`      // % perubahan profit vs periode sebelumnya
	TransactionGrowth float64     `json:"transaction_growth"` // % perubahan jumlah transaksi vs periode sebelumnya
	LowStockCount     int         `json:"low_stock_count"`    // Jumlah produk stok menipis (threshold dari query param)
}

// AssetReport represents summary of inventory capital and sales potential.
// Struct untuk rangkuman Aset/Modal berjalan Inventory
type AssetReport struct {
	TotalAssetCost   float64 `json:"total_asset_cost"`   // Modal tertanam dari HPP
	TotalAssetRetail float64 `json:"total_asset_retail"` // Potensi omset dari Harga Jual (Retail)
}
