# 🔐 Implementasi Authentication & Authorization

## 📋 Ringkasan Perubahan

Implementasi ini menambahkan sistem **Role-Based Access Control (RBAC)** dengan 2 role:
- **Admin**: Full access (CRUD produk, lihat profit, kelola user)
- **Kasir**: Limited access (lihat produk, proses transaksi, lihat laporan sendiri)

---

## 🗄️ Database Changes

### 1. Migration Files Created

#### `migration_create_users.sql`
- Membuat tabel `users` dengan kolom: id, username, password, nama_lengkap, role, is_active
- Menambahkan indexes untuk performa
- Trigger auto-update `updated_at`

#### `migration_add_harga_beli_to_products.sql`
- Menambahkan kolom `harga_beli` ke tabel `products` (nullable)
- Menambahkan kolom `created_by` untuk tracking
- Constraint validation (harga tidak boleh negatif)

#### `migration_add_kasir_to_transactions.sql`
- Menambahkan kolom `kasir_id` ke tabel `transactions`
- Menambahkan kolom `harga_beli` ke `transaction_details` (snapshot)

#### `seed_default_users.sql`
- Insert default users:
  - `admin` / `admin123` (role: admin)
  - `kasir1` / `kasir123` (role: kasir)
  - `kasir2` / `kasir123` (role: kasir)

### 2. Cara Menjalankan Migration

```bash
# 1. Login ke Supabase Dashboard
# 2. Buka SQL Editor
# 3. Jalankan migration satu per satu:

# Step 1: Create users table
# Copy-paste isi migration_create_users.sql

# Step 2: Add harga_beli to products
# Copy-paste isi migration_add_harga_beli_to_products.sql

# Step 3: Add kasir_id to transactions
# Copy-paste isi migration_add_kasir_to_transactions.sql

# Step 4: Seed default users
# Copy-paste isi seed_default_users.sql
```

---

## 🏗️ Code Structure

### New Files Created

```
kasir-api/
├── models/
│   ├── user.go              ✅ User model & login request/response
│   └── errors.go            ✅ Centralized error definitions
├── utils/
│   ├── password.go          ✅ Bcrypt password hashing
│   └── jwt.go               ✅ JWT token generation & validation
├── middleware/
│   ├── auth.go              ✅ JWT authentication middleware
│   ├── role.go              ✅ Role-based access control
│   ├── cors.go              ✅ CORS for frontend integration
│   └── logging.go           ✅ Request logging with slog
├── repositories/
│   └── user_repository.go   ✅ User database operations
├── services/
│   └── auth_service.go      ✅ Authentication business logic
└── handlers/
    └── auth_handler.go      ✅ Login & register endpoints
```

### Modified Files

```
├── models/
│   └── product.go           🔄 Added harga_beli, margin calculation
├── .env.example             🔄 Added JWT_SECRET, JWT_EXPIRE_HOURS
```

---

## 🔐 Authentication Flow

### 1. Login Flow

```
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
      "role": "admin",
      "is_active": true
    }
  }
}
```

### 2. Using JWT Token

Setiap request ke endpoint yang protected harus menyertakan token:

```
GET /api/produk
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🎭 Role-Based Access

### Admin Access

**Endpoints:**
- ✅ `GET /api/produk` - Lihat semua produk (dengan harga_beli & margin)
- ✅ `POST /api/produk` - Tambah produk baru
- ✅ `PUT /api/produk/:id` - Edit produk
- ✅ `DELETE /api/produk/:id` - Hapus produk
- ✅ `GET /api/reports/sales` - Laporan penjualan (dengan profit detail)
- ✅ `POST /api/auth/register` - Tambah user baru (kasir)

**Response Example (GET /api/produk):**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3000,
  "harga_beli": 2500,      // ← Admin bisa lihat
  "stok": 100,
  "margin": 16.67,         // ← Auto calculated
  "category_id": 1
}
```

### Kasir Access

**Endpoints:**
- ✅ `GET /api/produk` - Lihat produk (TANPA harga_beli)
- ✅ `POST /api/transaksi` - Proses penjualan
- ✅ `GET /api/reports/my-sales` - Laporan penjualan sendiri (tanpa profit)
- ❌ `POST /api/produk` - FORBIDDEN
- ❌ `PUT /api/produk/:id` - FORBIDDEN
- ❌ `DELETE /api/produk/:id` - FORBIDDEN

**Response Example (GET /api/produk):**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3000,           // Hanya harga jual
  "stok": 100,
  "category_id": 1
  // harga_beli TIDAK ditampilkan
}
```

---

## 🔧 Environment Variables

Update file `.env` Anda dengan:

```bash
# Security Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=8
API_KEY=your-api-key-here
```

**Generate secure JWT_SECRET:**
```bash
# Linux/Mac
openssl rand -base64 32

# Windows PowerShell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

---

## 📝 Next Steps

### TODO: Update main.go

Anda perlu update `main.go` untuk:
1. Initialize user repository & auth service
2. Setup auth handler
3. Apply middleware ke routes
4. Setup role-based routing

### TODO: Update Product Handler

Product handler perlu diupdate untuk:
1. Filter response berdasarkan role (hide harga_beli untuk kasir)
2. Validate role untuk create/update/delete
3. Track created_by saat create produk

### TODO: Update Transaction Handler

Transaction handler perlu:
1. Capture kasir_id dari JWT token
2. Snapshot harga_beli saat transaksi
3. Calculate profit untuk laporan admin

### TODO: Update Report Service

Report service perlu:
1. Endpoint `/api/reports/my-sales` untuk kasir
2. Endpoint `/api/reports/sales` dengan profit detail untuk admin
3. Filter by kasir_id

---

## 🧪 Testing Checklist

- [ ] Run all migrations di Supabase
- [ ] Verify users table created
- [ ] Verify default users inserted
- [ ] Test login dengan admin
- [ ] Test login dengan kasir
- [ ] Test JWT token validation
- [ ] Test role-based access (admin vs kasir)
- [ ] Update Postman collection dengan login endpoint
- [ ] Test product CRUD dengan role
- [ ] Test transaction dengan kasir_id
- [ ] Test reports dengan profit calculation

---

## 📚 Documentation

### How to Add New User (via Database)

```sql
-- Hash password dulu di aplikasi atau gunakan:
-- Password: newpassword123
-- Hash: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

INSERT INTO users (username, password, nama_lengkap, role) VALUES
('kasir3', '$2a$10$...hash...', 'Kasir Malam', 'kasir');
```

### How to Deactivate User

```sql
UPDATE users SET is_active = FALSE WHERE username = 'kasir1';
```

### How to Change Password

```sql
-- Hash password baru dulu, lalu:
UPDATE users SET password = '$2a$10$...new_hash...' WHERE username = 'admin';
```

---

## 🔒 Security Best Practices

1. ✅ Password di-hash dengan bcrypt (cost 10)
2. ✅ JWT token dengan expiration
3. ✅ Sensitive data (password) tidak di-return di JSON
4. ✅ Role validation di middleware
5. ✅ CORS configuration untuk frontend
6. ✅ Request logging untuk audit trail
7. ⚠️ **IMPORTANT**: Ganti JWT_SECRET di production!
8. ⚠️ **IMPORTANT**: Ganti default passwords!

---

## 📞 Support

Jika ada pertanyaan atau issue:
1. Cek error logs (slog akan print ke console)
2. Verify JWT_SECRET sudah di-set di .env
3. Verify migrations sudah dijalankan
4. Test dengan Postman collection yang sudah diupdate

---

**Status**: ✅ Database migrations ready, ✅ Core authentication ready, 🔄 Integration pending
