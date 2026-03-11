package repositories

import (
	"database/sql"     // Package standard Go untuk database SQL
	"kasir-api/models" // Import models untuk struct Category
)

// CategoryRepository handles database operations for categories
// Struct ini menyimpan koneksi database
type CategoryRepository struct {
	db *sql.DB // Pointer ke database connection
}

// NewCategoryRepository creates a new CategoryRepository
// Fungsi ini adalah "constructor" untuk membuat instance CategoryRepository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db} // Return struct dengan db yang sudah di-inject
}

// GetAll retrieves all categories from database
// Fungsi ini mengambil semua kategori dari table categories
func (r *CategoryRepository) GetAll(storeID int) ([]models.Category, error) {
	// SQL query untuk select semua kolom dari table categories
	query := "SELECT id, nama, description, COALESCE(discount_type, '') as discount_type, COALESCE(discount_value, 0) as discount_value FROM categories WHERE store_id = $1"

	// Execute query dan dapatkan rows (banyak baris)
	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err // Kalau error, return nil dan error
	}
	defer rows.Close() // Pastikan rows di-close setelah selesai

	// Buat slice kosong untuk menampung categories
	var categories []models.Category

	// Loop semua rows yang didapat dari database
	for rows.Next() {
		var category models.Category // Buat variable category untuk setiap row
		var discType string

		// Scan data dari row ke struct category
		err := rows.Scan(&category.ID, &category.Nama, &category.Description, &discType, &category.DiscountValue)
		if err != nil {
			return nil, err // Kalau scan error, return error
		}

		// Set discount_type jika tidak kosong
		if discType != "" {
			category.DiscountType = &discType
		}

		// Tambahkan category ke slice categories
		categories = append(categories, category)
	}

	return categories, nil // Return slice categories dan nil (no error)
}

// GetByID retrieves a category by ID with its products
// Fungsi ini mengambil 1 kategori berdasarkan ID beserta semua products dalam category tersebut
func (r *CategoryRepository) GetByID(id int, storeID int) (*models.Category, error) {
	// 1. Ambil category data
	query := "SELECT id, nama, description, COALESCE(discount_type, '') as discount_type, COALESCE(discount_value, 0) as discount_value FROM categories WHERE id = $1 AND store_id = $2"
	row := r.db.QueryRow(query, id, storeID)

	var category models.Category
	var discType string
	err := row.Scan(&category.ID, &category.Nama, &category.Description, &discType, &category.DiscountValue)
	if err != nil {
		return nil, err // Kalau tidak ketemu atau error, return nil
	}

	// Set discount_type jika tidak kosong
	if discType != "" {
		category.DiscountType = &discType
	}

	// 2. Ambil semua products yang punya category_id ini
	productsQuery := "SELECT id, nama, harga, stok FROM products WHERE category_id = $1 AND store_id = $2"
	rows, err := r.db.Query(productsQuery, id, storeID)
	if err != nil {
		// Kalau error query products, tetap return category (tanpa products)
		return &category, nil
	}
	defer rows.Close()

	// 3. Loop semua products dan tambahkan ke category.Products
	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Nama, &product.Harga, &product.Stok)
		if err != nil {
			continue // Skip product yang error
		}
		products = append(products, product)
	}

	// 4. Set products ke category
	category.Products = products

	return &category, nil
}

// Create adds a new category to database
// Fungsi ini menambahkan kategori baru ke database
func (r *CategoryRepository) Create(category *models.Category) error {
	// SQL query untuk INSERT
	// RETURNING id = return ID yang baru dibuat (auto-increment)
	query := "INSERT INTO categories (nama, description, discount_type, discount_value, store_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	// Execute query dan langsung scan ID yang di-return
	err := r.db.QueryRow(query, category.Nama, category.Description, category.DiscountType, category.DiscountValue, category.StoreID).Scan(&category.ID)

	return err // Return error (nil kalau sukses)
}

// Update updates an existing category
// Fungsi ini mengupdate kategori yang sudah ada
func (r *CategoryRepository) Update(category *models.Category) error {
	// SQL query untuk UPDATE
	// SET untuk set nilai baru termasuk discount
	// WHERE untuk kondisi (update kategori dengan id tertentu)
	query := "UPDATE categories SET nama = $1, description = $2, discount_type = $3, discount_value = $4 WHERE id = $5 AND store_id = $6"

	// Execute query
	_, err := r.db.Exec(query, category.Nama, category.Description, category.DiscountType, category.DiscountValue, category.ID, category.StoreID)

	return err // Return error (nil kalau sukses)
}

// Delete removes a category from database
// Fungsi ini menghapus kategori dari database
func (r *CategoryRepository) Delete(id int, storeID int) error {
	// SQL query untuk DELETE
	// WHERE untuk kondisi (hapus kategori dengan id tertentu)
	query := "DELETE FROM categories WHERE id = $1 AND store_id = $2"

	// Execute query
	// $1 = id
	_, err := r.db.Exec(query, id, storeID)

	return err // Return error (nil kalau sukses)
}
