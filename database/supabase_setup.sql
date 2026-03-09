-- ============================================
-- SUPABASE SETUP SCRIPT - KASIR API
-- ============================================
-- Script ini akan membuat semua table yang dibutuhkan
-- untuk aplikasi Kasir API di Supabase
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- ============================================

-- ============================================
-- 1. CREATE TABLE: PRODUCTS
-- ============================================
-- Table untuk menyimpan data produk
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL UNIQUE,  -- UNIQUE constraint untuk mencegah duplikasi nama
    harga INTEGER NOT NULL CHECK (harga >= 0),
    stok INTEGER NOT NULL CHECK (stok >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index UNIQUE untuk mempercepat pencarian berdasarkan nama
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_nama ON products(nama);

-- ============================================
-- 2. CREATE TABLE: CATEGORIES
-- ============================================
-- Table untuk menyimpan data kategori produk
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat pencarian berdasarkan nama
CREATE INDEX IF NOT EXISTS idx_categories_nama ON categories(nama);

-- ============================================
-- 3. INSERT SAMPLE DATA (OPTIONAL)
-- ============================================
-- Data contoh untuk testing
-- Hapus bagian ini jika tidak ingin data sample

-- Sample Products
INSERT INTO products (nama, harga, stok) VALUES
    ('Indomie Goreng', 3500, 100),
    ('Aqua 600ml', 3000, 50),
    ('Teh Pucuk', 4000, 75),
    ('Kopi Kapal Api', 2000, 200),
    ('Mie Sedaap', 3500, 80)
ON CONFLICT DO NOTHING;

-- Sample Categories
INSERT INTO categories (nama, description) VALUES
    ('Makanan', 'Produk makanan dan snack'),
    ('Minuman', 'Produk minuman kemasan'),
    ('Sembako', 'Kebutuhan pokok sehari-hari'),
    ('Alat Tulis', 'Perlengkapan tulis dan kantor')
ON CONFLICT DO NOTHING;

-- ============================================
-- 4. CREATE FUNCTION: AUTO UPDATE TIMESTAMP
-- ============================================
-- Function untuk auto-update kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 5. CREATE TRIGGERS
-- ============================================
-- Trigger untuk auto-update updated_at di table products
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger untuk auto-update updated_at di table categories
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 6. ENABLE ROW LEVEL SECURITY (OPTIONAL)
-- ============================================
-- Uncomment jika ingin mengaktifkan RLS
-- ALTER TABLE products ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE categories ENABLE ROW LEVEL SECURITY;

-- Policy untuk public access (untuk API)
-- CREATE POLICY "Allow public access" ON products FOR ALL USING (true);
-- CREATE POLICY "Allow public access" ON categories FOR ALL USING (true);

-- ============================================
-- SELESAI! 🎉
-- ============================================
-- Cek apakah table sudah dibuat dengan query:
-- SELECT * FROM products;
-- SELECT * FROM categories;
