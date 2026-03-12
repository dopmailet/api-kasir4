-- ============================================================
-- MIGRATION: Tambah batas maksimal transaksi harian ke paket
-- ============================================================

-- Tambah kolom max_daily_sales untuk membatasi transaksi (NULL = unlimited)
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS max_daily_sales INT DEFAULT NULL;

-- Update paket default sesuai spesifikasi
-- Gratis (ID 1): max 10 transaksi per hari
UPDATE subscription_packages 
SET max_daily_sales = 10
WHERE id = 1;

-- Basic (ID 2): max 100 transaksi per hari
UPDATE subscription_packages 
SET max_daily_sales = 100
WHERE id = 2;

-- Pro (ID 3): max 500 transaksi per hari
UPDATE subscription_packages 
SET max_daily_sales = 500
WHERE id = 3;

-- Business/unlimited bisa dibiarkan NULL (DEFAULT NULL)
