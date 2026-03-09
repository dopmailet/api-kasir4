-- Migration: Add discount columns to categories table
-- Tanggal: 2026-02-27
-- Deskripsi: Menambahkan kolom discount_type dan discount_value ke tabel categories
--            agar diskon per kategori bisa dikelola langsung

ALTER TABLE categories
  ADD COLUMN IF NOT EXISTS discount_type VARCHAR(20) DEFAULT NULL,
  ADD COLUMN IF NOT EXISTS discount_value DECIMAL(12,2) DEFAULT 0;

-- Catatan:
-- discount_type: 'percentage' atau 'fixed' (nullable, NULL = tidak ada diskon kategori)
-- discount_value: nilai diskon, misal 10.00 untuk 10% atau 5000.00 untuk Rp 5.000
