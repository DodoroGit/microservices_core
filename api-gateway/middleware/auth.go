package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 定義從 JWT payload 中讀取的欄位，需與 user-service 簽發時的結構相同
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RequireAuth 驗證請求的 Authorization header 是否帶有合法的 JWT token。
//
// 驗證流程：
//  1. 確認 header 存在且格式為 "Bearer <token>"
//  2. 解析 token，確認簽章演算法為 HS256
//  3. 用 JWT_SECRET 驗證簽章是否正確
//  4. 確認 token 尚未過期（jwt 套件自動處理）
//  5. 將 user_id 存入 gin.Context，讓後續 handler 可以使用
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── 1. 取出 header ─────────────────────────────────────────────────
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "缺少 Authorization header",
			})
			return
		}

		// ── 2. 確認格式為 "Bearer <token>" ────────────────────────────────
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "格式錯誤，預期為：Bearer <token>",
			})
			return
		}
		tokenString := parts[1]

		// ── 3. 解析並驗證 token ────────────────────────────────────────────
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 確認簽章演算法是預期的 HS256，防止演算法混淆攻擊
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("非預期的簽章演算法：%v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token 無效或已過期",
			})
			return
		}

		// ── 4. 將 user_id 存入 context，後續 handler 可透過 c.GetString("user_id") 取得
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
