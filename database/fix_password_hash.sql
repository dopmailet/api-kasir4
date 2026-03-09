-- ============================================
-- FIX PASSWORD HASH - Update User Passwords
-- ============================================
-- Script ini untuk update password hash yang correct
-- Jalankan di Supabase SQL Editor
-- ============================================

-- Update Admin password (admin123)
UPDATE users 
SET password = '$2a$10$EsHt4SEuTd8uLg9JoQ5jO.k90YbTmkGv0s0yZdNQ/T1QmqZowHRYa' 
WHERE username = 'admin';

-- Update Kasir1 password (kasir123)
UPDATE users 
SET password = '$2a$10$mTs5c9h0Nq7TVdscKY4E1OqHTSm92ec9vzOwMmyT.KcZmi2hBr33y' 
WHERE username = 'kasir1';

-- Update Kasir2 password (kasir123)
UPDATE users 
SET password = '$2a$10$mTs5c9h0Nq7TVdscKY4E1OqHTSm92ec9vzOwMmyT.KcZmi2hBr33y' 
WHERE username = 'kasir2';

-- Verify update
SELECT username, role, is_active, 
       LEFT(password, 20) as password_hash_preview
FROM users
ORDER BY id;

-- ============================================
-- Expected Result:
-- username | role  | is_active | password_hash_preview
-- ---------+-------+-----------+----------------------
-- admin    | admin | true      | $2a$10$EsHt4SEuTd8u
-- kasir1   | kasir | true      | $2a$10$mTs5c9h0Nq7T
-- kasir2   | kasir | true      | $2a$10$mTs5c9h0Nq7T
-- ============================================

-- Credentials after update:
-- admin / admin123
-- kasir1 / kasir123
-- kasir2 / kasir123
