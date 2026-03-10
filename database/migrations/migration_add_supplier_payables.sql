-- ============================================
-- ADD SUPPLIER & PAYABLES TABLES
-- ============================================

-- 1. Create suppliers table if not exists
CREATE TABLE IF NOT EXISTS suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(100),
    address TEXT,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    total_purchases INTEGER DEFAULT 0,
    total_spent NUMERIC(15,2) DEFAULT 0,
    total_payable NUMERIC(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create supplier_payables table (Hutang Supplier)
CREATE TABLE IF NOT EXISTS supplier_payables (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    purchase_id INTEGER, -- Optional, links to purchases table if you have it
    amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'partial', 'paid')),
    due_date DATE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create payable_payments table (Riwayat Pembayaran Hutang)
CREATE TABLE IF NOT EXISTS payable_payments (
    id SERIAL PRIMARY KEY,
    payable_id INTEGER NOT NULL REFERENCES supplier_payables(id) ON DELETE CASCADE,
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    payment_date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_supplier_payables_supplier_id ON supplier_payables(supplier_id);
CREATE INDEX IF NOT EXISTS idx_supplier_payables_status ON supplier_payables(status);
CREATE INDEX IF NOT EXISTS idx_payable_payments_payable_id ON payable_payments(payable_id);
CREATE INDEX IF NOT EXISTS idx_payable_payments_date ON payable_payments(payment_date);

-- Trigger to update updated_at on suppliers
CREATE OR REPLACE FUNCTION update_suppliers_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_suppliers_updated_at ON suppliers;
CREATE TRIGGER trigger_update_suppliers_updated_at
    BEFORE UPDATE ON suppliers
    FOR EACH ROW
    EXECUTE FUNCTION update_suppliers_updated_at();

-- Trigger to update updated_at on supplier_payables
CREATE OR REPLACE FUNCTION update_supplier_payables_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_supplier_payables_updated_at ON supplier_payables;
CREATE TRIGGER trigger_update_supplier_payables_updated_at
    BEFORE UPDATE ON supplier_payables
    FOR EACH ROW
    EXECUTE FUNCTION update_supplier_payables_updated_at();
