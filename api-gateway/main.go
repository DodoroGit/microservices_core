package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"api-gateway/config"
	"api-gateway/routes"
)

func main() {
	// 讀取設定（port、下游服務 URL）
	cfg := config.Load()

	// 初始化 Redis（用於 token 黑名單檢查）
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
	defer redisClient.Close()

	// 使用 gin.New() 而非 gin.Default()，
	// 因為 Recovery 與 Logger 已在 routes.Setup 中手動掛載，避免重複。
	router := gin.New()
	routes.Setup(router, cfg, redisClient)

	log.Printf("API Gateway 啟動，監聽 port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("API Gateway 啟動失敗：", err)
	}
}
