-- ============================================================
-- KasirPOS - Full Database Schema
-- Jalankan di Supabase SQL Editor untuk setup database baru
-- Generated from: postgres.rgokcvgkeixqjisurvlh (Supabase)
-- ============================================================

-- 1. USERS
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nama_lengkap VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'kasir',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. CATEGORIES
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    discount_type VARCHAR(20) DEFAULT NULL,
    discount_value NUMERIC DEFAULT 0
);

-- 3. PRODUCTS
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    harga INT NOT NULL,
    stok INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    harga_beli NUMERIC,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    default_discount_type VARCHAR(20),
    default_discount_value NUMERIC,
    barcode VARCHAR(100) UNIQUE,
    is_featured BOOLEAN DEFAULT false
);

-- 4. CUSTOMERS
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    address TEXT,
    notes TEXT,
    loyalty_points INT DEFAULT 0,
    total_transactions INT DEFAULT 0,
    total_spent NUMERIC DEFAULT 0,
    last_transaction_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    card_number VARCHAR(50) DEFAULT NULL
);

-- 5. CUSTOMER_SETTINGS
CREATE TABLE IF NOT EXISTS customer_settings (
    id SERIAL PRIMARY KEY,
    show_customer_in_pos BOOLEAN DEFAULT true,
    enable_loyalty_points BOOLEAN DEFAULT true
);

-- 6. DISCOUNTS
CREATE TABLE IF NOT EXISTS discounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    value NUMERIC NOT NULL,
    min_order_amount NUMERIC DEFAULT 0,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    product_id INT REFERENCES products(id) ON DELETE SET NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL
);

-- 7. TRANSACTIONS
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    kasir_id INT REFERENCES users(id) ON DELETE SET NULL,
    discount_id INT REFERENCES discounts(id) ON DELETE SET NULL,
    discount_amount NUMERIC DEFAULT 0,
    payment_amount NUMERIC DEFAULT 0,
    change_amount NUMERIC DEFAULT 0,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    customer_id INT REFERENCES customers(id) ON DELETE SET NULL
);

-- 8. TRANSACTION_DETAILS
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL,
    price NUMERIC NOT NULL,
    subtotal NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    harga_beli NUMERIC,
    discount_type VARCHAR(10) DEFAULT NULL,
    discount_value NUMERIC DEFAULT 0,
    discount_amount NUMERIC DEFAULT 0
);

-- 9. SUPPLIERS
CREATE TABLE IF NOT EXISTS suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(255),
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    notes TEXT,
    is_active BOOLEAN DEFAULT true,
    total_purchases INT DEFAULT 0,
    total_spent NUMERIC DEFAULT 0,
    total_payable NUMERIC DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- 10. PURCHASES
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    supplier_id INT REFERENCES suppliers(id) ON DELETE SET NULL,
    supplier_name VARCHAR(150),
    total_amount NUMERIC NOT NULL DEFAULT 0,
    notes TEXT,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_method VARCHAR(10) DEFAULT 'cash',
    payment_status VARCHAR(10) DEFAULT 'paid',
    paid_amount NUMERIC DEFAULT 0,
    remaining_amount NUMERIC DEFAULT 0,
    due_date DATE,
    payment_notes TEXT
);

-- 11. PURCHASE_ITEMS
CREATE TABLE IF NOT EXISTS purchase_items (
    id SERIAL PRIMARY KEY,
    purchase_id INT NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(150) NOT NULL,
    quantity INT NOT NULL,
    buy_price NUMERIC NOT NULL,
    sell_price NUMERIC,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    subtotal NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 12. SUPPLIER_PAYABLES
CREATE TABLE IF NOT EXISTS supplier_payables (
    id SERIAL PRIMARY KEY,
    supplier_id INT NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    purchase_id INT REFERENCES purchases(id) ON DELETE SET NULL,
    amount NUMERIC NOT NULL,
    paid_amount NUMERIC DEFAULT 0,
    status VARCHAR(20) DEFAULT 'unpaid',
    due_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- 13. PAYABLE_PAYMENTS
CREATE TABLE IF NOT EXISTS payable_payments (
    id SERIAL PRIMARY KEY,
    payable_id INT NOT NULL REFERENCES supplier_payables(id) ON DELETE CASCADE,
    amount NUMERIC NOT NULL,
    payment_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- 14. EMPLOYEES
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    posisi VARCHAR(50) NOT NULL,
    gaji_pokok NUMERIC NOT NULL DEFAULT 0,
    no_hp VARCHAR(20),
    alamat TEXT,
    tanggal_masuk DATE DEFAULT CURRENT_DATE,
    aktif BOOLEAN DEFAULT true,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 15. PAYROLL
CREATE TABLE IF NOT EXISTS payroll (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    periode VARCHAR(20),
    gaji_pokok NUMERIC NOT NULL,
    bonus NUMERIC DEFAULT 0,
    potongan NUMERIC DEFAULT 0,
    total NUMERIC NOT NULL,
    catatan TEXT,
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 16. EXPENSES
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    amount NUMERIC NOT NULL,
    expense_date DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT false,
    recurring_period VARCHAR(20),
    notes TEXT,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- 17. CASH_FUNDS
CREATE TABLE IF NOT EXISTS cash_funds (
    id SERIAL PRIMARY KEY,
    type VARCHAR(3) NOT NULL CHECK (type IN ('in', 'out')),
    amount NUMERIC NOT NULL,
    date DATE NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- 18. LOYALTY_TRANSACTIONS
CREATE TABLE IF NOT EXISTS loyalty_transactions (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    transaction_id INT REFERENCES transactions(id) ON DELETE SET NULL,
    points INT NOT NULL,
    type VARCHAR(10) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now()
);

-- 19. APP_SETTINGS
CREATE TABLE IF NOT EXISTS app_settings (
    key VARCHAR(100) PRIMARY KEY,
    value_json JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- SEED DATA: Default Admin User
-- Password: admin123 (bcrypt hash)
-- ============================================================
INSERT INTO users (username, password, nama_lengkap, role)
VALUES ('admin', '$2a$10$YhEzQCq9GVe7p.Hx0q0Jz.8mvJr6LKf4pIQx0iY6qRq5u9Dh3bxXu', 'Administrator', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Default customer settings
INSERT INTO customer_settings (show_customer_in_pos, enable_loyalty_points)
VALUES (true, true)
ON CONFLICT DO NOTHING;
