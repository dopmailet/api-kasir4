# 🗄️ Database Migration Guide - Session 4

## 📋 Ringkasan

Migration ini menambahkan fitur **Authentication & Profit Tracking** ke aplikasi Kasir API.

---

## ✅ Cara Termudah: 1 File Saja!

### **File yang Digunakan:**
```
database/migration_session_4_complete.sql
```

### **Langkah-langkah:**

1. **Login ke Supabase Dashboard**
   - Buka https://supabase.com
   - Login dengan akun Anda
   - Pilih project kasir-api

2. **Buka SQL Editor**
   - Klik menu "SQL Editor" di sidebar kiri
   - Klik tombol "New query"

3. **Copy-Paste Migration**
   - Buka file `database/migration_session_4_complete.sql`
   - Copy **SELURUH ISI FILE** (Ctrl+A, Ctrl+C)
   - Paste ke SQL Editor di Supabase (Ctrl+V)

4. **Run Migration**
   - Klik tombol "Run" (atau tekan Ctrl+Enter)
   - Tunggu sampai selesai (5-10 detik)
   - ✅ Jika tidak ada error, migration berhasil!

5. **Verify Migration**
   - Scroll ke bawah di SQL Editor
   - Akan ada beberapa verification queries
   - Jalankan satu per satu untuk memastikan semua berhasil

---

## 🔍 Apa yang Dilakukan Migration Ini?

### **1. Membuat Tabel `users`**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nama_lengkap VARCHAR(100) NOT NULL,
    role VARCHAR(20) CHECK (role IN ('admin', 'kasir')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Untuk apa?**
- Menyimpan data user (admin & kasir)
- Authentication (login)
- Role-based access control

### **2. Menambahkan Kolom ke Tabel `products`**
```sql
ALTER TABLE products 
ADD COLUMN harga_beli NUMERIC(10,2),
ADD COLUMN created_by INTEGER;
```

**Untuk apa?**
- `harga_beli`: Tracking harga modal untuk hitung profit
- `created_by`: Tracking siapa yang menambahkan produk

### **3. Menambahkan Kolom ke Tabel `transactions`**
```sql
ALTER TABLE transactions 
ADD COLUMN kasir_id INTEGER;

ALTER TABLE transaction_details 
ADD COLUMN harga_beli NUMERIC(10,2);
```

**Untuk apa?**
- `kasir_id`: Tracking kasir yang melakukan transaksi
- `harga_beli`: Snapshot harga beli saat transaksi (untuk laporan profit)

### **4. Insert Default Users**
```sql
INSERT INTO users (username, password, nama_lengkap, role) VALUES
('admin', '$2a$10$...', 'Administrator', 'admin'),
('kasir1', '$2a$10$...', 'Kasir Pagi', 'kasir'),
('kasir2', '$2a$10$...', 'Kasir Sore', 'kasir');
```

**Default Credentials:**
- `admin` / `admin123` (role: admin)
- `kasir1` / `kasir123` (role: kasir)
- `kasir2` / `kasir123` (role: kasir)

⚠️ **PENTING**: Ganti password ini di production!

---

## ✅ Verification Checklist

Setelah run migration, verify dengan query ini:

### **1. Cek Tabel Users**
```sql
SELECT id, username, nama_lengkap, role, is_active 
FROM users 
ORDER BY id;
```

**Expected Result:**
```
id | username | nama_lengkap    | role  | is_active
---+----------+-----------------+-------+-----------
1  | admin    | Administrator   | admin | true
2  | kasir1   | Kasir Pagi      | kasir | true
3  | kasir2   | Kasir Sore      | kasir | true
```

### **2. Cek Kolom Baru di Products**
```sql
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'products' 
  AND column_name IN ('harga_beli', 'created_by');
```

**Expected Result:**
```
column_name | data_type | is_nullable
------------+-----------+-------------
harga_beli  | numeric   | YES
created_by  | integer   | YES
```

### **3. Cek Kolom Baru di Transactions**
```sql
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'transactions' 
  AND column_name = 'kasir_id';
```

**Expected Result:**
```
column_name | data_type | is_nullable
------------+-----------+-------------
kasir_id    | integer   | YES
```

---

## 🐛 Troubleshooting

### **Error: "relation 'users' already exists"**
**Penyebab**: Migration sudah pernah dijalankan sebelumnya.

**Solusi**: 
- Jika ingin re-run, jalankan rollback dulu (ada di bagian bawah migration file)
- Atau skip error ini (tidak masalah)

### **Error: "column 'harga_beli' already exists"**
**Penyebab**: Kolom sudah ada dari migration sebelumnya.

**Solusi**: Skip error ini, tidak masalah.

### **Error: "duplicate key value violates unique constraint"**
**Penyebab**: Default users sudah ada.

**Solusi**: 
- Migration akan otomatis delete users lama sebelum insert
- Jika masih error, hapus manual: `DELETE FROM users WHERE username IN ('admin', 'kasir1', 'kasir2');`

---

## 🔄 Rollback (Jika Perlu)

Jika ingin membatalkan migration, ada rollback script di bagian bawah file `migration_session_4_complete.sql`.

**Cara rollback:**
1. Buka file migration
2. Scroll ke bagian "ROLLBACK"
3. Uncomment script rollback (hapus `/*` dan `*/`)
4. Copy-paste ke SQL Editor
5. Run

⚠️ **WARNING**: Rollback akan menghapus semua data users dan kolom baru!

---

## 📁 File Migration Lainnya (Opsional)

Jika Anda ingin menjalankan migration secara terpisah (untuk debugging), tersedia juga file individual:

1. `migration_create_users.sql` - Hanya tabel users
2. `migration_add_harga_beli_to_products.sql` - Hanya kolom harga_beli
3. `migration_add_kasir_to_transactions.sql` - Hanya kolom kasir_id
4. `seed_default_users.sql` - Hanya insert users

**Tapi disarankan pakai `migration_session_4_complete.sql` saja!**

---

## 📞 Need Help?

Jika ada masalah:
1. Cek error message di Supabase SQL Editor
2. Jalankan verification queries untuk cek status
3. Lihat dokumentasi lengkap di `docs/AUTH_IMPLEMENTATION.md`

---

**Status**: ✅ Ready to run  
**Estimated Time**: 5-10 detik  
**Risk Level**: Low (ada rollback script)
