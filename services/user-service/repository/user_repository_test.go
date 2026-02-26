//go:build integration

package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"user-service/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -------------------------------------------------------------------
// setupIntegrationDB：連接真實 DB，建立 table，並在 test 結束後清資料
//
// require.NoError 和 assert.NoError 的差別：
//   - require：失敗時立刻停止這個 test，後面的程式碼不會執行
//   - assert ：失敗時繼續執行，收集所有錯誤後一起回報
//
// 這裡用 require 是因為 DB 連不上的話後面所有操作都沒意義
// -------------------------------------------------------------------

func setupIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "admin")
	password := getEnvOrDefault("DB_PASSWORD", "admin123")
	dbname := getEnvOrDefault("DB_NAME", "userdb_test")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "failed to open db connection")
	require.NoError(t, db.Ping(), "failed to ping db — is postgres running?")

	// 建立 table（idempotent，重複跑不會壞）
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(100) NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`)
	require.NoError(t, err, "failed to create tables")

	// test 結束後清掉所有測試資料，保持 DB 乾淨
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE email LIKE '%@integration.test'")
		db.Close()
	})

	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// -------------------------------------------------------------------
// Create 測試
// -------------------------------------------------------------------

func TestUserRepository_Create(t *testing.T) {
	db := setupIntegrationDB(t)
	repo := NewUserRepository(db)

	tests := []struct {
		name        string
		user        *models.User
		expectError bool
	}{
		{
			name: "success",
			user: &models.User{
				ID:       "11111111-1111-1111-1111-111111111111",
				Email:    "create@integration.test",
				Username: "createuser",
				Password: "hashedpassword",
			},
			expectError: false,
		},
		{
			name: "duplicate email",
			user: &models.User{
				ID:       "22222222-2222-2222-2222-222222222222",
				Email:    "create@integration.test", // 同一個 email
				Username: "anotheruser",
				Password: "hashedpassword",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// DB 有回填 created_at / updated_at
				assert.False(t, tt.user.CreatedAt.IsZero(), "created_at should be set by DB")
				assert.False(t, tt.user.UpdatedAt.IsZero(), "updated_at should be set by DB")
			}
		})
	}
}

// -------------------------------------------------------------------
// FindByEmail 測試
// -------------------------------------------------------------------

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupIntegrationDB(t)
	repo := NewUserRepository(db)

	// 先建一筆資料供查詢用
	existing := &models.User{
		ID:       "33333333-3333-3333-3333-333333333333",
		Email:    "find@integration.test",
		Username: "finduser",
		Password: "hashedpassword",
	}
	require.NoError(t, repo.Create(existing))

	tests := []struct {
		name       string
		email      string
		expectUser bool
	}{
		{
			name:       "found",
			email:      "find@integration.test",
			expectUser: true,
		},
		{
			name:       "not found",
			email:      "nobody@integration.test",
			expectUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByEmail(tt.email)

			assert.NoError(t, err)
			if tt.expectUser {
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				// FindByEmail 會查 password（登入用），確認有撈回來
				assert.NotEmpty(t, user.Password)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

// -------------------------------------------------------------------
// FindByID 測試
// -------------------------------------------------------------------

func TestUserRepository_FindByID(t *testing.T) {
	db := setupIntegrationDB(t)
	repo := NewUserRepository(db)

	existing := &models.User{
		ID:       "44444444-4444-4444-4444-444444444444",
		Email:    "findbyid@integration.test",
		Username: "findbyiduser",
		Password: "hashedpassword",
	}
	require.NoError(t, repo.Create(existing))

	tests := []struct {
		name       string
		id         string
		expectUser bool
	}{
		{
			name:       "found",
			id:         "44444444-4444-4444-4444-444444444444",
			expectUser: true,
		},
		{
			name:       "not found",
			id:         "00000000-0000-0000-0000-000000000000",
			expectUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.FindByID(tt.id)

			assert.NoError(t, err)
			if tt.expectUser {
				assert.NotNil(t, user)
				assert.Equal(t, tt.id, user.ID)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

// -------------------------------------------------------------------
// Update 測試
// -------------------------------------------------------------------

func TestUserRepository_Update(t *testing.T) {
	db := setupIntegrationDB(t)
	repo := NewUserRepository(db)

	existing := &models.User{
		ID:       "55555555-5555-5555-5555-555555555555",
		Email:    "update@integration.test",
		Username: "oldname",
		Password: "hashedpassword",
	}
	require.NoError(t, repo.Create(existing))

	tests := []struct {
		name        string
		id          string
		newUsername string
		expectError bool
	}{
		{
			name:        "success",
			id:          "55555555-5555-5555-5555-555555555555",
			newUsername: "newname",
			expectError: false,
		},
		{
			name:        "user not found",
			id:          "00000000-0000-0000-0000-000000000000",
			newUsername: "newname",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.id, tt.newUsername)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 查回來確認真的有更新
				updated, _ := repo.FindByID(tt.id)
				assert.Equal(t, tt.newUsername, updated.Username)
			}
		})
	}
}

// -------------------------------------------------------------------
// Delete 測試
// -------------------------------------------------------------------

func TestUserRepository_Delete(t *testing.T) {
	db := setupIntegrationDB(t)
	repo := NewUserRepository(db)

	existing := &models.User{
		ID:       "66666666-6666-6666-6666-666666666666",
		Email:    "delete@integration.test",
		Username: "deleteuser",
		Password: "hashedpassword",
	}
	require.NoError(t, repo.Create(existing))

	tests := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name:        "success",
			id:          "66666666-6666-6666-6666-666666666666",
			expectError: false,
		},
		{
			name:        "already deleted",
			id:          "66666666-6666-6666-6666-666666666666", // 同一筆，已被刪
			expectError: true,
		},
		{
			name:        "not found",
			id:          "00000000-0000-0000-0000-000000000000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 查回來確認真的不見了
				deleted, _ := repo.FindByID(tt.id)
				assert.Nil(t, deleted)
			}
		})
	}
}
