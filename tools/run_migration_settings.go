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
		// Coba load dari current directory jika gagal
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

	fmt.Println("⏳ Memulai migrasi tabel app_settings untuk multi-tenant...")

	// 4. Skrip SQL untuk Migrasi
	// Drop old PK, add store_id, set new composite PK
	queries := []string{
		// Tambahkan kolom store_id jika belum ada (default ke Toko 1)
		`ALTER TABLE app_settings ADD COLUMN IF NOT EXISTS store_id INT NOT NULL DEFAULT 1;`,

		// Hapuskan Primary Key lama (key)
		`ALTER TABLE app_settings DROP CONSTRAINT IF EXISTS app_settings_pkey;`,

		// Buat Composite Primary Key baru (store_id, key)
		`ALTER TABLE app_settings ADD PRIMARY KEY (store_id, key);`,

		// Tambahkan Foreign Key ke tabel stores
		`ALTER TABLE app_settings ADD CONSTRAINT fk_app_settings_store FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE;`,
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

	fmt.Println("🚀 Migrasi app_settings selesai!")
}
