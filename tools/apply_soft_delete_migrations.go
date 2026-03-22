package main

import (
	"database/sql"
	"fmt"
	"kasir-api/config"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found or error loading it")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed reading config:", err)
	}

	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal("Failed connecting to DB:", err)
	}
	defer db.Close()

	commands := []string{
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;",
		"UPDATE users SET is_active = TRUE WHERE is_active IS NULL;",
		"ALTER TABLE transactions ADD COLUMN IF NOT EXISTS cashier_name VARCHAR(100);",
	}

	for _, cmd := range commands {
		fmt.Println("Executing:", cmd)
		_, err := db.Exec(cmd)
		if err != nil {
			log.Printf("Error executing %s: %v\n", cmd, err)
		} else {
			fmt.Println("Success.")
		}
	}
	fmt.Println("Migration completed.")
}
