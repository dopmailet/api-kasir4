-- Migration: Add store_id to cash_funds
ALTER TABLE cash_funds ADD COLUMN IF NOT EXISTS store_id INTEGER NOT NULL DEFAULT 1 REFERENCES stores(id);
CREATE INDEX IF NOT EXISTS idx_cash_funds_store ON cash_funds(store_id);
