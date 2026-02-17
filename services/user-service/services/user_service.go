package services

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"user-service/models"
	"user-service/repository"
)

// UserService 用戶業務邏輯層
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 創建用戶 Service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register 註冊新用戶
func (s *UserService) Register(req models.RegisterRequest) (*models.User, error) {
	// 檢查 email 是否已存在
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 創建用戶
	user := &models.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用戶登入
func (s *UserService) Login(req models.LoginRequest) (*models.User, error) {
	// 查找用戶
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 驗證密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// GetUsers 獲取所有用戶
func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

// GetUserByID 根據 ID 獲取用戶
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// UpdateUser 更新用戶
func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) error {
	return s.repo.Update(id, req.Username)
}

// DeleteUser 刪除用戶
func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
