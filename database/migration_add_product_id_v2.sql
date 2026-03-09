-- PERBAIKAN SCHEMA DISCOUNTS
-- Jalankan ini di SQL Editor Supabase/PGAdmin

-- 1. Tambahkan kolom product_id jika belum ada
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'discounts' AND column_name = 'product_id') THEN
        ALTER TABLE discounts ADD COLUMN product_id INT DEFAULT NULL;
    END IF;
END $$;

-- 2. Tambahkan Foreign Key Constraint ke tabel products
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_discounts_product') THEN
        ALTER TABLE discounts 
        ADD CONSTRAINT fk_discounts_product 
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;
    END IF;
END $$;

-- 3. Buat Index agar pencarian diskon produk cepat
CREATE INDEX IF NOT EXISTS idx_discounts_product_active ON discounts(product_id, is_active);

-- 4. Pastikan kolom discount_id dan discount_amount ada di transactions (jaga-jaga)
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'transactions' AND column_name = 'discount_id') THEN
        ALTER TABLE transactions ADD COLUMN discount_id INT DEFAULT NULL;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'transactions' AND column_name = 'discount_amount') THEN
        ALTER TABLE transactions ADD COLUMN discount_amount DECIMAL(15, 2) DEFAULT 0;
    END IF;
END $$;
