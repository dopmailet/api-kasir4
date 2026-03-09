# Panduan Integrasi Frontend: Customer Database & Loyalty (KasirPOS)

Dokumen ini berisi informasi krusial untuk AI Agent Frontend yang bertugas mengintegrasikan fitur Pelanggan & Loyalitas ke dalam UI KasirPOS.

## 1. Konfigurasi Fitur & Settings
Sebelum me-render UI, Frontend harus mengecek status pengaturan dari backend.

- **Endpoint**: `GET /api/settings/customer`
- **Response**:
```json
{
  "showCustomerInPOS": true,
  "enableLoyaltyPoints": true
}
```
**Instruksi FE:**
- Jika `showCustomerInPOS` == `false`, **SEMBUNYIKAN** semua elemen UI pencarian pelanggan di halaman kasir.
- Jika `enableLoyaltyPoints` == `false`, jangan tampilkan teks seputar *Point Reward / Loyalty* pada resi atau keranjang, namun fitur customer-nya sendiri (apabila `true`) tetap bisa dipakai untuk tracking histori belanja.

---

## 2. API Endpoints Pelanggan (Customer)

### A. Pencarian Cepat Pelanggan (Untuk Halaman POS / Kasir)
- **Endpoint**: `GET /api/customers/search?q={keyword}`
- **Catatan**: Selalu gunakan ini di *search bar* POS. Hanya mereturn pelanggan yang **aktif**.

### B. Manajemen Data (Menu Database Pelanggan)
Hanya role **Admin** yang dapat melihat Full List dan melakukan EDIT (PUT). Kasir hanya bisa menambah (POST) atau mengakses pencarian aktif.
- **GET All** : `GET /api/customers?page=1&limit=10&q={katakunci}&status={active|inactive|all}`
- **Create** : `POST /api/customers`
  ```json
  {
    "name": "John Doe",
    "phone": "08123456789",
    "address": "Jalan ABC",
    "notes": "Pelanggan VIP"
  }
  ```
- **Edit** : `PUT /api/customers/{id}` (Payload sama dengan Create, tambah field boolean `is_active`)
- **Get Detail** : `GET /api/customers/{id}`

### C. Riwayat Pelanggan (Detail Page)
- **Riwayat Belanja (Struk)**: `GET /api/customers/{id}/transactions`
- **Riwayat Poin (Earn/Redeem)**: `GET /api/customers/{id}/loyalty-transactions`

---

## 3. Integrasi pada Proses Checkout `/api/checkout`
Sistem pembayaran kini mendukung referensi `customer_id` secara **opsional**. 

**Payload Baru `POST /api/checkout`**:
```json
{
  "items": [ ... ],
  "payment_amount": 50000,
  "discount_amount": 0,
  "customer_id": 2  // <--- TAMBAHAN BARU (Kirim int ID jika ada pelanggan yang di-select, abaikan/hapus field jika guest)
}
```

**Response Checkout (Perubahan)**:
Ketika sukses (201/200 OK), response `data` akan memiliki obyek tambahan di dalamnya:
```json
{
  "id": 165,
  "total_amount": 40000,
  // ... field checkout lama ...
  "customer": {
     "id": 2,
     "name": "Budi",
     "phone": "0812xxx",
     "loyalty_points": 13 
  },
  "points_earned": 4
}
```

**Penting untuk Frontend Saat Checkout**:
1. FE **tidak perlu** menghitung poin secara manual, seluruh kalkulasi (1 Poin / Rp. 10.000) dan penambahan ke saldo pelanggan dilakukan secara **otomatis dan atomik** di Backend.
2. Gunakan nilai return `points_earned` dari respon sukses `/api/checkout` untuk menampilkannya di *UI Success* atau Struk ("Anda mendapatkan 4 poin dari transaksi ini").
