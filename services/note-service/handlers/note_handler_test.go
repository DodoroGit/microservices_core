package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"note-service/models"
	"note-service/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -------------------------------------------------------------------
// MockNoteService：手動實作 NoteServiceInterface 供測試用
// Handler 測試只需 mock service，完全不涉及 DB
// -------------------------------------------------------------------

type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) GetAll(category string) ([]models.Note, error) {
	args := m.Called(category)
	return args.Get(0).([]models.Note), args.Error(1)
}

func (m *MockNoteService) GetByID(id string) (*models.Note, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteService) Create(req models.CreateNoteRequest) (*models.Note, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteService) Update(id string, req models.UpdateNoteRequest) (*models.Note, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// -------------------------------------------------------------------
// 測試輔助：建立 gin test router
// -------------------------------------------------------------------

func setupTestRouter(svc services.NoteServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewNoteHandler(svc)
	r.GET("/health", h.Health)
	r.GET("/notes", h.GetNotes)
	r.GET("/notes/:id", h.GetNote)
	r.POST("/notes", h.CreateNote)
	r.PUT("/notes/:id", h.UpdateNote)
	r.DELETE("/notes/:id", h.DeleteNote)
	return r
}

func makeTestNote() *models.Note {
	return &models.Note{
		ID:        primitive.NewObjectID(),
		Title:     "Test Note",
		Category:  models.CategoryTech,
		Content:   "Test content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ===================================================================
// Health handler 測試
// Health 只有一個情境，不用 t.Run
// ===================================================================

func TestHealthHandler(t *testing.T) {
	mockSvc := new(MockNoteService)
	router := setupTestRouter(mockSvc)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "healthy", resp["status"])
	assert.Equal(t, "note-service", resp["service"])
}

// ===================================================================
// GetNotes handler 測試
// ===================================================================

func TestGetNotesHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		mockSvc.On("GetAll", "").Return([]models.Note{*note}, nil)

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []models.NoteResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Len(t, resp, 1)
		assert.Equal(t, note.ID.Hex(), resp[0].ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("with category query", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		mockSvc.On("GetAll", "技術筆記").Return([]models.Note{*note}, nil)

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/notes?category=技術筆記", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		mockSvc.On("GetAll", "").Return([]models.Note{}, nil)

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, "[]", w.Body.String())
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// GetNote handler 測試
// ===================================================================

func TestGetNoteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		mockSvc.On("GetByID", note.ID.Hex()).Return(note, nil)

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/notes/"+note.ID.Hex(), nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.NoteResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, note.ID.Hex(), resp.ID)
		assert.Equal(t, note.Title, resp.Title)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		mockSvc.On("GetByID", "no-such-id").Return(nil, errors.New("note not found"))

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/notes/no-such-id", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// CreateNote handler 測試
// ===================================================================

func TestCreateNoteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		reqBody := models.CreateNoteRequest{
			Title:    note.Title,
			Category: note.Category,
			Content:  note.Content,
		}
		mockSvc.On("Create", reqBody).Return(note, nil)

		router := setupTestRouter(mockSvc)
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.NoteResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, note.ID.Hex(), resp.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid body - missing required fields", func(t *testing.T) {
		mockSvc := new(MockNoteService)

		router := setupTestRouter(mockSvc)
		body := []byte(`{"title": "Only Title"}`)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		// body 不合法，service 不應該被呼叫
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid category", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		reqBody := models.CreateNoteRequest{
			Title:    "Test",
			Category: "無效分類",
			Content:  "Content",
		}
		mockSvc.On("Create", reqBody).Return(nil, errors.New("invalid category"))

		router := setupTestRouter(mockSvc)
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// UpdateNote handler 測試
// ===================================================================

func TestUpdateNoteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		newTitle := "Updated Title"
		reqBody := models.UpdateNoteRequest{Title: &newTitle}
		mockSvc.On("Update", note.ID.Hex(), reqBody).Return(note, nil)

		router := setupTestRouter(mockSvc)
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPut, "/notes/"+note.ID.Hex(), bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		newTitle := "Updated Title"
		reqBody := models.UpdateNoteRequest{Title: &newTitle}
		mockSvc.On("Update", "no-such-id", reqBody).Return(nil, errors.New("note not found"))

		router := setupTestRouter(mockSvc)
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPut, "/notes/no-such-id", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid category", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		invalidCat := models.Category("無效分類")
		reqBody := models.UpdateNoteRequest{Category: &invalidCat}
		mockSvc.On("Update", note.ID.Hex(), reqBody).Return(nil, errors.New("invalid category"))

		router := setupTestRouter(mockSvc)
		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodPut, "/notes/"+note.ID.Hex(), bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ===================================================================
// DeleteNote handler 測試
// ===================================================================

func TestDeleteNoteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		note := makeTestNote()
		mockSvc.On("Delete", note.ID.Hex()).Return(nil)

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodDelete, "/notes/"+note.ID.Hex(), nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockNoteService)
		mockSvc.On("Delete", "no-such-id").Return(errors.New("note not found"))

		router := setupTestRouter(mockSvc)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodDelete, "/notes/no-such-id", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
