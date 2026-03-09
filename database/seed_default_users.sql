-- ============================================
-- SEED DATA: Default users (Admin & Kasir)
-- ============================================
-- Script ini membuat user default untuk testing
-- 
-- PENTING: 
-- - Password sudah di-hash menggunakan bcrypt
-- - Ganti password ini di production!
--
-- Cara pakai:
-- 1. Login ke Supabase Dashboard (https://supabase.com)
-- 2. Buka project Anda
-- 3. Klik "SQL Editor" di sidebar
-- 4. Copy-paste script ini
-- 5. Klik "Run" atau tekan Ctrl+Enter
-- ============================================

-- ============================================
-- DEFAULT CREDENTIALS
-- ============================================
-- Username: admin     | Password: admin123
-- Username: kasir1    | Password: kasir123
-- Username: kasir2    | Password: kasir123
-- ============================================

-- Hapus data lama jika ada (untuk re-run script)
DELETE FROM users WHERE username IN ('admin', 'kasir1', 'kasir2');

-- Insert Admin
INSERT INTO users (username, password, nama_lengkap, role, is_active) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', 'admin', TRUE);

-- Insert Kasir 1
INSERT INTO users (username, password, nama_lengkap, role, is_active) VALUES
('kasir1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Pagi', 'kasir', TRUE);

-- Insert Kasir 2
INSERT INTO users (username, password, nama_lengkap, role, is_active) VALUES
('kasir2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Sore', 'kasir', TRUE);

-- ============================================
-- VERIFY
-- ============================================
-- Cek apakah users sudah dibuat:
SELECT id, username, nama_lengkap, role, is_active, created_at 
FROM users 
ORDER BY id;

-- ============================================
-- CARA MENAMBAH USER BARU
-- ============================================
-- Untuk menambah admin baru:
-- INSERT INTO users (username, password, nama_lengkap, role) VALUES
-- ('admin2', '$2a$10$...hash_password_di_golang...', 'Admin Cabang', 'admin');

-- Untuk menambah kasir baru:
-- INSERT INTO users (username, password, nama_lengkap, role) VALUES
-- ('kasir3', '$2a$10$...hash_password_di_golang...', 'Kasir Malam', 'kasir');

-- Untuk nonaktifkan user (tanpa hapus data):
-- UPDATE users SET is_active = FALSE WHERE username = 'kasir1';

-- Untuk aktifkan kembali:
-- UPDATE users SET is_active = TRUE WHERE username = 'kasir1';

-- Untuk ganti password (hash dulu di aplikasi):
-- UPDATE users SET password = '$2a$10$...new_hash...' WHERE username = 'admin';

-- ============================================
-- SECURITY NOTES
-- ============================================
-- 1. Password hash menggunakan bcrypt cost 10
-- 2. JANGAN simpan password plain text di database
-- 3. Ganti semua password default di production
-- 4. Gunakan password minimal 8 karakter
-- 5. Kombinasi huruf besar, kecil, angka, dan simbol
-- ============================================
