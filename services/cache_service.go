package services

import (
	"context"
	"encoding/json"
	"fmt"
	"kasir-api/config"
	"log"
	"time"
)

// CacheService handles Redis caching operations
type CacheService struct {
	defaultTTL time.Duration // Time To Live untuk cache
}

// NewCacheService creates a new CacheService
func NewCacheService() *CacheService {
	return &CacheService{
		defaultTTL: 5 * time.Minute, // Default cache 5 menit
	}
}

// Get mengambil data dari Redis cache
// key: Redis key
// dest: pointer ke variable untuk menampung hasil (akan di-unmarshal dari JSON)
// Return: true jika data ditemukan di cache, false jika tidak
func (c *CacheService) Get(key string, dest interface{}) bool {
	// Jika Redis client tidak tersedia, return false (cache miss)
	if config.RedisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ambil data dari Redis
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		// Cache miss atau error
		return false
	}

	// Unmarshal JSON ke dest
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		log.Printf("Error unmarshal cache for key %s: %v", key, err)
		return false
	}

	log.Printf("✅ Cache HIT: %s", key)
	return true
}

// Set menyimpan data ke Redis cache
// key: Redis key
// value: data yang akan disimpan (akan di-marshal ke JSON)
// ttl: Time To Live (0 = gunakan default)
func (c *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	// Jika Redis client tidak tersedia, skip caching
	if config.RedisClient == nil {
		return nil
	}

	// Marshal value ke JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshal cache: %v", err)
	}

	// Gunakan default TTL jika ttl = 0
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Simpan ke Redis
	err = config.RedisClient.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		log.Printf("Error set cache for key %s: %v", key, err)
		return err
	}

	log.Printf("📦 Cache SET: %s (TTL: %v)", key, ttl)
	return nil
}

// Delete menghapus cache berdasarkan key
func (c *CacheService) Delete(key string) error {
	// Jika Redis client tidak tersedia, skip
	if config.RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := config.RedisClient.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error delete cache for key %s: %v", key, err)
		return err
	}

	log.Printf("🗑️  Cache DELETE: %s", key)
	return nil
}

// DeletePattern menghapus semua cache yang match dengan pattern
// Contoh: DeletePattern("products:*") akan hapus semua cache yang dimulai dengan "products:"
func (c *CacheService) DeletePattern(pattern string) error {
	// Jika Redis client tidak tersedia, skip
	if config.RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Scan semua keys yang match dengan pattern
	iter := config.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	deletedCount := 0
	for iter.Next(ctx) {
		key := iter.Val()
		err := config.RedisClient.Del(ctx, key).Err()
		if err != nil {
			log.Printf("Error delete cache for key %s: %v", key, err)
		} else {
			deletedCount++
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	log.Printf("🗑️  Cache DELETE Pattern: %s (%d keys deleted)", pattern, deletedCount)
	return nil
}

// GenerateKey membuat cache key dengan format standar
// Contoh: GenerateKey("products", "list", "page:1") -> "products:list:page:1"
func (c *CacheService) GenerateKey(parts ...string) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}
