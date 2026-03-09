package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Ambil DATABASE_URL dari env atau default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:admin@localhost:5432/kasir_db?sslmode=disable"
		fmt.Println("⚠️ DATABASE_URL not set, using default:", dbURL)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run tools/run_migration.go <path/to/file.sql>")
	}

	filePath := os.Args[1]
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Failed to read file:", err)
	}

	fmt.Printf("Running migration from %s...\n", filePath)
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	fmt.Println("✅ Migration executed successfully!")
}
