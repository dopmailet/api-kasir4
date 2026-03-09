// reset_admin_password.go
// Tool untuk generate bcrypt hash password dan mencetak SQL query reset password.
//
// Cara pakai:
//   cd tools
//   go run reset_admin_password.go
//
// Salin output SQL ke Supabase SQL Editor lalu jalankan.

package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserReset struct {
	Username    string
	Password    string
	NamaLengkap string
	Role        string
}

func generateHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("gagal generate hash: %v", err))
	}
	return string(hash)
}

func main() {
	// -----------------------------------------------
	// Daftar user yang ingin di-reset passwordnya
	// Ubah sesuai kebutuhan
	// -----------------------------------------------
	users := []UserReset{
		{Username: "admin", Password: "admin123", NamaLengkap: "Administrator", Role: "admin"},
		{Username: "kasir1", Password: "kasir123", NamaLengkap: "Kasir Utama", Role: "kasir"},
	}

	fmt.Println("========================================")
	fmt.Println("🔑 KASIR API - PASSWORD RESET TOOL")
	fmt.Println("========================================")
	fmt.Println()

	for _, u := range users {
		hash := generateHash(u.Password)
		fmt.Printf("User     : %s\n", u.Username)
		fmt.Printf("Password : %s\n", u.Password)
		fmt.Printf("Hash     : %s\n", hash)
		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Println("📋 SQL UNTUK SUPABASE SQL EDITOR")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("-- LANGKAH 1: Cek user yang ada")
	fmt.Print("SELECT id, username, role, is_active FROM users")
	usernames := ""
	for i, u := range users {
		if i == 0 {
			usernames += fmt.Sprintf("'%s'", u.Username)
		} else {
			usernames += fmt.Sprintf(", '%s'", u.Username)
		}
	}
	fmt.Printf(" WHERE username IN (%s);\n\n", usernames)

	fmt.Println("-- LANGKAH 2: Jika user SUDAH ADA, gunakan UPDATE")
	for _, u := range users {
		hash := generateHash(u.Password)
		fmt.Printf("UPDATE users\n")
		fmt.Printf("SET password = '%s', is_active = true\n", hash)
		fmt.Printf("WHERE username = '%s';\n\n", u.Username)
	}

	fmt.Println("-- LANGKAH 3: Jika user BELUM ADA, gunakan INSERT")
	for _, u := range users {
		hash := generateHash(u.Password)
		fmt.Printf("INSERT INTO users (username, password, nama_lengkap, role, is_active)\n")
		fmt.Printf("VALUES ('%s', '%s', '%s', '%s', true);\n\n",
			u.Username, hash, u.NamaLengkap, u.Role)
	}

	fmt.Println("========================================")
	fmt.Println("✅ Setelah SQL dijalankan, login dengan:")
	for _, u := range users {
		fmt.Printf("   username: %-10s | password: %s\n", u.Username, u.Password)
	}
	fmt.Println("========================================")
}
