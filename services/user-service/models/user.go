package models

import "time"

// User 代表用戶資料模型
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 不在 JSON 中顯示
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest 註冊請求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登入請求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用戶請求
type UpdateUserRequest struct {
	Username string `json:"username"`
}

// LoginResponse 登入響應
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"` // JWT token，前端後續請求放進 Authorization header
	User    User   `json:"user"`
}
