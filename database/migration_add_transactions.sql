-- Migration untuk menambahkan tabel transactions dan transaction_details
-- Tabel ini digunakan untuk menyimpan data transaksi checkout

-- Tabel transactions: menyimpan header transaksi
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel transaction_details: menyimpan detail item yang dibeli
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat query
CREATE INDEX IF NOT EXISTS idx_transaction_details_transaction_id ON transaction_details(transaction_id);
CREATE INDEX IF NOT EXISTS idx_transaction_details_product_id ON transaction_details(product_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
