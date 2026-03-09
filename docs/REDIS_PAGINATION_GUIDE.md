# Redis Caching & Pagination Guide

## 📋 Overview

Kasir API sekarang dilengkapi dengan:
1. **Redis Caching** - Untuk meningkatkan performa dengan menyimpan data yang sering diakses di memory
2. **Pagination** - Untuk mengelola data dalam jumlah besar dengan efisien

## 🔴 Redis Caching

### Cara Kerja

#### GET Request (Read)
```
1. Client request GET /api/produk
2. Service cek Redis cache dulu
3. Jika ada di cache (CACHE HIT):
   → Return data dari Redis (sangat cepat!)
4. Jika tidak ada di cache (CACHE MISS):
   → Ambil dari database
   → Simpan ke Redis untuk request berikutnya
   → Return data ke client
```

#### POST/PUT/DELETE Request (Write)
```
1. Client request POST/PUT/DELETE
2. Service update/create/delete data di database
3. Service hapus cache yang relevan (invalidate)
4. Return response ke client

Kenapa hapus cache?
→ Supaya GET request berikutnya ambil data terbaru dari database
```

### Cache Keys Pattern

Cache keys menggunakan pattern hierarki:

```
products:list:search:{name}:page:{page}:limit:{limit}
products:detail:id:{id}
categories:list:page:{page}:limit:{limit}
```

### Cache TTL (Time To Live)

- **Default**: 5 menit
- Setelah 5 menit, cache otomatis dihapus
- Data terbaru akan diambil dari database

### Setup Redis

#### Local Development

1. **Install Redis** (Windows):
   ```bash
   # Download Redis dari: https://github.com/microsoftarchive/redis/releases
   # Atau gunakan Docker:
   docker run -d -p 6379:6379 redis:alpine
   ```

2. **Set environment variable**:
   ```bash
   # Di .env
   REDIS_URL=redis://localhost:6379/0
   ```

3. **Test Redis**:
   ```bash
   redis-cli ping
   # Output: PONG
   ```

#### Production (Railway/Upstash)

1. **Railway**:
   - Add Redis plugin di Railway dashboard
   - REDIS_URL akan otomatis di-set

2. **Upstash** (Alternative):
   - Buat database di https://upstash.com
   - Copy REDIS_URL
   - Set di environment variables

### Graceful Degradation

**Jika Redis tidak tersedia**, aplikasi tetap berjalan normal:
- Cache akan di-skip
- Data langsung diambil dari database
- Tidak ada error

Log akan menampilkan:
```
Warning: Gagal connect ke Redis, Redis caching akan dinonaktifkan
```

## 📄 Pagination

### Query Parameters

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `page` | integer | 1 | - | Halaman ke berapa |
| `limit` | integer | 10 | 100 | Jumlah items per halaman |

### Request Examples

#### Get halaman pertama (default)
```http
GET /api/produk
```

#### Get halaman 2, 10 items per halaman
```http
GET /api/produk?page=2&limit=10
```

#### Get halaman 1, 50 items per halaman
```http
GET /api/produk?page=1&limit=50
```

#### Kombinasi dengan search
```http
GET /api/produk?name=kopi&page=1&limit=20
```

### Response Format

Response sekarang include pagination metadata:

```json
{
  "data": [
    {
      "id": 1,
      "nama": "Kopi Arabica",
      "harga": 25000,
      "stok": 100,
      "category_id": 1,
      "category": {
        "id": 1,
        "nama": "Minuman",
        "description": "Minuman segar"
      }
    }
    // ... more products
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 45,
    "total_pages": 5
  }
}
```

### Pagination Metadata

| Field | Description |
|-------|-------------|
| `page` | Halaman saat ini |
| `limit` | Items per halaman |
| `total_items` | Total semua items (dari database) |
| `total_pages` | Total halaman yang tersedia |

### Implementasi di Frontend

```javascript
// Fetch halaman pertama
const response = await fetch('/api/produk?page=1&limit=10');
const data = await response.json();

console.log(`Showing ${data.data.length} items`);
console.log(`Total items: ${data.pagination.total_items}`);
console.log(`Total pages: ${data.pagination.total_pages}`);

// Fetch halaman berikutnya
const nextPage = data.pagination.page + 1;
if (nextPage <= data.pagination.total_pages) {
  const nextResponse = await fetch(`/api/produk?page=${nextPage}&limit=10`);
  // ...
}
```

## 🚀 Performance Benefits

### Tanpa Redis & Pagination
```
Request 1: Database query (100ms)
Request 2: Database query (100ms)
Request 3: Database query (100ms)
Total: 300ms

Jika 1000 products → return semua (slow!)
```

### Dengan Redis & Pagination
```
Request 1: Database query (100ms) + Set cache
Request 2: Redis cache (5ms) ✅
Request 3: Redis cache (5ms) ✅
Total: 110ms (73% faster!)

Jika 1000 products → return 10 per page (fast!)
```

## 🧪 Testing

### Test Pagination

```bash
# Test default pagination
curl http://localhost:8080/api/produk

# Test custom pagination
curl "http://localhost:8080/api/produk?page=2&limit=5"

# Test dengan search
curl "http://localhost:8080/api/produk?name=kopi&page=1&limit=20"
```

### Test Redis Caching

```bash
# Request 1 (Cache MISS - dari database)
curl http://localhost:8080/api/produk
# Check logs: "📦 Cache SET: products:list:..."

# Request 2 (Cache HIT - dari Redis)
curl http://localhost:8080/api/produk
# Check logs: "✅ Cache HIT: products:list:..."

# Create new product (invalidate cache)
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{"nama":"Test","harga":1000,"stok":10}'
# Check logs: "🗑️ Cache DELETE Pattern: products:list:*"

# Request 3 (Cache MISS lagi - cache sudah dihapus)
curl http://localhost:8080/api/produk
# Check logs: "📦 Cache SET: products:list:..."
```

### Monitor Redis

```bash
# Monitor semua commands di Redis
redis-cli monitor

# Check semua keys
redis-cli keys "*"

# Check specific key
redis-cli get "products:list:search::page:1:limit:10"

# Clear all cache
redis-cli flushall
```

## 📊 Cache Invalidation Strategy

| Operation | Cache Invalidation |
|-----------|-------------------|
| **GET** `/api/produk` | No invalidation (read only) |
| **GET** `/api/produk/{id}` | No invalidation (read only) |
| **POST** `/api/produk` | Delete `products:list:*` |
| **PUT** `/api/produk/{id}` | Delete `products:detail:id:{id}` + `products:list:*` |
| **DELETE** `/api/produk/{id}` | Delete `products:detail:id:{id}` + `products:list:*` |

## 🔧 Configuration

### Cache TTL

Edit `services/cache_service.go`:

```go
func NewCacheService() *CacheService {
	return &CacheService{
		defaultTTL: 5 * time.Minute, // Ubah sesuai kebutuhan
	}
}
```

### Pagination Limits

Edit `models/pagination.go`:

```go
func NewPaginationParams(page, limit int) PaginationParams {
	// Default limit = 10, max = 100
	if limit < 1 {
		limit = 10  // Ubah default
	}
	if limit > 100 {
		limit = 100 // Ubah max
	}
	// ...
}
```

## 🎯 Best Practices

### 1. Cache Strategy
- ✅ Cache data yang sering dibaca (GET)
- ✅ Invalidate cache setelah write (POST/PUT/DELETE)
- ✅ Gunakan TTL yang reasonable (5-10 menit)
- ❌ Jangan cache data yang sangat dinamis

### 2. Pagination
- ✅ Gunakan limit yang reasonable (10-50)
- ✅ Set max limit untuk prevent abuse
- ✅ Include pagination metadata di response
- ❌ Jangan return semua data tanpa pagination

### 3. Error Handling
- ✅ Graceful degradation jika Redis down
- ✅ Log cache hits/misses untuk monitoring
- ✅ Handle edge cases (page > total_pages)

## 📝 Migration Notes

### Breaking Changes
⚠️ **Response format berubah untuk GET /api/produk**

**Before:**
```json
[
  {"id": 1, "nama": "Product 1"},
  {"id": 2, "nama": "Product 2"}
]
```

**After:**
```json
{
  "data": [
    {"id": 1, "nama": "Product 1"},
    {"id": 2, "nama": "Product 2"}
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 2,
    "total_pages": 1
  }
}
```

### Frontend Update Required

Update frontend code untuk handle new response format:

```javascript
// Before
const products = await response.json();

// After
const result = await response.json();
const products = result.data;
const pagination = result.pagination;
```

## 🐛 Troubleshooting

### Redis Connection Failed
```
Warning: Gagal connect ke Redis, Redis caching akan dinonaktifkan
```
**Solution**: Check Redis is running, check REDIS_URL

### Cache Not Working
**Check**:
1. Redis is running: `redis-cli ping`
2. REDIS_URL is set correctly
3. Check logs for cache operations

### Pagination Not Working
**Check**:
1. Query params format: `?page=1&limit=10`
2. Page and limit are positive integers
3. Check response format includes `pagination` field

## 📚 References

- [Redis Documentation](https://redis.io/documentation)
- [go-redis Client](https://github.com/redis/go-redis)
- [Pagination Best Practices](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/)
