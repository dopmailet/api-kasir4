package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("⚠️  Warning: Error loading .env file, relying on system environment variables")
		}
	}

	// 2. Dapatkan koneksi string
	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		dbConn = os.Getenv("DATABASE_URL")
	}
	if dbConn == "" {
		log.Fatal("❌ DB_CONN or DATABASE_URL is not set in environment variables")
	}

	// 3. Connect ke database
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal("❌ Gagal koneksi ke database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Gagal ping database:", err)
	}

	fmt.Println("⏳ Memulai migrasi tabel subscription_packages untuk Landing Page SaaS...")

	// 4. Skrip SQL untuk Migrasi
	queries := []string{
		// Tambahkan 7 Kolom Baru (TIDAK ADA store_id karena ini tabel Global)
		`ALTER TABLE subscription_packages
			ADD COLUMN IF NOT EXISTS description VARCHAR(500) DEFAULT NULL,
			ADD COLUMN IF NOT EXISTS features JSONB DEFAULT '[]'::jsonb,
			ADD COLUMN IF NOT EXISTS period VARCHAR(20) DEFAULT '/bulan',
			ADD COLUMN IF NOT EXISTS discount_percent DECIMAL(5,2) DEFAULT 0,
			ADD COLUMN IF NOT EXISTS discount_label VARCHAR(100) DEFAULT NULL,
			ADD COLUMN IF NOT EXISTS is_popular BOOLEAN DEFAULT FALSE,
			ADD COLUMN IF NOT EXISTS sort_order INTEGER DEFAULT 0,
			ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW(),
			ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();`,

		// Seed Data (Opsional, diinject ulang kalau masih kosong fitur-fiturnya)
		`INSERT INTO subscription_packages 
			(id, name, max_kasir, max_products, price, is_active, description, features, period, discount_percent, discount_label, is_popular, sort_order)
		VALUES
			(1, 'Free', 1, 100, 0, true, NULL, '["1 Toko", "1 Akses Kasir", "Laporan Harian Dasar", "Manajemen Produk"]'::jsonb, 'selamanya', 0, NULL, false, 1),
			(2, 'Basic', 1, 500, 100000, true, NULL, '["1 Toko", "1 Kasir", "Laporan & Analitik Lanjut", "Manajemen Gaji Karyawan", "Export Excel / PDF"]'::jsonb, '/bulan', 0, NULL, false, 2),
			(3, 'Pro', 3, 1000, 150000, true, NULL, '["1 Toko", "3 Kasir", "Laporan & Analitik Lanjut", "Manajemen Gaji Karyawan", "Export Excel / PDF", "Prioritas Dukungan"]'::jsonb, '/bulan', 0, NULL, true, 3)
		ON CONFLICT (id) DO UPDATE SET 
			description = EXCLUDED.description,
			features = EXCLUDED.features,
			period = EXCLUDED.period,
			discount_percent = EXCLUDED.discount_percent,
			discount_label = EXCLUDED.discount_label,
			is_popular = EXCLUDED.is_popular,
			sort_order = EXCLUDED.sort_order;`,

		// Mengkalibrasi sequence numbering (Mencegah Duplicate Key ID pasca seed)
		`SELECT setval(pg_get_serial_sequence('subscription_packages', 'id'), coalesce(max(id),0) + 1, false) FROM subscription_packages;`,
	}

	// 5. Eksekusi setiap query
	for i, q := range queries {
		fmt.Printf("Mengeksekusi query %d...\n", i+1)
		_, err := db.Exec(q)
		if err != nil {
			log.Printf("⚠️ Gagal mengeksekusi query %d: %v\n", i+1, err)
			log.Println("Catatan: Jika errornya karena constraint sudah ada, Anda dapat mengabaikannya.")
		} else {
			fmt.Printf("✅ Query %d berhasil.\n", i+1)
		}
	}

	fmt.Println("🚀 Migrasi subscription_packages selesai!")
}
