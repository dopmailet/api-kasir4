package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		log.Fatal("DB_CONN is not set in environment or .env file")
	}

	// Buat connection pool menggunakan pgxpool
	pool, err := pgxpool.New(context.Background(), dbConn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Baca file SQL
	sqlFile := "database/migrations/add_expenses_table.sql"
	sqlBytes, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v\n", err)
	}
	sqlQuery := string(sqlBytes)

	// Eksekusi SQL
	_, err = pool.Exec(context.Background(), sqlQuery)
	if err != nil {
		log.Fatalf("Failed to execute migration: %v\n", err)
	}

	fmt.Println("✅ Migration applied successfully!")
}
