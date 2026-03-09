# Upstash Redis Setup Guide (FREE)

## 🔴 Apa itu Upstash?

Upstash adalah **serverless Redis** dengan free tier yang generous:
- ✅ 10,000 commands per day (FREE)
- ✅ Global edge network (fast)
- ✅ No credit card required
- ✅ Perfect untuk production apps

---

## 📦 Setup Upstash Redis

### **Step 1: Create Account**

1. Buka https://upstash.com
2. Click **"Get Started"**
3. Sign up dengan GitHub/Google/Email
4. Verify email (jika perlu)

### **Step 2: Create Redis Database**

1. **Dashboard** → Click **"Create Database"**

2. **Configure Database:**
   ```
   Name: kasir-api-cache
   Type: Regional (pilih region terdekat, e.g., ap-southeast-1)
   Primary Region: Singapore (atau terdekat dengan user Anda)
   Read Regions: (optional, biarkan kosong untuk free tier)
   Eviction: No eviction (recommended)
   ```

3. Click **"Create"**

4. Wait ~30 seconds untuk database ready

### **Step 3: Get Redis URL**

1. Click database yang baru dibuat
2. Scroll ke **"REST API"** section
3. Copy **"UPSTASH_REDIS_REST_URL"** atau
4. Scroll ke **"Connection"** → Copy **Redis URL**

   Format:
   ```
   redis://default:YOUR_PASSWORD@YOUR_ENDPOINT.upstash.io:6379
   ```

   Example:
   ```
   redis://default:AYNxASQgYmU5ZjE3YjYt...@gusc1-merry-firefly-12345.upstash.io:6379
   ```

### **Step 4: Set Environment Variable**

#### **Local Development (.env)**
```bash
# .env
DATABASE_URL=postgresql://...your-supabase-url...
REDIS_URL=redis://default:YOUR_PASSWORD@YOUR_ENDPOINT.upstash.io:6379
PORT=8080
```

#### **Railway Deployment**
1. Go to Railway dashboard
2. Select your service
3. Go to **"Variables"** tab
4. Click **"New Variable"**
5. Add:
   ```
   Variable: REDIS_URL
   Value: redis://default:YOUR_PASSWORD@YOUR_ENDPOINT.upstash.io:6379
   ```
6. Click **"Add"**
7. Railway will auto-redeploy

#### **Vercel Deployment**
1. Go to Vercel dashboard
2. Select your project
3. Go to **"Settings"** → **"Environment Variables"**
4. Add:
   ```
   Name: REDIS_URL
   Value: redis://default:YOUR_PASSWORD@YOUR_ENDPOINT.upstash.io:6379
   Environment: Production, Preview, Development
   ```
5. Click **"Save"**
6. Redeploy

### **Step 5: Test Connection**

#### **Local Test**
```bash
# Jalankan app
go run main.go

# Check logs
# Should see: ✅ Redis connected successfully!
```

#### **Production Test**
```bash
# Check deployment logs
# Should see: ✅ Redis connected successfully!

# Test API
curl https://your-app.com/api/produk
# Check logs for: 📦 Cache SET

curl https://your-app.com/api/produk
# Check logs for: ✅ Cache HIT
```

---

## 📊 Monitor Redis Usage

### **Upstash Dashboard**

1. Go to https://console.upstash.com
2. Click your database
3. View **"Metrics"**:
   - Commands per day
   - Storage used
   - Latency

### **Free Tier Limits**
```
Daily Commands: 10,000 (resets every 24h)
Storage: 256 MB
Connections: 1,000
```

**Apakah cukup?**
- ✅ Small apps (< 1000 users/day): **YES**
- ✅ Medium apps (< 5000 users/day): **YES** (with good TTL)
- ⚠️ Large apps (> 10000 users/day): Consider paid tier

---

## 🔧 Optimize for Free Tier

### **1. Set Appropriate TTL**

Edit `services/cache_service.go`:
```go
func NewCacheService() *CacheService {
	return &CacheService{
		defaultTTL: 10 * time.Minute, // Increase to 10 minutes
	}
}
```

**Why?** Longer TTL = fewer cache refreshes = fewer commands

### **2. Cache Only Popular Endpoints**

Only cache frequently accessed data:
- ✅ GET /api/produk (list)
- ✅ GET /api/produk/{id} (detail)
- ❌ Don't cache reports (changes frequently)

### **3. Use Cache Patterns Wisely**

```go
// Good: Specific cache keys
products:list:page:1:limit:10

// Bad: Too many variations
products:list:search:a:page:1:limit:10
products:list:search:b:page:1:limit:10
// ... creates too many keys
```

---

## 🚀 Production Checklist

Before deploying to production with Upstash:

- [ ] Upstash account created
- [ ] Redis database created (region: closest to users)
- [ ] REDIS_URL copied
- [ ] Environment variable set (Railway/Vercel)
- [ ] Local test passed (✅ Redis connected)
- [ ] Production deployment successful
- [ ] Cache HIT/MISS working in logs
- [ ] Monitor Upstash dashboard (commands/day)

---

## 🐛 Troubleshooting

### **Error: "dial tcp: i/o timeout"**

**Cause**: Network issue or wrong URL

**Solution**:
1. Check REDIS_URL format
2. Ensure no firewall blocking port 6379
3. Try REST API instead (Upstash supports both)

### **Error: "NOAUTH Authentication required"**

**Cause**: Missing password in URL

**Solution**:
```bash
# Wrong
redis://your-endpoint.upstash.io:6379

# Correct
redis://default:YOUR_PASSWORD@your-endpoint.upstash.io:6379
```

### **Warning: "Gagal connect ke Redis"**

**Cause**: Redis URL not set or invalid

**Solution**:
1. Check `.env` file has `REDIS_URL`
2. Restart app: `go run main.go`
3. Check logs for connection success

### **Commands Limit Exceeded**

**Cause**: More than 10,000 commands/day

**Solution**:
1. Increase TTL (cache longer)
2. Reduce cache keys (fewer variations)
3. Upgrade to paid tier ($10/month for 100K commands)

---

## 💡 Best Practices

### **1. Connection Pooling**
Already implemented in `config/redis.go`:
```go
opt, err := redis.ParseURL(redisURL)
// Uses connection pooling by default
```

### **2. Graceful Degradation**
Already implemented:
```go
if config.RedisClient == nil {
    return false // Skip cache, use DB
}
```

### **3. Monitor Performance**
Check Upstash dashboard weekly:
- Commands usage (should be < 10K/day)
- Latency (should be < 50ms)
- Storage (should be < 256MB)

---

## 📚 Resources

- **Upstash Docs**: https://docs.upstash.com/redis
- **Upstash Console**: https://console.upstash.com
- **Pricing**: https://upstash.com/pricing
- **Status**: https://status.upstash.com

---

## 🎉 Summary

**Upstash Redis** adalah pilihan perfect untuk:
- ✅ Development & testing
- ✅ Small to medium production apps
- ✅ Apps dengan budget terbatas
- ✅ Serverless deployments (Vercel, Railway)

**Setup time**: ~5 minutes
**Cost**: FREE (up to 10K commands/day)
**Reliability**: 99.99% uptime SLA

**Selamat mencoba!** 🚀
