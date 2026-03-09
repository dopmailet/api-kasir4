-- Menambahkan kolom created_by pada tabel transactions
ALTER TABLE transactions ADD COLUMN created_by INT REFERENCES users(id);

-- Optional: Jika ingin data lama tidak error/memiliki default (misal admin punya ID 1)
-- ALTER TABLE transactions ADD COLUMN created_by INT DEFAULT 1 REFERENCES users(id);
-- Namun demi konsistensi, biarkan NULL untuk transaksi lama.
