package services

import (
	"errors"
	"testing"
	"time"

	"note-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -------------------------------------------------------------------
// MockNoteRepository：手動實作 NoteRepositoryInterface 供測試用
// Service 測試只需 mock repository，完全不涉及 DB
// -------------------------------------------------------------------

type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) FindAll(category string) ([]models.Note, error) {
	args := m.Called(category)
	return args.Get(0).([]models.Note), args.Error(1)
}

func (m *MockNoteRepository) FindByID(id string) (*models.Note, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) Create(note *models.Note) (*models.Note, error) {
	args := m.Called(note)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) Update(id string, updates map[string]interface{}) (*models.Note, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *MockNoteRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// makeNote：建立一個測試用的 Note，避免每個 test 重複寫
func makeNote() *models.Note {
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
// GetAll 測試
// ===================================================================

func TestGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		mockRepo.On("FindAll", "").Return([]models.Note{*makeNote()}, nil)

		svc := NewNoteService(mockRepo)
		notes, err := svc.GetAll("")

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("with category filter", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		mockRepo.On("FindAll", string(models.CategoryTech)).Return([]models.Note{*makeNote()}, nil)

		svc := NewNoteService(mockRepo)
		notes, err := svc.GetAll(string(models.CategoryTech))

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		mockRepo.On("FindAll", "").Return([]models.Note{}, nil)

		svc := NewNoteService(mockRepo)
		notes, err := svc.GetAll("")

		assert.NoError(t, err)
		assert.Empty(t, notes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		mockRepo.On("FindAll", "").Return([]models.Note(nil), errors.New("db connection failed"))

		svc := NewNoteService(mockRepo)
		notes, err := svc.GetAll("")

		assert.Error(t, err)
		assert.Nil(t, notes)
		assert.Contains(t, err.Error(), "db connection failed")
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// GetByID 測試
// ===================================================================

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		note := makeNote()
		mockRepo.On("FindByID", note.ID.Hex()).Return(note, nil)

		svc := NewNoteService(mockRepo)
		result, err := svc.GetByID(note.ID.Hex())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, note.ID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		// DB 查無此 ID → 回傳 nil, nil（不是 error，只是找不到）
		mockRepo.On("FindByID", "nonexistent-id").Return(nil, nil)

		svc := NewNoteService(mockRepo)
		result, err := svc.GetByID("nonexistent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "note not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		mockRepo.On("FindByID", "error-id").Return(nil, errors.New("db error"))

		svc := NewNoteService(mockRepo)
		result, err := svc.GetByID("error-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// Create 測試
// ===================================================================

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		created := makeNote()
		// 接受任意 *models.Note，因為 ID 會在 repository 裡才產生
		mockRepo.On("Create", mock.AnythingOfType("*models.Note")).Return(created, nil)

		svc := NewNoteService(mockRepo)
		result, err := svc.Create(models.CreateNoteRequest{
			Title:    "New Note",
			Category: models.CategoryTech,
			Content:  "Content here",
		})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, created.ID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid category", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)

		svc := NewNoteService(mockRepo)
		result, err := svc.Create(models.CreateNoteRequest{
			Title:    "New Note",
			Category: "無效分類",
			Content:  "Content here",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid category")
		// category 驗證失敗，repository 不應該被呼叫
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		mockRepo.On("Create", mock.AnythingOfType("*models.Note")).Return(nil, errors.New("db error"))

		svc := NewNoteService(mockRepo)
		result, err := svc.Create(models.CreateNoteRequest{
			Title:    "New Note",
			Category: models.CategoryTech,
			Content:  "Content here",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// Update 測試
// ===================================================================

func TestUpdate(t *testing.T) {
	t.Run("success - update title", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		note := makeNote()
		newTitle := "Updated Title"

		mockRepo.On("FindByID", note.ID.Hex()).Return(note, nil)
		mockRepo.On("Update", note.ID.Hex(), mock.MatchedBy(func(m map[string]interface{}) bool {
			v, ok := m["title"]
			return ok && v == newTitle
		})).Return(note, nil)

		svc := NewNoteService(mockRepo)
		result, err := svc.Update(note.ID.Hex(), models.UpdateNoteRequest{Title: &newTitle})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		newTitle := "Updated Title"
		// FindByID 查無此筆 → 回傳 nil, nil
		mockRepo.On("FindByID", "nonexistent-id").Return(nil, nil)

		svc := NewNoteService(mockRepo)
		result, err := svc.Update("nonexistent-id", models.UpdateNoteRequest{Title: &newTitle})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "note not found")
		// FindByID 回傳 not found，Update 不應該被呼叫
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid category", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		invalidCat := models.Category("無效分類")

		svc := NewNoteService(mockRepo)
		result, err := svc.Update("any-id", models.UpdateNoteRequest{Category: &invalidCat})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid category")
		// category 驗證在 FindByID 之前，repository 完全不應該被呼叫
		mockRepo.AssertExpectations(t)
	})

	t.Run("no fields to update - returns existing note", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		note := makeNote()
		// 空 request → service 直接回傳現有的 note，不呼叫 Update
		mockRepo.On("FindByID", note.ID.Hex()).Return(note, nil)

		svc := NewNoteService(mockRepo)
		result, err := svc.Update(note.ID.Hex(), models.UpdateNoteRequest{})

		assert.NoError(t, err)
		assert.Equal(t, note, result)
		mockRepo.AssertExpectations(t)
	})
}

// ===================================================================
// Delete 測試
// ===================================================================

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		note := makeNote()
		mockRepo.On("FindByID", note.ID.Hex()).Return(note, nil)
		mockRepo.On("Delete", note.ID.Hex()).Return(nil)

		svc := NewNoteService(mockRepo)
		err := svc.Delete(note.ID.Hex())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		// FindByID 查無此筆 → 不應該呼叫 Delete
		mockRepo.On("FindByID", "nonexistent-id").Return(nil, nil)

		svc := NewNoteService(mockRepo)
		err := svc.Delete("nonexistent-id")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "note not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on delete", func(t *testing.T) {
		mockRepo := new(MockNoteRepository)
		note := makeNote()
		mockRepo.On("FindByID", note.ID.Hex()).Return(note, nil)
		mockRepo.On("Delete", note.ID.Hex()).Return(errors.New("db error"))

		svc := NewNoteService(mockRepo)
		err := svc.Delete(note.ID.Hex())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}
