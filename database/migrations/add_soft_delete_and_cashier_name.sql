-- Menambahkan kolom is_active untuk soft delete kasir
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
UPDATE users SET is_active = TRUE WHERE is_active IS NULL;

-- Menambahkan kolom cashier_name untuk menyimpan data nama kasir yang paten
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS cashier_name VARCHAR(100);
