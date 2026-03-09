-- ==========================================
-- MIGRATION: Tambah Modul Pembelian (Purchase)
-- Tanggal: 2026-02-16
-- Deskripsi: Menambah tabel purchases dan purchase_items
--            untuk mencatat pembelian/pengadaan barang dari supplier
-- ==========================================

-- 1. TABLE: PURCHASES (Header Pembelian)
-- Menyimpan informasi header setiap pembelian
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    supplier_name VARCHAR(150),                         -- Nama supplier (optional)
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,     -- Total harga pembelian
    notes TEXT,                                          -- Catatan tambahan (optional)
    created_by INT,                                      -- Admin yang mencatat
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_purchases_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Index untuk query riwayat pembelian berdasarkan tanggal
CREATE INDEX IF NOT EXISTS idx_purchases_created_at ON purchases(created_at);

-- 2. TABLE: PURCHASE_ITEMS (Detail Item Pembelian)
-- Menyimpan detail item yang dibeli dalam setiap pembelian
CREATE TABLE IF NOT EXISTS purchase_items (
    id SERIAL PRIMARY KEY,
    purchase_id INT NOT NULL,                            -- FK ke header pembelian
    product_id INT,                                      -- FK ke produk (NULL jika produk baru)
    product_name VARCHAR(150) NOT NULL,                  -- Nama produk (disimpan untuk riwayat)
    quantity INT NOT NULL CHECK (quantity > 0),           -- Jumlah beli (harus > 0)
    buy_price DECIMAL(15, 2) NOT NULL CHECK (buy_price >= 0), -- Harga beli per unit
    sell_price DECIMAL(15, 2),                           -- Harga jual (hanya untuk produk baru)
    category_id INT,                                     -- Kategori (hanya untuk produk baru)
    subtotal DECIMAL(15, 2) NOT NULL,                    -- quantity × buy_price
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_purchase_items_purchase FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    CONSTRAINT fk_purchase_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    CONSTRAINT fk_purchase_items_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Index untuk query detail per pembelian
CREATE INDEX IF NOT EXISTS idx_purchase_items_purchase_id ON purchase_items(purchase_id);
