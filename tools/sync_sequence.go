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
	// Load file dari root directory (.) atau sekitarnya
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		log.Println("Peringatan: Error loading .env file")
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Reset sequence ID toko karena kita sudah insert manual Toko ID=1
	_, err = db.Exec(`SELECT setval(pg_get_serial_sequence('stores', 'id'), coalesce(max(id),0) + 1, false) FROM stores;`)
	if err != nil {
		log.Fatalf("Gagal mereset sequence tabel stores: %v", err)
	}

	// Lakukan hal yang sama untuk tabel users karena kita insert superadmin juga
	_, err = db.Exec(`SELECT setval(pg_get_serial_sequence('users', 'id'), coalesce(max(id),0) + 1, false) FROM users;`)
	if err != nil {
		log.Fatalf("Gagal mereset sequence tabel users: %v", err)
	}

	fmt.Println("Berhasil sinkronisasi ID sequence untuk tabel stores dan users!")
}
