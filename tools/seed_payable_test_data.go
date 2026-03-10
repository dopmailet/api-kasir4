package main

import (
	"database/sql"
	"fmt"
	"kasir-api/config"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Config Error:", err)
	}

	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Connection failed:", err)
	}

	fmt.Println("🛠️ Seeding Test Data...")

	// 1. Insert a supplier
	var supplierID int
	err = db.QueryRow(`
		INSERT INTO suppliers (name, contact_person, is_active)
		VALUES ('Test Supplier PT', 'Budi', true)
		RETURNING id
	`).Scan(&supplierID)
	if err != nil {
		log.Fatal("❌ Insert supplier failed:", err)
	}

	// 2. Insert a supplier_payable
	var payableID int
	err = db.QueryRow(`
		INSERT INTO supplier_payables (supplier_id, amount, paid_amount, status, due_date)
		VALUES ($1, 500000, 100000, 'partial', '2026-12-31')
		RETURNING id
	`, supplierID).Scan(&payableID)
	if err != nil {
		log.Fatal("❌ Insert supplier_payable failed:", err)
	}

	// 3. Insert a payable_payment
	_, err = db.Exec(`
		INSERT INTO payable_payments (payable_id, amount, payment_date, notes)
		VALUES ($1, 100000, '2026-03-01', 'Cicilan pertama')
	`, payableID)
	if err != nil {
		log.Fatal("❌ Insert payable_payments failed:", err)
	}

	fmt.Println("✅ Test Data Seeded Successfully!")
}
