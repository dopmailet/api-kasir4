# 📋 UPDATE BACKEND — Instruksi untuk Frontend

**Tanggal:** 16 Februari 2026  
**Dari:** AI Backend (Go API)  
**Untuk:** AI Frontend  

---

## 🔴 PERUBAHAN PENTING (Breaking Changes)

### 1. POST /api/produk → DIHAPUS
Endpoint `POST /api/produk` untuk membuat produk baru sudah **dihapus**. 
Produk baru sekarang dibuat lewat modul **Pembelian** (`POST /api/purchases`).

**Aksi frontend:** Hapus tombol "Tambah Produk" atau redirect ke halaman Pembelian.

### 2. PUT /api/produk/{id} → DIBATASI
Sekarang hanya bisa mengubah: **nama**, **harga jual (harga)**, dan **category_id**.
**Tidak bisa** mengubah stok atau harga_beli lewat endpoint ini.

```json
// Request PUT /api/produk/{id}
{
  "nama": "Nama Produk Baru",
  "harga": 50000,
  "category_id": 2
}
// ❌ "stok" dan "harga_beli" diabaikan jika dikirim
```

**Aksi frontend:** Di form edit produk, hilangkan input "Stok" dan "Harga Beli". Hanya tampilkan input Nama, Harga Jual, dan Kategori.

---

## 🆕 ENDPOINT BARU: Modul Pembelian

### POST /api/purchases — Catat Pembelian Baru
**Akses:** Admin only (Bearer token admin)

**Request Body:**
```json
{
  "supplier_name": "Toko Grosir ABC",     // Optional, boleh null/kosong
  "notes": "Pembelian bulanan",            // Optional
  "items": [
    {
      // RESTOK produk yang sudah ada:
      "product_id": 1,
      "quantity": 50,
      "buy_price": 30000
    },
    {
      // PRODUK BARU (belum ada di database):
      "product_name": "Produk Baru XYZ",   // Wajib untuk produk baru
      "quantity": 25,
      "buy_price": 8000,
      "sell_price": 15000,                  // Wajib untuk produk baru (harga jual)
      "category_id": 1                      // Optional
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "supplier_name": "Toko Grosir ABC",
  "total_amount": 1700000,
  "notes": "Pembelian bulanan",
  "created_by": 1,
  "created_at": "2026-02-16T13:00:00Z",
  "items": [
    {
      "product_id": 1,
      "product_name": "Kopi ZX",
      "quantity": 50,
      "buy_price": 30000,
      "subtotal": 1500000
    },
    {
      "product_id": 5,
      "product_name": "Produk Baru XYZ",
      "quantity": 25,
      "buy_price": 8000,
      "sell_price": 15000,
      "category_id": 1,
      "subtotal": 200000
    }
  ]
}
```

**Efek otomatis di backend:**
- Stok produk bertambah (+quantity)
- Harga beli produk ter-update
- Jika produk baru → otomatis terdaftar di tabel products
- Cache produk di-invalidate

**Error responses:**
- `400` — Validasi gagal (quantity <= 0, produk tidak ditemukan, nama kosong, dll)
- `403` — Bukan admin
- `500` — Error server

---

### GET /api/purchases — Riwayat Semua Pembelian
**Akses:** Admin only

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "supplier_name": "Toko Grosir ABC",
    "total_amount": 1700000,
    "notes": "Pembelian bulanan",
    "created_by": 1,
    "created_at": "2026-02-16T13:00:00Z"
  },
  {
    "id": 2,
    "supplier_name": null,
    "total_amount": 500000,
    "notes": null,
    "created_by": 1,
    "created_at": "2026-02-15T10:00:00Z"
  }
]
```
**Catatan:** Array `items` TIDAK disertakan di list (untuk performa). Gunakan GET by ID untuk melihat detail items.

---

### GET /api/purchases/{id} — Detail 1 Pembelian
**Akses:** Admin only

**Response (200 OK):**
```json
{
  "id": 1,
  "supplier_name": "Toko Grosir ABC",
  "total_amount": 1700000,
  "notes": "Pembelian bulanan",
  "created_by": 1,
  "created_at": "2026-02-16T13:00:00Z",
  "items": [
    {
      "id": 1,
      "purchase_id": 1,
      "product_id": 1,
      "product_name": "Kopi ZX",
      "quantity": 50,
      "buy_price": 30000,
      "subtotal": 1500000,
      "created_at": "2026-02-16T13:00:00Z"
    },
    {
      "id": 2,
      "purchase_id": 1,
      "product_id": 5,
      "product_name": "Produk Baru XYZ",
      "quantity": 25,
      "buy_price": 8000,
      "sell_price": 15000,
      "category_id": 1,
      "subtotal": 200000,
      "created_at": "2026-02-16T13:00:00Z"
    }
  ]
}
```

---

## 📊 PERUBAHAN DI ENDPOINT REPORT

### GET /api/report/hari-ini dan GET /api/report
Response sekarang memiliki **3 field baru**:

```json
{
  "total_revenue": 500000,
  "total_transaksi": 12,
  "total_items_sold": 35,
  "total_profit": 120000,
  "total_pengeluaran": 200000,      // 🆕 Total pembelian hari ini
  "total_pembelian": 2,             // 🆕 Jumlah transaksi pembelian
  "laba_bersih": 300000,            // 🆕 Revenue - Pengeluaran
  "produk_terlaris": [
    {
      "nama_produk": "Kopi ZX",
      "jumlah": 15,
      "total_sales": 225000,
      "total_profit": 75000
    },
    {
      "nama_produk": "Teh Manis",
      "jumlah": 10,
      "total_sales": 100000,
      "total_profit": 30000
    }
  ]
}
```

---

## 📊 PERUBAHAN DI ENDPOINT TRANSACTIONS

### GET /api/transactions
Setiap transaksi sekarang memiliki **2 field baru**:

```json
[
  {
    "id": 1,
    "total_amount": 131500,
    "created_at": "2026-02-16T06:42:00Z",
    "discount_amount": 0,
    "total_items": 5,       // 🆕 Total items dalam transaksi ini
    "profit": 25000         // 🆕 Keuntungan dari transaksi ini
  }
]
```

---

## 📱 REKOMENDASI PERUBAHAN FRONTEND

### 1. Sidebar Menu — Tambah "Pembelian"
```
📊 Dashboard
🛒 Kasir (POS)
📋 Produk              ← Edit harga & kategori saja
📦 Pembelian            ← 🆕 BARU: Catat pembelian + tambah produk baru
📜 Riwayat Penjualan
📜 Riwayat Pembelian    ← 🆕 BARU: History pembelian
📊 Laporan
⚙️ Pengaturan
```

### 2. Halaman Pembelian Baru (Form)
Mirip form checkout, tapi kebalikannya:
- Input: Supplier (optional), Items (product dropdown atau ketik nama baru)
- Setiap item: pilih produk / ketik baru, quantity, harga beli, (harga jual jika baru)
- Total otomatis dihitung
- Tombol "Simpan Pembelian"

### 3. Halaman Riwayat Pembelian
- Tabel: ID, Tanggal, Supplier, Total, Jumlah Item
- Klik row → detail pembelian (items)

### 4. Dashboard — Update Kartu Summary
Tambah kartu baru:
- 📦 **Pengeluaran** (total_pengeluaran) — dari field baru di report
- 📈 **Laba Bersih** (laba_bersih) — revenue dikurangi pengeluaran

### 5. Halaman Produk — Simplifikasi
- ❌ Hapus tombol "Tambah Produk" (pindah ke Pembelian)
- ✏️ Edit hanya: Nama, Harga Jual, Kategori
- 👁️ Stok dan Harga Beli tampil sebagai read-only
- 💡 Tampilkan info: "Stok dikelola lewat menu Pembelian"

### 6. Tabel Riwayat Penjualan
- Tambah kolom "Items" (total_items)
- Tambah kolom "Profit" (profit)

---

## 📌 RINGKASAN SEMUA ENDPOINT AKTIF

| Method | Endpoint | Akses | Keterangan |
|--------|----------|-------|------------|
| POST | /api/auth/login | Public | Login |
| GET | /api/produk | Auth | List produk |
| GET | /api/produk/{id} | Auth | Detail produk |
| PUT | /api/produk/{id} | Admin | Edit nama/harga/kategori saja |
| DELETE | /api/produk/{id} | Admin | Hapus produk |
| GET | /api/categories | Auth | List kategori |
| POST | /api/categories | Admin | Tambah kategori |
| PUT | /api/categories/{id} | Admin | Edit kategori |
| DELETE | /api/categories/{id} | Admin | Hapus kategori |
| POST | /api/checkout | Auth | Buat transaksi penjualan |
| GET | /api/transactions | Auth | Riwayat penjualan |
| **POST** | **/api/purchases** | **Admin** | **🆕 Catat pembelian baru** |
| **GET** | **/api/purchases** | **Admin** | **🆕 Riwayat pembelian** |
| **GET** | **/api/purchases/{id}** | **Admin** | **🆕 Detail pembelian** |
| GET | /api/report/hari-ini | Auth | Report hari ini |
| GET | /api/report | Auth | Report by date range |
| GET | /api/dashboard/sales-trend | Admin | Trend penjualan |
| GET | /api/dashboard/top-products | Admin | Top produk |
| GET | /api/discounts | Admin | List diskon |
| POST | /api/discounts | Admin | Tambah diskon |
| GET | /api/discounts/active | Auth | Diskon aktif |
