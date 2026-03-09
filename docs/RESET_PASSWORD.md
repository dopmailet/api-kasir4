# 🔑 Panduan Reset Password - KasirPOS API

Dokumen ini menjelaskan cara reset password user jika login gagal.

## Kapan Perlu Reset Password?

- Login gagal dengan error `"username atau password salah"`
- `JWT_SECRET` di `.env` baru diganti (token lama tidak valid → **cukup login ulang**)
- Lupa password yang tersimpan di database

> **Penting:** Mengganti `JWT_SECRET` di `.env` **tidak** mengubah password di database.
> Token lama yang sudah digenerate memang tidak valid, tapi login ulang akan langsung berhasil.
> Reset password hanya diperlukan jika hash di database memang salah/tidak cocok.

---

## Langkah Reset Password

### 1. Generate Hash Password Baru

Jalankan tool berikut dari folder `tools/`:

```bash
cd tools
go run reset_admin_password.go
```

Output akan menampilkan hash bcrypt dan SQL query siap pakai.

### 2. Jalankan SQL di Supabase

Buka [supabase.com](https://supabase.com) → pilih project → **SQL Editor**

**Cek dulu user yang ada:**
```sql
SELECT id, username, role, is_active FROM users
WHERE username IN ('admin', 'kasir1');
```

**Jika user SUDAH ADA → gunakan UPDATE:**
```sql
UPDATE users
SET password = '<hash_dari_output_tool>', is_active = true
WHERE username = 'admin';
```

**Jika user BELUM ADA → gunakan INSERT:**
```sql
INSERT INTO users (username, password, nama_lengkap, role, is_active)
VALUES ('admin', '<hash_dari_output_tool>', 'Administrator', 'admin', true);
```

> ⚠️ **Jangan gunakan `ON CONFLICT (username)`** kecuali kolom username sudah punya UNIQUE constraint.
> Gunakan SELECT dulu untuk cek, lalu pilih UPDATE atau INSERT.

---

## Default Credentials

| Username | Password  | Role  |
|----------|-----------|-------|
| `admin`  | `admin123`| admin |
| `kasir1` | `kasir123`| kasir |

---

## Mengganti Password Default

Edit bagian `users` di file `tools/reset_admin_password.go`:

```go
users := []UserReset{
    {Username: "admin", Password: "password_baru_anda", NamaLengkap: "Administrator", Role: "admin"},
    {Username: "kasir1", Password: "password_baru_kasir", NamaLengkap: "Kasir Utama", Role: "kasir"},
}
```

Lalu jalankan ulang dan gunakan SQL output-nya.

---

## Catatan Teknis

- Password di-hash menggunakan **bcrypt** dengan cost factor 10
- Hash bersifat **one-way** — tidak bisa di-decode, hanya bisa diverifikasi
- Setiap kali tool dijalankan, hash yang dihasilkan **berbeda** (salt acak) tapi tetap valid
- JWT token menggunakan secret dari `.env` → ganti `JWT_SECRET` = semua token lama invalid
