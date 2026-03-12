-- ============================================================
-- MIGRATION: Patch subscription_packages untuk mendukung fitur
-- PUT/POST hardening (tambah kolom yang mungkin belum ada)
-- 
-- Jalankan di Supabase SQL Editor
-- AMAN: Semua menggunakan IF NOT EXISTS / DEFAULT values
-- ============================================================

-- Pastikan kolom description ada (nullable TEXT)
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS description TEXT DEFAULT NULL;

-- Pastikan kolom features ada sebagai JSONB
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS features JSONB DEFAULT '[]'::jsonb;

-- Pastikan kolom period ada (mis: "/bulan", "/tahun")
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS period VARCHAR(50) DEFAULT '/bulan';

-- Pastikan kolom discount_percent ada
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS discount_percent NUMERIC DEFAULT 0;

-- Pastikan kolom discount_label ada (nullable)
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS discount_label VARCHAR(100) DEFAULT NULL;

-- Pastikan kolom is_popular ada
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS is_popular BOOLEAN DEFAULT false;

-- Pastikan kolom sort_order ada (untuk urutan tampilan)
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS sort_order INT DEFAULT 0;

-- Pastikan kolom updated_at ada
ALTER TABLE subscription_packages 
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();

-- Isi sort_order untuk data yang sudah ada (berdasarkan id)
UPDATE subscription_packages 
SET sort_order = id 
WHERE sort_order = 0 OR sort_order IS NULL;

-- Pastikan features tidak NULL (ubah NULL → empty array)
UPDATE subscription_packages 
SET features = '[]'::jsonb 
WHERE features IS NULL;

-- Verifikasi hasil
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'subscription_packages'
ORDER BY ordinal_position;
