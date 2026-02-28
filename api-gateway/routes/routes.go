package routes

import (
	"net/http"
	"time"

	"api-gateway/config"
	"api-gateway/middleware"
	"api-gateway/proxy"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup 將所有 middleware 與路由掛載到 Gin engine 上。
func Setup(r *gin.Engine, cfg *config.Config) {
	p := proxy.New()

	// ── 全域 Middleware ──────────────────────────────────────────────────────
	r.Use(gin.Recovery()) // 攔截 panic，回傳 500，避免整個服務崩潰
	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ── Health Check ─────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "api-gateway",
		})
	})

	// ── User Service 路由 ─────────────────────────────────────────────────────
	//
	// 公開路由：不需要驗證身份（登入、註冊不可能先有 token）
	public := r.Group("/api")
	{
		public.POST("/users/login",    p.Forward(cfg.UserServiceURL, "/api"))
		public.POST("/users/register", p.Forward(cfg.UserServiceURL, "/api"))
	}

	// 受保護路由：需要帶 Bearer token（透過 middleware/auth.go 驗證）
	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret))
	{
		protected.GET("/users",        p.Forward(cfg.UserServiceURL, "/api"))
		protected.GET("/users/:id",    p.Forward(cfg.UserServiceURL, "/api"))
		protected.POST("/users",       p.Forward(cfg.UserServiceURL, "/api"))
		protected.PUT("/users/:id",    p.Forward(cfg.UserServiceURL, "/api"))
		protected.DELETE("/users/:id", p.Forward(cfg.UserServiceURL, "/api"))
	}
}
