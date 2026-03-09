-- MIGRATION: Add Category Support to Discounts
-- Description: Menambahkan kolom category_id untuk memungkinkan diskon per kategori

-- 1. Tambahkan kolom category_id
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'discounts' AND column_name = 'category_id') THEN
        ALTER TABLE discounts ADD COLUMN category_id INT DEFAULT NULL;
    END IF;
END $$;

-- 2. Tambahkan Foreign Key ke Categories
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_discounts_category') THEN
        ALTER TABLE discounts 
        ADD CONSTRAINT fk_discounts_category 
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;
    END IF;
END $$;

-- 3. Update Index (tambahkan category_id ke index pencarian)
DROP INDEX IF EXISTS idx_discounts_product_active; -- Hapus index lama
CREATE INDEX idx_discounts_lookup ON discounts(product_id, category_id, is_active);
