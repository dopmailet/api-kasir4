-- Menambahkan kolom is_featured ke tabel products
ALTER TABLE products 
ADD COLUMN is_featured BOOLEAN DEFAULT FALSE;
