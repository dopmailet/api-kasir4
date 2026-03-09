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
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("DB_CONN environment variable not set")
	}

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Unable to parse DB_CONN:", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	sqlBytes, err := ioutil.ReadFile("database/migrations/add_payroll_system.sql")
	if err != nil {
		log.Fatal("Unable to read sql file:", err)
	}

	_, err = pool.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	fmt.Println("✅ Migration applied successfully!")
}
