# 📊 Dashboard Analytics API Documentation

Dokumentasi ini berisi daftar endpoint untuk fitur Dashboard & Analisis Bisnis.
Semua endpoint di bawah ini memerlukan **Authentication Token** dan **Role Admin**.

## 🔐 Authentication
Semua request harus menyertakan header:
`Authorization: Bearer <token_admin>`

---

## 1. Get Sales Trend (Grafik Penjualan)
Endpoint ini digunakan untuk menampilkan grafik penjualan dan profit berdasarkan periode waktu.

- **URL**: `/api/dashboard/sales-trend`
- **Method**: `GET`
- **Query Params**:
  - `period`: `day` (Default, 7 hari terakhir), `month` (12 bulan terakhir), `year` (5 tahun terakhir).

**Contoh Request:**
`GET /api/dashboard/sales-trend?period=month`

**Contoh Response:**
```json
{
  "period": "month",
  "data": [
    {
      "date": "2024-01",
      "total_sales": 5000000,
      "total_profit": 1500000,
      "transaction_count": 50
    },
    {
      "date": "2024-02",
      "total_sales": 6200000,
      "total_profit": 1800000,
      "transaction_count": 65
    }
  ]
}
```

---

## 2. Get Top Selling Products (Produk Terlaris)
Endpoint ini menampilkan ranking produk berdasarkan jumlah terjual (Quantity) dan keuntungan (Profit).

- **URL**: `/api/dashboard/top-products`
- **Method**: `GET`
- **Query Params**:
  - `limit`: Jumlah data yang ditampilkan (Default: 5).

**Contoh Request:**
`GET /api/dashboard/top-products?limit=5`

**Contoh Response:**
```json
{
  "by_quantity": [
    {
      "nama_produk": "Kopi Susu",
      "jumlah": 150,
      "total_sales": 2250000,
      "total_profit": 750000
    },
    {
      "nama_produk": "Roti Bakar",
      "jumlah": 120,
      "total_sales": 1800000,
      "total_profit": 600000
    }
  ],
  "by_profit": [
    {
      "nama_produk": "Steak Sapi", // Walau jumlah dikit, profit gede
      "jumlah": 50,
      "total_sales": 5000000,
      "total_profit": 2500000
    },
    {
      "nama_produk": "Kopi Susu",
      "jumlah": 150,
      "total_sales": 2250000,
      "total_profit": 750000
    }
  ]
}
```

---

## ⚠️ Error Responses
- **401 Unauthorized**: Token tidak valid atau tidak ada.
- **403 Forbidden**: Token valid tapi bukan role Admin (misal: Kasir).
- **500 Internal Server Error**: Kesalahan teknis di server.
