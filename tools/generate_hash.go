package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Cara pakai: go run tools/generate_hash.go <password_baru>")
		fmt.Println("💡 Contoh: go run tools/generate_hash.go rahasia123")
		os.Exit(1)
	}

	password := os.Args[1]

	// Generate hash menggunakan bcrypt cost 10 (default standar aplikasi)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal membuat hash:", err)
	}

	fmt.Println("========= HASIL BCRYPT HASH =========")
	fmt.Println("Password Asli :", password)
	fmt.Println("Hasil Hash    :", string(hash))
	fmt.Println("=====================================")
	fmt.Println("✓ Copy 'Hasil Hash' di atas (dimulai dari $2a$...)")
	fmt.Println("✓ Paste ke kolom 'password' di Supabase")
}
