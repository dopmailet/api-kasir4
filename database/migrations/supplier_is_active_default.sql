-- Migration: Set default value for is_active in suppliers table
ALTER TABLE suppliers ALTER COLUMN is_active SET DEFAULT TRUE;

-- Data fix: Activate existing suppliers without status (including those accidentally set to FALSE if they were created before default TRUE was introduced)
UPDATE suppliers SET is_active = TRUE WHERE is_active IS NULL OR is_active = FALSE;
