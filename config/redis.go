package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// RedisClient adalah instance global Redis client
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis menginisialisasi koneksi ke Redis
func InitRedis() {
	// Ambil Redis URL dari environment variable
	// Format: redis://username:password@host:port/db
	// Contoh: redis://localhost:6379/0
	redisURL := os.Getenv("REDIS_URL")

	// Jika REDIS_URL tidak di-set, gunakan default localhost
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
		log.Println("REDIS_URL tidak di-set, menggunakan default: redis://localhost:6379/0")
	}

	// Parse Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: Gagal parse REDIS_URL, Redis caching akan dinonaktifkan. Error: %v", err)
		RedisClient = nil
		return
	}

	// Buat Redis client
	RedisClient = redis.NewClient(opt)

	// Test koneksi dengan ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Gagal connect ke Redis, Redis caching akan dinonaktifkan. Error: %v", err)
		RedisClient = nil
		return
	}

	log.Println("✅ Redis connected successfully!")
}

// CloseRedis menutup koneksi Redis
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis connection closed")
	}
}
