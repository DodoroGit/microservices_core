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

// -------------------------------------------------------------------
// Register 測試
// -------------------------------------------------------------------

func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		req           models.RegisterRequest
		setupMock     func(*MockUserRepository)
		expectUser    bool
		expectedError string
		// 額外驗證：成功時才會執行
		checkUser func(*testing.T, *models.User)
	}{
		{
			name: "success",
			req: models.RegisterRequest{
				Email:    "new@example.com",
				Username: "newuser",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", "new@example.com").Return(nil, nil)
				m.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectUser: true,
			checkUser: func(t *testing.T, u *models.User) {
				assert.Equal(t, "new@example.com", u.Email)
				assert.Equal(t, "newuser", u.Username)
				// 確認密碼已被 hash，不是明文
				assert.NotEqual(t, "password123", u.Password)
			},
		},
		{
			name: "email already exists",
			req: models.RegisterRequest{
				Email:    "exist@example.com",
				Username: "someone",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				existing := &models.User{Email: "exist@example.com"}
				m.On("FindByEmail", "exist@example.com").Return(existing, nil)
			},
			expectUser:    false,
			expectedError: "email already exists",
		},
		{
			name: "db error on find",
			req: models.RegisterRequest{
				Email:    "error@example.com",
				Username: "someone",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", "error@example.com").Return(nil, fmt.Errorf("db connection failed"))
			},
			expectUser:    false,
			expectedError: "db connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)
			svc := NewUserService(mockRepo)

			user, err := svc.Register(tt.req)

			if tt.expectUser {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				if tt.checkUser != nil {
					tt.checkUser(t, user)
				}
			} else {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// -------------------------------------------------------------------
// Login 測試
// Login 的特殊性：success / wrong password 都需要先有 bcrypt hash 過的 user
// 所以用 setupUser helper 產生，不重複寫 register 流程
// -------------------------------------------------------------------

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

func TestLogin(t *testing.T) {
	hashedUser := setupHashedUser(t, "user@example.com", "user", "correctpassword")

	tests := []struct {
		name          string
		req           models.LoginRequest
		setupMock     func(*MockUserRepository)
		expectUser    bool
		expectedError string
	}{
		{
			name: "success",
			req:  models.LoginRequest{Email: "user@example.com", Password: "correctpassword"},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", "user@example.com").Return(hashedUser, nil)
			},
			expectUser: true,
		},
		{
			name: "user not found",
			req:  models.LoginRequest{Email: "ghost@example.com", Password: "somepassword"},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", "ghost@example.com").Return(nil, nil)
			},
			expectUser:    false,
			expectedError: "invalid credentials",
		},
		{
			name: "wrong password",
			req:  models.LoginRequest{Email: "user@example.com", Password: "wrongpassword"},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", "user@example.com").Return(hashedUser, nil)
			},
			expectUser:    false,
			expectedError: "invalid credentials",
		},
		{
			name: "db error",
			req:  models.LoginRequest{Email: "user@example.com", Password: "correctpassword"},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", "user@example.com").Return(nil, fmt.Errorf("db connection failed"))
			},
			expectUser:    false,
			expectedError: "db connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)
			svc := NewUserService(mockRepo)

			user, err := svc.Login(tt.req)

			if tt.expectUser {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.req.Email, user.Email)
			} else {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// -------------------------------------------------------------------
// GetUserByID 測試
// -------------------------------------------------------------------

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		setupMock     func(*MockUserRepository)
		expectUser    bool
		expectedError string
	}{
		{
			name: "success",
			id:   "abc-123",
			setupMock: func(m *MockUserRepository) {
				m.On("FindByID", "abc-123").Return(&models.User{ID: "abc-123", Email: "u@example.com"}, nil)
			},
			expectUser: true,
		},
		{
			name: "not found",
			id:   "not-exist",
			setupMock: func(m *MockUserRepository) {
				m.On("FindByID", "not-exist").Return(nil, nil)
			},
			expectUser:    false,
			expectedError: "user not found",
		},
		{
			name: "db error",
			id:   "error-id",
			setupMock: func(m *MockUserRepository) {
				m.On("FindByID", "error-id").Return(nil, fmt.Errorf("db error"))
			},
			expectUser:    false,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)
			svc := NewUserService(mockRepo)

			user, err := svc.GetUserByID(tt.id)

			if tt.expectUser {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, user.ID)
			} else {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// -------------------------------------------------------------------
// DeleteUser 測試
// -------------------------------------------------------------------

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		setupMock     func(*MockUserRepository)
		expectedError string
	}{
		{
			name: "success",
			id:   "abc-123",
			setupMock: func(m *MockUserRepository) {
				m.On("Delete", "abc-123").Return(nil)
			},
		},
		{
			name: "not found",
			id:   "ghost-id",
			setupMock: func(m *MockUserRepository) {
				m.On("Delete", "ghost-id").Return(fmt.Errorf("user not found"))
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)
			svc := NewUserService(mockRepo)

			err := svc.DeleteUser(tt.id)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
