# 📚 API CRUD Guide - Kasir API

Panduan lengkap untuk melakukan operasi **CRUD (Create, Read, Update, Delete)** pada Kasir API.

---

## 📋 Table of Contents

- [Base URL](#base-url)
- [Products CRUD](#products-crud)
- [Categories CRUD](#categories-crud)
- [Relasi Product-Category](#relasi-product-category)
- [Error Handling](#error-handling)
- [Tips & Best Practices](#tips--best-practices)

---

## 🌐 Base URL

### Local Development
```
http://localhost:8080
```

### Production (Railway)
```
https://your-app.railway.app
```

---

## 📦 Products CRUD

### 1️⃣ CREATE - Tambah Product Baru

#### **Endpoint:**
```http
POST /api/produk
```

#### **Request Body:**
```json
{
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": 1
}
```

#### **Response (201 Created):**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": 1
}
```

#### **Contoh dengan cURL:**
```bash
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Indomie Goreng",
    "harga": 3500,
    "stok": 100,
    "category_id": 1
  }'
```

#### **Contoh Tanpa Category:**
```json
{
  "nama": "Produk Tanpa Kategori",
  "harga": 5000,
  "stok": 50
}
```

#### **Catatan:**
- ✅ `id` dibuat otomatis oleh database
- ✅ `category_id` bersifat **opsional** (boleh NULL)
- ✅ `nama` harus **unik** (tidak boleh duplikat)
- ⚠️ Jika nama sudah ada, stok akan **ditambahkan** (UPSERT logic)

---

### 2️⃣ READ - Lihat Semua Products

#### **Endpoint:**
```http
GET /api/produk
```

#### **Response (200 OK):**
```json
[
  {
    "id": 1,
    "nama": "Indomie Goreng",
    "harga": 3500,
    "stok": 100,
    "category_id": 1,
    "category": {
      "id": 1,
      "nama": "Makanan",
      "deskription": "Produk makanan dan snack"
    }
  },
  {
    "id": 2,
    "nama": "Aqua 600ml",
    "harga": 3000,
    "stok": 50,
    "category_id": 2,
    "category": {
      "id": 2,
      "nama": "Minuman",
      "deskription": "Produk minuman kemasan"
    }
  }
]
```

#### **Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/produk
```

#### **Catatan:**
- ✅ Return semua products dengan **category name** (via JOIN)
- ✅ Product tanpa category tetap muncul (LEFT JOIN)

---

### 3️⃣ READ - Lihat 1 Product by ID

#### **Endpoint:**
```http
GET /api/produk/{id}
```

#### **Contoh:**
```http
GET /api/produk/1
```

#### **Response (200 OK):**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": 1,
  "category": {
    "id": 1,
    "nama": "Makanan",
    "deskription": "Produk makanan dan snack"
  }
}
```

#### **Response (404 Not Found):**
```json
{
  "error": "Product tidak ditemukan"
}
```

#### **Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/produk/1
```

---

### 4️⃣ UPDATE - Update Product

#### **Endpoint:**
```http
PUT /api/produk/{id}
```

#### **Request Body:**
```json
{
  "nama": "Indomie Goreng Special",
  "harga": 4000,
  "stok": 90,
  "category_id": 1
}
```

#### **Response (200 OK):**
```json
{
  "id": 1,
  "nama": "Indomie Goreng Special",
  "harga": 4000,
  "stok": 90,
  "category_id": 1
}
```

#### **Contoh dengan cURL:**
```bash
curl -X PUT http://localhost:8080/api/produk/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Indomie Goreng Special",
    "harga": 4000,
    "stok": 90,
    "category_id": 1
  }'
```

#### **Update - Hapus Category (Set NULL):**
```json
{
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": null
}
```

#### **Catatan:**
- ✅ Semua field harus dikirim (nama, harga, stok, category_id)
- ✅ `category_id` bisa di-set `null` untuk hapus category

---

### 5️⃣ DELETE - Hapus Product

#### **Endpoint:**
```http
DELETE /api/produk/{id}
```

#### **Contoh:**
```http
DELETE /api/produk/1
```

#### **Response (200 OK):**
```json
{
  "message": "sukses delete"
}
```

#### **Response (404 Not Found):**
```json
{
  "error": "Product tidak ditemukan"
}
```

#### **Contoh dengan cURL:**
```bash
curl -X DELETE http://localhost:8080/api/produk/1
```

#### **Catatan:**
- ⚠️ **Permanent delete** - data tidak bisa dikembalikan
- ✅ ID yang dihapus tidak akan dipakai lagi (ada gap)

---

## 🏷️ Categories CRUD

### 1️⃣ CREATE - Tambah Category Baru

#### **Endpoint:**
```http
POST /api/categories
```

#### **Request Body:**
```json
{
  "nama": "Elektronik",
  "deskription": "Produk elektronik dan gadget"
}
```

#### **Response (201 Created):**
```json
{
  "id": 5,
  "nama": "Elektronik",
  "deskription": "Produk elektronik dan gadget"
}
```

#### **Contoh dengan cURL:**
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Elektronik",
    "deskription": "Produk elektronik dan gadget"
  }'
```

#### **Catatan:**
- ✅ `id` dibuat otomatis oleh database
- ✅ `deskription` bersifat opsional (boleh kosong)

---

### 2️⃣ READ - Lihat Semua Categories

#### **Endpoint:**
```http
GET /api/categories
```

#### **Response (200 OK):**
```json
[
  {
    "id": 1,
    "nama": "Makanan",
    "deskription": "Produk makanan dan snack"
  },
  {
    "id": 2,
    "nama": "Minuman",
    "deskription": "Produk minuman kemasan"
  },
  {
    "id": 3,
    "nama": "Sembako",
    "deskription": "Kebutuhan pokok sehari-hari"
  }
]
```

#### **Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/categories
```

#### **Catatan:**
- ✅ Return semua categories **tanpa products** (untuk performa)

---

### 3️⃣ READ - Lihat 1 Category by ID (dengan Products)

#### **Endpoint:**
```http
GET /api/categories/{id}
```

#### **Contoh:**
```http
GET /api/categories/1
```

#### **Response (200 OK):**
```json
{
  "id": 1,
  "nama": "Makanan",
  "deskription": "Produk makanan dan snack",
  "products": [
    {
      "id": 1,
      "nama": "Indomie Goreng",
      "harga": 3500,
      "stok": 100
    },
    {
      "id": 5,
      "nama": "Mie Sedaap",
      "harga": 3500,
      "stok": 80
    }
  ]
}
```

#### **Response (404 Not Found):**
```json
{
  "error": "Category tidak ditemukan"
}
```

#### **Contoh dengan cURL:**
```bash
curl http://localhost:8080/api/categories/1
```

#### **Catatan:**
- ✅ Return category **dengan semua products** yang termasuk dalam category tersebut
- ✅ Jika category tidak punya products, `products` akan array kosong `[]`

---

### 4️⃣ UPDATE - Update Category

#### **Endpoint:**
```http
PUT /api/categories/{id}
```

#### **Request Body:**
```json
{
  "nama": "Makanan & Minuman",
  "deskription": "Produk makanan, minuman, dan snack"
}
```

#### **Response (200 OK):**
```json
{
  "id": 1,
  "nama": "Makanan & Minuman",
  "deskription": "Produk makanan, minuman, dan snack"
}
```

#### **Contoh dengan cURL:**
```bash
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Makanan & Minuman",
    "deskription": "Produk makanan, minuman, dan snack"
  }'
```

#### **Catatan:**
- ✅ Semua field harus dikirim (nama, deskription)
- ✅ Update category **tidak mempengaruhi** products yang sudah ada

---

### 5️⃣ DELETE - Hapus Category

#### **Endpoint:**
```http
DELETE /api/categories/{id}
```

#### **Contoh:**
```http
DELETE /api/categories/1
```

#### **Response (200 OK):**
```json
{
  "message": "sukses delete"
}
```

#### **Response (404 Not Found):**
```json
{
  "error": "Category tidak ditemukan"
}
```

#### **Contoh dengan cURL:**
```bash
curl -X DELETE http://localhost:8080/api/categories/1
```

#### **Catatan:**
- ⚠️ **Permanent delete** - data tidak bisa dikembalikan
- ✅ Products yang punya `category_id` ini akan di-set **NULL** (ON DELETE SET NULL)
- ✅ Products **tidak ikut terhapus**, hanya category_id-nya yang jadi NULL

---

## 🔗 Relasi Product-Category

### Workflow: Tambah Category → Tambah Product

#### **Step 1: Tambah Category**
```http
POST /api/categories
Content-Type: application/json

{
  "nama": "Elektronik",
  "deskription": "Produk elektronik dan gadget"
}
```

**Response:**
```json
{
  "id": 5,  // ← Catat ID ini!
  "nama": "Elektronik",
  "deskription": "Produk elektronik dan gadget"
}
```

---

#### **Step 2: Tambah Product dengan Category ID**
```http
POST /api/produk
Content-Type: application/json

{
  "nama": "Powerbank 10000mAh",
  "harga": 150000,
  "stok": 25,
  "category_id": 5  // ← Gunakan ID dari Step 1
}
```

**Response:**
```json
{
  "id": 10,
  "nama": "Powerbank 10000mAh",
  "harga": 150000,
  "stok": 25,
  "category_id": 5
}
```

---

#### **Step 3: Verify - Lihat Category dengan Products**
```http
GET /api/categories/5
```

**Response:**
```json
{
  "id": 5,
  "nama": "Elektronik",
  "deskription": "Produk elektronik dan gadget",
  "products": [
    {
      "id": 10,
      "nama": "Powerbank 10000mAh",
      "harga": 150000,
      "stok": 25
    }
  ]
}
```

---

### Workflow: Pindah Product ke Category Lain

#### **Step 1: Lihat Product Sekarang**
```http
GET /api/produk/1
```

**Response:**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": 1,
  "category": {
    "nama": "Makanan"
  }
}
```

---

#### **Step 2: Update Category Product**
```http
PUT /api/produk/1
Content-Type: application/json

{
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": 3  // ← Pindah ke Sembako
}
```

---

#### **Step 3: Verify - Lihat Product Setelah Update**
```http
GET /api/produk/1
```

**Response:**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100,
  "category_id": 3,
  "category": {
    "nama": "Sembako"  // ← Category berubah!
  }
}
```

---

## ❌ Error Handling

### Common Errors

#### **1. Product/Category Not Found (404)**
```json
{
  "error": "Product tidak ditemukan"
}
```

**Penyebab:** ID yang dicari tidak ada di database

---

#### **2. Invalid Request Body (400)**
```json
{
  "error": "Invalid request body"
}
```

**Penyebab:** JSON format salah atau field required tidak ada

---

#### **3. Foreign Key Violation (500)**
```json
{
  "error": "violates foreign key constraint"
}
```

**Penyebab:** `category_id` yang diberikan tidak ada di table categories

**Solusi:** Cek ID categories yang tersedia:
```http
GET /api/categories
```

---

#### **4. Duplicate Product Name (Upsert)**
Jika product dengan nama yang sama sudah ada, stok akan **ditambahkan**:

**Existing Product:**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "stok": 100
}
```

**POST dengan nama sama:**
```http
POST /api/produk
{
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 50
}
```

**Result:**
```json
{
  "id": 1,
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 150  // ← 100 + 50 = 150
}
```

---

## 💡 Tips & Best Practices

### 1. **Gunakan Postman Collection**
Import file `Kasir-API.postman_collection.json` untuk testing yang lebih mudah.

### 2. **Cek Categories Terlebih Dahulu**
Sebelum tambah product, cek dulu ID categories yang tersedia:
```http
GET /api/categories
```

### 3. **Gunakan Environment Variables**
Di Postman, buat variable `{{base_url}}`:
- Local: `http://localhost:8080`
- Production: `https://your-app.railway.app`

### 4. **Backup Data Sebelum Delete**
DELETE bersifat permanent. Backup data penting sebelum hapus.

### 5. **Validasi Input di Client**
- Harga dan stok harus **angka positif**
- Nama product sebaiknya **tidak kosong**
- Category ID harus **valid** (ada di database)

### 6. **Handle NULL Category**
Product tanpa category (category_id = NULL) tetap valid dan bisa ditampilkan.

### 7. **Monitor Auto-Increment ID**
ID tidak akan di-recycle setelah delete. Jika delete ID 5, ID berikutnya adalah 6 (bukan 5).

---

## 📊 Quick Reference Table

| Operation | Products | Categories |
|-----------|----------|------------|
| **Create** | `POST /api/produk` | `POST /api/categories` |
| **Read All** | `GET /api/produk` | `GET /api/categories` |
| **Read One** | `GET /api/produk/{id}` | `GET /api/categories/{id}` |
| **Update** | `PUT /api/produk/{id}` | `PUT /api/categories/{id}` |
| **Delete** | `DELETE /api/produk/{id}` | `DELETE /api/categories/{id}` |

---

## 🎯 Testing Checklist

### Products
- [ ] POST product dengan category
- [ ] POST product tanpa category
- [ ] GET all products (lihat category name muncul)
- [ ] GET product by ID (lihat category detail)
- [ ] PUT update harga dan stok
- [ ] PUT pindah category
- [ ] PUT set category_id = null
- [ ] DELETE product

### Categories
- [ ] POST category baru
- [ ] GET all categories
- [ ] GET category by ID (lihat products muncul)
- [ ] PUT update nama dan deskripsi
- [ ] DELETE category (cek products jadi NULL)

### Relasi
- [ ] Tambah category → Tambah product dengan category_id
- [ ] Lihat category → Products muncul di response
- [ ] Lihat product → Category name muncul di response
- [ ] Delete category → Products tetap ada (category_id jadi NULL)

---

## 🚀 Ready to Use!

API Anda sudah siap digunakan! Silakan test semua endpoint dan selamat coding! 🎉

**Need Help?**
- Cek file `README.md` untuk setup dan deployment
- Cek file `CHALLENGE_SESSION_2.md` untuk penjelasan teknis relasi
- Lihat `database/supabase_setup.sql` untuk database schema

---

**Last Updated:** 2026-02-01  
**Version:** 1.0.0
