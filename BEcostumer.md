Tambahkan fitur baru Customer Database + Loyalty sederhana pada backend aplikasi POS berbasis Go + PostgreSQL + JWT auth + REST API.

KONTEKS:
Backend POS sudah memiliki:
- auth JWT
- roles admin & kasir
- products
- categories
- purchases
- employees
- payroll
- expenses
- dashboard
- reports
- transactions / checkout
- settings dasar jika ada

TUJUAN:
Tambahkan modul customer sederhana yang dapat terhubung ke transaksi POS dan mendukung loyalty points versi awal.

ATURAN BISNIS:
1. Customer bersifat opsional pada transaksi.
2. Customer bisa aktif/nonaktif.
3. Customer nonaktif tidak boleh dipilih untuk transaksi baru.
4. Loyalty sederhana:
   - setiap Rp10.000 dari total akhir setelah diskon = 1 poin
   - floor pembulatan ke bawah
   - poin hanya bertambah setelah checkout berhasil
   - redeem poin belum diimplementasikan
5. Jika customer_id null, transaksi tetap valid.
6. Fitur customer di POS bisa disembunyikan dari frontend lewat settings, tetapi backend tetap menyediakan endpoint customer.
7. Admin bisa CRUD customer.
8. Kasir boleh mencari customer aktif dan menambah customer cepat.
9. Semua perubahan harus backward compatible.

DATABASE CHANGES:

1. Buat tabel customers
Kolom minimal:
- id UUID / bigserial
- customer_code varchar unique
- name varchar not null
- phone varchar not null
- address text null
- notes text null
- loyalty_points integer not null default 0
- total_spent numeric(18,2) not null default 0
- total_transactions integer not null default 0
- last_transaction_at timestamptz null
- is_active boolean not null default true
- created_at timestamptz not null default now()
- updated_at timestamptz not null default now()

Aturan:
- phone boleh unique jika memungkinkan, tapi jika existing bisnis berpotensi duplikat, minimal beri index
- buat trigger updated_at seperti tabel lain

2. Buat tabel loyalty_transactions
Kolom:
- id UUID / bigserial
- customer_id FK references customers(id) on delete cascade
- transaction_id FK references transactions(id) on delete set null
- type varchar not null check in ('earn','adjust')
- points integer not null
- description text null
- created_by_user_id FK references users(id) on delete set null
- created_at timestamptz not null default now()

Catatan:
- versi awal hanya pakai type='earn'
- adjust disiapkan untuk masa depan

3. Tambahkan customer_id nullable ke tabel transactions / sales header
- customer_id FK references customers(id) on delete set null
- index pada customer_id

4. Tambahkan settings sederhana jika belum ada penyimpanan server-side
Minimal key:
- show_customer_in_pos boolean default true
- enable_loyalty_points boolean default true

Jika project sudah punya sistem settings generik, gunakan itu.
Jika belum ada, buat tabel app_settings sederhana:
- key varchar primary key
- value_json jsonb not null
- updated_at timestamptz not null default now()

API YANG HARUS DIBUAT:

1. CUSTOMER CRUD
- GET /api/customers
  Query params:
  - q
  - status=active|inactive|all
  - page
  - limit
  - sort_by
  - sort_order
  Role: admin untuk full list, kasir boleh akses list sederhana aktif saja bila dipakai POS search

- GET /api/customers/:id
  Return detail customer + summary

- POST /api/customers
  Role: admin, kasir
  Payload:
  - name
  - phone
  - address optional
  - notes optional
  - is_active optional
  Saat create:
  - generate customer_code otomatis, misal CUST-000001

- PUT /api/customers/:id
  Role: admin
  Update:
  - name
  - phone
  - address
  - notes
  - is_active

2. CUSTOMER SEARCH
- GET /api/customers/search?q=
  Return hanya customer aktif
  Search by name or phone
  Limit hasil agar ringan untuk POS
  Role: admin, kasir

3. CUSTOMER TRANSACTION HISTORY
- GET /api/customers/:id/transactions
  Return daftar transaksi customer
  fields minimal:
  - transaction_id
  - date
  - total
  - cashier_name

4. CUSTOMER LOYALTY HISTORY
- GET /api/customers/:id/loyalty-transactions
  Return ledger poin customer
  fields:
  - id
  - type
  - points
  - description
  - transaction_id
  - created_at

5. SETTINGS ENDPOINT
Jika settings existing sudah ada, tambahkan key berikut:
- showCustomerInPOS
- enableLoyaltyPoints

Jika belum ada, buat:
- GET /api/settings/customer
- PUT /api/settings/customer

Payload:
{
  "showCustomerInPOS": true,
  "enableLoyaltyPoints": true
}

Role:
- GET: admin
- PUT: admin

CHECKOUT INTEGRATION:

Update endpoint checkout existing:
- POST /api/checkout

Tambahkan customer_id nullable pada payload transaksi.
Flow checkout:
1. validasi customer_id jika ada
2. pastikan customer aktif
3. simpan transaksi seperti biasa
4. jika customer_id ada dan loyalty aktif:
   - hitung points_earned = floor(final_total / 10000)
   - jika points_earned > 0:
     - insert ke loyalty_transactions type='earn'
     - update customers.loyalty_points += points_earned
   - update customers.total_spent += final_total
   - update customers.total_transactions += 1
   - update customers.last_transaction_at = now()
5. jika customer_id ada tetapi loyalty dimatikan:
   - tetap update total_spent, total_transactions, last_transaction_at
   - jangan tambah poin

Checkout response tambahkan:
- customer_summary nullable
- points_earned integer
Contoh:
{
  "transaction_id": "...",
  "invoice_no": "...",
  "total": 85000,
  "customer_summary": {
    "id": "...",
    "name": "Andi",
    "phone": "08123456789",
    "loyalty_points": 120
  },
  "points_earned": 8
}

VALIDATION RULES:
- name wajib
- phone wajib
- customer nonaktif tidak bisa dipakai untuk checkout
- customer nonaktif tetap bisa dilihat di admin page
- customer boleh diaktifkan kembali
- transaksi lama dengan customer nonaktif tetap valid
- jika customer dihapus secara logis, gunakan is_active=false, jangan hard delete dari UI
- hard delete endpoint tidak perlu dibuat

AUTHORIZATION:
- Admin:
  - full customer management
  - update settings
  - lihat seluruh detail customer
- Kasir:
  - search customer aktif
  - create customer cepat
  - lihat detail customer sederhana jika diperlukan oleh POS
  - tidak boleh nonaktifkan customer
  - tidak boleh ubah settings

IMPLEMENTATION NOTES:
- gunakan pattern project yang sudah ada
- pisahkan layer: handler, service, repository
- gunakan transaction database untuk proses checkout + update loyalty agar atomic
- tambahkan migration SQL yang rapi
- tambahkan index untuk phone, is_active, last_transaction_at
- gunakan response JSON konsisten dengan modul lain
- tambahkan unit/integration test minimal untuk:
  - create customer
  - search customer
  - checkout dengan customer aktif
  - checkout dengan customer nonaktif gagal
  - loyalty points bertambah sesuai aturan
  - loyalty disabled tidak menambah poin

BONUS YANG BOLEH DITAMBAHKAN JIKA MUDAH:
- endpoint quick create customer yang sama dengan POST /api/customers
- field customer_code di response list/detail
- sorting customer berdasarkan total_spent atau loyalty_points

YANG JANGAN DITAMBAHKAN DULU:
- redeem poin
- tier member bronze/silver/gold
- referral
- coupon
- birthday promo
- WhatsApp blast
- multi store customer segmentation

HASIL AKHIR YANG DIHARAPKAN:
- Admin dapat mengelola customer
- Kasir dapat memilih customer saat transaksi
- Customer bisa disembunyikan dari POS melalui settings
- Customer bisa dinonaktifkan
- Checkout dapat mencatat transaksi customer
- Loyalty points otomatis bertambah setelah transaksi sukses
- Sistem tetap sederhana, stabil, dan mudah dikembangkan nanti