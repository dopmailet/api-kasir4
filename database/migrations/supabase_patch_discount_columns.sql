-- ============================================================
-- SUPABASE PATCH: Tambah kolom diskon & pembayaran
-- Tanggal: 2026-02-23
-- Jalankan script ini di: Supabase → SQL Editor → Run
-- ============================================================
-- Script ini aman dijalankan berkali-kali (IF NOT EXISTS / IF EXISTS)
-- ============================================================

-- ─── 1. Tabel TRANSACTIONS ─────────────────────────────────
-- Tambah kolom discount_id, discount_amount, payment_amount, change_amount

ALTER TABLE transactions
  ADD COLUMN IF NOT EXISTS discount_id     INTEGER       DEFAULT NULL,
  ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(15,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS payment_amount  DECIMAL(15,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS change_amount   DECIMAL(15,2) DEFAULT 0;

-- ─── 2. Tabel TRANSACTION_DETAILS ──────────────────────────
-- Tambah kolom diskon per item

ALTER TABLE transaction_details
  ADD COLUMN IF NOT EXISTS discount_type   VARCHAR(10)   DEFAULT NULL,
  ADD COLUMN IF NOT EXISTS discount_value  DECIMAL(12,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0;

-- ─── 3. Verifikasi kolom sudah ada ─────────────────────────
-- Jalankan query di bawah untuk memastikan semua kolom sudah ada:

SELECT column_name, data_type, column_default
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name   = 'transactions'
  AND column_name  IN ('discount_id', 'discount_amount', 'payment_amount', 'change_amount')
ORDER BY column_name;

SELECT column_name, data_type, column_default
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name   = 'transaction_details'
  AND column_name  IN ('discount_type', 'discount_value', 'discount_amount')
ORDER BY column_name;
