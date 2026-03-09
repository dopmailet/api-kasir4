-- ============================================
-- CLEANUP DUPLICATE PRODUCTS
-- ============================================
-- Script ini akan menggabungkan produk duplikat
-- dengan menjumlahkan stok dan mengambil harga terbaru
--
-- PENTING: Script ini akan menghapus data duplikat!
-- Pastikan Anda sudah backup data sebelum menjalankan
-- ============================================

-- STEP 1: Cek data duplikat yang ada
-- Jalankan query ini dulu untuk melihat data duplikat
SELECT 
    nama,
    COUNT(*) as jumlah_duplikat,
    STRING_AGG(id::text, ', ') as id_list,
    SUM(stok) as total_stok,
    MAX(harga) as harga_terbaru
FROM products
GROUP BY nama
HAVING COUNT(*) > 1
ORDER BY nama;

-- STEP 2: Backup data (PENTING!)
-- Buat backup table sebelum melakukan perubahan
CREATE TABLE IF NOT EXISTS products_backup AS 
SELECT * FROM products;

-- Verifikasi backup berhasil
SELECT COUNT(*) as total_backup FROM products_backup;

-- STEP 3: Update produk pertama dengan total stok
-- Untuk setiap nama produk yang duplikat:
-- - Ambil ID terkecil (produk pertama)
-- - Update stoknya dengan total stok dari semua duplikat
-- - Update harganya dengan harga terbaru
WITH duplicates AS (
    SELECT 
        nama,
        MIN(id) as keep_id,
        SUM(stok) as total_stok,
        MAX(harga) as latest_harga
    FROM products
    GROUP BY nama
    HAVING COUNT(*) > 1
)
UPDATE products p
SET 
    stok = d.total_stok,
    harga = d.latest_harga,
    updated_at = CURRENT_TIMESTAMP
FROM duplicates d
WHERE p.id = d.keep_id;

-- STEP 4: Hapus produk duplikat (yang bukan ID pertama)
-- Hanya keep produk dengan ID terkecil untuk setiap nama
DELETE FROM products p
WHERE EXISTS (
    SELECT 1 
    FROM (
        SELECT nama, MIN(id) as keep_id
        FROM products
        GROUP BY nama
        HAVING COUNT(*) > 1
    ) d
    WHERE p.nama = d.nama AND p.id != d.keep_id
);

-- STEP 5: Verifikasi tidak ada duplikat lagi
-- Query ini seharusnya return 0 rows
SELECT 
    nama,
    COUNT(*) as jumlah
FROM products
GROUP BY nama
HAVING COUNT(*) > 1;

-- STEP 6: Sekarang bisa tambahkan UNIQUE constraint
-- Hapus index biasa yang sudah ada (jika ada)
DROP INDEX IF EXISTS idx_products_nama;

-- Tambahkan UNIQUE constraint
ALTER TABLE products 
ADD CONSTRAINT products_nama_unique UNIQUE (nama);

-- Buat index UNIQUE
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_nama_unique ON products(nama);

-- ============================================
-- SELESAI! 🎉
-- ============================================
-- Verifikasi hasil akhir
SELECT 
    COUNT(*) as total_products,
    SUM(stok) as total_stok
FROM products;

-- Jika ada masalah, restore dari backup:
-- DROP TABLE products;
-- ALTER TABLE products_backup RENAME TO products;
