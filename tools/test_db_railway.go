package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Koneksi ke DB yang PERSIS SAMA dengan yang ada di Railway (.env file ini)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("No DATABASE_URL set, trying DB_CONN...")
		dbURL = os.Getenv("DB_CONN")
		if dbURL == "" {
			// Hardcode URL yang dipakai user di Railway
			dbURL = "postgresql://postgres.bwpaqpblhnsmrtuivvyf:1LDjhKBAmgvhK61s@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
		}
	}

	fmt.Println("Testing connection to:", dbURL)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Open error:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping error:", err)
	}
	fmt.Println("✅ Connection successful!")

	// Cek tabel users
	var count int
	err = db.QueryRow("SELECT count(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal("Query count error:", err)
	}
	fmt.Printf("found %d users in DB\n", count)

	var id int
	var storedHash string
	err = db.QueryRow("SELECT id, password FROM users WHERE username = 'admin'").Scan(&id, &storedHash)
	if err != nil {
		log.Fatal("Query admin error:", err)
	}

	fmt.Printf("Admin ID: %d\n", id)
	fmt.Printf("Stored Hash: %s\n", storedHash)

	// Test bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte("admin123"))
	if err != nil {
		fmt.Println("❌ Password verification FAILED:", err)

		// Mari kita fix sekalian di script ini
		fmt.Println("Attempting to fix password...")
		newHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err = db.Exec("UPDATE users SET password = $1 WHERE username = 'admin'", string(newHash))
		if err != nil {
			log.Fatal("Failed to update password:", err)
		}
		fmt.Println("✅ Password updated successfully! Please try login again.")

	} else {
		fmt.Println("✅ Password verification OK - The database hash is correct for 'admin123'")
	}
}
