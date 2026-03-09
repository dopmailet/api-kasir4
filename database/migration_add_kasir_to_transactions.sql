-- ============================================
-- MIGRATION: Add kasir_id to transactions table
-- ============================================
-- Migration ini menambahkan kolom kasir_id untuk tracking
-- siapa kasir yang melakukan transaksi penjualan
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- ============================================

-- ============================================
-- 1. ADD COLUMN: kasir_id
-- ============================================
-- Tambah kolom kasir_id ke table transactions
-- NULL = untuk transaksi lama yang dilakukan sebelum sistem user
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS kasir_id INTEGER;

-- ============================================
-- 2. ADD FOREIGN KEY CONSTRAINT
-- ============================================
-- Tambah foreign key constraint ke table users
-- ON DELETE SET NULL = kalau user kasir dihapus, kasir_id jadi NULL (data tetap ada)
-- ON UPDATE CASCADE = kalau user.id diupdate, kasir_id ikut update
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_kasir 
FOREIGN KEY (kasir_id) 
REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- ============================================
-- 3. CREATE INDEX
-- ============================================
-- Index untuk mempercepat query filter by kasir_id
CREATE INDEX IF NOT EXISTS idx_transactions_kasir_id 
ON transactions(kasir_id);

-- ============================================
-- 4. ADD COLUMN: harga_beli to transaction_details
-- ============================================
-- Tambah kolom harga_beli untuk snapshot harga beli saat transaksi
-- Ini penting untuk laporan profit yang akurat (harga beli bisa berubah)
ALTER TABLE transaction_details 
ADD COLUMN IF NOT EXISTS harga_beli NUMERIC(10,2);

-- ============================================
-- ROLLBACK (jika perlu)
-- ============================================
-- Jika ingin rollback migration ini, jalankan:
-- ALTER TABLE transaction_details DROP COLUMN IF EXISTS harga_beli;
-- DROP INDEX IF EXISTS idx_transactions_kasir_id;
-- ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_kasir;
-- ALTER TABLE transactions DROP COLUMN IF EXISTS kasir_id;

-- ============================================
-- VERIFY
-- ============================================
-- Cek apakah kolom sudah ditambahkan:
-- SELECT column_name, data_type, is_nullable, column_default
-- FROM information_schema.columns 
-- WHERE table_name = 'transactions'
-- ORDER BY ordinal_position;

-- Test query dengan join ke users:
-- SELECT 
--     t.id,
--     t.total_amount,
--     t.created_at,
--     u.nama_lengkap as kasir_nama,
--     u.role as kasir_role
-- FROM transactions t
-- LEFT JOIN users u ON t.kasir_id = u.id
-- ORDER BY t.created_at DESC
-- LIMIT 10;
