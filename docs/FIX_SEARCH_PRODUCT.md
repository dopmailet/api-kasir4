# Fix: Search Product by Name

## 🐛 Bug yang Diperbaiki

**Problem:**
```
GET /api/produk?nama=te
```
Menampilkan **semua produk** yang mengandung "te" di mana saja:
- ✅ "Teh Manis" 
- ✅ "Es Teh"
- ✅ "Kopi Latte" ← Tidak seharusnya muncul!

**Expected:**
Hanya menampilkan produk yang **dimulai dengan** "te":
- ✅ "Teh Manis"
- ❌ "Es Teh"
- ❌ "Kopi Latte"

---

## 🔍 Penyebab

### **Sebelum (Bug):**
```go
// Search pattern: %searchName%
args = append(args, "%"+searchName+"%")
```

**SQL Query:**
```sql
WHERE p.nama ILIKE '%te%'
```

**Behavior:** **CONTAINS** - match di mana saja dalam string
- "**Te**h Manis" ✅
- "Es **Te**h" ✅
- "Kopi Lat**te**" ✅ ← Bug!

---

## ✅ Solusi

### **Sesudah (Fixed):**
```go
// Search pattern: searchName%
args = append(args, searchName+"%")
```

**SQL Query:**
```sql
WHERE p.nama ILIKE 'te%'
```

**Behavior:** **STARTS WITH** (prefix match) - hanya match di awal string
- "**Te**h Manis" ✅
- "Es Teh" ❌
- "Kopi Latte" ❌

---

## 📊 Perbandingan Search Patterns

| Pattern | SQL | Behavior | Example Match |
|---------|-----|----------|---------------|
| `%te%` | `ILIKE '%te%'` | **Contains** (anywhere) | "Teh", "Es Teh", "Latte" |
| `te%` | `ILIKE 'te%'` | **Starts with** (prefix) | "Teh", "Teh Manis" |
| `%te` | `ILIKE '%te'` | **Ends with** (suffix) | "Latte", "Chocolate" |
| `te` | `ILIKE 'te'` | **Exact match** | "te" only |

---

## 🧪 Testing

### **Test Case 1: Search "te"**

**Request:**
```bash
GET /api/produk?nama=te
```

**Expected Results:**
```json
[
  {
    "id": 1,
    "nama": "Teh Manis",
    "harga": 8000,
    "stok": 100
  },
  {
    "id": 2,
    "nama": "Teh Tarik",
    "harga": 10000,
    "stok": 50
  }
]
```

**Should NOT include:**
- "Es Teh" (tidak dimulai dengan "te")
- "Kopi Latte" (tidak dimulai dengan "te")

---

### **Test Case 2: Search "ko"**

**Request:**
```bash
GET /api/produk?nama=ko
```

**Expected Results:**
```json
[
  {
    "id": 3,
    "nama": "Kopi Hitam",
    "harga": 10000,
    "stok": 100
  },
  {
    "id": 4,
    "nama": "Kopi Susu",
    "harga": 15000,
    "stok": 50
  }
]
```

**Should NOT include:**
- "Chocolate" (tidak dimulai dengan "ko")

---

### **Test Case 3: Case Insensitive**

**Request:**
```bash
GET /api/produk?nama=TE
GET /api/produk?nama=Te
GET /api/produk?nama=te
```

**Expected:** Semua 3 request di atas return hasil yang **sama** karena menggunakan `ILIKE` (case-insensitive)

---

## 💡 Alternative: Support Multiple Search Modes

Jika di masa depan kamu mau support berbagai mode search, bisa tambahkan query parameter:

### **Option 1: Search Mode Parameter**

```go
// GET /api/produk?nama=te&mode=starts_with
// GET /api/produk?nama=te&mode=contains
// GET /api/produk?nama=te&mode=exact

searchMode := r.URL.Query().Get("mode")

switch searchMode {
case "exact":
    args = append(args, searchName)
case "contains":
    args = append(args, "%"+searchName+"%")
case "starts_with":
    fallthrough
default:
    args = append(args, searchName+"%")
}
```

### **Option 2: Wildcard Support**

```go
// User bisa pakai wildcard sendiri
// GET /api/produk?nama=te%      → starts with
// GET /api/produk?nama=%te%     → contains
// GET /api/produk?nama=%te      → ends with

// Jika user tidak pakai wildcard, default ke starts_with
if !strings.Contains(searchName, "%") {
    searchName = searchName + "%"
}
args = append(args, searchName)
```

---

## 📝 Best Practices

### **Untuk Search Product:**

1. **Default: Starts With** ✅
   - Paling umum untuk search product
   - User ketik "ko" → expect "Kopi", bukan "Chocolate"
   - Lebih cepat (bisa pakai index)

2. **Case Insensitive** ✅
   - Gunakan `ILIKE` bukan `LIKE`
   - User tidak perlu exact case

3. **Trim Whitespace** ✅
   ```go
   searchName = strings.TrimSpace(searchName)
   ```

4. **Add Index** ✅
   ```sql
   CREATE INDEX idx_products_nama ON products(nama);
   ```

---

## 🚀 Performance Tips

### **1. Use Index for Prefix Search**

Prefix search (`te%`) bisa pakai index:
```sql
CREATE INDEX idx_products_nama ON products(nama);
```

### **2. Avoid Leading Wildcard**

Leading wildcard (`%te`) **tidak bisa pakai index** → slow!
```sql
-- ❌ Slow (full table scan)
WHERE nama ILIKE '%te%'

-- ✅ Fast (can use index)
WHERE nama ILIKE 'te%'
```

### **3. Full-Text Search (Advanced)**

Untuk search yang lebih advanced, gunakan PostgreSQL Full-Text Search:
```sql
-- Add tsvector column
ALTER TABLE products ADD COLUMN search_vector tsvector;

-- Create index
CREATE INDEX idx_products_search ON products USING GIN(search_vector);

-- Update search vector
UPDATE products SET search_vector = to_tsvector('indonesian', nama);

-- Search
SELECT * FROM products 
WHERE search_vector @@ to_tsquery('indonesian', 'kopi');
```

---

## ✅ Summary

**Fixed:**
- ✅ Search sekarang pakai **prefix match** (`te%`)
- ✅ Hanya return products yang **dimulai dengan** search term
- ✅ Case insensitive tetap work
- ✅ Lebih akurat dan sesuai ekspektasi user

**Before:**
```
GET /api/produk?nama=te
→ Returns: "Teh", "Es Teh", "Latte" (semua yang mengandung "te")
```

**After:**
```
GET /api/produk?nama=te
→ Returns: "Teh Manis", "Teh Tarik" (hanya yang dimulai dengan "te")
```

**Perfect!** ✅
