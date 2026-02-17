package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	userServiceURL = getEnv("USER_SERVICE_URL", "http://localhost:8081")
	port           = getEnv("PORT", "8080")
)

func main() {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "api-gateway",
		})
	})

	// User service routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("", proxyRequest(userServiceURL, "/users"))
		userRoutes.GET("/:id", proxyRequest(userServiceURL, "/users"))
		userRoutes.POST("", proxyRequest(userServiceURL, "/users"))
		userRoutes.PUT("/:id", proxyRequest(userServiceURL, "/users"))
		userRoutes.DELETE("/:id", proxyRequest(userServiceURL, "/users"))
		userRoutes.POST("/login", proxyRequest(userServiceURL, "/users/login"))
		userRoutes.POST("/register", proxyRequest(userServiceURL, "/users/register"))
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func proxyRequest(serviceURL, basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build target URL
		targetURL := serviceURL + basePath
		if c.Param("id") != "" {
			targetURL += "/" + c.Param("id")
		}

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// Create new request
		req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Send request
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Read response
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Forward response
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
