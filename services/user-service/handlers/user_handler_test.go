package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"user-service/models"
)

// -------------------------------------------------------------------
// MockUserService：手動實作 UserServiceInterface 供測試用
// Handler 測試只需 mock service，完全不涉及 DB
// -------------------------------------------------------------------

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(req models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Login(req models.LoginRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, req models.UpdateUserRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// -------------------------------------------------------------------
// 測試輔助：建立 gin test router
// -------------------------------------------------------------------

func setupTestRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/users/register", handler.Register)
	r.POST("/users/login", handler.Login)
	r.GET("/users", handler.GetUsers)
	r.GET("/users/:id", handler.GetUser)
	r.PUT("/users/:id", handler.UpdateUser)
	r.DELETE("/users/:id", handler.DeleteUser)
	r.GET("/health", handler.Health)
	return r
}

// ===================================================================
// Register handler 測試
// ===================================================================

func TestRegisterHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("Register", models.RegisterRequest{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
		}).Return(&models.User{ID: "uuid-001", Email: "test@example.com", Username: "testuser"}, nil)

		router := setupTestRouter(NewUserHandler(mockSvc))

		body, _ := json.Marshal(models.RegisterRequest{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
		})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.User
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "test@example.com", resp.Email)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid body - missing fields", func(t *testing.T) {
		mockSvc := new(MockUserService)

		router := setupTestRouter(NewUserHandler(mockSvc))

		body, _ := json.Marshal(map[string]string{"email": "not-valid-email"})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// body 不合法，service 不應該被呼叫
		mockSvc.AssertExpectations(t)
	})

	t.Run("service error - email exists", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("Register", models.RegisterRequest{
			Email:    "exist@example.com",
			Username: "someone",
			Password: "password123",
		}).Return(nil, fmt.Errorf("email already exists"))

		router := setupTestRouter(NewUserHandler(mockSvc))

		body, _ := json.Marshal(models.RegisterRequest{
			Email:    "exist@example.com",
			Username: "someone",
			Password: "password123",
		})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// Login handler 測試
// ===================================================================

func TestLoginHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("Login", models.LoginRequest{
			Email:    "user@example.com",
			Password: "password123",
		}).Return(&models.User{ID: "uuid-001", Email: "user@example.com"}, nil)

		router := setupTestRouter(NewUserHandler(mockSvc))

		body, _ := json.Marshal(models.LoginRequest{
			Email:    "user@example.com",
			Password: "password123",
		})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Login successful", resp.Message)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := new(MockUserService)

		router := setupTestRouter(NewUserHandler(mockSvc))

		body, _ := json.Marshal(map[string]string{"email": "no-password"})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("Login", models.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpass",
		}).Return(nil, fmt.Errorf("invalid credentials"))

		router := setupTestRouter(NewUserHandler(mockSvc))

		body, _ := json.Marshal(models.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpass",
		})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// GetUser handler 測試
// ===================================================================

func TestGetUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("GetUserByID", "abc-123").Return(&models.User{ID: "abc-123", Email: "u@example.com"}, nil)

		router := setupTestRouter(NewUserHandler(mockSvc))

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/users/abc-123", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.User
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "abc-123", resp.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("GetUserByID", "no-such-id").Return(nil, fmt.Errorf("user not found"))

		router := setupTestRouter(NewUserHandler(mockSvc))

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/users/no-such-id", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// DeleteUser handler 測試
// ===================================================================

func TestDeleteUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("DeleteUser", "abc-123").Return(nil)

		router := setupTestRouter(NewUserHandler(mockSvc))

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("DELETE", "/users/abc-123", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		mockSvc.On("DeleteUser", "ghost-id").Return(fmt.Errorf("user not found"))

		router := setupTestRouter(NewUserHandler(mockSvc))

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("DELETE", "/users/ghost-id", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// Health handler 測試
// Health 只有一個情境，不用 t.Run
// ===================================================================

func TestHealthHandler(t *testing.T) {
	mockSvc := new(MockUserService)
	router := setupTestRouter(NewUserHandler(mockSvc))

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "user-service", response["service"])
}
