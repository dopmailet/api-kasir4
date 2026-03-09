-- ============================================
-- MIGRATION: Add harga_beli to products table
-- ============================================
-- Migration ini menambahkan kolom harga_beli untuk tracking
-- harga modal/pembelian produk, sehingga bisa menghitung profit
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- ============================================

-- ============================================
-- 1. ADD COLUMN: harga_beli
-- ============================================
-- Tambah kolom harga_beli ke table products
-- NULL = untuk produk lama yang belum ada data harga beli
-- NUMERIC(10,2) = support angka sampai 99,999,999.99
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS harga_beli NUMERIC(10,2);

-- ============================================
-- 2. ADD COLUMN: created_by (tracking)
-- ============================================
-- Tambah kolom untuk tracking siapa yang menambahkan produk
-- NULL = untuk produk lama yang ditambahkan sebelum sistem user
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS created_by INTEGER;

-- ============================================
-- 3. ADD FOREIGN KEY CONSTRAINT
-- ============================================
-- Tambah foreign key constraint ke table users
-- ON DELETE SET NULL = kalau user dihapus, created_by jadi NULL
-- ON UPDATE CASCADE = kalau user.id diupdate, created_by ikut update
ALTER TABLE products
ADD CONSTRAINT fk_products_created_by 
FOREIGN KEY (created_by) 
REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- ============================================
-- 4. CREATE INDEX
-- ============================================
-- Index untuk mempercepat query filter by created_by
CREATE INDEX IF NOT EXISTS idx_products_created_by 
ON products(created_by);

-- ============================================
-- 5. ADD CHECK CONSTRAINT (Business Logic)
-- ============================================
-- Pastikan harga_beli tidak negatif (jika ada nilai)
ALTER TABLE products
ADD CONSTRAINT check_harga_beli_positive 
CHECK (harga_beli IS NULL OR harga_beli >= 0);

-- Pastikan harga jual tidak negatif
ALTER TABLE products
ADD CONSTRAINT check_harga_positive 
CHECK (harga >= 0);

-- ============================================
-- ROLLBACK (jika perlu)
-- ============================================
-- Jika ingin rollback migration ini, jalankan:
-- ALTER TABLE products DROP CONSTRAINT IF EXISTS check_harga_positive;
-- ALTER TABLE products DROP CONSTRAINT IF EXISTS check_harga_beli_positive;
-- DROP INDEX IF EXISTS idx_products_created_by;
-- ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_created_by;
-- ALTER TABLE products DROP COLUMN IF EXISTS created_by;
-- ALTER TABLE products DROP COLUMN IF EXISTS harga_beli;

-- ============================================
-- VERIFY
-- ============================================
-- Cek apakah kolom sudah ditambahkan:
-- SELECT column_name, data_type, is_nullable, column_default
-- FROM information_schema.columns 
-- WHERE table_name = 'products'
-- ORDER BY ordinal_position;

-- Test query dengan profit calculation:
-- SELECT 
--     id, 
--     nama, 
--     harga as harga_jual,
--     harga_beli,
--     CASE 
--         WHEN harga_beli IS NOT NULL THEN harga - harga_beli 
--         ELSE NULL 
--     END as profit_per_unit,
--     CASE 
--         WHEN harga_beli IS NOT NULL THEN 
--             ROUND(((harga - harga_beli) / harga * 100)::numeric, 2)
--         ELSE NULL 
--     END as margin_persen
-- FROM products;
