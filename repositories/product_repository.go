package repositories

import (
	"database/sql"     // Package standard Go untuk database SQL
	"fmt"              // Package untuk formatting string
	"kasir-api/models" // Import models untuk struct Product
)

// ProductRepository handles database operations for products
// Struct ini menyimpan koneksi database
type ProductRepository struct {
	db *sql.DB // Pointer ke database connection
}

// NewProductRepository creates a new ProductRepository
// Fungsi ini adalah "constructor" untuk membuat instance ProductRepository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db} // Return struct dengan db yang sudah di-inject
}

// GetAll retrieves all products from database with pagination
// Fungsi ini mengambil produk dari table products dengan pagination
// Parameter searchName untuk filter by name (kosong = ambil semua)
// Parameter pagination untuk limit dan offset
// Return: products, total count, error
func (r *ProductRepository) GetAll(searchName string, searchBarcode string, pagination *models.PaginationParams) ([]models.Product, int, error) {
	// SQL query dengan LEFT JOIN ke table categories
	// LEFT JOIN = ambil semua products, meskipun tidak punya category
	query := `
		SELECT 
			p.id, 
			p.nama, 
			p.harga, 
			p.stok, 
			p.barcode,
			p.category_id,
			p.harga_beli,
			p.default_discount_type,
			p.default_discount_value,
			p.is_featured,
			p.created_by,
			c.id as category_id_full,
			c.nama as category_name,
			c.description as category_description,
			COALESCE(c.discount_type, '') as category_discount_type,
			COALESCE(c.discount_value, 0) as category_discount_value
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`

	// Query untuk count total items
	countQuery := `SELECT COUNT(*) FROM products p`

	// Buat slice untuk menampung arguments query
	var args []interface{}
	var countArgs []interface{}
	paramIndex := 1
	countParamIndex := 1
	hasWhere := false

	// Filter by barcode (exact match, prioritas tertinggi)
	if searchBarcode != "" {
		query += fmt.Sprintf(" WHERE p.barcode = $%d", paramIndex)
		countQuery += fmt.Sprintf(" WHERE p.barcode = $%d", countParamIndex)
		args = append(args, searchBarcode)
		countArgs = append(countArgs, searchBarcode)
		paramIndex++
		countParamIndex++
		hasWhere = true
	}

	// Filter by name (prefix match)
	if searchName != "" && !hasWhere {
		query += fmt.Sprintf(" WHERE p.nama ILIKE $%d", paramIndex)
		countQuery += fmt.Sprintf(" WHERE p.nama ILIKE $%d", countParamIndex)
		args = append(args, searchName+"%")
		countArgs = append(countArgs, searchName+"%")
		paramIndex++
		countParamIndex++
	}

	// Hitung total items untuk pagination metadata
	var totalItems int
	if len(countArgs) == 0 {
		countArgs = nil
	}
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Tambahkan ORDER BY untuk konsistensi
	query += " ORDER BY p.id DESC"

	// Tambahkan LIMIT dan OFFSET untuk pagination
	if pagination != nil {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
		args = append(args, pagination.Limit, pagination.GetOffset())
	}

	// Execute query dan dapatkan rows (banyak baris)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err // Kalau error, return nil dan error
	}
	defer rows.Close() // Pastikan rows di-close setelah selesai (penting!)

	// Buat slice kosong untuk menampung products
	var products []models.Product

	// Loop semua rows yang didapat dari database
	for rows.Next() {
		var product models.Product           // Buat variable product untuk setiap row
		var categoryID sql.NullInt64         // Untuk handle NULL dari LEFT JOIN
		var categoryName sql.NullString      // Untuk handle NULL dari LEFT JOIN
		var categoryDesc sql.NullString      // Untuk handle NULL dari LEFT JOIN
		var catDiscType string               // Diskon kategori
		var catDiscValue float64             // Nilai diskon kategori
		var hargaBeli sql.NullFloat64        // Untuk handle NULL dari harga_beli
		var barcode sql.NullString           // Untuk handle NULL dari barcode
		var defaultDiscType sql.NullString   // Untuk handle NULL dari default_discount_type
		var defaultDiscValue sql.NullFloat64 // Untuk handle NULL dari default_discount_value
		var createdBy sql.NullInt64          // Untuk handle NULL dari created_by

		// Scan data dari row ke struct product
		// Urutan harus sama dengan SELECT
		err := rows.Scan(
			&product.ID,
			&product.Nama,
			&product.Harga,
			&product.Stok,
			&barcode,
			&product.CategoryID,
			&hargaBeli,
			&defaultDiscType,
			&defaultDiscValue,
			&product.IsFeatured,
			&createdBy,
			&categoryID,
			&categoryName,
			&categoryDesc,
			&catDiscType,
			&catDiscValue,
		)
		if err != nil {
			return nil, 0, err // Kalau scan error, return error
		}

		// Set harga_beli jika valid
		if hargaBeli.Valid {
			product.HargaBeli = &hargaBeli.Float64
		}

		// Set barcode jika valid
		if barcode.Valid {
			product.Barcode = &barcode.String
		}

		// Set default discount jika valid
		if defaultDiscType.Valid {
			product.DefaultDiscountType = &defaultDiscType.String
		}
		if defaultDiscValue.Valid {
			product.DefaultDiscountValue = &defaultDiscValue.Float64
		}

		// Set created_by jika valid
		if createdBy.Valid {
			id := int(createdBy.Int64)
			product.CreatedBy = &id
		}

		// Jika ada category, populate Category struct dengan semua field termasuk diskon
		if categoryName.Valid && categoryID.Valid {
			cat := &models.Category{
				ID:            int(categoryID.Int64),
				Nama:          categoryName.String,
				Description:   categoryDesc.String,
				DiscountValue: catDiscValue,
			}
			if catDiscType != "" {
				cat.DiscountType = &catDiscType
			}
			product.Category = cat
		}

		// Tambahkan product ke slice products
		products = append(products, product)
	}

	return products, totalItems, nil // Return slice products, total count, dan nil (no error)
}

// GetByID retrieves a product by ID
// Fungsi ini mengambil 1 produk berdasarkan ID
func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	// SQL query dengan LEFT JOIN ke categories untuk mendapatkan category lengkap
	// $1 akan diganti dengan value id saat execute
	query := `
		SELECT 
			p.id, 
			p.nama, 
			p.harga, 
			p.stok, 
			p.barcode,
			p.category_id,
			p.harga_beli,
			p.default_discount_type,
			p.default_discount_value,
			p.is_featured,
			p.created_by,
			c.id as category_id_full,
			c.nama as category_name,
			c.description as category_description,
			COALESCE(c.discount_type, '') as category_discount_type,
			COALESCE(c.discount_value, 0) as category_discount_value
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	// QueryRow untuk query yang return 1 row saja
	row := r.db.QueryRow(query, id) // id akan replace $1

	var product models.Product           // Buat variable untuk menampung hasil
	var categoryID sql.NullInt64         // Untuk handle NULL dari LEFT JOIN
	var categoryName sql.NullString      // Untuk handle NULL dari LEFT JOIN
	var categoryDesc sql.NullString      // Untuk handle NULL dari LEFT JOIN
	var catDiscType string               // Diskon kategori
	var catDiscValue float64             // Nilai diskon kategori
	var hargaBeli sql.NullFloat64        // Untuk handle NULL dari harga_beli
	var barcode sql.NullString           // Untuk handle NULL dari barcode
	var defaultDiscType sql.NullString   // Untuk handle NULL dari default_discount_type
	var defaultDiscValue sql.NullFloat64 // Untuk handle NULL dari default_discount_value
	var createdBy sql.NullInt64          // Untuk handle NULL dari created_by

	// Scan hasil query ke struct product
	err := row.Scan(
		&product.ID,
		&product.Nama,
		&product.Harga,
		&product.Stok,
		&barcode,
		&product.CategoryID,
		&hargaBeli,
		&defaultDiscType,
		&defaultDiscValue,
		&product.IsFeatured,
		&createdBy,
		&categoryID,
		&categoryName,
		&categoryDesc,
		&catDiscType,
		&catDiscValue,
	)
	if err != nil {
		return nil, err // Kalau tidak ketemu atau error, return nil
	}

	// Set harga_beli jika valid
	if hargaBeli.Valid {
		product.HargaBeli = &hargaBeli.Float64
	}

	// Set barcode jika valid
	if barcode.Valid {
		product.Barcode = &barcode.String
	}

	// Set default discount jika valid
	if defaultDiscType.Valid {
		product.DefaultDiscountType = &defaultDiscType.String
	}
	if defaultDiscValue.Valid {
		product.DefaultDiscountValue = &defaultDiscValue.Float64
	}

	// Set created_by jika valid
	if createdBy.Valid {
		id := int(createdBy.Int64)
		product.CreatedBy = &id
	}

	// Jika ada category, populate Category struct dengan semua field
	if categoryName.Valid && categoryID.Valid {
		cat := &models.Category{
			ID:            int(categoryID.Int64),
			Nama:          categoryName.String,
			Description:   categoryDesc.String,
			DiscountValue: catDiscValue,
		}
		if catDiscType != "" {
			cat.DiscountType = &catDiscType
		}
		product.Category = cat
	}

	return &product, nil // Return pointer ke product
}

// Create adds a new product to database or updates stock if product exists
// Fungsi ini menambahkan produk baru ATAU menambah stok jika produk dengan nama sama sudah ada
func (r *ProductRepository) Create(product *models.Product) error {
	// SQL query dengan UPSERT logic (INSERT ... ON CONFLICT ... DO UPDATE)
	// Jika produk dengan nama yang sama sudah ada, maka:
	// - Stok akan ditambahkan (stok lama + stok baru)
	// - Harga akan diupdate dengan harga terbaru
	// - CategoryID akan diupdate jika diberikan
	// - HargaBeli akan diupdate (EXCLUDED.harga_beli)
	// Jika belum ada, akan insert produk baru
	query := `
		INSERT INTO products (nama, harga, stok, category_id, harga_beli, created_by, barcode, default_discount_type, default_discount_value, is_featured) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (nama) 
		DO UPDATE SET 
			harga = EXCLUDED.harga,
			stok = products.stok + EXCLUDED.stok,
			category_id = EXCLUDED.category_id,
			harga_beli = EXCLUDED.harga_beli,
			barcode = EXCLUDED.barcode,
			default_discount_type = EXCLUDED.default_discount_type,
			default_discount_value = EXCLUDED.default_discount_value,
			is_featured = EXCLUDED.is_featured
		RETURNING id, stok
	`

	// Execute query dan scan ID + stok terbaru yang di-return
	err := r.db.QueryRow(query, product.Nama, product.Harga, product.Stok, product.CategoryID, product.HargaBeli, product.CreatedBy, product.Barcode, product.DefaultDiscountType, product.DefaultDiscountValue, product.IsFeatured).Scan(&product.ID, &product.Stok)

	return err // Return error (nil kalau sukses)
}

// Update updates product info (nama, harga jual, kategori only)
// Stok dikelola lewat pembelian (POST /api/purchases) dan penjualan (POST /api/checkout)
// Harga beli dikelola lewat pembelian (POST /api/purchases)
func (r *ProductRepository) Update(product *models.Product) error {
	// SQL query untuk UPDATE — nama, harga jual, kategori, barcode, diskon default, dan is_featured
	// Stok dan harga_beli TIDAK bisa diubah dari sini
	query := "UPDATE products SET nama = $1, harga = $2, category_id = $3, barcode = $4, default_discount_type = $5, default_discount_value = $6, is_featured = $7 WHERE id = $8"

	_, err := r.db.Exec(query, product.Nama, product.Harga, product.CategoryID, product.Barcode, product.DefaultDiscountType, product.DefaultDiscountValue, product.IsFeatured, product.ID)

	return err
}

// GetByBarcode retrieves a product by barcode (exact match)
// Fungsi ini mengambil 1 produk berdasarkan barcode
func (r *ProductRepository) GetByBarcode(barcode string) (*models.Product, error) {
	query := `
		SELECT 
			p.id, 
			p.nama, 
			p.harga, 
			p.stok, 
			p.barcode,
			p.category_id,
			p.harga_beli,
			p.default_discount_type,
			p.default_discount_value,
			p.is_featured,
			p.created_by,
			c.id as category_id_full,
			c.nama as category_name,
			c.description as category_description,
			COALESCE(c.discount_type, '') as category_discount_type,
			COALESCE(c.discount_value, 0) as category_discount_value
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.barcode = $1
	`

	row := r.db.QueryRow(query, barcode)

	var product models.Product
	var categoryID sql.NullInt64
	var categoryName sql.NullString
	var categoryDesc sql.NullString
	var catDiscType string
	var catDiscValue float64
	var hargaBeli sql.NullFloat64
	var barcodeVal sql.NullString
	var defaultDiscType sql.NullString
	var defaultDiscValue sql.NullFloat64
	var createdBy sql.NullInt64

	err := row.Scan(
		&product.ID,
		&product.Nama,
		&product.Harga,
		&product.Stok,
		&barcodeVal,
		&product.CategoryID,
		&hargaBeli,
		&defaultDiscType,
		&defaultDiscValue,
		&product.IsFeatured,
		&createdBy,
		&categoryID,
		&categoryName,
		&categoryDesc,
		&catDiscType,
		&catDiscValue,
	)
	if err != nil {
		return nil, err
	}

	if hargaBeli.Valid {
		product.HargaBeli = &hargaBeli.Float64
	}
	if barcodeVal.Valid {
		product.Barcode = &barcodeVal.String
	}
	if defaultDiscType.Valid {
		product.DefaultDiscountType = &defaultDiscType.String
	}
	if defaultDiscValue.Valid {
		product.DefaultDiscountValue = &defaultDiscValue.Float64
	}
	if createdBy.Valid {
		id := int(createdBy.Int64)
		product.CreatedBy = &id
	}
	if categoryName.Valid && categoryID.Valid {
		cat := &models.Category{
			ID:            int(categoryID.Int64),
			Nama:          categoryName.String,
			Description:   categoryDesc.String,
			DiscountValue: catDiscValue,
		}
		if catDiscType != "" {
			cat.DiscountType = &catDiscType
		}
		product.Category = cat
	}

	return &product, nil
}

// Delete removes a product from database
// Fungsi ini menghapus produk dari database
func (r *ProductRepository) Delete(id int) error {
	// SQL query untuk DELETE
	// WHERE untuk kondisi (hapus produk dengan id tertentu)
	query := "DELETE FROM products WHERE id = $1"

	// Execute query
	// $1 = id
	_, err := r.db.Exec(query, id)

	return err // Return error (nil kalau sukses)
}
