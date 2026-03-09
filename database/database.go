package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// InitDB menginisialisasi koneksi database menggunakan database/sql
func InitDB(dsn string) *sql.DB {
	log.Println("🔌 Connecting to database...")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Connection pool settings untuk performa optimal dengan cloud DB
	db.SetMaxOpenConns(10)                 // Max koneksi aktif
	db.SetMaxIdleConns(5)                  // Max koneksi idle (siap pakai)
	db.SetConnMaxLifetime(5 * time.Minute) // Recycle koneksi setiap 5 menit
	db.SetConnMaxIdleTime(2 * time.Minute) // Tutup koneksi idle > 2 menit

	// Test koneksi dengan Ping
	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}

	// Buat table jika belum ada
	createTables(db)

	fmt.Println("✅ Database connected successfully")
	return db
}

// createTables membuat table dasar jika belum ada
func createTables(db *sql.DB) {
	productsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(255) NOT NULL,
		harga INTEGER NOT NULL,
		stok INTEGER NOT NULL
	)`

	_, err := db.Exec(productsTable)
	if err != nil {
		log.Fatal("❌ Failed to create products table:", err)
	}

	categoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(255) NOT NULL,
		description TEXT
	)`

	_, err = db.Exec(categoriesTable)
	if err != nil {
		log.Fatal("❌ Failed to create categories table:", err)
	}

	transactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		total_amount DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(transactionsTable)
	if err != nil {
		log.Fatal("❌ Failed to create transactions table:", err)
	}

	transactionDetailsTable := `
	CREATE TABLE IF NOT EXISTS transaction_details (
		id SERIAL PRIMARY KEY,
		transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
		product_id INTEGER NOT NULL REFERENCES products(id),
		quantity INTEGER NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		subtotal DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(transactionDetailsTable)
	if err != nil {
		log.Fatal("❌ Failed to create transaction_details table:", err)
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_transaction_details_transaction_id ON transaction_details(transaction_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transaction_details_product_id ON transaction_details(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at)`,
	}

	for _, indexSQL := range indexes {
		_, err = db.Exec(indexSQL)
		if err != nil {
			log.Fatal("❌ Failed to create index:", err)
		}
	}

	fmt.Println("✅ Tables created/verified successfully")
}
