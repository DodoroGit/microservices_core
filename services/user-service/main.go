package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"user-service/config"
	"user-service/database"
	"user-service/handlers"
	"user-service/repository"
	"user-service/routes"
	"user-service/services"
)

func main() {
	// 載入配置
	cfg := config.Load()

	// 初始化資料庫
	db, err := database.InitPostgres(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 創建資料表
	if err := database.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	// 初始化 Redis
	redisClient := database.InitRedis(cfg.Redis)
	defer redisClient.Close()

	// 初始化各層
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, cfg.JWTSecret)

	// 設定路由
	router := gin.Default()
	routes.SetupRoutes(router, userHandler)

	// 啟動服務
	log.Printf("User Service starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
