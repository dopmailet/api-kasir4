-- ============================================================
-- Migration: Tambah max_daily_sales ke subscription_packages
-- Jalankan di Supabase SQL Editor
-- ============================================================

-- 1. Tambah kolom max_daily_sales (NULL = unlimited/tidak terbatas)
ALTER TABLE subscription_packages 
ADD COLUMN IF NOT EXISTS max_daily_sales INT DEFAULT NULL;

-- 2. Set batas untuk paket Gratis (10 transaksi/hari)
UPDATE subscription_packages 
SET max_daily_sales = 10 
WHERE LOWER(name) LIKE '%gratis%' OR LOWER(name) LIKE '%free%';

-- 3. Paket berbayar = unlimited (NULL)
UPDATE subscription_packages 
SET max_daily_sales = NULL 
WHERE price > 0;

-- Verifikasi hasilnya
SELECT id, name, price, max_kasir, max_products, max_daily_sales 
FROM subscription_packages 
ORDER BY id;
