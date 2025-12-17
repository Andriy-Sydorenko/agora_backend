package database

import (
	"context"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var redisClient *redis.Client

func ConnectRedisClient(cfg *config.RedisConfig) *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalln("failed to connect to Redis:", err)
		return nil
	}

	log.Println("✅ Redis connected successfully")
	return redisClient
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatalln("Redis client not initialized. Call ConnectRedis() first")
		return nil
	}
	return redisClient
}
