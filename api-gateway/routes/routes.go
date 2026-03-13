package routes

import (
	"net/http"
	"strings"
	"time"

	"api-gateway/config"
	"api-gateway/middleware"
	"api-gateway/proxy"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const aiTimeout = 120 * time.Second

// Setup 將所有 middleware 與路由掛載到 Gin engine 上。
func Setup(r *gin.Engine, cfg *config.Config, redisClient *redis.Client) {
	p := proxy.New()
	aiProxy := proxy.NewWithTimeout(aiTimeout)

	// ── 全域 Middleware ──────────────────────────────────────────────────────
	r.Use(gin.Recovery()) // 攔截 panic，回傳 500，避免整個服務崩潰
	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") ||
				strings.HasSuffix(origin, ".ngrok-free.app") ||
				strings.HasSuffix(origin, ".ngrok.io")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
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
		public.POST("/users/login", p.Forward(cfg.UserServiceURL, "/api"))
		public.POST("/users/register", p.Forward(cfg.UserServiceURL, "/api"))
	}

	// 受保護路由：需要帶 Bearer token（透過 middleware/auth.go 驗證）
	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret, redisClient))
	{
		protected.POST("/users/logout", p.Forward(cfg.UserServiceURL, "/api"))
		protected.GET("/users", p.Forward(cfg.UserServiceURL, "/api"))
		protected.GET("/users/:id", p.Forward(cfg.UserServiceURL, "/api"))
		protected.PUT("/users/:id", p.Forward(cfg.UserServiceURL, "/api"))
		protected.DELETE("/users/:id", p.Forward(cfg.UserServiceURL, "/api"))
	}

	// AI Service 路由：需要登入，120 秒 timeout（LLM 推理較慢）
	ai := r.Group("/api/ai")
	ai.Use(middleware.RequireAuth(cfg.JWTSecret, redisClient))
	{
		ai.POST("/background", aiProxy.Forward(cfg.AIServiceURL, "/api"))
		ai.POST("/projects", aiProxy.Forward(cfg.AIServiceURL, "/api"))
		ai.POST("/skills", aiProxy.Forward(cfg.AIServiceURL, "/api"))
		ai.POST("/daily", aiProxy.Forward(cfg.AIServiceURL, "/api"))
		ai.POST("/project-story", aiProxy.Forward(cfg.AIServiceURL, "/api"))
	}

	// Note Service 路由 - 所有登入用戶可讀
	notes := r.Group("/api")
	notes.Use(middleware.RequireAuth(cfg.JWTSecret, redisClient))
	{
		notes.GET("/notes", p.Forward(cfg.NoteServiceURL, "/api"))
		notes.GET("/notes/:id", p.Forward(cfg.NoteServiceURL, "/api"))
	}

	// Admin-only 路由：需要 admin 角色
	admin := r.Group("/api")
	admin.Use(middleware.RequireAuth(cfg.JWTSecret, redisClient))
	admin.Use(middleware.RequireAdmin())
	{
		admin.PATCH("/users/:id/role", p.Forward(cfg.UserServiceURL, "/api"))
		admin.POST("/notes", p.Forward(cfg.NoteServiceURL, "/api"))
		admin.PUT("/notes/:id", p.Forward(cfg.NoteServiceURL, "/api"))
		admin.DELETE("/notes/:id", p.Forward(cfg.NoteServiceURL, "/api"))
	}
}
