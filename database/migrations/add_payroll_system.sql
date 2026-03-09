-- Migration: Add Payroll System (employees and payroll tables)
-- Date: 2026-03-02

-- ==========================================
-- 1. TABLE: EMPLOYEES
-- ==========================================
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    posisi VARCHAR(50) NOT NULL,          -- misal: kasir, cleaning, kurir, admin
    gaji_pokok DECIMAL(15,2) NOT NULL DEFAULT 0,
    no_hp VARCHAR(20),
    alamat TEXT,
    tanggal_masuk DATE DEFAULT CURRENT_DATE,
    aktif BOOLEAN DEFAULT TRUE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,  -- nullable, link ke user jika ada
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk updated_at
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
CREATE TRIGGER update_employees_updated_at
BEFORE UPDATE ON employees
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Index untuk filter berdasarkan aktif (Soft Delete)
CREATE INDEX IF NOT EXISTS idx_employees_aktif ON employees(aktif);

-- ==========================================
-- 2. TABLE: PAYROLL
-- ==========================================
CREATE TABLE IF NOT EXISTS payroll (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    periode VARCHAR(20),                  -- label bebas: "Maret 2026", "Minggu 1 Maret", dll
    gaji_pokok DECIMAL(15,2) NOT NULL,
    bonus DECIMAL(15,2) DEFAULT 0,
    potongan DECIMAL(15,2) DEFAULT 0,
    total DECIMAL(15,2) NOT NULL,         -- gaji_pokok + bonus - potongan
    catatan TEXT,
    paid_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk updated_at
DROP TRIGGER IF EXISTS update_payroll_updated_at ON payroll;
CREATE TRIGGER update_payroll_updated_at
BEFORE UPDATE ON payroll
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Index untuk querying / report timezone-aware 
CREATE INDEX IF NOT EXISTS idx_payroll_employee_id ON payroll(employee_id);
CREATE INDEX IF NOT EXISTS idx_payroll_paid_at ON payroll(paid_at);
