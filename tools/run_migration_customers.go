package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// run_migration_customers.go menambahkan tabel customers, loyalty_transactions, app_settings
// dan memodifikasi tabel transactions.
func main() {
	log.Println("Memulai migrasi fitur Customer Database & Loyalty...")

	// 1. Load env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("⚠️  Warning: Error loading .env file (mengabaikan jika berjalan di production atau Docker)")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("❌ DB_CONN environment variable is required")
	}

	// 2. Connect
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Database tidak merespon: %v", err)
	}
	log.Println("✅ Berhasil terhubung ke database")

	// 3. Menjalankan DDL
	queries := []string{
		`CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			customer_code VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(150) NOT NULL,
			phone VARCHAR(50) NOT NULL,
			address TEXT,
			notes TEXT,
			loyalty_points INTEGER NOT NULL DEFAULT 0,
			total_spent NUMERIC(18,2) NOT NULL DEFAULT 0,
			total_transactions INTEGER NOT NULL DEFAULT 0,
			last_transaction_at TIMESTAMPTZ,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);`,
		`CREATE INDEX IF NOT EXISTS idx_customers_active ON customers(is_active);`,
		`CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name);`,

		`CREATE TABLE IF NOT EXISTS loyalty_transactions (
			id SERIAL PRIMARY KEY,
			customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
			transaction_id INTEGER REFERENCES transactions(id) ON DELETE SET NULL,
			type VARCHAR(20) NOT NULL CHECK (type IN ('earn', 'adjust')),
			points INTEGER NOT NULL,
			description TEXT,
			created_by_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_loyalty_transactions_customer_id ON loyalty_transactions(customer_id);`,

		`CREATE TABLE IF NOT EXISTS app_settings (
			key VARCHAR(100) PRIMARY KEY,
			value_json JSONB NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,

		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS customer_id INTEGER REFERENCES customers(id) ON DELETE SET NULL;`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_customer_id ON transactions(customer_id);`,

		// Trigger updated_at customers
		`CREATE OR REPLACE FUNCTION update_customers_modtime()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;`,

		`DROP TRIGGER IF EXISTS trg_customers_updated_at ON customers;`,
		`CREATE TRIGGER trg_customers_updated_at
		BEFORE UPDATE ON customers
		FOR EACH ROW EXECUTE FUNCTION update_customers_modtime();`,

		// Seed initial setup
		`INSERT INTO app_settings (key, value_json) 
		VALUES ('customer_settings', '{"showCustomerInPOS": true, "enableLoyaltyPoints": true}')
		ON CONFLICT (key) DO NOTHING;`,
	}

	for i, q := range queries {
		log.Printf("⏳ Menjalankan query %d/%d...", i+1, len(queries))
		_, err := pool.Exec(ctx, q)
		if err != nil {
			log.Fatalf("❌ Gagal eksekusi query %d: %v\nQuery: %s", i+1, err, q)
		}
	}

	log.Println("🎉 Migrasi fitur Customer & Loyalty selesai dengan sukses!")
}
