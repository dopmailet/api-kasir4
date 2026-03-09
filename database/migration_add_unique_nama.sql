-- ============================================
-- MIGRATION: ADD UNIQUE CONSTRAINT TO PRODUCTS.NAMA
-- ============================================
-- Script ini menambahkan UNIQUE constraint pada kolom nama
-- di table products agar fitur UPSERT bisa bekerja
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- ============================================

-- Hapus index biasa yang sudah ada (jika ada)
DROP INDEX IF EXISTS idx_products_nama;

-- Tambahkan UNIQUE constraint pada kolom nama
-- Ini akan memastikan tidak ada produk dengan nama yang sama
ALTER TABLE products 
ADD CONSTRAINT products_nama_unique UNIQUE (nama);

-- Buat index UNIQUE untuk mempercepat pencarian
-- (Index ini otomatis dibuat oleh UNIQUE constraint, tapi kita explicit)
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_nama_unique ON products(nama);

-- ============================================
-- SELESAI! 🎉
-- ============================================
-- Sekarang kolom nama di table products sudah UNIQUE
-- Fitur UPSERT (INSERT ON CONFLICT) akan bekerja dengan baik
