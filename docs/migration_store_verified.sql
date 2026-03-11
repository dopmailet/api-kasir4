-- Migrasi: Tambahkan kolom is_verified ke tabel stores
-- Jalankan sekali di Supabase SQL Editor

ALTER TABLE stores ADD COLUMN IF NOT EXISTS is_verified BOOLEAN NOT NULL DEFAULT false;

-- Toko yang sudah pernah pakai paket berbayar (package_id > 1) langsung ditandai verified
UPDATE stores SET is_verified = true WHERE subscription_package_id > 1;

-- Toko #1 (Toko Default / Perintis) langsung verified
UPDATE stores SET is_verified = true WHERE id = 1;
