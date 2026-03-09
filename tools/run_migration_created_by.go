package main

import (
	"fmt"
	"kasir-api/config"
	"kasir-api/database"
	"log"
	"os"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect ke DB
	db := database.InitDB(cfg.GetDatabaseURL())
	defer db.Close()

	// Baca file SQL
	sqlScript, err := os.ReadFile("database/migrations/add_created_by_to_transactions.sql")
	if err != nil {
		log.Fatal("Gagal membaca file SQL:", err)
	}

	// Jalankan migrasi
	_, err = db.Exec(string(sqlScript))
	if err != nil {
		log.Fatal("Gagal menjalankan migrasi:", err)
	}

	fmt.Println("✅ Migrasi add_created_by_to_transactions.sql berhasil dijalankan!")
}
