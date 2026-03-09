package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"log"
	"time"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction creates a new transaction with details (OPTIMIZED - batch queries)
// API contract tidak berubah - request dan response format tetap sama
func (r *TransactionRepository) CreateTransaction(req *models.CheckoutRequest) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// ─── STEP 1: Batch fetch ALL products in 1 query (bukan per-item) ───
	productIDPlaceholders := ""
	productIDArgs := make([]interface{}, len(req.Items))
	for i, item := range req.Items {
		productIDArgs[i] = item.ProductID
		if i > 0 {
			productIDPlaceholders += ", "
		}
		productIDPlaceholders += fmt.Sprintf("$%d", i+1)
	}

	productRows, err := tx.Query(
		fmt.Sprintf("SELECT id, harga, stok, category_id, harga_beli FROM products WHERE id IN (%s)", productIDPlaceholders),
		productIDArgs...,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data produk: %w", err)
	}

	type productInfo struct {
		Price      float64
		Stok       int
		CategoryID sql.NullInt64
		HargaBeli  sql.NullFloat64
	}
	productMap := make(map[int]*productInfo)
	for productRows.Next() {
		var id int
		var p productInfo
		if scanErr := productRows.Scan(&id, &p.Price, &p.Stok, &p.CategoryID, &p.HargaBeli); scanErr != nil {
			productRows.Close()
			return nil, scanErr
		}
		productMap[id] = &p
	}
	productRows.Close()

	// Validasi
	for _, item := range req.Items {
		p, exists := productMap[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", item.ProductID)
		}
		if p.Stok < item.Quantity {
			return nil, fmt.Errorf("stok produk ID %d tidak mencukupi (sisa: %d, diminta: %d)", item.ProductID, p.Stok, item.Quantity)
		}
	}

	// ─── STEP 2: Batch fetch ALL active item discounts in 1 query ───
	discountRows, err := tx.Query(`
		SELECT id, type, value, product_id, category_id FROM discounts 
		WHERE is_active = TRUE AND NOW() BETWEEN start_date AND end_date
		AND (product_id IS NOT NULL OR category_id IS NOT NULL)
		ORDER BY value DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data diskon: %w", err)
	}

	type discInfo struct {
		Type  string
		Value float64
	}
	productDiscounts := make(map[int]*discInfo)
	categoryDiscounts := make(map[int]*discInfo)
	for discountRows.Next() {
		var id int
		var dType string
		var dValue float64
		var productID, categoryID sql.NullInt64
		if scanErr := discountRows.Scan(&id, &dType, &dValue, &productID, &categoryID); scanErr != nil {
			discountRows.Close()
			return nil, scanErr
		}
		// Jika diskon punya KEDUA product_id DAN category_id → data tidak valid, skip
		if productID.Valid && categoryID.Valid {
			log.Printf("⚠️ Diskon ID %d punya product_id DAN category_id sekaligus — diskip (data tidak valid)", id)
			continue
		}
		if productID.Valid {
			pid := int(productID.Int64)
			if _, exists := productDiscounts[pid]; !exists {
				productDiscounts[pid] = &discInfo{Type: dType, Value: dValue}
			}
		} else if categoryID.Valid {
			cid := int(categoryID.Int64)
			if _, exists := categoryDiscounts[cid]; !exists {
				categoryDiscounts[cid] = &discInfo{Type: dType, Value: dValue}
			}
		}
	}
	discountRows.Close()

	// ─── STEP 3: Process items in-memory ───
	// Prioritas: gunakan diskon dari frontend jika ada (discount_amount > 0),
	// fallback ke diskon dari DB jika tidak ada diskon frontend.
	var totalAmount float64
	var totalDiscount float64
	details := make([]models.TransactionDetail, 0, len(req.Items))

	for _, item := range req.Items {
		p := productMap[item.ProductID]

		// Harga satuan SELALU dari database — jangan izinkan frontend override
		// untuk mencegah harga naik/turun tanpa sepengetahuan admin
		unitPrice := p.Price
		if item.Price > 0 && item.Price != p.Price {
			log.Printf("⚠️ Produk ID %d: frontend kirim harga %.0f, DB harga %.0f — gunakan harga DB",
				item.ProductID, item.Price, p.Price)
		}

		var discountAmount float64
		var discountType string
		var discountValue float64

		if item.DiscountAmount > 0 {
			// Gunakan diskon yang sudah dihitung frontend
			discountAmount = item.DiscountAmount
			discountType = item.DiscountType
			discountValue = item.DiscountValue
		} else {
			// Fallback: cari diskon dari DB (per produk / kategori)
			var itemDiscountPerUnit float64
			disc, found := productDiscounts[item.ProductID]
			if !found && p.CategoryID.Valid {
				disc, _ = categoryDiscounts[int(p.CategoryID.Int64)]
			}
			if disc != nil {
				if disc.Type == "PERCENTAGE" {
					itemDiscountPerUnit = unitPrice * (disc.Value / 100)
					discountType = "percentage"
				} else {
					itemDiscountPerUnit = disc.Value
					discountType = "fixed"
				}
				if itemDiscountPerUnit > unitPrice {
					itemDiscountPerUnit = unitPrice
				}
				discountValue = disc.Value
				discountAmount = itemDiscountPerUnit * float64(item.Quantity)
			}
		}

		// Subtotal = (harga × qty) - total diskon item
		subtotal := (unitPrice * float64(item.Quantity)) - discountAmount
		if subtotal < 0 {
			subtotal = 0
		}
		totalAmount += subtotal
		totalDiscount += discountAmount

		var hargaBeliSnapshot float64
		if p.HargaBeli.Valid {
			hargaBeliSnapshot = p.HargaBeli.Float64
		} else {
			hargaBeliSnapshot = unitPrice
		}

		details = append(details, models.TransactionDetail{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			Price:          unitPrice,
			Subtotal:       subtotal,
			DiscountType:   discountType,
			DiscountValue:  discountValue,
			DiscountAmount: discountAmount,
			HargaBeli:      hargaBeliSnapshot,
		})
	}

	// ─── STEP 4: Batch UPDATE stock in 1 query ───
	stockQuery := "UPDATE products SET stok = CASE "
	stockIDs := ""
	stockArgs := make([]interface{}, 0, len(req.Items)*2)
	argIdx := 1
	for i, item := range req.Items {
		stockQuery += fmt.Sprintf("WHEN id = $%d THEN stok - $%d ", argIdx, argIdx+1)
		stockArgs = append(stockArgs, item.ProductID, item.Quantity)
		if i > 0 {
			stockIDs += ", "
		}
		stockIDs += fmt.Sprintf("$%d", argIdx)
		argIdx += 2
	}
	stockQuery += fmt.Sprintf("END WHERE id IN (%s)", stockIDs)
	_, err = tx.Exec(stockQuery, stockArgs...)
	if err != nil {
		return nil, fmt.Errorf("gagal update stok: %w", err)
	}

	// ─── STEP 5: Process Global Discount (if any) ───
	// Catatan alur diskon:
	// - totalAmount   = jumlah subtotal per item (sudah dikurangi diskon per-item)
	// - totalDiscount = akumulasi diskon per-item yang sudah dihitung di STEP 3
	// - req.DiscountAmount (root) = total diskon keseluruhan dari frontend (INFORMASI SAJA,
	//   sudah termasuk diskon per-item → JANGAN dikurangi lagi dari totalAmount)
	// - globalDiscountAmount = diskon tambahan dari tabel discounts (via DiscountID)
	var usedDiscountID *int
	var globalDiscountAmount float64

	if req.DiscountID != nil {
		// Diskon global dari tabel discounts (berdasarkan ID) — ini diskon TAMBAHAN
		var d models.Discount
		var isValid bool
		err = tx.QueryRow(`
			SELECT id, type, value, min_order_amount, (NOW() BETWEEN start_date AND end_date) as is_valid 
			FROM discounts WHERE id = $1 AND is_active = TRUE 
			AND product_id IS NULL AND category_id IS NULL`, *req.DiscountID).Scan(
			&d.ID, &d.Type, &d.Value, &d.MinOrderAmount, &isValid,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("diskon global ID %d tidak valid atau (mungkin diskon produk/kategori)", *req.DiscountID)
		}
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, fmt.Errorf("diskon global sudah kedaluwarsa")
		}
		if totalAmount < d.MinOrderAmount {
			return nil, fmt.Errorf("min order %.0f tidak terpenuhi", d.MinOrderAmount)
		}
		if d.Type == "PERCENTAGE" {
			globalDiscountAmount = totalAmount * (d.Value / 100)
		} else {
			globalDiscountAmount = d.Value
		}
		if globalDiscountAmount > totalAmount {
			globalDiscountAmount = totalAmount
		}
		usedDiscountID = req.DiscountID
	}
	// req.DiscountAmount TIDAK dikurangi lagi — diskon per-item sudah masuk
	// ke totalAmount via pengurangan subtotal di STEP 3.

	// finalTotal = total setelah semua diskon
	// Jika ada DiscountID (diskon global tambahan), kurangi dari totalAmount.
	// Jika tidak, totalAmount sudah bersih.
	finalTotal := totalAmount - globalDiscountAmount
	if finalTotal < 0 {
		finalTotal = 0
	}
	totalDiscount += globalDiscountAmount

	// change_amount = uang bayar - total yang harus dibayar (setelah SEMUA diskon)
	paymentAmount := req.PaymentAmount
	changeAmount := 0.0
	if paymentAmount > 0 {
		changeAmount = paymentAmount - finalTotal
		if changeAmount < 0 {
			changeAmount = 0
		}
	}

	// ─── STEP 6: Insert transaction header ───
	// Gunakan req.DiscountAmount sebagai discount_amount yang disimpan ke DB
	// (total diskon keseluruhan dari frontend, sudah include diskon per-item).
	// Jika frontend tidak mengirim req.DiscountAmount, fallback ke totalDiscount internal.
	savedDiscountAmount := totalDiscount
	if req.DiscountAmount > 0 {
		savedDiscountAmount = req.DiscountAmount
	}

	var transactionID int
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount, discount_id, discount_amount, payment_amount, change_amount, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		finalTotal, usedDiscountID, savedDiscountAmount, paymentAmount, changeAmount, req.CreatedBy,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// ─── STEP 7: Batch insert details (termasuk discount per item) ───
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, price, subtotal, harga_beli, discount_type, discount_value, discount_amount) VALUES "
		values := make([]interface{}, 0, len(details)*9)
		for i, detail := range details {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9)
			// Simpan NULL jika discount_type kosong
			var discType interface{}
			if detail.DiscountType != "" {
				discType = detail.DiscountType
			}
			values = append(values,
				transactionID, detail.ProductID, detail.Quantity, detail.Price, detail.Subtotal,
				detail.HargaBeli, discType, detail.DiscountValue, detail.DiscountAmount,
			)
		}
		_, err = tx.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	}

	// ─── STEP 8: Commit ───
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:             transactionID,
		TotalAmount:    finalTotal,
		DiscountID:     usedDiscountID,
		DiscountAmount: totalDiscount,
		PaymentAmount:  paymentAmount,
		ChangeAmount:   changeAmount,
		CreatedBy:      &req.CreatedBy,
	}, nil
}

// GetAll retrieves all transactions ordered by date descending
// Fungsi ini mengambil semua data transaksi untuk history, termasuk profit per transaksi
func (r *TransactionRepository) GetAll(userID *int) ([]models.Transaction, error) {
	// Profit = (total_amount - discount_amount) - HPP
	// total_amount = total setelah diskon item
	// discount_amount = diskon transaksi (terpisah, perlu dikurangi)
	// Subquery HPP per transaksi untuk hindari duplikasi
	query := `
		SELECT 
			t.id, 
			t.total_amount, 
			t.discount_id, 
			t.discount_amount, 
			COALESCE(t.payment_amount, 0) as payment_amount,
			COALESCE(t.change_amount, 0) as change_amount,
			t.created_at,
			COALESCE(hpp.total_qty, 0) as total_items,
			t.total_amount - COALESCE(hpp.total_hpp, 0) as profit,
			t.created_by,
			u.username
		FROM transactions t
		LEFT JOIN (
			SELECT 
				td.transaction_id,
				SUM(td.quantity) as total_qty,
				SUM(COALESCE(td.harga_beli, 0) * td.quantity) as total_hpp
			FROM transaction_details td
			GROUP BY td.transaction_id
		) hpp ON hpp.transaction_id = t.id
		LEFT JOIN users u ON t.created_by = u.id
	`

	args := []interface{}{}
	if userID != nil {
		query += " WHERE t.created_by = $1 "
		args = append(args, *userID)
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var t models.Transaction
		var discountID sql.NullInt64
		var createdBy sql.NullInt64
		var username sql.NullString

		err := rows.Scan(&t.ID, &t.TotalAmount, &discountID, &t.DiscountAmount, &t.PaymentAmount, &t.ChangeAmount, &t.CreatedAt, &t.TotalItems, &t.Profit, &createdBy, &username)
		if err != nil {
			return nil, err
		}

		if discountID.Valid {
			id := int(discountID.Int64)
			t.DiscountID = &id
		}
		if createdBy.Valid {
			id := int(createdBy.Int64)
			t.CreatedBy = &id
		}
		if username.Valid {
			t.Username = username.String
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetByDateRange retrieves transactions within a date range
// startDate dan endDate sudah mengandung timezone yang benar dari handler
func (r *TransactionRepository) GetByDateRange(startDate, endDate time.Time, userID *int) ([]models.Transaction, error) {
	query := `
		SELECT 
			t.id, 
			t.total_amount, 
			t.discount_id, 
			t.discount_amount, 
			COALESCE(t.payment_amount, 0) as payment_amount,
			COALESCE(t.change_amount, 0) as change_amount,
			t.created_at,
			COALESCE(hpp.total_qty, 0) as total_items,
			t.total_amount - COALESCE(hpp.total_hpp, 0) as profit,
			t.created_by,
			u.username
		FROM transactions t
		LEFT JOIN (
			SELECT 
				td.transaction_id,
				SUM(td.quantity) as total_qty,
				SUM(COALESCE(td.harga_beli, 0) * td.quantity) as total_hpp
			FROM transaction_details td
			GROUP BY td.transaction_id
		) hpp ON hpp.transaction_id = t.id
		LEFT JOIN users u ON t.created_by = u.id
		WHERE t.created_at BETWEEN $1 AND $2
	`

	args := []interface{}{startDate, endDate}

	if userID != nil {
		query += " AND t.created_by = $3 "
		args = append(args, *userID)
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var t models.Transaction
		var discountID sql.NullInt64
		var createdBy sql.NullInt64
		var username sql.NullString

		err := rows.Scan(&t.ID, &t.TotalAmount, &discountID, &t.DiscountAmount, &t.PaymentAmount, &t.ChangeAmount, &t.CreatedAt, &t.TotalItems, &t.Profit, &createdBy, &username)
		if err != nil {
			return nil, err
		}

		if discountID.Valid {
			id := int(discountID.Int64)
			t.DiscountID = &id
		}
		if createdBy.Valid {
			id := int(createdBy.Int64)
			t.CreatedBy = &id
		}
		if username.Valid {
			t.Username = username.String
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetByID retrieves a transaction by ID with all its items
func (r *TransactionRepository) GetByID(id int) (*models.TransactionWithItems, error) {
	queryHeader := `
		SELECT 
			t.id, 
			t.total_amount, 
			t.discount_amount, 
			COALESCE(t.payment_amount, 0) as payment_amount,
			COALESCE(t.change_amount, 0) as change_amount,
			t.created_at,
			COALESCE(hpp.total_qty, 0) as total_items,
			t.total_amount - COALESCE(hpp.total_hpp, 0) as profit,
			t.created_by,
			u.username
		FROM transactions t
		LEFT JOIN (
			SELECT 
				td.transaction_id,
				SUM(td.quantity) as total_qty,
				SUM(COALESCE(td.harga_beli, 0) * td.quantity) as total_hpp
			FROM transaction_details td
			GROUP BY td.transaction_id
		) hpp ON hpp.transaction_id = t.id
		LEFT JOIN users u ON t.created_by = u.id
		WHERE t.id = $1
	`

	var result models.TransactionWithItems
	var createdBy sql.NullInt64
	var username sql.NullString

	err := r.db.QueryRow(queryHeader, id).Scan(
		&result.ID,
		&result.TotalAmount,
		&result.DiscountAmount,
		&result.PaymentAmount,
		&result.ChangeAmount,
		&result.CreatedAt,
		&result.TotalItems,
		&result.Profit,
		&createdBy,
		&username,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaksi dengan ID %d tidak ditemukan", id)
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data transaksi: %w", err)
	}

	if createdBy.Valid {
		id := int(createdBy.Int64)
		result.CreatedBy = &id
	}
	if username.Valid {
		result.Username = username.String
	}

	queryItems := `
		SELECT 
			td.id,
			td.product_id,
			COALESCE(p.nama, 'Produk Dihapus') as product_name,
			td.quantity,
			td.price,
			td.subtotal,
			COALESCE(td.discount_type, '') as discount_type,
			COALESCE(td.discount_value, 0) as discount_value,
			COALESCE(td.discount_amount, 0) as discount_amount
		FROM transaction_details td
		LEFT JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id = $1
		ORDER BY td.id
	`

	rows, err := r.db.Query(queryItems, id)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail items: %w", err)
	}
	defer rows.Close()

	var items []models.TransactionDetail
	for rows.Next() {
		var item models.TransactionDetail
		err := rows.Scan(
			&item.ID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.Price, &item.Subtotal,
			&item.DiscountType, &item.DiscountValue, &item.DiscountAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca detail item: %w", err)
		}
		items = append(items, item)
	}

	if items == nil {
		items = []models.TransactionDetail{}
	}
	result.Items = items

	return &result, nil
}
