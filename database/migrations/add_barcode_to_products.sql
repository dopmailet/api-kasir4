-- Migration: Add barcode column to products table
-- Tanggal: 2026-02-26
-- Deskripsi: Menambahkan kolom barcode (VARCHAR(100), nullable, UNIQUE)
--            agar produk bisa di-scan via barcode scanner

ALTER TABLE products 
ADD COLUMN IF NOT EXISTS barcode VARCHAR(100) UNIQUE;

-- Catatan:
-- barcode: kode barcode unik per produk (nullable, NULL = belum ada barcode)
-- UNIQUE constraint: tidak boleh ada 2 produk dengan barcode yang sama
