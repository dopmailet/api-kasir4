# Kasir API - Layered Architecture

REST API untuk sistem kasir dengan arsitektur berlapis (Layered Architecture).

## 🏗️ Arsitektur

Aplikasi ini menggunakan **Layered Architecture** dengan pemisahan tanggung jawab yang jelas:

```
┌─────────────────────────────────────┐
│         HTTP Request/Response        │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│           HANDLER LAYER              │  ← Terima request, kirim response
│     (product_handler.go, etc)        │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│           SERVICE LAYER              │  ← Business logic
│     (product_service.go, etc)        │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│         REPOSITORY LAYER             │  ← Database queries
│   (product_repository.go, etc)       │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│            DATABASE                  │  ← PostgreSQL (Supabase/Railway)
│         (products, categories)       │
└─────────────────────────────────────┘
```

## 📁 Struktur Folder

```
kasir-api/
├── main.go                    # Entry point & routing
├── go.mod                     # Dependencies
├── .env.example               # Template environment variables
│
├── config/                    # Konfigurasi
│   └── database.go            # Database connection
│
├── model/                     # Data structures
│   ├── product.go             # Product model
│   └── category.go            # Category model
│
├── repository/                # Database layer
│   ├── product_repository.go  # Product DB queries
│   └── category_repository.go # Category DB queries
│
├── service/                   # Business logic layer
│   ├── product_service.go     # Product business logic
│   └── category_service.go    # Category business logic
│
└── handler/                   # HTTP layer
    ├── product_handler.go     # Product HTTP handlers
    └── category_handler.go    # Category HTTP handlers
```

## 🚀 Teknologi

- **Go** 1.24+
- **PostgreSQL** - Database (Supabase/Railway)
- **Redis** - Caching layer (Optional)
- **lib/pq** - PostgreSQL driver
- **go-redis** - Redis client

## ✨ Features

- ✅ **Layered Architecture** - Clean separation of concerns
- ✅ **Redis Caching** - Fast response dengan in-memory cache
- ✅ **Pagination** - Efficient data handling untuk large datasets
- ✅ **Search** - Filter products by name
- ✅ **Transactions** - Checkout dengan multiple items
- ✅ **Sales Reports** - Daily dan date range reports
- ✅ **Category Management** - Product categorization dengan JOIN

### 🔴 Redis Caching

API ini menggunakan Redis untuk caching:
- **Cache Strategy**: Cache-aside pattern
- **TTL**: 5 menit (configurable)
- **Invalidation**: Automatic pada POST/PUT/DELETE
- **Graceful Degradation**: App tetap jalan tanpa Redis

📖 **Panduan lengkap**: [REDIS_PAGINATION_GUIDE.md](./docs/REDIS_PAGINATION_GUIDE.md)

### 📄 Pagination

Semua list endpoints support pagination:
```
GET /api/produk?page=1&limit=10
```

Response format:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 45,
    "total_pages": 5
  }
}
```

## 📦 Instalasi

### 1. Clone Repository
```bash
git clone https://github.com/zainalzxc/kasir-api.git
cd kasir-api
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Database

**Opsi A: Local PostgreSQL**
```bash
# Install PostgreSQL
# Buat database baru
createdb kasir_db

# Copy .env.example ke .env
cp .env.example .env

# Edit .env dan isi DATABASE_URL
DATABASE_URL=host=localhost user=postgres password=postgres dbname=kasir_db port=5432 sslmode=disable
```

**Opsi B: Supabase (Cloud) - RECOMMENDED** ⭐
1. **Jalankan setup wizard:**
   ```bash
   # Windows
   setup-supabase.bat
   
   # Atau baca panduan lengkap
   ```
2. **Ikuti panduan lengkap di:** [SUPABASE_SETUP.md](./SUPABASE_SETUP.md)
3. **Quick steps:**
   - Buat project di https://supabase.com
   - Jalankan SQL script: `database/supabase_setup.sql`
   - Copy connection string dari Settings > Database
   - Paste ke file `.env`

### 4. Jalankan Aplikasi
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Endpoints

### Health Check
- `GET /health` - Cek status API

### Products
- `GET /api/produk` - Get all products (with pagination)
  - Query params: `?page=1&limit=10&name=search`
- `POST /api/produk` - Create new product (auto-update stock if product name exists)
- `GET /api/produk/{id}` - Get product by ID
- `PUT /api/produk/{id}` - Update product
- `DELETE /api/produk/{id}` - Delete product

### Categories
- `GET /api/categories` - Get all categories
- `POST /api/categories` - Create new category
- `GET /api/categories/{id}` - Get category by ID
- `PUT /api/categories/{id}` - Update category
- `DELETE /api/categories/{id}` - Delete category

### Transactions
- `POST /api/checkout` - Create transaction with multiple items

### Reports
- `GET /api/report/hari-ini` - Get today's sales report
- `GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Get sales report by date range

## 📮 Testing dengan Postman

### Import Postman Collection

1. **Buka Postman**
2. **Import Collection**:
   - Klik **Import** di Postman
   - Pilih file `Kasir-API.postman_collection.json`
   - Klik **Import**

3. **Import Environment** (Optional):
   - Import file `Kasir-API.postman_environment.json`
   - Pilih environment "Kasir API - Local"

### Quick Test

1. **Health Check**: `GET /health`
2. **Get All Products**: `GET /api/produk`
3. **Create Product**: `POST /api/produk`
   ```json
   {
     "nama": "Indomie Goreng",
     "harga": 3500,
     "stok": 100
   }
   ```

📖 **Panduan lengkap**: Lihat [POSTMAN_GUIDE.md](./POSTMAN_GUIDE.md)

## 🧪 Testing dengan cURL

### Create Product
```bash
POST /api/produk
Content-Type: application/json

{
  "nama": "Indomie Goreng",
  "harga": 3500,
  "stok": 100
}
```

### Get All Products
```bash
GET /api/produk
```

## 🚢 Deployment ke Production

### **Opsi 1: Railway (Recommended) ⭐**

Railway menyediakan PostgreSQL **DAN** Redis dalam satu platform!

#### **1. Push ke GitHub**
```bash
git add .
git commit -m "feat: add Redis caching and pagination"
git push origin main
```

#### **2. Deploy di Railway**
1. Login ke https://railway.app
2. New Project → Deploy from GitHub
3. Pilih repository `kasir-api`
4. Railway akan auto-detect Go project

#### **3. Add PostgreSQL**
- Click "New" → "Database" → "Add PostgreSQL"
- `DATABASE_URL` otomatis di-set ✅

#### **4. Add Redis** 🔴
- Click "New" → "Database" → "Add Redis"
- `REDIS_URL` otomatis di-set ✅

#### **5. Deploy!**
Railway akan otomatis build dan deploy!

Check logs untuk:
```
✅ Database connected successfully
✅ Redis connected successfully!
🚀 Server running on port: XXXX
```

---

### **Opsi 2: Supabase + Upstash Redis (FREE)**

Jika sudah menggunakan Supabase untuk PostgreSQL:

#### **1. Setup Upstash Redis (FREE)**
1. Buka https://upstash.com
2. Create account (free, no credit card)
3. Create Redis database
4. Copy Redis URL

📖 **Panduan lengkap**: [UPSTASH_REDIS_SETUP.md](./docs/UPSTASH_REDIS_SETUP.md)

#### **2. Set Environment Variables**
Di Railway/Vercel/Render:
```bash
DATABASE_URL=postgresql://...supabase...
REDIS_URL=redis://default:xxx@xxx.upstash.io:6379
```

#### **3. Deploy**
Push ke GitHub dan deploy seperti biasa!

---

### **Tanpa Redis (Tetap Berfungsi)**

App akan tetap berjalan normal tanpa Redis:
- ✅ Semua fitur tetap bekerja
- ✅ Data langsung dari database
- ⚠️ Sedikit lebih lambat (no caching)

Log akan menampilkan:
```
Warning: Gagal connect ke Redis, Redis caching akan dinonaktifkan
```

## 🔧 Development

### Menambah Fitur Baru

1. **Buat Model** di `model/`
2. **Buat Repository** di `repository/`
3. **Buat Service** di `service/`
4. **Buat Handler** di `handler/`
5. **Update Routing** di `main.go`

### Contoh: Menambah Resource "Supplier"

```go
// 1. model/supplier.go
type Supplier struct {
    ID   int    `json:"id" gorm:"primaryKey"`
    Nama string `json:"nama"`
}

// 2. repository/supplier_repository.go
type SupplierRepository interface {
    FindAll() ([]model.Supplier, error)
    // ... methods lainnya
}

// 3. service/supplier_service.go
type SupplierService interface {
    GetAll() ([]model.Supplier, error)
    // ... methods lainnya
}

// 4. handler/supplier_handler.go
type SupplierHandler struct {
    service service.SupplierService
}

// 5. main.go - tambahkan routing
http.HandleFunc("/api/suppliers", ...)
```

## 📝 Notes

- Database tables akan dibuat otomatis (auto-migration)
- Data persisten di database (tidak hilang saat restart)
- Gunakan environment variables untuk konfigurasi
- Ikuti pattern yang sama untuk konsistensi

## 👨‍💻 Author

Bootcamp Golang - KodingWorks

## 📄 License

MIT
