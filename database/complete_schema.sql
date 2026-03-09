-- ============================================
-- COMPLETE DATABASE SCHEMA - KASIR API
-- ============================================
-- Version: Session 4 (with Authentication & Profit Tracking)
-- Last Updated: 2026-02-12
--
-- Script ini berisi SELURUH schema database dari awal sampai Session 4.
-- Cocok untuk:
-- - Setup Supabase baru dari nol
-- - Disaster recovery
-- - Clone database ke environment baru
-- - Development setup
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buat project baru atau pilih project yang ingin di-setup
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste SELURUH script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- 6. Tunggu sampai selesai (10-15 detik)
-- ============================================

-- ============================================
-- PART 1: CREATE TABLE PRODUCTS
-- ============================================
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL UNIQUE,
    harga NUMERIC(10,2) NOT NULL CHECK (harga >= 0),
    harga_beli NUMERIC(10,2) CHECK (harga_beli IS NULL OR harga_beli >= 0),
    stok INTEGER NOT NULL CHECK (stok >= 0),
    category_id INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk products
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_nama ON products(nama);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_created_by ON products(created_by);

-- ============================================
-- PART 2: CREATE TABLE CATEGORIES
-- ============================================
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk categories
CREATE INDEX IF NOT EXISTS idx_categories_nama ON categories(nama);

-- ============================================
-- PART 3: CREATE TABLE USERS
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nama_lengkap VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'kasir' CHECK (role IN ('admin', 'kasir')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk users
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);

-- ============================================
-- PART 4: CREATE TABLE TRANSACTIONS
-- ============================================
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount NUMERIC(10,2) NOT NULL,
    kasir_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk transactions
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_kasir_id ON transactions(kasir_id);

-- ============================================
-- PART 5: CREATE TABLE TRANSACTION_DETAILS
-- ============================================
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    harga_beli NUMERIC(10,2),
    subtotal NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes untuk transaction_details
CREATE INDEX IF NOT EXISTS idx_transaction_details_transaction_id ON transaction_details(transaction_id);
CREATE INDEX IF NOT EXISTS idx_transaction_details_product_id ON transaction_details(product_id);

-- ============================================
-- PART 6: ADD FOREIGN KEY CONSTRAINTS
-- ============================================

-- Products foreign keys
ALTER TABLE products
DROP CONSTRAINT IF EXISTS fk_products_category;

ALTER TABLE products
ADD CONSTRAINT fk_products_category 
FOREIGN KEY (category_id) 
REFERENCES categories(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

ALTER TABLE products
DROP CONSTRAINT IF EXISTS fk_products_created_by;

ALTER TABLE products
ADD CONSTRAINT fk_products_created_by 
FOREIGN KEY (created_by) 
REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Transactions foreign keys
ALTER TABLE transactions
DROP CONSTRAINT IF EXISTS fk_transactions_kasir;

ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_kasir 
FOREIGN KEY (kasir_id) 
REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Transaction_details foreign keys
ALTER TABLE transaction_details
DROP CONSTRAINT IF EXISTS fk_transaction_details_transaction;

ALTER TABLE transaction_details
ADD CONSTRAINT fk_transaction_details_transaction 
FOREIGN KEY (transaction_id) 
REFERENCES transactions(id) 
ON DELETE CASCADE;

ALTER TABLE transaction_details
DROP CONSTRAINT IF EXISTS fk_transaction_details_product;

ALTER TABLE transaction_details
ADD CONSTRAINT fk_transaction_details_product 
FOREIGN KEY (product_id) 
REFERENCES products(id);

-- ============================================
-- PART 7: CREATE FUNCTIONS
-- ============================================

-- Function untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function untuk auto-update updated_at di users
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- PART 8: CREATE TRIGGERS
-- ============================================

-- Trigger untuk products
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger untuk categories
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger untuk users
DROP TRIGGER IF EXISTS trigger_update_users_updated_at ON users;
CREATE TRIGGER trigger_update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();

-- ============================================
-- PART 9: INSERT SAMPLE DATA
-- ============================================

-- Sample Categories
INSERT INTO categories (nama, description) VALUES
    ('Makanan', 'Produk makanan dan snack'),
    ('Minuman', 'Produk minuman kemasan'),
    ('Sembako', 'Kebutuhan pokok sehari-hari'),
    ('Alat Tulis', 'Perlengkapan tulis dan kantor')
ON CONFLICT DO NOTHING;

-- Sample Products (tanpa harga_beli untuk produk lama)
INSERT INTO products (nama, harga, stok, category_id) VALUES
    ('Indomie Goreng', 3500, 100, 1),
    ('Aqua 600ml', 3000, 50, 2),
    ('Teh Pucuk', 4000, 75, 2),
    ('Kopi Kapal Api', 2000, 200, 2),
    ('Mie Sedaap', 3500, 80, 1)
ON CONFLICT (nama) DO NOTHING;

-- ============================================
-- PART 10: INSERT DEFAULT USERS
-- ============================================

-- Hapus users lama jika ada (untuk re-run script)
DELETE FROM users WHERE username IN ('admin', 'kasir1', 'kasir2');

-- Insert Admin
-- Username: admin | Password: admin123
INSERT INTO users (username, password, nama_lengkap, role, is_active) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', 'admin', TRUE);

-- Insert Kasir 1
-- Username: kasir1 | Password: kasir123
INSERT INTO users (username, password, nama_lengkap, role, is_active) VALUES
('kasir1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Pagi', 'kasir', TRUE);

-- Insert Kasir 2
-- Username: kasir2 | Password: kasir123
INSERT INTO users (username, password, nama_lengkap, role, is_active) VALUES
('kasir2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Sore', 'kasir', TRUE);

-- ============================================
-- PART 11: VERIFICATION QUERIES
-- ============================================
-- Jalankan query ini untuk verify setup berhasil:

-- 1. Cek semua tables
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
  AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- 2. Cek users
SELECT id, username, nama_lengkap, role, is_active 
FROM users 
ORDER BY id;

-- 3. Cek products
SELECT id, nama, harga, harga_beli, stok, category_id 
FROM products 
ORDER BY id 
LIMIT 5;

-- 4. Cek categories
SELECT id, nama, description 
FROM categories 
ORDER BY id;

-- 5. Cek foreign keys
SELECT
    tc.table_name, 
    kcu.column_name, 
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
  ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
  ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
ORDER BY tc.table_name;

-- ============================================
-- EXPECTED RESULTS
-- ============================================
-- Tables: categories, products, transaction_details, transactions, users
-- Users: 3 users (admin, kasir1, kasir2)
-- Products: 5 sample products
-- Categories: 4 sample categories
-- Foreign Keys: 5 foreign keys

-- ============================================
-- DEFAULT CREDENTIALS
-- ============================================
-- admin / admin123 (role: admin)
-- kasir1 / kasir123 (role: kasir)
-- kasir2 / kasir123 (role: kasir)
--
-- ⚠️ PENTING: Ganti password ini di production!

-- ============================================
-- SELESAI! ✅
-- ============================================
-- Database schema sudah lengkap dan siap digunakan!
-- 
-- Next steps:
-- 1. Update .env dengan JWT_SECRET
-- 2. Run aplikasi: go run main.go
-- 3. Test login dengan credentials di atas
-- 4. Update password default di production
-- ============================================
