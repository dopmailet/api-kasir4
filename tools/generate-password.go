package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Password yang ingin kita hash untuk KASIR
	password := "kasir123"
	username := "kasir1"
	role := "kasir"
	fullname := "Kasir Utama"

	// Generate hash dari password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating hash:", err)
		return
	}

	hashString := string(hash)

	fmt.Println("========================================")
	fmt.Println("🔑 PASSWORD HASH GENERATOR (KASIR)")
	fmt.Println("========================================")
	fmt.Println("Username          :", username)
	fmt.Println("Password Original :", password)
	fmt.Println("Bcrypt Hash       :", hashString)
	fmt.Println("========================================")
	fmt.Println("📋 SQL QUERY UNTUK SUPABASE:")
	fmt.Println()
	fmt.Printf("-- 1. Reset password user '%s' menjadi '%s'\n", username, password)
	fmt.Printf("UPDATE users SET password = '%s' WHERE username = '%s';\n", hashString, username)
	fmt.Println()
	fmt.Println("-- 2. ATAU buat user baru jika belum ada:")
	fmt.Printf("INSERT INTO users (username, password, nama_lengkap, role) VALUES ('%s', '%s', '%s', '%s') ON CONFLICT (username) DO NOTHING;\n", username, hashString, fullname, role)
	fmt.Println("========================================")
}
