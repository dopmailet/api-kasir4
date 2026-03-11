package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"kasir-api/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Peringatan: Error loading .env file, continuing with existing env vars")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// 2. Koneksi ke Database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	fmt.Println("Berhasil terhubung ke database. Memulai pembuatan akun Superadmin...")

	// 3. Setup Detail Akun Superadmin
	username := "superadmin"
	rawPassword := "SuperAdmin123!"
	fullName := "Super Administrator"
	role := "superadmin"
	storeID := 1
	isActive := true
	isSuperadmin := true

	// 4. Hash Password
	hashedPassword, err := utils.HashPassword(rawPassword)
	if err != nil {
		log.Fatalf("Gagal melakukan hashing password: %v", err)
	}

	// 5. Cek apakah user superadmin sudah ada (menghindari duplikat)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Fatalf("Gagal mengecek username: %v", err)
	}

	if count > 0 {
		log.Printf("Gagal: Akun dengan username '%s' sudah ada di database.\n", username)
		return
	}

	// 6. Masukkan ke tabel users
	// Catatan: Pastikan kolom is_superadmin ada di tabel users Anda
	query := `
		INSERT INTO users (username, password, nama_lengkap, role, is_active, store_id, is_superadmin)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var newUserID int
	err = db.QueryRow(
		query,
		username,
		hashedPassword,
		fullName,
		role,
		isActive,
		storeID,
		isSuperadmin,
	).Scan(&newUserID)

	if err != nil {
		log.Fatalf("Gagal menyimpan akun Superadmin ke database: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Println("✅ AKUN SUPERADMIN BERHASIL DIBUAT!")
	fmt.Println("==================================================")
	fmt.Printf("ID User     : %d\n", newUserID)
	fmt.Printf("Username    : %s\n", username)
	fmt.Printf("Password    : %s\n", rawPassword)
	fmt.Printf("Store ID    : %d (Toko Default)\n", storeID)
	fmt.Println("==================================================")
	fmt.Println("⚠️  SANGAT PENTING: Segera login dan ganti password Anda setelah deploy ke production!")
}
