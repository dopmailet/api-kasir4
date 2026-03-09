# 🚀 Quick Start - Testing Production API

## Production URL
**Base URL:** `https://kasir-api-production-zainalzxc.up.railway.app`

---

## ⚡ Quick Test (No Postman Required)

### **1. Health Check**
```bash
curl https://kasir-api-production-zainalzxc.up.railway.app/health
```

**Expected:**
```json
{
  "status": "OK",
  "message": "API Running"
}
```

### **2. Get All Products**
```bash
curl https://kasir-api-production-zainalzxc.up.railway.app/api/produk
```

### **3. Create Product**
```bash
curl -X POST https://kasir-api-production-zainalzxc.up.railway.app/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Kopi Susu",
    "harga": 15000,
    "stok": 50
  }'
```

### **4. Checkout**
```bash
curl -X POST https://kasir-api-production-zainalzxc.up.railway.app/api/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2}
    ]
  }'
```

---

## 📮 Postman Setup

### **Step 1: Import Collection**
1. Buka Postman
2. Import file: `Kasir-API.postman_collection.json`

### **Step 2: Import Environment**
1. Import file: `Kasir-API-Production.postman_environment.json`
2. Import file: `Kasir-API-Local.postman_environment.json`

### **Step 3: Select Environment**
1. Klik dropdown di pojok kanan atas
2. Pilih **"Kasir API - Production"**
3. Done! ✅

### **Step 4: Test!**
1. Pilih request **"Health Check"**
2. Klik **"Send"**
3. Verify response OK

---

## 🧪 Testing Checklist

### **Production Deployment Verification**

- [ ] ✅ Health check returns OK
- [ ] ✅ Database connection working
- [ ] ✅ CRUD operations berfungsi
- [ ] ✅ Checkout flow berfungsi
- [ ] ✅ Report endpoint berfungsi
- [ ] ✅ HTTPS working (secure connection)
- [ ] ✅ No prepared statement errors
- [ ] ✅ Response time acceptable (< 500ms)

---

## 🎯 Test Scenarios untuk Production

### **Scenario 1: Basic Functionality**

```bash
# 1. Health check
curl https://kasir-api-production-zainalzxc.up.railway.app/health

# 2. Create product
curl -X POST https://kasir-api-production-zainalzxc.up.railway.app/api/produk \
  -H "Content-Type: application/json" \
  -d '{"nama":"Kopi","harga":10000,"stok":100}'

# 3. Get all products
curl https://kasir-api-production-zainalzxc.up.railway.app/api/produk

# 4. Checkout
curl -X POST https://kasir-api-production-zainalzxc.up.railway.app/api/checkout \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":2}]}'

# 5. Check report
curl https://kasir-api-production-zainalzxc.up.railway.app/api/report/hari-ini
```

### **Scenario 2: Performance Test**

Test batch INSERT optimization dengan 10 items:

```bash
curl -X POST https://kasir-api-production-zainalzxc.up.railway.app/api/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 1},
      {"product_id": 2, "quantity": 2},
      {"product_id": 3, "quantity": 1},
      {"product_id": 4, "quantity": 3},
      {"product_id": 5, "quantity": 1},
      {"product_id": 6, "quantity": 2},
      {"product_id": 7, "quantity": 1},
      {"product_id": 8, "quantity": 4},
      {"product_id": 9, "quantity": 1},
      {"product_id": 10, "quantity": 2}
    ]
  }'
```

---

## 📊 Benchmark Production

### **Using `hey` (HTTP load generator)**

Install hey:
```bash
# Windows (with Chocolatey)
choco install hey

# Or download from: https://github.com/rakyll/hey/releases
```

Test checkout performance:
```bash
hey -n 100 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":1}]}' \
  https://kasir-api-production-zainalzxc.up.railway.app/api/checkout
```

**Expected Results:**
- Total requests: 100
- Concurrency: 10
- Success rate: 100%
- Average response time: < 500ms

---

## 🔍 Monitoring & Debugging

### **Check Railway Logs**

1. Go to: https://railway.app
2. Open project: **kasir-api-production**
3. Click **"Deployments"**
4. Click **"View Logs"**

### **Common Issues**

#### **Issue 1: Connection Timeout**
**Cause:** Railway cold start  
**Solution:** First request might be slow (5-10s), subsequent requests will be fast

#### **Issue 2: Database Connection Error**
**Cause:** Database not configured  
**Solution:** Check `DB_CONN` environment variable di Railway

#### **Issue 3: 404 Not Found**
**Cause:** Wrong endpoint  
**Solution:** Verify URL: `https://kasir-api-production-zainalzxc.up.railway.app/api/...`

---

## 🎥 Video Demo Preparation (Sesi 4)

### **Checklist untuk Video Demo:**

1. **Show Health Check** ✅
   ```
   GET /health → Status OK
   ```

2. **Show CRUD Operations** ✅
   ```
   POST /api/produk → Create product
   GET /api/produk → List products
   PUT /api/produk/1 → Update product
   DELETE /api/produk/1 → Delete product
   ```

3. **Show Checkout Flow** ✅
   ```
   GET /api/produk → Check stock
   POST /api/checkout → Process checkout
   GET /api/produk → Verify stock berkurang
   ```

4. **Show Reports** ✅
   ```
   GET /api/report/hari-ini → Daily report
   ```

5. **Show Performance** ✅
   ```
   POST /api/checkout (10 items) → Fast response
   ```

### **Recording Tips:**

- Use **Postman** untuk visual yang bagus
- Show **response time** di Postman
- Demonstrate **batch INSERT** dengan 10 items
- Show **no errors** pada concurrent requests
- Highlight **optimasi** yang sudah dilakukan

---

## 📝 Production URL Reference

### **All Endpoints:**

```
Base URL: https://kasir-api-production-zainalzxc.up.railway.app

Health:
  GET  /health

Products:
  GET    /api/produk
  GET    /api/produk/{id}
  GET    /api/produk?nama={search}
  POST   /api/produk
  PUT    /api/produk/{id}
  DELETE /api/produk/{id}

Categories:
  GET    /api/categories
  GET    /api/categories/{id}
  POST   /api/categories
  PUT    /api/categories/{id}
  DELETE /api/categories/{id}

Checkout:
  POST   /api/checkout

Reports:
  GET    /api/report/hari-ini
  GET    /api/report?start_date={date}&end_date={date}
```

---

## 🎉 Production Ready!

Your API is now live at:
**https://kasir-api-production-zainalzxc.up.railway.app**

Features:
- ✅ HTTPS enabled (secure)
- ✅ Database connected (PostgreSQL)
- ✅ Batch INSERT optimization
- ✅ No prepared statement errors
- ✅ Ready for Sesi 4 competition!

**Good luck!** 🚀
