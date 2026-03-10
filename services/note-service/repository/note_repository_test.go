//go:build integration

package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"note-service/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// -------------------------------------------------------------------
// setupIntegrationDB：連接真實 MongoDB，並在 test 結束後清資料
//
// require.NoError 和 assert.NoError 的差別：
//   - require：失敗時立刻停止這個 test，後面的程式碼不會執行
//   - assert ：失敗時繼續執行，收集所有錯誤後一起回報
//
// 這裡用 require 是因為 MongoDB 連不上的話後面所有操作都沒意義
// -------------------------------------------------------------------

func setupIntegrationDB(t *testing.T) *NoteRepository {
	t.Helper()

	mongoURL := getTestEnv("MONGODB_URL", "mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	require.NoError(t, err, "failed to connect to MongoDB")
	require.NoError(t, client.Ping(context.Background(), nil), "failed to ping MongoDB — is it running?")

	db := client.Database("test_notedb")
	repo := NewNoteRepository(db)

	// test 結束後清掉 notes collection，保持 DB 乾淨
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Collection("notes").Drop(ctx)
		client.Disconnect(ctx)
	})

	return repo
}

func getTestEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// ===================================================================
// Create 測試
// ===================================================================

func TestNoteRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		note := &models.Note{
			Title:    "Integration Test Note",
			Category: models.CategoryTech,
			Content:  "Test content",
		}

		created, err := repo.Create(note)

		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.NotEqual(t, primitive.NilObjectID, created.ID, "ID should be assigned by repository")
		assert.Equal(t, "Integration Test Note", created.Title)
		assert.False(t, created.CreatedAt.IsZero(), "created_at should be set")
		assert.False(t, created.UpdatedAt.IsZero(), "updated_at should be set")
	})
}

// ===================================================================
// FindAll 測試
// ===================================================================

func TestNoteRepository_FindAll(t *testing.T) {
	t.Run("returns all notes", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		repo.Create(&models.Note{Title: "Note 1", Category: models.CategoryTech, Content: "Content 1"})
		repo.Create(&models.Note{Title: "Note 2", Category: models.CategoryProject, Content: "Content 2"})

		notes, err := repo.FindAll("")

		assert.NoError(t, err)
		assert.Len(t, notes, 2)
	})

	t.Run("empty collection", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		notes, err := repo.FindAll("")

		assert.NoError(t, err)
		// 空集合應回傳空 slice，不是 nil
		assert.NotNil(t, notes)
		assert.Empty(t, notes)
	})

	t.Run("filter by category", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		repo.Create(&models.Note{Title: "Tech Note", Category: models.CategoryTech, Content: "Content"})
		repo.Create(&models.Note{Title: "Project Note", Category: models.CategoryProject, Content: "Content"})

		notes, err := repo.FindAll(string(models.CategoryTech))

		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, models.CategoryTech, notes[0].Category)
	})

	t.Run("sorted by created_at descending", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		first, _ := repo.Create(&models.Note{Title: "First", Category: models.CategoryTech, Content: "Content"})
		time.Sleep(10 * time.Millisecond)
		repo.Create(&models.Note{Title: "Second", Category: models.CategoryTech, Content: "Content"})

		notes, err := repo.FindAll("")

		require.NoError(t, err)
		require.Len(t, notes, 2)
		// 最新的排在第一筆
		assert.NotEqual(t, first.ID, notes[0].ID, "newest note should be first")
		assert.Equal(t, "Second", notes[0].Title)
	})
}

// ===================================================================
// FindByID 測試
// ===================================================================

func TestNoteRepository_FindByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		created, _ := repo.Create(&models.Note{Title: "Find Me", Category: models.CategoryDaily, Content: "Content"})

		found, err := repo.FindByID(created.ID.Hex())

		assert.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Find Me", found.Title)
	})

	t.Run("not found", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		found, err := repo.FindByID(primitive.NewObjectID().Hex())

		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("invalid id format - treated as not found", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		found, err := repo.FindByID("not-a-valid-objectid")

		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

// ===================================================================
// Update 測試
// ===================================================================

func TestNoteRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		created, _ := repo.Create(&models.Note{Title: "Original Title", Category: models.CategoryTech, Content: "Content"})

		updates := map[string]interface{}{"title": "Updated Title"}
		updated, err := repo.Update(created.ID.Hex(), updates)

		assert.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, "Updated Title", updated.Title)
		// updated_at 應該被刷新
		assert.True(t, !updated.UpdatedAt.Before(created.UpdatedAt), "updated_at should be refreshed")
	})

	t.Run("not found", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		updates := map[string]interface{}{"title": "Updated"}
		result, err := repo.Update(primitive.NewObjectID().Hex(), updates)

		assert.NoError(t, err)
		// 找不到時回傳 nil, nil（不是 error）
		assert.Nil(t, result)
	})
}

// ===================================================================
// Delete 測試
// ===================================================================

func TestNoteRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		created, _ := repo.Create(&models.Note{Title: "Delete Me", Category: models.CategoryTech, Content: "Content"})

		err := repo.Delete(created.ID.Hex())

		assert.NoError(t, err)
		// 查回來確認真的不見了
		deleted, _ := repo.FindByID(created.ID.Hex())
		assert.Nil(t, deleted)
	})

	t.Run("already deleted", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		created, _ := repo.Create(&models.Note{Title: "Delete Me", Category: models.CategoryTech, Content: "Content"})
		// 先刪一次
		require.NoError(t, repo.Delete(created.ID.Hex()))

		// 再刪一次，應該要失敗
		err := repo.Delete(created.ID.Hex())

		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := setupIntegrationDB(t)

		err := repo.Delete(primitive.NewObjectID().Hex())

		assert.Error(t, err)
	})
}
