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
		CREATE TABLE IF NOT EXISTS loyalty_transactions (
			id SERIAL PRIMARY KEY,
			customer_id INT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
			transaction_id INT REFERENCES transactions(id) ON DELETE SET NULL,
			points INT NOT NULL,              
			type VARCHAR(10) NOT NULL,        
			description TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);

		ALTER TABLE loyalty_transactions ADD COLUMN IF NOT EXISTS store_id INTEGER NOT NULL DEFAULT 1 REFERENCES stores(id);
		CREATE INDEX IF NOT EXISTS idx_loyalty_transactions_store ON loyalty_transactions(store_id);

		ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_points INT DEFAULT 0;
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("✅ Migration successful: loyalty_transactions table and loyalty_points column added.")
}
