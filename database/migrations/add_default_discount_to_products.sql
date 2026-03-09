-- Migration: Add default discount columns to products table
-- Tanggal: 2026-02-24
-- Deskripsi: Menambahkan kolom default_discount_type dan default_discount_value
--            agar diskon default per produk tersimpan di database (konsisten di semua device)

ALTER TABLE products 
ADD COLUMN IF NOT EXISTS default_discount_type VARCHAR(20),
ADD COLUMN IF NOT EXISTS default_discount_value DECIMAL(15,2);

-- Catatan:
-- default_discount_type: 'percentage' atau 'fixed' (nullable, NULL = tidak ada diskon default)
-- default_discount_value: nilai diskon, misal 5.00 untuk 5% atau 5000.00 untuk Rp 5.000
