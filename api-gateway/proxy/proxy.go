package proxy

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Proxy 持有一個共用的 HTTP client，用來將請求轉發給下游服務。
// 共用同一個 client 是為了讓 TCP connection pool 能夠被重複利用，避免每次請求都重新建立連線。
type Proxy struct {
	client *http.Client
}

func New() *Proxy {
	return &Proxy{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Forward 回傳一個 Gin handler，將收到的請求轉發到 targetBaseURL，
// 並在轉發前將 pathPrefix 從路徑中去除。
//
// 範例：
//
//	收到請求：GET  /api/users/123
//	去除前綴：/api
//	轉發目標：GET  http://user-service:8081/users/123
func (p *Proxy) Forward(targetBaseURL, pathPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── 1. 重寫路徑：去掉 gateway 前綴 ────────────────────────────────
		servicePath := strings.TrimPrefix(c.Request.URL.Path, pathPrefix)
		targetURL := targetBaseURL + servicePath
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		// ── 2. 讀取請求 body ───────────────────────────────────────────────
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// ── 3. 建立對下游服務的新請求 ──────────────────────────────────────
		req, err := http.NewRequestWithContext(
			c.Request.Context(),
			c.Request.Method,
			targetURL,
			bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "建立請求失敗"})
			return
		}

		// 轉送原始 headers（Authorization、Content-Type 等）
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// ── 4. 發送請求到下游服務 ──────────────────────────────────────────
		resp, err := p.client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "下游服務無法連線"})
			return
		}
		defer resp.Body.Close()

		// ── 5. 將下游的 response 原封不動回傳給前端 ────────────────────────
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取回應失敗"})
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
	}
}
