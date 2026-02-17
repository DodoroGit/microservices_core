package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/handlers"
)

// SetupRoutes 設定所有路由
func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	// 健康檢查
	router.GET("/health", userHandler.Health)

	// 用戶路由
	router.POST("/users/register", userHandler.Register)
	router.POST("/users/login", userHandler.Login)
	router.GET("/users", userHandler.GetUsers)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
}
