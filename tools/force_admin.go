package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("DB_CONN is missing")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	hashBytes, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	hash := string(hashBytes)

	query := `
		INSERT INTO users (username, password, nama_lengkap, role, is_active) 
		VALUES ($1, $2, 'Administrator', 'admin', true)
		ON CONFLICT (username) DO UPDATE 
		SET password = EXCLUDED.password, is_active = true
	`

	_, err = pool.Exec(ctx, query, "admin", hash)
	if err != nil {
		log.Fatal("Failed to update admin:", err)
	}

	fmt.Println("Admin password forcefully updated and hashed to admin123")
}
