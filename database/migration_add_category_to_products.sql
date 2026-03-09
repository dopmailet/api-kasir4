-- ============================================
-- MIGRATION: Add category_id to products table
-- ============================================
-- Migration ini menambahkan relasi antara products dan categories
-- Setiap product bisa memiliki 1 category (optional)
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- ============================================

-- ============================================
-- 1. ADD COLUMN: category_id
-- ============================================
-- Tambah kolom category_id ke table products
-- NULL = category bersifat optional (boleh tidak punya category)
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS category_id INTEGER;

-- ============================================
-- 2. ADD FOREIGN KEY CONSTRAINT
-- ============================================
-- Tambah foreign key constraint ke table categories
-- ON DELETE SET NULL = kalau category dihapus, category_id di product jadi NULL
-- ON UPDATE CASCADE = kalau category.id diupdate, category_id di product ikut update
ALTER TABLE products
ADD CONSTRAINT fk_products_category 
FOREIGN KEY (category_id) 
REFERENCES categories(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- ============================================
-- 3. CREATE INDEX
-- ============================================
-- Index untuk mempercepat query JOIN dan WHERE category_id
CREATE INDEX IF NOT EXISTS idx_products_category_id 
ON products(category_id);

-- ============================================
-- 4. UPDATE SAMPLE DATA (OPTIONAL)
-- ============================================
-- Update beberapa product dengan category_id
-- Sesuaikan dengan ID category yang ada di database Anda

-- Ambil ID categories yang ada
-- SELECT id, nama FROM categories;

-- Contoh: Update products dengan category
-- UPDATE products SET category_id = 1 WHERE nama = 'Indomie Goreng';  -- Makanan
-- UPDATE products SET category_id = 2 WHERE nama = 'Aqua 600ml';      -- Minuman
-- UPDATE products SET category_id = 2 WHERE nama = 'Teh Pucuk';       -- Minuman
-- UPDATE products SET category_id = 2 WHERE nama = 'Kopi Kapal Api';  -- Minuman
-- UPDATE products SET category_id = 1 WHERE nama = 'Mie Sedaap';      -- Makanan

-- ============================================
-- ROLLBACK (jika perlu)
-- ============================================
-- Jika ingin rollback migration ini, jalankan:
-- DROP INDEX IF EXISTS idx_products_category_id;
-- ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_category;
-- ALTER TABLE products DROP COLUMN IF EXISTS category_id;

-- ============================================
-- VERIFY
-- ============================================
-- Cek apakah kolom sudah ditambahkan:
-- SELECT column_name, data_type, is_nullable 
-- FROM information_schema.columns 
-- WHERE table_name = 'products';

-- Test JOIN query:
-- SELECT p.id, p.nama, p.harga, p.stok, p.category_id, c.nama as category_name
-- FROM products p
-- LEFT JOIN categories c ON p.category_id = c.id;
