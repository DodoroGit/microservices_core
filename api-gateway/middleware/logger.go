package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 記錄每個進入的請求，包含 method、path、status code 與處理耗時。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf("[Gateway] %s %s | status=%d | latency=%s | ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
		)
	}
}
