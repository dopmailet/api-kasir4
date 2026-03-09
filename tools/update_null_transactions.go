package main

import (
	"fmt"
	"kasir-api/config"
	"kasir-api/database"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Warning: Gagal load config: %v", err)
	}

	// Initialize Database Connection
	db := database.InitDB(cfg.GetDatabaseURL())

	// Cari ID Admin pertama
	var adminID int
	err = db.QueryRow("SELECT id FROM users WHERE role = 'admin' ORDER BY id ASC LIMIT 1").Scan(&adminID)
	if err != nil {
		log.Fatalf("Gagal mencari ID admin: %v", err)
	}

	fmt.Printf("Admin ID ditemukan: %d\n", adminID)

	// Update transactions yang masih NULL
	query := "UPDATE transactions SET created_by = $1 WHERE created_by IS NULL"
	res, err := db.Exec(query, adminID)
	if err != nil {
		log.Fatalf("Gagal mengupdate transaksi: %v", err)
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("✅ Berhasil mengupdate %d transaksi lama dengan created_by = %d\n", rowsAffected, adminID)
}
