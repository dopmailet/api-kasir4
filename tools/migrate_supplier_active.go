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
		ALTER TABLE suppliers ALTER COLUMN is_active SET DEFAULT TRUE;
		UPDATE suppliers SET is_active = TRUE WHERE is_active IS NULL OR is_active = FALSE;
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("✅ Migration successful: suppliers is_active default set to TRUE and existing data updated.")
}
