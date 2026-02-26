package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"user-service/models"
)

// -------------------------------------------------------------------
// MockUserRepository：手動實作 UserRepositoryInterface 供測試用
// testify/mock 會幫我們追蹤每個 method 被呼叫的次數與參數
// -------------------------------------------------------------------

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(id string, username string) error {
	args := m.Called(id, username)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// setupHashedUser：呼叫真實 Register 取得已 bcrypt hash 過的 user
// Login 測試需要 hash 過的密碼，用這個 helper 產生，避免重複寫 register 流程
func setupHashedUser(t *testing.T, email, username, password string) *models.User {
	t.Helper()
	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByEmail", email).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	svc := NewUserService(mockRepo)
	user, err := svc.Register(models.RegisterRequest{Email: email, Username: username, Password: password})
	assert.NoError(t, err)
	return user
}

// ===================================================================
// Register 測試
// ===================================================================

func TestRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		// email 不存在 → 回傳 nil, nil
		mockRepo.On("FindByEmail", "new@example.com").Return(nil, nil)
		// Create 被呼叫時，接受任意 *models.User，成功不回錯誤
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		svc := NewUserService(mockRepo)
		user, err := svc.Register(models.RegisterRequest{
			Email:    "new@example.com",
			Username: "newuser",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "new@example.com", user.Email)
		assert.Equal(t, "newuser", user.Username)
		// 確認密碼已被 hash，不是明文
		assert.NotEqual(t, "password123", user.Password)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		// email 已存在 → 回傳一個有資料的 user
		existing := &models.User{Email: "exist@example.com"}
		mockRepo.On("FindByEmail", "exist@example.com").Return(existing, nil)

		svc := NewUserService(mockRepo)
		user, err := svc.Register(models.RegisterRequest{
			Email:    "exist@example.com",
			Username: "someone",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "email already exists")
		// Create 應該完全沒被呼叫
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error on find", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		// FindByEmail 本身就出錯（DB 連線問題等）
		mockRepo.On("FindByEmail", "error@example.com").Return(nil, fmt.Errorf("db connection failed"))

		svc := NewUserService(mockRepo)
		user, err := svc.Register(models.RegisterRequest{
			Email:    "error@example.com",
			Username: "someone",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "db connection failed")
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// Login 測試
// ===================================================================

func TestLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// 先透過 Register 產生 hash 過的 user，模擬 DB 裡存的狀態
		hashedUser := setupHashedUser(t, "user@example.com", "user", "correctpassword")

		mockRepo := new(MockUserRepository)
		mockRepo.On("FindByEmail", "user@example.com").Return(hashedUser, nil)

		svc := NewUserService(mockRepo)
		user, err := svc.Login(models.LoginRequest{
			Email:    "user@example.com",
			Password: "correctpassword",
		})

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "user@example.com", user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		// email 查不到 → 回傳 nil, nil（不是 error，只是找不到）
		mockRepo.On("FindByEmail", "ghost@example.com").Return(nil, nil)

		svc := NewUserService(mockRepo)
		user, err := svc.Login(models.LoginRequest{
			Email:    "ghost@example.com",
			Password: "somepassword",
		})

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		// 有這個 user，但密碼錯
		hashedUser := setupHashedUser(t, "user@example.com", "user", "correctpassword")

		mockRepo := new(MockUserRepository)
		mockRepo.On("FindByEmail", "user@example.com").Return(hashedUser, nil)

		svc := NewUserService(mockRepo)
		user, err := svc.Login(models.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpassword",
		})

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("FindByEmail", "user@example.com").Return(nil, fmt.Errorf("db connection failed"))

		svc := NewUserService(mockRepo)
		user, err := svc.Login(models.LoginRequest{
			Email:    "user@example.com",
			Password: "correctpassword",
		})

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "db connection failed")
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// GetUserByID 測試
// ===================================================================

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("FindByID", "abc-123").Return(&models.User{ID: "abc-123", Email: "u@example.com"}, nil)

		svc := NewUserService(mockRepo)
		user, err := svc.GetUserByID("abc-123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "abc-123", user.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		// DB 查無此 ID → 回傳 nil, nil
		mockRepo.On("FindByID", "not-exist").Return(nil, nil)

		svc := NewUserService(mockRepo)
		user, err := svc.GetUserByID("not-exist")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("FindByID", "error-id").Return(nil, fmt.Errorf("db error"))

		svc := NewUserService(mockRepo)
		user, err := svc.GetUserByID("error-id")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// DeleteUser 測試
// ===================================================================

func TestDeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("Delete", "abc-123").Return(nil)

		svc := NewUserService(mockRepo)
		err := svc.DeleteUser("abc-123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockRepo.On("Delete", "ghost-id").Return(fmt.Errorf("user not found"))

		svc := NewUserService(mockRepo)
		err := svc.DeleteUser("ghost-id")

		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
		mockRepo.AssertExpectations(t)
	})
}
