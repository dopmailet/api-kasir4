package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"log"
	"time"
)

// PurchaseRepository handles database operations for purchases
// Repository untuk operasi pembelian/pengadaan barang
type PurchaseRepository struct {
	db *sql.DB
}

// NewPurchaseRepository creates a new PurchaseRepository
func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

// Create creates a new purchase with items
// Fungsi ini mencatat pembelian baru:
// - Jika product_id NULL → buat produk baru di tabel products
// - Jika product_id ada → update stok dan harga_beli produk yang sudah ada
// Semua dalam 1 database transaction (atomic)
func (r *PurchaseRepository) Create(req *models.PurchaseRequest, createdBy int) (*models.Purchase, error) {
	// Begin database transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("gagal memulai transaksi: %w", err)
	}

	// Defer rollback jika terjadi error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var totalAmount float64
	processedItems := make([]models.PurchaseItem, 0, len(req.Items))

	// ─── PROSES SETIAP ITEM ───
	for i, item := range req.Items {
		subtotal := float64(item.Quantity) * item.BuyPrice
		totalAmount += subtotal

		var productID int
		var productName string

		if item.ProductID != nil {
			// ═══ RESTOK: Produk sudah ada ═══
			// 1. Ambil nama produk dan validasi produk ada
			err = tx.QueryRow("SELECT id, nama FROM products WHERE id = $1", *item.ProductID).Scan(&productID, &productName)
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("item #%d: produk dengan ID %d tidak ditemukan", i+1, *item.ProductID)
			}
			if err != nil {
				return nil, fmt.Errorf("item #%d: gagal mengambil data produk: %w", i+1, err)
			}

			// 2. Update stok (tambah) dan harga_beli
			_, err = tx.Exec(
				"UPDATE products SET stok = stok + $1, harga_beli = $2 WHERE id = $3",
				item.Quantity, item.BuyPrice, productID,
			)
			if err != nil {
				return nil, fmt.Errorf("item #%d: gagal update stok produk '%s': %w", i+1, productName, err)
			}

			log.Printf("📦 Restok: %s +%d unit (harga beli: %.0f)", productName, item.Quantity, item.BuyPrice)

		} else {
			// ═══ PRODUK BARU: Buat produk dan set stok awal ═══
			if item.ProductName == nil || *item.ProductName == "" {
				return nil, fmt.Errorf("item #%d: nama produk wajib diisi untuk produk baru", i+1)
			}
			if item.SellPrice == nil || *item.SellPrice <= 0 {
				return nil, fmt.Errorf("item #%d: harga jual wajib diisi dan harus > 0 untuk produk baru", i+1)
			}

			productName = *item.ProductName

			// Cek apakah produk dengan nama yang sama sudah ada
			var existingID int
			errCheck := tx.QueryRow("SELECT id FROM products WHERE nama = $1", productName).Scan(&existingID)
			if errCheck == nil {
				// Produk dengan nama yang sama sudah ada → restok saja
				productID = existingID
				_, err = tx.Exec(
					"UPDATE products SET stok = stok + $1, harga_beli = $2 WHERE id = $3",
					item.Quantity, item.BuyPrice, productID,
				)
				if err != nil {
					return nil, fmt.Errorf("item #%d: gagal update stok produk '%s': %w", i+1, productName, err)
				}
				log.Printf("📦 Produk '%s' sudah ada, restok +%d unit", productName, item.Quantity)
			} else {
				// Produk benar-benar baru → insert ke tabel products
				err = tx.QueryRow(
					`INSERT INTO products (nama, harga, harga_beli, stok, category_id, created_by) 
					 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
					productName, *item.SellPrice, item.BuyPrice, item.Quantity, item.CategoryID, createdBy,
				).Scan(&productID)
				if err != nil {
					return nil, fmt.Errorf("item #%d: gagal membuat produk baru '%s': %w", i+1, productName, err)
				}
				log.Printf("✅ Produk baru: '%s' (ID: %d, stok: %d, beli: %.0f, jual: %.0f)",
					productName, productID, item.Quantity, item.BuyPrice, *item.SellPrice)
			}
		}

		// Simpan item yang sudah diproses
		processedItems = append(processedItems, models.PurchaseItem{
			ProductID:   &productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			BuyPrice:    item.BuyPrice,
			SellPrice:   item.SellPrice,
			CategoryID:  item.CategoryID,
			Subtotal:    subtotal,
		})
	}

	// ─── INSERT HEADER PURCHASE ───
	// Tentukan payment_method, paid_amount, payment_status, remaining_amount
	paymentMethod := req.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "cash"
	}
	var paidAmount float64
	if req.PaidAmount != nil {
		paidAmount = *req.PaidAmount
	} else if paymentMethod == "cash" {
		paidAmount = totalAmount
	}
	remainingAmount := totalAmount - paidAmount
	paymentStatus := "paid"
	if remainingAmount > 0 && paidAmount > 0 {
		paymentStatus = "partial"
	} else if paidAmount == 0 {
		paymentStatus = "unpaid"
	}

	var purchaseID int
	err = tx.QueryRow(
		`INSERT INTO purchases (supplier_id, supplier_name, total_amount, payment_method, payment_status, paid_amount, remaining_amount, due_date, payment_notes, notes, created_by) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8::date, $9, $10, $11) RETURNING id`,
		req.SupplierID, req.SupplierName, totalAmount, paymentMethod, paymentStatus, paidAmount, remainingAmount, req.DueDate, req.PaymentNotes, req.Notes, createdBy,
	).Scan(&purchaseID)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan pembelian: %w", err)
	}

	// ─── BATCH INSERT PURCHASE ITEMS ───
	if len(processedItems) > 0 {
		query := `INSERT INTO purchase_items 
			(purchase_id, product_id, product_name, quantity, buy_price, sell_price, category_id, subtotal) VALUES `
		values := make([]interface{}, 0, len(processedItems)*8)

		for i, item := range processedItems {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)
			values = append(values,
				purchaseID, item.ProductID, item.ProductName,
				item.Quantity, item.BuyPrice, item.SellPrice,
				item.CategoryID, item.Subtotal,
			)
		}

		_, err = tx.Exec(query, values...)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan detail pembelian: %w", err)
		}
	}

	// ─── INTEGRASI SUPPLIER (jika supplier_id dikirim) ───
	if req.SupplierID != nil {
		// Verifikasi supplier ada
		var supplierExists bool
		err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1)", *req.SupplierID).Scan(&supplierExists)
		if err != nil {
			return nil, fmt.Errorf("gagal memverifikasi supplier: %w", err)
		}
		if !supplierExists {
			return nil, fmt.Errorf("supplier dengan ID %d tidak ditemukan", *req.SupplierID)
		}

		// Update total_purchases dan total_spent di suppliers
		_, err = tx.Exec(`
			UPDATE suppliers 
			SET total_purchases = total_purchases + 1,
			    total_spent     = total_spent + $1,
			    updated_at      = NOW()
			WHERE id = $2
		`, totalAmount, *req.SupplierID)
		if err != nil {
			return nil, fmt.Errorf("gagal update statistik supplier: %w", err)
		}

		// Buat record hutang (payable) dengan status 'unpaid'
		_, err = tx.Exec(`
			INSERT INTO supplier_payables (supplier_id, purchase_id, amount, paid_amount, status)
			VALUES ($1, $2, $3, 0, 'unpaid')
		`, *req.SupplierID, purchaseID, totalAmount)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat payable supplier: %w", err)
		}

		log.Printf("📋 Payable dibuat untuk supplier ID=%d, purchase ID=%d, amount=%.0f", *req.SupplierID, purchaseID, totalAmount)
	}

	// ─── COMMIT TRANSACTION ───
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	log.Printf("✅ Pembelian berhasil: ID=%d, Total=%.0f, Items=%d", purchaseID, totalAmount, len(processedItems))

	// Build response
	dueDateStr := req.DueDate
	purchase := &models.Purchase{
		ID:              purchaseID,
		SupplierID:      req.SupplierID,
		SupplierName:    req.SupplierName,
		TotalAmount:     totalAmount,
		PaymentMethod:   paymentMethod,
		PaymentStatus:   paymentStatus,
		PaidAmount:      paidAmount,
		RemainingAmount: remainingAmount,
		DueDate:         dueDateStr,
		PaymentNotes:    req.PaymentNotes,
		Notes:           req.Notes,
		CreatedBy:       &createdBy,
		CreatedAt:       time.Now(),
		Items:           processedItems,
	}

	return purchase, nil
}

// GetAll retrieves all purchases ordered by date descending
// Fungsi ini mengambil riwayat semua pembelian
func (r *PurchaseRepository) GetAll() ([]models.Purchase, error) {
	query := `
		SELECT 
			p.id, p.supplier_id, p.supplier_name,
			p.total_amount,
			COALESCE(p.payment_method, 'cash'),
			COALESCE(p.payment_status, 'paid'),
			COALESCE(p.paid_amount, p.total_amount),
			COALESCE(p.remaining_amount, 0),
			to_char(p.due_date, 'YYYY-MM-DD'),
			p.payment_notes, p.notes, p.created_by, p.created_at
		FROM purchases p
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat pembelian: %w", err)
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var p models.Purchase
		var supplierID sql.NullInt64
		var supplierName sql.NullString
		var dueDate sql.NullString
		var paymentNotes sql.NullString
		var notes sql.NullString
		var createdBy sql.NullInt64

		err := rows.Scan(
			&p.ID, &supplierID, &supplierName,
			&p.TotalAmount,
			&p.PaymentMethod, &p.PaymentStatus,
			&p.PaidAmount, &p.RemainingAmount,
			&dueDate, &paymentNotes, &notes, &createdBy, &p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data pembelian: %w", err)
		}
		if supplierID.Valid {
			id := int(supplierID.Int64)
			p.SupplierID = &id
		}
		if supplierName.Valid {
			p.SupplierName = &supplierName.String
		}
		if dueDate.Valid {
			p.DueDate = &dueDate.String
		}
		if paymentNotes.Valid {
			p.PaymentNotes = &paymentNotes.String
		}
		if notes.Valid {
			p.Notes = &notes.String
		}
		if createdBy.Valid {
			id := int(createdBy.Int64)
			p.CreatedBy = &id
		}

		purchases = append(purchases, p)
	}

	if purchases == nil {
		purchases = []models.Purchase{}
	}

	return purchases, nil
}

// GetByID retrieves a purchase by ID with its items
// Fungsi ini mengambil detail 1 pembelian beserta item-itemnya
func (r *PurchaseRepository) GetByID(id int) (*models.Purchase, error) {
	var p models.Purchase
	var supplierID sql.NullInt64
	var supplierName sql.NullString
	var dueDate sql.NullString
	var paymentNotes sql.NullString
	var notes sql.NullString
	var createdBy sql.NullInt64

	err := r.db.QueryRow(`
		SELECT id, supplier_id, supplier_name,
		       total_amount,
		       COALESCE(payment_method, 'cash'),
		       COALESCE(payment_status, 'paid'),
		       COALESCE(paid_amount, total_amount),
		       COALESCE(remaining_amount, 0),
		       to_char(due_date, 'YYYY-MM-DD'),
		       payment_notes, notes, created_by, created_at
		FROM purchases WHERE id = $1`, id,
	).Scan(&p.ID, &supplierID, &supplierName,
		&p.TotalAmount,
		&p.PaymentMethod, &p.PaymentStatus,
		&p.PaidAmount, &p.RemainingAmount,
		&dueDate, &paymentNotes, &notes, &createdBy, &p.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pembelian dengan ID %d tidak ditemukan", id)
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pembelian: %w", err)
	}

	if supplierID.Valid {
		id2 := int(supplierID.Int64)
		p.SupplierID = &id2
	}
	if supplierName.Valid {
		p.SupplierName = &supplierName.String
	}
	if dueDate.Valid {
		p.DueDate = &dueDate.String
	}
	if paymentNotes.Valid {
		p.PaymentNotes = &paymentNotes.String
	}
	if notes.Valid {
		p.Notes = &notes.String
	}
	if createdBy.Valid {
		id2 := int(createdBy.Int64)
		p.CreatedBy = &id2
	}

	// 2. Ambil detail items
	queryItems := `
		SELECT id, purchase_id, product_id, product_name, quantity, buy_price, sell_price, category_id, subtotal, created_at
		FROM purchase_items 
		WHERE purchase_id = $1 
		ORDER BY id
	`

	rows, err := r.db.Query(queryItems, p.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail pembelian: %w", err)
	}
	defer rows.Close()

	var items []models.PurchaseItem
	for rows.Next() {
		var item models.PurchaseItem
		var productID sql.NullInt64
		var sellPrice sql.NullFloat64
		var categoryID sql.NullInt64

		err := rows.Scan(
			&item.ID, &item.PurchaseID, &productID, &item.ProductName,
			&item.Quantity, &item.BuyPrice, &sellPrice, &categoryID,
			&item.Subtotal, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca detail item: %w", err)
		}

		if productID.Valid {
			id := int(productID.Int64)
			item.ProductID = &id
		}
		if sellPrice.Valid {
			item.SellPrice = &sellPrice.Float64
		}
		if categoryID.Valid {
			id := int(categoryID.Int64)
			item.CategoryID = &id
		}

		items = append(items, item)
	}

	if items == nil {
		items = []models.PurchaseItem{}
	}
	p.Items = items

	return &p, nil
}

// GetTotalPengeluaran retrieves total purchase amount for a date range
// Fungsi ini menghitung total pengeluaran (pembelian) untuk laporan
func (r *PurchaseRepository) GetTotalPengeluaran(startDate, endDate time.Time) (float64, int, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_pengeluaran,
			COUNT(*) as total_pembelian
		FROM purchases
		WHERE created_at BETWEEN $1 AND $2
	`

	var totalPengeluaran float64
	var totalPembelian int
	err := r.db.QueryRow(query, startDate, endDate).Scan(&totalPengeluaran, &totalPembelian)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal menghitung total pengeluaran: %w", err)
	}

	return totalPengeluaran, totalPembelian, nil
}
