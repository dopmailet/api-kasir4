# 🔄 Database Migration & Deployment Guide

## 📋 Skenario Penggunaan

### **Skenario 1: Update Database yang Sudah Ada (Session 4)**
Anda sudah punya database dengan tabel products, categories, transactions.  
**Gunakan:** `migration_session_4_complete.sql`

### **Skenario 2: Setup Supabase Baru dari Nol**
Anda buat project Supabase baru atau ingin clone database.  
**Gunakan:** `complete_schema.sql`

---

## 🆕 Skenario 1: Update Database Existing

### **File yang Digunakan:**
```
database/migration_session_4_complete.sql
```

### **Kapan Digunakan:**
- ✅ Database sudah ada (products, categories, transactions)
- ✅ Hanya ingin menambahkan fitur Session 4 (auth & profit tracking)
- ✅ Upgrade dari Session 3 ke Session 4

### **Apa yang Dilakukan:**
- ✅ Membuat tabel `users`
- ✅ Menambahkan kolom `harga_beli` & `created_by` ke `products`
- ✅ Menambahkan kolom `kasir_id` ke `transactions`
- ✅ Menambahkan kolom `harga_beli` ke `transaction_details`
- ✅ Insert 3 default users

### **Langkah-langkah:**
1. Login ke Supabase Dashboard
2. Buka SQL Editor
3. Copy-paste isi `migration_session_4_complete.sql`
4. Run
5. ✅ Selesai!

---

## 🌟 Skenario 2: Setup Database Baru (Recommended untuk Clone)

### **File yang Digunakan:**
```
database/complete_schema.sql
```

### **Kapan Digunakan:**
- ✅ Setup Supabase project baru dari nol
- ✅ Clone database ke environment baru (dev, staging, production)
- ✅ Disaster recovery (restore database)
- ✅ Testing di local PostgreSQL
- ✅ Onboarding developer baru

### **Apa yang Dilakukan:**
File ini berisi **SELURUH schema database** dari awal sampai Session 4:

**Tables:**
- ✅ `products` (dengan harga_beli, created_by)
- ✅ `categories`
- ✅ `users`
- ✅ `transactions` (dengan kasir_id)
- ✅ `transaction_details` (dengan harga_beli)

**Indexes:**
- ✅ 10+ indexes untuk performa optimal

**Foreign Keys:**
- ✅ 5 foreign key constraints

**Triggers:**
- ✅ Auto-update `updated_at` untuk products, categories, users

**Sample Data:**
- ✅ 4 categories
- ✅ 5 sample products
- ✅ 3 default users (admin, kasir1, kasir2)

### **Langkah-langkah:**
1. **Buat Supabase Project Baru**
   - Login ke https://supabase.com
   - Klik "New Project"
   - Isi nama project, password, region
   - Tunggu project dibuat (2-3 menit)

2. **Run Complete Schema**
   - Klik "SQL Editor"
   - Copy-paste **SELURUH ISI** `complete_schema.sql`
   - Klik "Run"
   - Tunggu 10-15 detik

3. **Verify Setup**
   - Scroll ke bawah, ada verification queries
   - Jalankan untuk cek semua table sudah ada
   - Cek users, products, categories

4. **Update Connection String**
   - Copy DATABASE_URL dari Supabase Settings
   - Update di `.env` file Anda

5. **Test Connection**
   - Run aplikasi: `go run main.go`
   - Test login dengan default credentials

---

## 📊 Perbandingan File Migration

| File | Ukuran | Untuk | Isi |
|------|--------|-------|-----|
| `migration_session_4_complete.sql` | ~7 KB | Update existing DB | Hanya perubahan Session 4 |
| `complete_schema.sql` | ~12 KB | Setup baru dari nol | SEMUA schema + data |

---

## 🔄 Use Cases Detail

### **Use Case 1: Development → Staging → Production**

**Development (Local):**
```bash
# Setup database lokal
psql -U postgres -d kasir_dev -f database/complete_schema.sql
```

**Staging (Supabase):**
1. Buat project Supabase baru untuk staging
2. Run `complete_schema.sql`
3. Update `.env.staging` dengan DATABASE_URL baru

**Production (Supabase):**
1. Buat project Supabase baru untuk production
2. Run `complete_schema.sql`
3. Update `.env.production` dengan DATABASE_URL baru
4. **GANTI PASSWORD DEFAULT!**

### **Use Case 2: Disaster Recovery**

**Backup:**
```bash
# Backup data dari Supabase
# (gunakan Supabase Dashboard → Database → Backups)
```

**Restore:**
1. Buat project Supabase baru
2. Run `complete_schema.sql` (schema)
3. Restore data dari backup (data)

### **Use Case 3: Clone untuk Testing**

**Scenario:** Anda ingin test fitur baru tanpa ganggu production

1. Buat project Supabase baru "kasir-api-test"
2. Run `complete_schema.sql`
3. (Optional) Import data production untuk testing realistis
4. Test dengan aman tanpa ganggu production

### **Use Case 4: Onboarding Developer Baru**

**Scenario:** Developer baru join tim

1. Clone repository
2. Buat Supabase project sendiri
3. Run `complete_schema.sql`
4. Update `.env` dengan DATABASE_URL sendiri
5. Langsung bisa development tanpa setup manual

---

## 🎯 Rekomendasi Best Practice

### **Untuk Update Existing Database:**
✅ Gunakan `migration_session_4_complete.sql`  
✅ Backup database dulu sebelum migration  
✅ Test di staging dulu sebelum production  

### **Untuk Setup Baru:**
✅ Gunakan `complete_schema.sql`  
✅ Langsung dapat semua schema + sample data  
✅ Ganti password default setelah setup  

### **Untuk Version Control:**
✅ Simpan semua migration files di git  
✅ Buat migration file baru untuk setiap perubahan  
✅ Update `complete_schema.sql` setiap ada perubahan major  

---

## 📁 File Structure

```
database/
├── complete_schema.sql                    ← MASTER: Setup baru dari nol
├── migration_session_4_complete.sql       ← Update existing DB
├── README_MIGRATION.md                    ← Panduan migration
├── README_DEPLOYMENT.md                   ← Panduan ini
│
├── (Optional - untuk referensi)
├── supabase_setup.sql                     ← Setup awal (Session 1-3)
├── migration_create_users.sql             ← Partial migration
├── migration_add_harga_beli_to_products.sql
├── migration_add_kasir_to_transactions.sql
└── seed_default_users.sql
```

---

## ✅ Verification Checklist

Setelah run migration, verify dengan checklist ini:

### **Tables Created:**
- [ ] `products` (dengan kolom harga_beli, created_by)
- [ ] `categories`
- [ ] `users`
- [ ] `transactions` (dengan kolom kasir_id)
- [ ] `transaction_details` (dengan kolom harga_beli)

### **Indexes Created:**
- [ ] idx_products_nama (UNIQUE)
- [ ] idx_products_category_id
- [ ] idx_products_created_by
- [ ] idx_users_username
- [ ] idx_transactions_kasir_id
- [ ] (dan lainnya)

### **Foreign Keys Created:**
- [ ] products → categories
- [ ] products → users (created_by)
- [ ] transactions → users (kasir_id)
- [ ] transaction_details → transactions
- [ ] transaction_details → products

### **Sample Data Inserted:**
- [ ] 3 users (admin, kasir1, kasir2)
- [ ] 4 categories
- [ ] 5 sample products

### **Triggers Working:**
- [ ] Auto-update updated_at di products
- [ ] Auto-update updated_at di categories
- [ ] Auto-update updated_at di users

---

## 🐛 Troubleshooting

### **Error: "relation already exists"**
**Penyebab:** Table sudah ada dari setup sebelumnya.

**Solusi:**
- Jika ingin fresh install, drop tables dulu
- Atau gunakan `migration_session_4_complete.sql` untuk update saja

### **Error: "duplicate key value"**
**Penyebab:** Sample data sudah ada.

**Solusi:** Skip error ini, tidak masalah.

### **Error: "foreign key constraint"**
**Penyebab:** Urutan table creation salah.

**Solusi:** Gunakan `complete_schema.sql` yang sudah benar urutannya.

---

## 📞 Need Help?

- **Migration Issues:** Lihat `README_MIGRATION.md`
- **Setup Issues:** Lihat `QUICK_START_AUTH.md`
- **Full Documentation:** Lihat `AUTH_IMPLEMENTATION.md`

---

**Created:** 2026-02-12  
**Version:** Session 4  
**Status:** ✅ Production Ready
