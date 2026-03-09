-- Migration: Add Expenses Table
-- Date: 2026-03-05

-- ==========================================
-- 1. TABLE: EXPENSES
-- ==========================================
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    category VARCHAR(50) NOT NULL,        -- 'listrik','wifi','sewa','marketing','maintenance','lainnya'
    description VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    expense_date DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_period VARCHAR(20),          -- 'monthly','quarterly','yearly' (nullable)
    notes TEXT,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ==========================================
-- 2. FUNCTION: update_timestamp (jika belum ada)
-- ==========================================
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==========================================
-- 3. TRIGGER: Update updated_at on expenses
-- ==========================================
DROP TRIGGER IF EXISTS trigger_update_timestamp_expenses ON expenses;

CREATE TRIGGER trigger_update_timestamp_expenses
BEFORE UPDATE ON expenses
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
