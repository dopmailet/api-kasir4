package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"kasir-api/config"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load env explicitly
	godotenv.Load(".env")

	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Config Error:", err)
	}

	// Connect to Database
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Connection failed:", err)
	}

	fmt.Println("🛠️ Updating Database for Supplier Payables...")

	// Read SQL File
	sqlFile := "database/migrations/migration_add_supplier_payables.sql"
	content, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("❌ Failed to read SQL file:", err)
	}

	// Execute SQL
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal("❌ Migration Failed:", err)
	} else {
		fmt.Println("✅ Database Schema Updated Successfully!")
	}
}
