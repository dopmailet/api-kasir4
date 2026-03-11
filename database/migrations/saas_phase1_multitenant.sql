-- ============================================================
-- FASE 1: MIGRASI MULTI-TENANT SAAS
-- Jalankan di Supabase SQL Editor secara berurutan
-- ============================================================

-- LANGKAH 1: Buat tabel Paket Langganan
CREATE TABLE IF NOT EXISTS subscription_packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,          -- contoh: 'Gratis', 'Basic', 'Pro'
    max_kasir INT NOT NULL DEFAULT 1,   -- batas maksimal kasir per toko
    max_products INT NOT NULL DEFAULT 100, -- batas maksimal produk per toko
    price NUMERIC NOT NULL DEFAULT 0,   -- harga per bulan
    features_json JSONB DEFAULT '{}',   -- fitur tambahan (fleksibel)
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Seed data paket default
INSERT INTO subscription_packages (name, max_kasir, max_products, price)
VALUES 
    ('Gratis', 1, 50, 0),
    ('Basic', 3, 300, 99000),
    ('Pro', 10, 999999, 199000)
ON CONFLICT DO NOTHING;

-- ============================================================
-- LANGKAH 2: Buat tabel Toko (Tenant)
CREATE TABLE IF NOT EXISTS stores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    subscription_package_id INT REFERENCES subscription_packages(id) DEFAULT 1,
    subscription_end_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Buat Toko Default (untuk data lama yang sudah ada)
INSERT INTO stores (id, name, subscription_package_id, is_active)
VALUES (1, 'Toko Utama (Default)', 1, true)
ON CONFLICT DO NOTHING;

-- ============================================================
-- LANGKAH 3: Tambahkan store_id ke semua tabel operasional
-- AMAN: Defaultnya store_id = 1, jadi data lama tidak rusak

ALTER TABLE users ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE products ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE discounts ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE purchases ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE cash_funds ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE employees ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE payroll ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;
ALTER TABLE cash_funds ADD COLUMN IF NOT EXISTS store_id INT REFERENCES stores(id) DEFAULT 1;

-- ============================================================
-- LANGKAH 4: Update semua data lama ke store_id = 1
UPDATE users SET store_id = 1 WHERE store_id IS NULL;
UPDATE products SET store_id = 1 WHERE store_id IS NULL;
UPDATE categories SET store_id = 1 WHERE store_id IS NULL;
UPDATE customers SET store_id = 1 WHERE store_id IS NULL;
UPDATE discounts SET store_id = 1 WHERE store_id IS NULL;
UPDATE transactions SET store_id = 1 WHERE store_id IS NULL;
UPDATE purchases SET store_id = 1 WHERE store_id IS NULL;
UPDATE suppliers SET store_id = 1 WHERE store_id IS NULL;
UPDATE expenses SET store_id = 1 WHERE store_id IS NULL;
UPDATE cash_funds SET store_id = 1 WHERE store_id IS NULL;
UPDATE employees SET store_id = 1 WHERE store_id IS NULL;
UPDATE payroll SET store_id = 1 WHERE store_id IS NULL;
UPDATE expenses SET store_id = 1 WHERE store_id IS NULL;
UPDATE cash_funds SET store_id = 1 WHERE store_id IS NULL;

-- ============================================================
-- LANGKAH 5: Buat kolom superadmin di users
-- superadmin tidak punya store_id, dia pemilik platform
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_superadmin BOOLEAN DEFAULT false;

-- Verifikasi hasil
SELECT 'stores' AS tabel, count(*) FROM stores
UNION ALL SELECT 'subscription_packages', count(*) FROM subscription_packages
UNION ALL SELECT 'users dengan store_id', count(*) FROM users WHERE store_id IS NOT NULL;
