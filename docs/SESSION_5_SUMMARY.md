# 🚀 Session 5 Summary: Dashboard Analytics & Advanced Discount System

Sesi ini berfokus pada fitur "Professional Business Tools": Analisis Data dan Promosi Fleksibel (Per Produk & Global).

---

## 📊 1. Dashboard Analytics
Fitur ini memberikan wawasan bisnis bagi Admin.

### ✅ Endpoint Baru:
1.  **Sales Trend** (`GET /api/dashboard/sales-trend?period=day|month|year`)
    *   Menampilkan grafik penjualan dan profit.
2.  **Top Products** (`GET /api/dashboard/top-products?limit=5`)
    *   Menampilkan produk terlaris.

---

## 🏷️ 2. Advanced Discount System
Sistem diskon bertingkat yang mendukung promosi produk spesifik maupun total belanja.

### ✅ Database Updates:
*   New Table: `discounts` dengan kolom `product_id` (Nullable).
*   Logic: Jika `product_id` diisi, diskon hanya untuk produk tersebut. Jika NULL, diskon untuk total transaksi.

### ✅ Logic Checkout (Transaction Engine):
Saat Checkout (`POST /api/checkout`), sistem bekerja dengan urutan:
1.  **Product Discount (Otomatis)**:
    *   Sistem mengecek setiap item di keranjang.
    *   Jika ada diskon aktif untuk produk tersebut (misal "Kopi Susu 50%"), harga langsung dipotong.
    *   Subtotal item = (Harga Asli - Diskon) * Qty.
2.  **Global Discount (Manual/Selected)**:
    *   Jika Kasir memilih ID Diskon Global (misal "Promo Kemerdekaan 10%").
    *   Diskon dihitung dari Total Belanja (yang sudah dipotong diskon produk).
    *   Total Bayar = Total Belanja - Diskon Global.

### ✅ API Management Diskon:
*   `GET /api/discounts`: List semua diskon (termasuk per produk).
*   `POST /api/discounts`: Buat promo baru.
    *   Isi `product_id` untuk diskon produk.
    *   Kosongkan `product_id` untuk diskon global.
*   `GET /api/discounts/active`: List diskon global untuk dipilih kasir (Diskon produk tidak perlu dipilih karena otomatis).

---

## 🧪 Cara Testing Fitur Baru

### Langkah 1: Setup Database (REVISI PENTING!)
Anda harus menjalankan script SQL berikut di SQL Editor Supabase/pgAdmin:
`database/migration_create_discounts.sql`
*(Jika tabel `discounts` sudah ada, harap `DROP TABLE discounts CASCADE;` terlebih dahulu).*

### Langkah 2: Buat Diskon Produk (Sebagai Admin)
```bash
POST /api/discounts
Body:
{
    "name": "Kopi Susu 50% Off",
    "type": "PERCENTAGE",
    "value": 50,
    "product_id": 1,  <-- Hanya berlaku untuk Product ID 1
    "start_date": "...", "end_date": "...", "is_active": true
}
```

### Langkah 3: Checkout (Sebagai Kasir)
Beli produk ID 1. Harga otomatis terpotong 50%.
Jika pakai diskon global juga, diskon bertumpuk.

---

**Status Project: Backend Lengkap & Professional!** 🚀
Siap untuk pengembangan Frontend di sesi berikutnya.
