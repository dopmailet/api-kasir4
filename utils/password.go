package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword meng-hash password menggunakan bcrypt
// Cost 10 adalah default yang recommended (balance antara security & performance)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash membandingkan password plain text dengan hash
// Returns true jika password cocok
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
