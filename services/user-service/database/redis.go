package database

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"user-service/config"
)

// InitRedis 初始化 Redis 連接
func InitRedis(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	})
	log.Println("Redis client initialized")
	return client
}
