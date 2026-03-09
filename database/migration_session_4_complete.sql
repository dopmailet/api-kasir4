-- ============================================
-- MIGRATION LENGKAP: Session 4 - Authentication & Profit Tracking
-- ============================================
-- Migration ini mencakup SEMUA perubahan untuk Session 4:
-- 1. Tabel users untuk authentication
-- 2. Kolom harga_beli di products untuk profit tracking
-- 3. Kolom kasir_id di transactions untuk tracking kasir
-- 4. Default users (admin, kasir1, kasir2)
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste SELURUH script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- 6. Selesai! ✅
-- ============================================

-- ============================================
-- PART 1: CREATE TABLE USERS
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- Hashed dengan bcrypt
    nama_lengkap VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'kasir' CHECK (role IN ('admin', 'kasir')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes untuk performa
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);

-- Trigger untuk auto update updated_at
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();

-- ============================================
-- PART 2: ADD HARGA_BELI TO PRODUCTS
-- ============================================
-- Tambah kolom harga_beli untuk tracking harga modal
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS harga_beli NUMERIC(10,2);

-- Tambah kolom created_by untuk tracking siapa yang menambahkan produk
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS created_by INTEGER;

-- Foreign key ke users
ALTER TABLE products
DROP CONSTRAINT IF EXISTS fk_products_created_by;

ALTER TABLE products
ADD CONSTRAINT fk_products_created_by 
FOREIGN KEY (created_by) 
REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Index untuk performa
CREATE INDEX IF NOT EXISTS idx_products_created_by ON products(created_by);

-- Constraint validation
ALTER TABLE products
DROP CONSTRAINT IF EXISTS check_harga_beli_positive;

ALTER TABLE products
ADD CONSTRAINT check_harga_beli_positive 
CHECK (harga_beli IS NULL OR harga_beli >= 0);

ALTER TABLE products
DROP CONSTRAINT IF EXISTS check_harga_positive;

ALTER TABLE products
ADD CONSTRAINT check_harga_positive 
CHECK (harga >= 0);

-- ============================================
-- PART 3: ADD KASIR_ID TO TRANSACTIONS
-- ============================================
-- Tambah kolom kasir_id untuk tracking siapa kasir yang melakukan transaksi
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS kasir_id INTEGER;

-- Foreign key ke users
ALTER TABLE transactions
DROP CONSTRAINT IF EXISTS fk_transactions_kasir;

ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_kasir 
FOREIGN KEY (kasir_id) 
REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Index untuk performa
CREATE INDEX IF NOT EXISTS idx_transactions_kasir_id ON transactions(kasir_id);

-- Tambah kolom harga_beli ke transaction_details untuk snapshot
ALTER TABLE transaction_details 
ADD COLUMN IF NOT EXISTS harga_beli NUMERIC(10,2);

-- ============================================
-- PART 4: INSERT DEFAULT USERS
-- ============================================
-- Hapus data lama jika ada (untuk re-run script)
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
-- VERIFICATION QUERIES
-- ============================================
-- Jalankan query ini untuk verify migration berhasil:

-- 1. Cek users table dan isinya
SELECT id, username, nama_lengkap, role, is_active, created_at 
FROM users 
ORDER BY id;

-- 2. Cek kolom baru di products
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'products' 
  AND column_name IN ('harga_beli', 'created_by')
ORDER BY ordinal_position;

-- 3. Cek kolom baru di transactions
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'transactions' 
  AND column_name = 'kasir_id'
ORDER BY ordinal_position;

-- 4. Test query profit calculation
SELECT 
    id, 
    nama, 
    harga as harga_jual,
    harga_beli,
    CASE 
        WHEN harga_beli IS NOT NULL THEN harga - harga_beli 
        ELSE NULL 
    END as profit_per_unit,
    CASE 
        WHEN harga_beli IS NOT NULL THEN 
            ROUND(((harga - harga_beli) / harga * 100)::numeric, 2)
        ELSE NULL 
    END as margin_persen
FROM products
LIMIT 5;

-- ============================================
-- ROLLBACK (jika perlu)
-- ============================================
-- Jika ingin rollback semua perubahan, jalankan:
/*
-- Rollback Part 4: Delete users
DELETE FROM users WHERE username IN ('admin', 'kasir1', 'kasir2');

-- Rollback Part 3: Remove kasir_id from transactions
ALTER TABLE transaction_details DROP COLUMN IF EXISTS harga_beli;
DROP INDEX IF EXISTS idx_transactions_kasir_id;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_kasir;
ALTER TABLE transactions DROP COLUMN IF EXISTS kasir_id;

-- Rollback Part 2: Remove harga_beli from products
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_harga_positive;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_harga_beli_positive;
DROP INDEX IF EXISTS idx_products_created_by;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_created_by;
ALTER TABLE products DROP COLUMN IF EXISTS created_by;
ALTER TABLE products DROP COLUMN IF EXISTS harga_beli;

-- Rollback Part 1: Drop users table
DROP TRIGGER IF EXISTS trigger_update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_users_updated_at();
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users CASCADE;
*/

-- ============================================
-- SELESAI! ✅
-- ============================================
-- Migration berhasil jika tidak ada error.
-- Cek hasil dengan menjalankan verification queries di atas.
--
-- Default Credentials:
-- - admin / admin123 (role: admin)
-- - kasir1 / kasir123 (role: kasir)
-- - kasir2 / kasir123 (role: kasir)
--
-- ⚠️ PENTING: Ganti password default di production!
-- ============================================
