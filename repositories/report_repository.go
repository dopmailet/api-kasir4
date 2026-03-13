package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

// Default timezone removed since it's dynamically parsed now

// ReportRepository handles database operations for reports
type ReportRepository struct {
	db *sql.DB
}

// NewReportRepository creates a new ReportRepository
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetDailySalesReport retrieves sales report for today (fallback, uses default timezone)
func (r *ReportRepository) GetDailySalesReport(userID *int, storeID int, tzName string) (*models.SalesReport, error) {
	loc, _ := time.LoadLocation(tzName)
	if loc == nil {
		loc = time.FixedZone("WITA", 8*60*60)
	}
	now := time.Now().In(loc)
	return r.getSalesReportByDateRange(now, now, userID, storeID, tzName)
}

// GetSalesReportByDateRange retrieves sales report for a date range
// startDate dan endDate sudah mengandung timezone yang benar dari handler/caller
func (r *ReportRepository) GetSalesReportByDateRange(startDate, endDate time.Time, userID *int, storeID int, tzName string) (*models.SalesReport, error) {
	return r.getSalesReportByDateRange(startDate, endDate, userID, storeID, tzName)
}

// getSalesReportByDateRange is a private helper function
func (r *ReportRepository) getSalesReportByDateRange(startDate, endDate time.Time, userID *int, storeID int, tzName string) (*models.SalesReport, error) {
	var report models.SalesReport

	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	// Prepare arguments context — store_id selalu ada di $4
	argsBase := []interface{}{startStr, endStr, tzName, storeID}
	userFilterStr := ""
	if userID != nil {
		userFilterStr = " AND created_by = $5 "
		argsBase = append(argsBase, *userID)
	}

	// Prepare join user Filter String context
	userJoinFilterStr := ""
	if userID != nil {
		userJoinFilterStr = " AND t.created_by = $5 "
	}

	dateFilterSql := " AND created_at >= ($1::timestamp AT TIME ZONE $3) AND created_at < (($2::timestamp + INTERVAL '1 day') AT TIME ZONE $3) "
	dateFilterSqlT := " AND t.created_at >= ($1::timestamp AT TIME ZONE $3) AND t.created_at < (($2::timestamp + INTERVAL '1 day') AT TIME ZONE $3) "

	// Query 1A: Total revenue (nett) dan total transaksi
	queryRevenue := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE store_id = $4 ` + dateFilterSql + userFilterStr + `
	`
	err := r.db.QueryRow(queryRevenue, argsBase...).Scan(
		&report.TotalRevenue,
		&report.TotalTransaksi,
	)
	if err != nil {
		return nil, err
	}

	// Query 1B: Total items terjual dan profit
	queryItems := `
		SELECT 
			COALESCE(SUM(hpp.total_qty), 0) as total_items_sold,
			COALESCE(SUM(t.total_amount) - SUM(hpp.total_hpp), 0) as total_profit
		FROM transactions t
		JOIN (
			SELECT 
				td.transaction_id,
				SUM(td.quantity) as total_qty,
				SUM(COALESCE(td.harga_beli, td.price) * td.quantity) as total_hpp
			FROM transaction_details td
			GROUP BY td.transaction_id
		) hpp ON hpp.transaction_id = t.id
		WHERE t.store_id = $4 ` + dateFilterSqlT + userJoinFilterStr + `
	`
	err = r.db.QueryRow(queryItems, argsBase...).Scan(
		&report.TotalItemsSold,
		&report.TotalProfit,
	)
	if err != nil {
		return nil, err
	}

	// Query 2: Total pengeluaran (pembelian) dalam periode yang sama
	queryPengeluaran := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_pengeluaran,
			COUNT(*) as total_pembelian
		FROM purchases
		WHERE store_id = $4 ` + dateFilterSql + `
	`

	err = r.db.QueryRow(queryPengeluaran, argsBase[0], argsBase[1], argsBase[2], argsBase[3]).Scan(
		&report.TotalPengeluaran,
		&report.TotalPembelian,
	)
	if err != nil {
		return nil, err
	}

	// Query 3: Total pengeluaran Gaji (Payroll) dalam periode yang sama
	queryPayroll := `
		SELECT 
			COALESCE(SUM(total), 0) as total_payroll
		FROM payroll
		WHERE store_id = $4 
		  AND paid_at >= ($1::timestamp AT TIME ZONE $3) 
		  AND paid_at < (($2::timestamp + INTERVAL '1 day') AT TIME ZONE $3)
	`
	err = r.db.QueryRow(queryPayroll, argsBase[0], argsBase[1], argsBase[2], argsBase[3]).Scan(
		&report.TotalPayroll,
	)
	if err != nil {
		return nil, err
	}

	// Query 4: Total pengeluaran Operasional (Expenses) dalam periode yang sama
	queryExpenses := `
		SELECT 
			COALESCE(SUM(amount), 0) as total_expenses
		FROM expenses
		WHERE store_id = $4 
		  AND expense_date >= ($1::timestamp AT TIME ZONE $3) 
		  AND expense_date < (($2::timestamp + INTERVAL '1 day') AT TIME ZONE $3)
	`
	err = r.db.QueryRow(queryExpenses, argsBase[0], argsBase[1], argsBase[2], argsBase[3]).Scan(
		&report.TotalExpenses,
	)
	if err != nil {
		return nil, err
	}

	// Hitung laba bersih = laba kotor (total_profit) - pengeluaran_gaji - pengeluaran_operasional
	report.LabaBersih = report.TotalProfit - report.TotalPayroll - report.TotalExpenses

	// Query 5: Semua produk terjual (sorted by total_sales DESC)
	queryProducts := `
		SELECT 
			p.nama as nama_produk,
			SUM(td.quantity) as jumlah,
			COALESCE(SUM(td.subtotal), 0) as total_sales,
			COALESCE(SUM(
				td.subtotal
				- (COALESCE(td.harga_beli, 0) * td.quantity)
				- (
					td.subtotal
					/ NULLIF((SELECT SUM(s.subtotal) FROM transaction_details s WHERE s.transaction_id = td.transaction_id), 0)
					* COALESCE(t.discount_amount - (SELECT COALESCE(SUM(s.discount_amount),0) FROM transaction_details s WHERE s.transaction_id = td.transaction_id), 0)
				)
			), 0) as total_profit
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.store_id = $4 ` + dateFilterSqlT + userJoinFilterStr + `
		GROUP BY p.id, p.nama
		ORDER BY total_sales DESC
	`

	rows, err := r.db.Query(queryProducts, argsBase...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produkTerlaris []models.TopProduct
	for rows.Next() {
		var p models.TopProduct
		if err := rows.Scan(&p.NamaProduk, &p.Jumlah, &p.TotalSales, &p.TotalProfit); err != nil {
			return nil, err
		}
		produkTerlaris = append(produkTerlaris, p)
	}

	if produkTerlaris == nil {
		produkTerlaris = []models.TopProduct{}
	}
	report.ProdukTerlaris = produkTerlaris

	return &report, nil
}

// GetDashboardAssets retrieves total asset cost and total asset retail from products table
// where stock > 0
func (r *ReportRepository) GetDashboardAssets(storeID int) (*models.AssetReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(stok * harga_beli), 0) AS total_asset_cost,
			COALESCE(SUM(stok * harga), 0)      AS total_asset_retail
		FROM products
		WHERE stok > 0 AND store_id = $1
	`

	var report models.AssetReport
	err := r.db.QueryRow(query, storeID).Scan(
		&report.TotalAssetCost,
		&report.TotalAssetRetail,
	)

	if err != nil {
		return nil, err
	}

	return &report, nil
}

// GetSalesTrend retrieves sales trend data for chart
func (r *ReportRepository) GetSalesTrend(startDate, endDate time.Time, interval string, tzName string, storeID int) ([]models.SalesTrend, error) {
	var trends []models.SalesTrend

	if tzName == "" {
		tzName = "Asia/Jakarta"
	}

	var dateFormat string
	switch interval {
	case "month":
		dateFormat = "YYYY-MM"
	case "year":
		dateFormat = "YYYY"
	default: // "day"
		dateFormat = "YYYY-MM-DD"
	}

	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")

	query := `
		SELECT 
			TO_CHAR((t.created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) as period,
			COALESCE(SUM(t.total_amount), 0) as total_sales,
			COALESCE(SUM(t.total_amount) - SUM(hpp.total_hpp), 0) as total_profit,
			COUNT(DISTINCT t.id) as transaction_count
		FROM transactions t
		JOIN (
			SELECT 
				td.transaction_id,
				SUM(COALESCE(td.harga_beli, 0) * td.quantity) as total_hpp
			FROM transaction_details td
			GROUP BY td.transaction_id
		) hpp ON hpp.transaction_id = t.id
		WHERE DATE(t.created_at AT TIME ZONE 'UTC' AT TIME ZONE $1) >= $3
		  AND DATE(t.created_at AT TIME ZONE 'UTC' AT TIME ZONE $1) <= $4
		  AND t.store_id = $5
		GROUP BY TO_CHAR((t.created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2)
		ORDER BY TO_CHAR((t.created_at AT TIME ZONE 'UTC' AT TIME ZONE $1), $2) ASC
	`

	rows, err := r.db.Query(query, tzName, dateFormat, startStr, endStr, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.SalesTrend
		err := rows.Scan(&t.Date, &t.TotalSales, &t.TotalProfit, &t.TransactionCount)
		if err != nil {
			return nil, err
		}
		trends = append(trends, t)
	}

	return trends, nil
}

// GetTopProducts returns top selling products by quantity and by profit
func (r *ReportRepository) GetTopProducts(startDate, endDate time.Time, limit int, storeID int, tzName string) ([]models.TopProduct, []models.TopProduct, error) {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")
	dateFilterSqlT := " AND t.created_at >= ($1::timestamp AT TIME ZONE $3) AND t.created_at < (($2::timestamp + INTERVAL '1 day') AT TIME ZONE $3) "

	// 1. Top by Quantity
	queryQty := `
		SELECT 
			p.nama,
			COALESCE(SUM(td.quantity), 0) as jumlah,
			COALESCE(SUM(td.subtotal), 0) as total_sales,
			COALESCE(SUM(
				td.subtotal
				- (COALESCE(td.harga_beli, 0) * td.quantity)
				- (
					td.subtotal
					/ NULLIF((SELECT SUM(s.subtotal) FROM transaction_details s WHERE s.transaction_id = td.transaction_id), 0)
					* COALESCE(t.discount_amount - (SELECT COALESCE(SUM(s.discount_amount),0) FROM transaction_details s WHERE s.transaction_id = td.transaction_id), 0)
				)
			), 0) as total_profit
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.store_id = $4 ` + dateFilterSqlT + `
		GROUP BY p.id, p.nama
		ORDER BY jumlah DESC
		LIMIT $5
	`

	rowsQty, err := r.db.Query(queryQty, startStr, endStr, tzName, storeID, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rowsQty.Close()

	var topQty []models.TopProduct
	for rowsQty.Next() {
		var p models.TopProduct
		if err := rowsQty.Scan(&p.NamaProduk, &p.Jumlah, &p.TotalSales, &p.TotalProfit); err != nil {
			return nil, nil, err
		}
		topQty = append(topQty, p)
	}

	// 2. Top by Profit
	queryProfit := `
		SELECT 
			p.nama,
			COALESCE(SUM(td.quantity), 0) as jumlah,
			COALESCE(SUM(td.subtotal), 0) as total_sales,
			COALESCE(SUM(
				td.subtotal
				- (COALESCE(td.harga_beli, 0) * td.quantity)
				- (
					td.subtotal
					/ NULLIF((SELECT SUM(s.subtotal) FROM transaction_details s WHERE s.transaction_id = td.transaction_id), 0)
					* COALESCE(t.discount_amount - (SELECT COALESCE(SUM(s.discount_amount),0) FROM transaction_details s WHERE s.transaction_id = td.transaction_id), 0)
				)
			), 0) as total_profit
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.store_id = $4 ` + dateFilterSqlT + `
		GROUP BY p.id, p.nama
		ORDER BY total_profit DESC
		LIMIT $5
	`

	rowsProfit, err := r.db.Query(queryProfit, startStr, endStr, tzName, storeID, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rowsProfit.Close()

	var topProfit []models.TopProduct
	for rowsProfit.Next() {
		var p models.TopProduct
		if err := rowsProfit.Scan(&p.NamaProduk, &p.Jumlah, &p.TotalSales, &p.TotalProfit); err != nil {
			return nil, nil, err
		}
		topProfit = append(topProfit, p)
	}

	return topQty, topProfit, nil
}

// CountLowStockProducts menghitung jumlah produk yang stoknya <= threshold
// Digunakan untuk widget peringatan stok menipis di dashboard
func (r *ReportRepository) CountLowStockProducts(threshold int, storeID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products WHERE stok <= $1 AND store_id = $2`
	err := r.db.QueryRow(query, threshold, storeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
