# 📦 TODO: Modul Pembelian (Purchase Module)

**Tanggal Mulai:** 16 Februari 2026  
**Status:** 🔲 Belum Dimulai

---

## TAHAP 1: Database Schema
> Buat tabel baru di Supabase

- [ ] **1.1** Buat tabel `purchases` (header pembelian)
  - `id` (PK, auto increment)
  - `supplier_name` (varchar, nullable/optional)
  - `total_amount` (decimal, total pembelian)
  - `notes` (text, nullable, catatan tambahan)
  - `created_by` (FK ke users, siapa yang input)
  - `created_at` (timestamp)

- [ ] **1.2** Buat tabel `purchase_items` (detail item pembelian)
  - `id` (PK, auto increment)
  - `purchase_id` (FK ke purchases)
  - `product_id` (FK ke products, nullable — null jika produk baru)
  - `product_name` (varchar, nama produk — untuk produk baru)
  - `quantity` (int, jumlah beli)
  - `buy_price` (decimal, harga beli per unit)
  - `sell_price` (decimal, harga jual — hanya untuk produk baru)
  - `category_id` (FK ke categories, nullable — hanya untuk produk baru)
  - `subtotal` (decimal, quantity × buy_price)
  - `created_at` (timestamp)

- [ ] **1.3** Generate SQL migration script
- [ ] **1.4** Update `schema_complete.sql` dengan tabel baru

---

## TAHAP 2: Backend — Model
> Buat struct Go untuk pembelian

- [ ] **2.1** Buat file `models/purchase.go`
  - Struct `Purchase` (header pembelian)
  - Struct `PurchaseItem` (detail item)
  - Struct `PurchaseRequest` (request body dari frontend)
  - Struct `PurchaseItemRequest` (item dalam request)

---

## TAHAP 3: Backend — Repository
> Buat layer database untuk pembelian

- [ ] **3.1** Buat file `repositories/purchase_repository.go`
  - Fungsi `Create()` — Simpan pembelian baru ke database
    - Insert ke tabel `purchases` (header)
    - Insert ke tabel `purchase_items` (detail)
    - Jika produk baru → insert ke tabel `products`
    - Jika produk lama → update stok (stok += quantity)
    - Update `harga_beli` di tabel `products`
    - Semua dalam 1 database transaction (rollback jika error)
  - Fungsi `GetAll()` — Ambil riwayat semua pembelian
  - Fungsi `GetByID()` — Ambil detail 1 pembelian + items-nya

---

## TAHAP 4: Backend — Service
> Buat layer business logic untuk pembelian

- [ ] **4.1** Buat file `services/purchase_service.go`
  - Fungsi `Create()` — Validasi input + panggil repository
    - Validasi: minimal 1 item
    - Validasi: quantity > 0
    - Validasi: buy_price >= 0
    - Validasi: jika produk baru → nama & sell_price wajib diisi
  - Fungsi `GetAll()` — Ambil riwayat pembelian
  - Fungsi `GetByID()` — Ambil detail pembelian

---

## TAHAP 5: Backend — Handler
> Buat layer HTTP untuk pembelian

- [ ] **5.1** Buat file `handlers/purchase_handler.go`
  - `POST /api/purchases` — Catat pembelian baru (Admin only)
  - `GET /api/purchases` — Riwayat pembelian (Admin only)
  - `GET /api/purchases/{id}` — Detail 1 pembelian (Admin only)

---

## TAHAP 6: Backend — Register Routes
> Daftarkan endpoint baru di main.go

- [ ] **6.1** Tambahkan dependency injection di `main.go`
  - purchaseRepo → purchaseService → purchaseHandler
- [ ] **6.2** Register routes:
  - `/api/purchases` → GET (list), POST (create)
  - `/api/purchases/` → GET by ID
- [ ] **6.3** Update console log (daftar endpoint)

---

## TAHAP 7: Backend — Ubah Modul Produk
> Sesuaikan produk agar stok tidak bisa diubah manual

- [ ] **7.1** Ubah `POST /api/produk` → HAPUS (produk baru lewat pembelian saja)
- [ ] **7.2** Ubah `PUT /api/produk/{id}` → Hanya bisa ubah: nama, harga jual, kategori
  - ❌ Tidak bisa ubah stok
  - ❌ Tidak bisa ubah harga beli
- [ ] **7.3** Invalidate cache produk setelah pembelian berhasil

---

## TAHAP 8: Backend — Report Pengeluaran
> Tambahkan data pengeluaran ke laporan

- [ ] **8.1** Tambah endpoint `GET /api/report/pengeluaran`
  - Total pengeluaran hari ini
  - Total pengeluaran per periode
- [ ] **8.2** Update `GET /api/report/hari-ini`
  - Tambah field `total_pengeluaran` (total pembelian hari ini)
  - Tambah field `laba_bersih` (revenue - pengeluaran)

---

## TAHAP 9: Testing & Deploy
> Test semua dan deploy ke Railway

- [ ] **9.1** Build & pastikan tidak ada error
- [ ] **9.2** Test manual via Postman:
  - [ ] Catat pembelian dengan produk baru
  - [ ] Catat pembelian restok produk lama
  - [ ] Catat pembelian tanpa supplier (optional)
  - [ ] Cek stok bertambah setelah pembelian
  - [ ] Cek harga_beli ter-update setelah pembelian
  - [ ] Cek riwayat pembelian
  - [ ] Cek report pengeluaran
  - [ ] Cek produk tidak bisa ubah stok manual
- [ ] **9.3** Commit & push ke GitHub
- [ ] **9.4** Jalankan SQL migration di Supabase
- [ ] **9.5** Verifikasi Railway deploy sukses

---

## TAHAP 10: Dokumentasi untuk Frontend
> Berikan spesifikasi ke AI frontend

- [ ] **10.1** Buat dokumen spesifikasi API pembelian
  - Request/Response format lengkap
  - Contoh penggunaan
  - Flow frontend yang direkomendasikan
- [ ] **10.2** Buat SQL migration yang perlu dijalankan di Supabase

---

## 📊 Ringkasan File yang Akan Dibuat/Diubah

### File BARU:
| No | File | Keterangan |
|----|------|-----------|
| 1 | `models/purchase.go` | Struct Purchase, PurchaseItem |
| 2 | `repositories/purchase_repository.go` | Database operations |
| 3 | `services/purchase_service.go` | Business logic & validasi |
| 4 | `handlers/purchase_handler.go` | HTTP handlers |
| 5 | `database/migrations/add_purchases.sql` | SQL migration |

### File yang DIUBAH:
| No | File | Perubahan |
|----|------|-----------|
| 1 | `main.go` | Tambah DI & routes pembelian |
| 2 | `handlers/product_handler.go` | Hapus POST, batasi PUT |
| 3 | `repositories/product_repository.go` | Hapus Create, batasi Update |
| 4 | `services/product_service.go` | Hapus Create, batasi Update |
| 5 | `models/report.go` | Tambah field pengeluaran & laba |
| 6 | `repositories/report_repository.go` | Query pengeluaran |
| 7 | `database/schema_complete.sql` | Tambah tabel purchases |

---

## ⚠️ Catatan Penting
- Tahap 7 (ubah modul produk) dikerjakan TERAKHIR setelah modul pembelian sudah jalan
- Setiap tahap di-test dulu sebelum lanjut ke tahap berikutnya
- Frontend BARU dikerjakan setelah semua backend selesai & tested
