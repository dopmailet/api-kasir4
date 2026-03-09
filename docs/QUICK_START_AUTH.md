# 🚀 Quick Start Guide - Authentication System

## ⚡ 5-Minute Setup

### Step 1: Run Database Migration (Supabase) - 1 File Saja! ✅

1. Login ke https://supabase.com
2. Buka project Anda
3. Klik "SQL Editor"
4. Buka file `database/migration_session_4_complete.sql`
5. Copy-paste **SELURUH ISI FILE** ke SQL Editor
6. Klik "Run" atau tekan Ctrl+Enter
7. Tunggu sampai selesai (sekitar 5-10 detik)
8. ✅ Selesai! Migration berhasil jika tidak ada error

**File yang digunakan:**
- ✅ `database/migration_session_4_complete.sql` (ALL-IN-ONE)

**Apa yang dilakukan migration ini:**
- ✅ Membuat tabel `users`
- ✅ Menambahkan `harga_beli` ke tabel `products`
- ✅ Menambahkan `kasir_id` ke tabel `transactions`
- ✅ Insert 3 default users (admin, kasir1, kasir2)

### Step 2: Update .env File

```bash
# Tambahkan ke .env Anda:
JWT_SECRET=ganti-dengan-random-string-yang-panjang-dan-aman
JWT_EXPIRE_HOURS=8
```

**Generate secure JWT_SECRET:**
```powershell
# Windows PowerShell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

### Step 3: Build & Run

```bash
go build -o kasir-api.exe
./kasir-api.exe
```

---

## 🔑 Default Credentials

| Username | Password | Role | Nama Lengkap |
|----------|----------|------|--------------|
| `admin` | `admin123` | admin | Administrator |
| `kasir1` | `kasir123` | kasir | Kasir Pagi |
| `kasir2` | `kasir123` | kasir | Kasir Sore |

⚠️ **IMPORTANT**: Ganti password ini di production!

---

## 📡 API Endpoints

### 1. Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "nama_lengkap": "Administrator",
      "role": "admin"
    }
  }
}
```

### 2. Access Protected Endpoint

```http
GET /api/produk
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🎭 Role Differences

### Admin Response (GET /api/produk)
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3000,
  "harga_beli": 2500,      // ← Admin bisa lihat
  "stok": 100,
  "margin": 16.67          // ← Auto calculated
}
```

### Kasir Response (GET /api/produk)
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3000,           // Hanya harga jual
  "stok": 100
  // harga_beli TIDAK ditampilkan
}
```

---

## 🛠️ Common Tasks

### Add New Kasir (via SQL)

```sql
-- Password: kasir123 (already hashed)
INSERT INTO users (username, password, nama_lengkap, role) VALUES
('kasir3', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Kasir Malam', 'kasir');
```

### Deactivate User

```sql
UPDATE users SET is_active = FALSE WHERE username = 'kasir1';
```

### Activate User

```sql
UPDATE users SET is_active = TRUE WHERE username = 'kasir1';
```

### List All Users

```sql
SELECT id, username, nama_lengkap, role, is_active FROM users;
```

---

## 🧪 Testing Checklist

- [ ] Run all 4 migrations
- [ ] Verify 3 users created (admin, kasir1, kasir2)
- [ ] Test login dengan admin → get token
- [ ] Test login dengan kasir → get token
- [ ] Test GET /api/produk dengan admin token → lihat harga_beli
- [ ] Test GET /api/produk dengan kasir token → tidak lihat harga_beli
- [ ] Test POST /api/produk dengan admin token → success
- [ ] Test POST /api/produk dengan kasir token → 403 Forbidden
- [ ] Test request tanpa token → 401 Unauthorized

---

## 🐛 Troubleshooting

### Error: "JWT_SECRET tidak ditemukan"
**Solution**: Tambahkan `JWT_SECRET=...` ke file `.env`

### Error: "Invalid or expired token"
**Solution**: Login ulang untuk mendapatkan token baru

### Error: "Forbidden: Insufficient permissions"
**Solution**: Endpoint ini hanya untuk admin, login dengan akun admin

### Error: "User tidak aktif"
**Solution**: Aktifkan user dengan SQL: `UPDATE users SET is_active = TRUE WHERE username = '...'`

### Build Error: "could not import golang.org/x/crypto/bcrypt"
**Solution**: Run `go get golang.org/x/crypto/bcrypt`

---

## 📚 Documentation

- **Full Implementation Guide**: `docs/AUTH_IMPLEMENTATION.md`
- **Complete Summary**: `docs/SESSION_4_SUMMARY.md`
- **Migration Files**: `database/migration_*.sql`
- **Seed Data**: `database/seed_default_users.sql`

---

## 🎯 Next Steps

1. ✅ Database migrations done
2. ✅ Core authentication done
3. 🔄 **TODO**: Update `main.go` untuk integrate auth
4. 🔄 **TODO**: Update handlers untuk role-based filtering
5. 🔄 **TODO**: Update Postman collection

---

**Need Help?** Check the full documentation in `docs/` folder!
