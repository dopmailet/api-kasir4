package main

import (
	"database/sql"
	"fmt"
	"kasir-api/config"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Load config loaded from environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Failed to load config:", err)
	}

	dbURL := cfg.GetDatabaseURL()
	fmt.Println("🔌 Connecting to:", dbURL)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Connection failed:", err)
	}

	fmt.Println("🛠️ Fixing Database Schema...")

	// 1. Add product_id column
	_, err = db.Exec("ALTER TABLE discounts ADD COLUMN IF NOT EXISTS product_id INT DEFAULT NULL;")
	if err != nil {
		fmt.Println("⚠️ Error adding column (might exist):", err)
	} else {
		fmt.Println("✅ Column 'product_id' added.")
	}

	// 2. Add Foreign Key
	// Must drop constraint first if exists to avoid error
	db.Exec("ALTER TABLE discounts DROP CONSTRAINT IF EXISTS fk_discounts_product;")
	_, err = db.Exec("ALTER TABLE discounts ADD CONSTRAINT fk_discounts_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;")
	if err != nil {
		fmt.Println("⚠️ Error adding FK:", err)
	} else {
		fmt.Println("✅ FK 'fk_discounts_product' added.")
	}

	// 3. Add Index
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_discounts_product ON discounts(product_id);")
	if err != nil {
		fmt.Println("⚠️ Error adding index:", err)
	} else {
		fmt.Println("✅ Index created.")
	}

	fmt.Println("🎉 Database Fix Complete!")
}
