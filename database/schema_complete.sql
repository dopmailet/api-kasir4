-- MASTER SCHEMA KASIR API v1.1 (Synced with Production)
-- Created Date: 2026-02-12
-- Compatible with: PostgreSQL (Supabase)

-- ==========================================
-- 1. UTILITIES
-- ==========================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ==========================================
-- 2. TABLE: USERS
-- ==========================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, 
    nama_lengkap VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'kasir' CHECK (role IN ('admin', 'kasir')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ==========================================
-- 3. TABLE: CATEGORIES
-- ==========================================
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ==========================================
-- 4. TABLE: PRODUCTS
-- ==========================================
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(150) NOT NULL UNIQUE,
    harga DECIMAL(15, 2) NOT NULL CHECK (harga >= 0),
    stok INT NOT NULL DEFAULT 0 CHECK (stok >= 0),
    harga_beli DECIMAL(15, 2) CHECK (harga_beli IS NULL OR harga_beli >= 0),
    category_id INT, 
    created_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_products_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_products_nama ON products(nama);
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ==========================================
-- 5. TABLE: DISCOUNTS
-- ==========================================
CREATE TABLE IF NOT EXISTS discounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,              
    value DECIMAL(10, 2) NOT NULL,          
    min_order_amount DECIMAL(15, 2) DEFAULT 0, 
    product_id INT DEFAULT NULL,            
    category_id INT DEFAULT NULL, -- New: Discount per Category
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_discounts_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    CONSTRAINT fk_discounts_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_discounts_active ON discounts(is_active, start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_discounts_lookup ON discounts(product_id, category_id, is_active);
CREATE TRIGGER update_discounts_updated_at
BEFORE UPDATE ON discounts
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ==========================================
-- 6. TABLE: TRANSACTIONS
-- ==========================================
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount DECIMAL(15, 2) NOT NULL,
    kasir_id INT,
    discount_id INT DEFAULT NULL,
    discount_amount DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transactions_kasir FOREIGN KEY (kasir_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_transactions_discount FOREIGN KEY (discount_id) REFERENCES discounts(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);

-- ==========================================
-- 7. TABLE: TRANSACTION DETAILS
-- ==========================================
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price DECIMAL(15, 2) NOT NULL,    
    subtotal DECIMAL(15, 2) NOT NULL, 
    harga_beli DECIMAL(15, 2), -- Snapshot harga beli saat transaksi (Profit Analysis)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_details_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    CONSTRAINT fk_details_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL
);

-- ==========================================
-- 8. TABLE: PURCHASES (Header Pembelian)
-- ==========================================
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    supplier_name VARCHAR(150),
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_purchases_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_purchases_created_at ON purchases(created_at);

-- ==========================================
-- 9. TABLE: PURCHASE_ITEMS (Detail Item Pembelian)
-- ==========================================
CREATE TABLE IF NOT EXISTS purchase_items (
    id SERIAL PRIMARY KEY,
    purchase_id INT NOT NULL,
    product_id INT,
    product_name VARCHAR(150) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    buy_price DECIMAL(15, 2) NOT NULL CHECK (buy_price >= 0),
    sell_price DECIMAL(15, 2),
    category_id INT,
    subtotal DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_purchase_items_purchase FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    CONSTRAINT fk_purchase_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    CONSTRAINT fk_purchase_items_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_purchase_items_purchase_id ON purchase_items(purchase_id);

-- ==========================================
INSERT INTO users (username, password, nama_lengkap, role) VALUES
('admin', '$2a$10$Rz/u.W2/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3example', 'Administrator', 'admin'),
('kasir1', '$2a$10$Rz/u.W2/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3/tF3example', 'Kasir Utama', 'kasir')
ON CONFLICT (username) DO NOTHING;

INSERT INTO categories (name, description) VALUES 
('Minuman', 'Aneka minuman segar'), 
('Makanan', 'Makanan berat'), 
('Snack', 'Camilan ringan')
ON CONFLICT DO NOTHING;
