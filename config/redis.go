package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Membuat konteks untuk operasi Redis
	ctx := context.Background()

	// Memeriksa koneksi dengan perintah Ping
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connect to Redis:\n%v", err)
	}

	log.Println("Redis connection success")
	return rdb
}
