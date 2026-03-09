-- ==========================================
-- MIGRATION: Tambah kolom diskon per item ke transaction_details
-- Tanggal: 2026-02-22
-- Deskripsi: Menyimpan discount_type, discount_value, dan discount_amount per item
-- ==========================================

ALTER TABLE transaction_details
  ADD COLUMN IF NOT EXISTS discount_type VARCHAR(10) DEFAULT NULL,
  ADD COLUMN IF NOT EXISTS discount_value DECIMAL(12,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0;
