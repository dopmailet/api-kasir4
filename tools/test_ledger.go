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
	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated hash:", string(hash))

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres.bwpaqpblhnsmrtuivvyf:1LDjhKBAmgvhK61s@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete and re-insert admin
	db.Exec("DELETE FROM users WHERE username = 'admin'")
	_, err = db.Exec(`
		INSERT INTO users (username, password, nama_lengkap, role, is_active)
		VALUES ('admin', $1, 'Administrator', 'admin', true)
	`, string(hash))
	if err != nil {
		log.Fatal("Insert error:", err)
	}
	fmt.Println("✅ Admin user created successfully!")

	// Verify
	var storedHash string
	err = db.QueryRow("SELECT password FROM users WHERE username = 'admin'").Scan(&storedHash)
	if err != nil {
		log.Fatal("Query error:", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte("admin123"))
	if err != nil {
		fmt.Println("❌ Password verification FAILED:", err)
	} else {
		fmt.Println("✅ Password verification OK - login admin/admin123 should work!")
	}
}
