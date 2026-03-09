# 📮 Postman Collection - Kasir API

## 🚀 Cara Import ke Postman

### **Metode 1: Import File JSON**

1. **Buka Postman**
2. **Klik "Import"** di pojok kiri atas
3. **Pilih "Upload Files"**
4. **Pilih file:** `Kasir-API.postman_collection.json`
5. **Klik "Import"**
6. **Done!** ✅ Collection sudah siap digunakan

### **Metode 2: Drag & Drop**

1. **Buka Postman**
2. **Drag file** `Kasir-API.postman_collection.json` ke window Postman
3. **Done!** ✅

---

## 📁 Struktur Collection

Collection ini berisi **6 folder utama**:

### **1. Health Check** 🏥
- `GET /health` - Check if API is running

### **2. Products** 📦
- `GET /api/produk` - Get all products
- `GET /api/produk/{id}` - Get product by ID
- `GET /api/produk?nama=xxx` - Search product by name
- `POST /api/produk` - Create new product
- `PUT /api/produk/{id}` - Update product
- `DELETE /api/produk/{id}` - Delete product

### **3. Categories** 🏷️
- `GET /api/categories` - Get all categories
- `GET /api/categories/{id}` - Get category by ID
- `POST /api/categories` - Create new category
- `PUT /api/categories/{id}` - Update category
- `DELETE /api/categories/{id}` - Delete category

### **4. Checkout / Transaction** 🛒
- `POST /api/checkout` - Single item checkout
- `POST /api/checkout` - Multiple items checkout
- `POST /api/checkout` - Large order (10 items) - untuk test optimasi

### **5. Reports** 📊
- `GET /api/report/hari-ini` - Daily sales report
- `GET /api/report?start_date=xxx&end_date=xxx` - Date range report
- `GET /api/report` - This month report

### **6. Test Scenarios** 🧪
- **Scenario 1:** Setup Test Data (create 3 products)
- **Scenario 2:** Complete Checkout Flow (check → checkout → verify)

---

## 🎯 Cara Menggunakan

### **Step 1: Start Server**

```bash
# Di terminal
cd c:\Users\Admin\Desktop\kasir-api
go run main.go
```

**Expected output:**
```
✅ Database connected successfully
🚀 Server running on port: 8080
```

### **Step 2: Test Health Check**

1. Buka Postman
2. Pilih folder **"Health Check"**
3. Klik request **"Health Check"**
4. Klik **"Send"**

**Expected response:**
```json
{
  "status": "OK",
  "message": "API Running"
}
```

### **Step 3: Setup Test Data**

Gunakan **Test Scenarios → Scenario 1: Setup Test Data**

Jalankan secara berurutan:
1. Create Product - Kopi
2. Create Product - Teh
3. Create Product - Roti

**Tip:** Klik kanan pada folder → **"Run folder"** untuk run semua sekaligus!

### **Step 4: Test Checkout**

Gunakan **Test Scenarios → Scenario 2: Test Checkout Flow**

Jalankan secara berurutan:
1. Check Products (lihat stok awal)
2. Checkout Order
3. Check Stock After Checkout (verify stok berkurang)
4. Check Daily Report (verify transaksi tercatat)

---

## 🔧 Environment Setup

Collection ini mendukung **2 environment**:

### **1. Production (Railway)** 🚀
- **URL:** `https://kasir-api-production-zainalzxc.up.railway.app`
- **File:** `Kasir-API-Production.postman_environment.json`
- **Use case:** Testing di production server

### **2. Local Development** 💻
- **URL:** `http://localhost:8080`
- **File:** `Kasir-API-Local.postman_environment.json`
- **Use case:** Testing saat development

---

### **Cara Import Environment:**

1. **Buka Postman**
2. **Klik "Import"**
3. **Import kedua file environment:**
   - `Kasir-API-Production.postman_environment.json`
   - `Kasir-API-Local.postman_environment.json`
4. **Done!** ✅

### **Cara Switch Environment:**

1. **Klik dropdown** di pojok kanan atas (biasanya tertulis "No Environment")
2. **Pilih environment:**
   - **"Kasir API - Production"** → untuk test di Railway
   - **"Kasir API - Local"** → untuk test di localhost
3. **Done!** Semua request akan otomatis pakai URL yang sesuai

### **Manual Setup (Alternative):**

Jika tidak mau import environment file, bisa set manual:

1. Klik **"Kasir API - Bootcamp Golang"** collection
2. Pilih tab **"Variables"**
3. Ganti value `base_url`:
   - **Production:** `https://kasir-api-production-zainalzxc.up.railway.app`
   - **Local:** `http://localhost:8080`
4. Klik **"Save"**

---

## 📝 Contoh Request Body

### **Create Product**
```json
{
  "nama": "Kopi Susu",
  "harga": 15000,
  "stok": 50
}
```

### **Update Product**
```json
{
  "nama": "Kopi Susu Premium",
  "harga": 18000,
  "stok": 45
}
```

### **Create Category**
```json
{
  "nama": "Minuman",
  "description": "Kategori untuk semua jenis minuman"
}
```

### **Checkout - Single Item**
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ]
}
```

### **Checkout - Multiple Items**
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 2,
      "quantity": 3
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```

---

## 🧪 Testing Checklist

### **Basic CRUD Testing**

- [ ] ✅ Health check returns OK
- [ ] ✅ Create product berhasil
- [ ] ✅ Get all products menampilkan data
- [ ] ✅ Get product by ID berhasil
- [ ] ✅ Search by name berfungsi
- [ ] ✅ Update product berhasil
- [ ] ✅ Delete product berhasil

### **Checkout Testing**

- [ ] ✅ Checkout single item berhasil
- [ ] ✅ Checkout multiple items berhasil
- [ ] ✅ Stock berkurang setelah checkout
- [ ] ✅ Total amount dihitung dengan benar
- [ ] ✅ Error jika stock tidak cukup
- [ ] ✅ Error jika product tidak ada

### **Report Testing**

- [ ] ✅ Daily report menampilkan transaksi hari ini
- [ ] ✅ Date range report berfungsi
- [ ] ✅ Total penjualan dihitung dengan benar

### **Performance Testing**

- [ ] ✅ Checkout dengan 10 items berhasil (test batch INSERT)
- [ ] ✅ Tidak ada error "prepared statement already exists"
- [ ] ✅ Response time < 100ms untuk single checkout

---

## 🎯 Test Scenarios

### **Scenario 1: Happy Path**

1. Create 3 products (Kopi, Teh, Roti)
2. Get all products → verify 3 products exist
3. Checkout dengan 2 items
4. Check stock → verify stock berkurang
5. Check report → verify transaction recorded

**Expected:** ✅ All pass

### **Scenario 2: Error Handling**

1. Checkout dengan product_id yang tidak ada
   - **Expected:** Error "produk tidak ditemukan"

2. Checkout dengan quantity > stock
   - **Expected:** Error "stok tidak mencukupi"

3. Get product dengan ID yang tidak ada
   - **Expected:** 404 Not Found

### **Scenario 3: Performance Test**

1. Create 10 products
2. Checkout dengan 10 items sekaligus
3. Measure response time

**Expected:** 
- ✅ Response time < 200ms
- ✅ Tidak ada error prepared statement
- ✅ Semua items ter-insert dengan benar

---

## 🚀 Advanced: Run Collection dengan Newman

Newman adalah CLI tool untuk run Postman collection dari terminal.

### **Install Newman:**
```bash
npm install -g newman
```

### **Run Collection:**
```bash
newman run Kasir-API.postman_collection.json
```

### **Run dengan Environment:**
```bash
newman run Kasir-API.postman_collection.json \
  --env-var "base_url=http://localhost:8080"
```

### **Export Results:**
```bash
newman run Kasir-API.postman_collection.json \
  --reporters cli,json \
  --reporter-json-export results.json
```

---

## 📊 Response Examples

### **Success Response - Create Product**
```json
{
  "id": 1,
  "nama": "Kopi Susu",
  "harga": 15000,
  "stok": 50
}
```

### **Success Response - Checkout**
```json
{
  "id": 1,
  "total_amount": 44000,
  "created_at": "2024-02-07T20:30:00Z"
}
```

### **Success Response - Daily Report**
```json
{
  "date": "2024-02-07",
  "total_transactions": 5,
  "total_amount": 250000,
  "transactions": [
    {
      "id": 1,
      "total_amount": 44000,
      "created_at": "2024-02-07T10:30:00Z"
    }
  ]
}
```

### **Error Response - Product Not Found**
```json
{
  "error": "produk dengan ID 999 tidak ditemukan"
}
```

### **Error Response - Insufficient Stock**
```json
{
  "error": "stok produk ID 1 tidak mencukupi (tersedia: 5, diminta: 10)"
}
```

---

## 💡 Tips & Tricks

### **1. Save Responses**
Klik **"Save Response"** untuk save example responses yang bisa dijadikan dokumentasi.

### **2. Use Tests Tab**
Tambahkan test scripts di tab "Tests" untuk automated testing:

```javascript
// Test status code
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

// Test response body
pm.test("Response has id", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('id');
});
```

### **3. Use Pre-request Scripts**
Generate dynamic data di tab "Pre-request Script":

```javascript
// Generate random product name
pm.variables.set("product_name", "Product " + Math.floor(Math.random() * 1000));
```

### **4. Chain Requests**
Save response data untuk request berikutnya:

```javascript
// Save product_id from response
var jsonData = pm.response.json();
pm.environment.set("product_id", jsonData.id);
```

Lalu gunakan di request berikutnya:
```
{{product_id}}
```

---

## 🎉 Kesimpulan

Dengan Postman Collection ini, kamu bisa:

- ✅ **Test semua endpoint** dengan mudah
- ✅ **Tidak perlu ketik manual** request body
- ✅ **Run test scenarios** secara otomatis
- ✅ **Share dengan tim** (export/import collection)
- ✅ **Dokumentasi API** yang interaktif

**Happy Testing!** 🚀
