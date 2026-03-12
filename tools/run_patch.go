package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("DB_CONN is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Read SQL file
	content, err := ioutil.ReadFile("database/migrations/patch_add_max_daily_sales.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("Error executing SQL: %v", err)
	}

	fmt.Println("✅ SQL migration for max_daily_sales applied successfully!")
}
