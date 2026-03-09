package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"kasir-api/config"
	"log"

	_ "github.com/lib/pq"
)

func main() {
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

	fmt.Println("🛠️ Updating Database for Category Discounts...")

	// Read SQL File
	sqlFile := "database/migration_add_category_discount.sql"
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
