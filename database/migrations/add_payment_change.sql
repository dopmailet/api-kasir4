-- ==========================================
-- MIGRATION: Tambah payment_amount dan change_amount ke transactions
-- Tanggal: 2026-02-17
-- Deskripsi: Menyimpan uang bayar dan kembalian customer
-- ==========================================

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_amount DECIMAL(15, 2) DEFAULT 0;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS change_amount DECIMAL(15, 2) DEFAULT 0;
