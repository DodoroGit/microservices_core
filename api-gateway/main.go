package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"api-gateway/config"
	"api-gateway/routes"
)

func main() {
	// 讀取設定（port、下游服務 URL）
	cfg := config.Load()

	// 使用 gin.New() 而非 gin.Default()，
	// 因為 Recovery 與 Logger 已在 routes.Setup 中手動掛載，避免重複。
	router := gin.New()
	routes.Setup(router, cfg)

	log.Printf("API Gateway 啟動，監聽 port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("API Gateway 啟動失敗：", err)
	}
}
