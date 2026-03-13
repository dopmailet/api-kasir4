-- Migration: Add loyalty_transactions table and loyalty_points to customers
CREATE TABLE IF NOT EXISTS loyalty_transactions (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    transaction_id INT REFERENCES transactions(id) ON DELETE SET NULL,
    points INT NOT NULL,              
    type VARCHAR(10) NOT NULL,        
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    store_id INTEGER NOT NULL DEFAULT 1 REFERENCES stores(id)
);

CREATE INDEX IF NOT EXISTS idx_loyalty_transactions_store ON loyalty_transactions(store_id);

ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_points INT DEFAULT 0;
