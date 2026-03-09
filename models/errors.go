package models

import "errors"

// Product errors
var (
	ErrInvalidPrice         = errors.New("harga jual tidak boleh negatif")
	ErrInvalidPurchasePrice = errors.New("harga beli tidak boleh negatif")
	ErrNegativeMargin       = errors.New("peringatan: harga beli tidak boleh lebih besar atau sama dengan harga jual (margin negatif)")
	ErrProductNotFound      = errors.New("produk tidak ditemukan")
	ErrInsufficientStock    = errors.New("stok tidak mencukupi")
)

// User errors
var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrUserInactive       = errors.New("user tidak aktif")
	ErrInvalidRole        = errors.New("role tidak valid")
	ErrUnauthorized       = errors.New("tidak memiliki akses")
	ErrForbidden          = errors.New("akses ditolak")
	ErrCannotDeleteSelf   = errors.New("tidak dapat menghapus diri sendiri")
)

// Auth errors
var (
	ErrInvalidToken  = errors.New("token tidak valid")
	ErrExpiredToken  = errors.New("token sudah kadaluarsa")
	ErrMissingToken  = errors.New("token tidak ditemukan")
	ErrInvalidAPIKey = errors.New("API key tidak valid")
)

// Transaction errors
var (
	ErrEmptyCart         = errors.New("keranjang belanja kosong")
	ErrInvalidQuantity   = errors.New("jumlah tidak valid")
	ErrTransactionFailed = errors.New("transaksi gagal")
)
