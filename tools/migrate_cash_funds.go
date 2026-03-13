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
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found or error loading, relying on explicit connection string if needed")
	}

	dbUrl := os.Getenv("DB_CONN")
	if dbUrl == "" {
		dbUrl = "postgresql://postgres.bwpaqpblhnsmrtuivvyf:1LDjhKBAmgvhK61s@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	query := `
		ALTER TABLE cash_funds ADD COLUMN IF NOT EXISTS store_id INTEGER NOT NULL DEFAULT 1 REFERENCES stores(id);
		CREATE INDEX IF NOT EXISTS idx_cash_funds_store ON cash_funds(store_id);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("✅ Migration successful: Added store_id and index to cash_funds")
}
