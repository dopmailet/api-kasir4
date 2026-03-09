# 🚀 Quick Reference - Database Files

## 📁 File Mana yang Harus Digunakan?

### ✅ **Untuk Setup Supabase BARU dari NOL**
```
File: database/complete_schema.sql
```
**Kapan:** Clone database, disaster recovery, setup dev baru  
**Isi:** SEMUA schema + sample data (lengkap!)  
**Waktu:** ~10-15 detik  

---

### ✅ **Untuk UPDATE Database yang SUDAH ADA**
```
File: database/migration_session_4_complete.sql
```
**Kapan:** Upgrade dari Session 3 ke Session 4  
**Isi:** Hanya perubahan Session 4 (auth & profit)  
**Waktu:** ~5-10 detik  

---

## 🎯 Decision Tree

```
Apakah database Anda sudah ada?
│
├─ TIDAK (database baru/kosong)
│  └─ Gunakan: complete_schema.sql ✅
│
└─ YA (sudah ada products, categories, transactions)
   └─ Gunakan: migration_session_4_complete.sql ✅
```

---

## 📊 Comparison Table

| Kriteria | complete_schema.sql | migration_session_4_complete.sql |
|----------|---------------------|----------------------------------|
| **Untuk** | Setup baru | Update existing |
| **Tables** | Semua (5 tables) | Hanya update existing |
| **Sample Data** | ✅ Ya | ✅ Ya (users only) |
| **Ukuran** | ~12 KB | ~7 KB |
| **Waktu** | 10-15 detik | 5-10 detik |
| **Use Case** | Clone, DR, New Dev | Upgrade Session 3→4 |

---

## 🔧 Quick Commands

### **Setup Baru (Complete Schema)**
```bash
# 1. Login ke Supabase
# 2. SQL Editor
# 3. Copy-paste complete_schema.sql
# 4. Run
# 5. Done! ✅
```

### **Update Existing (Migration)**
```bash
# 1. Login ke Supabase
# 2. SQL Editor
# 3. Copy-paste migration_session_4_complete.sql
# 4. Run
# 5. Done! ✅
```

---

## 📝 Default Credentials (Both Files)

| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | admin |
| kasir1 | kasir123 | kasir |
| kasir2 | kasir123 | kasir |

⚠️ **Ganti di production!**

---

## ✅ Verification

Run query ini setelah migration:

```sql
-- Cek tables
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' ORDER BY table_name;

-- Expected: categories, products, transaction_details, transactions, users

-- Cek users
SELECT username, role FROM users;

-- Expected: admin, kasir1, kasir2
```

---

**Need More Info?**
- Full guide: `README_DEPLOYMENT.md`
- Migration guide: `README_MIGRATION.md`
