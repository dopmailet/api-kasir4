# 🚀 SESSION 5 ROADMAP: Professional Dashboard & Marketing

## 🎯 Objective
Mengubah aplikasi dari sekadar "pencatat transaksi" menjadi **alat analisis bisnis profesional** dengan fitur Dashboard Analystics dan Manajemen Diskon.

---

## 📦 Part 1: Advanced Analytics (Dashboard)

Fitur ini akan memberikan "mata dewa" bagi Admin untuk melihat kesehatan bisnis.

### **Fitur Baru:**
1.  **Top Products API**
    *   `GET /api/dashboard/top-products?limit=5`
    *   Menampilkan 5 produk terlaris berdasarkan Quantity terjual.
    *   Menampilkan 5 produk dengan Profit tertinggi.

2.  **Sales Trend Chart**
    *   `GET /api/dashboard/trend?days=7`
    *   Data grafik penjualan per hari selama seminggu terakhir.

3.  **Low Stock Alert**
    *   `GET /api/dashboard/low-stock`
    *   List produk yang stoknya hampir habis (misal: stok < 10).

4.  **Summary Cards**
    *   Total Omzet Hari Ini
    *   Total Profit Hari Ini
    *   Jumlah Transaksi Hari Ini

---

## 🏷️ Part 2: Discount System (Marketing)

Fitur ini memungkinkan Admin mengatur strategi harga tanpa mengubah kode program.

### **Database Schema:**
Tabel baru `discounts`:
- `id`: PK
- `name`: "Promo Kemerdekaan"
- `type`: "PERCENT" | "FIXED"
- `value`: 10 (10%) | 5000 (Rp 5.000)
- `product_id`: NULL (jika diskon global) atau Specific Product
- `start_date`: Kapan mulai
- `end_date`: Kapan berakhir
- `is_active`: Toggle on/off manual

### **Logic Update:**
1.  **Admin Dashboard**: CRUD (Create, Read, Update, Delete) diskon.
2.  **Checkout Process**: Modifikasi `TransactionService` untuk menghitung `final_price` setelah diskon.

---

## 🛠️ Prioritas Pengerjaan

Saya sarankan kita mulai dari **Part 1 (Analytics)** dulu karena lebih aman (hanya membaca data, tidak mengubah logika uang). Setelah itu baru masuk ke **Part 2 (Diskon)** yang lebih kompleks logikanya.

Bagaimana menurut Anda?
