package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"user-service/models"
	"user-service/services"
)

// Claims 定義 JWT payload 的內容
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// UserHandler 用戶 HTTP 處理層
type UserHandler struct {
	service   services.UserServiceInterface
	jwtSecret string
}

// NewUserHandler 創建用戶 Handler
func NewUserHandler(service services.UserServiceInterface, jwtSecret string) *UserHandler {
	return &UserHandler{service: service, jwtSecret: jwtSecret}
}

// Register 註冊處理
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login 登入處理：驗證帳密，成功後簽發 JWT token
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ── 產生 JWT token ────────────────────────────────────────────────────
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 1 小時後過期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "產生 token 失敗"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Message: "Login successful",
		Token:   tokenString,
		User:    *user,
	})
}

// GetUsers 獲取用戶列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser 獲取單個用戶
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用戶
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateUser(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser 刪除用戶
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateRole 更新用戶角色
func (h *UserHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateRole(id, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// Logout 登出：解析 token 取得剩餘有效時間，交由 service 寫入黑名單
func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 token"})
		return
	}
	tokenString := parts[1]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 token"})
		return
	}

	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
		return
	}

	if err := h.service.Logout(tokenString, remaining); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// Health 健康檢查
func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "user-service",
	})
}
