-- Migration: Create Discounts Table
-- Created: 2026-02-12
-- Description: Menambahkan tabel discounts untuk fitur promo/diskon

CREATE TABLE IF NOT EXISTS discounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,              -- Nama promo, misal "Diskon Akhir Tahun"
    type VARCHAR(20) NOT NULL,               -- Jenis: 'PERCENTAGE' atau 'FIXED'
    value DECIMAL(10, 2) NOT NULL,           -- Nilai: 10 (10%) atau 5000 (Rp 5.000)
    min_order_amount DECIMAL(15, 2) DEFAULT 0, -- Syarat min. belanja (jika diskon global)
    product_id INT DEFAULT NULL,             -- NULL = Diskon Global. Jika diisi = Diskon Per Produk
    start_date TIMESTAMP NOT NULL,           -- Promo mulai berlaku
    end_date TIMESTAMP NOT NULL,             -- Promo berakhir
    is_active BOOLEAN DEFAULT TRUE,          -- Toggle aktif/nonaktif
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_discounts_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL
);

-- Index update untuk query cepat discount per product
CREATE INDEX idx_discounts_product ON discounts(product_id, is_active, start_date, end_date);

-- Tambahkan kolom discount_id dan discount_amount di tabel transactions
-- Agar kita tahu transaksi mana yang pakai diskon apa
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS discount_id INT,
ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(15, 2) DEFAULT 0;

-- Tambahkan Foreign Key ke tabel discounts (optional, set NULL on delete agar history aman)
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_discounts
FOREIGN KEY (discount_id) REFERENCES discounts(id) ON DELETE SET NULL;

-- Trigger update_updated_at untuk discounts
CREATE OR REPLACE TRIGGER update_discounts_updated_at
BEFORE UPDATE ON discounts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Seed Data: Contoh Diskon
INSERT INTO discounts (name, type, value, min_order_amount, start_date, end_date, is_active)
VALUES 
('Opening Promo', 'PERCENTAGE', 10.00, 0, NOW(), NOW() + INTERVAL '1 month', TRUE),
('Potongan 5rb', 'FIXED', 5000.00, 50000, NOW(), NOW() + INTERVAL '1 month', TRUE);
