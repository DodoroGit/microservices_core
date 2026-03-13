package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"user-service/models"
)

// UserRepositoryInterface 定義 repository 層的契約，讓 service 層依賴 interface 而非具體實作
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(id string, username string) error
	UpdateRole(id string, role string) error
	Delete(id string) error
	AddToBlacklist(token string, ttl time.Duration) error
}

// UserRepository 用戶資料訪問層，同時持有 DB 與 Redis
type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
}

// NewUserRepository 創建用戶 Repository
func NewUserRepository(db *sql.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{db: db, redis: redis}
}

// Create 創建用戶
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (id, email, username, password, role)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING created_at, updated_at`

	err := r.db.QueryRow(query, user.ID, user.Email, user.Username, user.Password, user.Role).
		Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByEmail 根據 email 查找用戶
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, username, password, role, created_at, updated_at
	          FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByID 根據 ID 查找用戶
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, username, role, created_at, updated_at
	          FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindAll 獲取所有用戶
func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `SELECT id, email, username, role, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// Update 更新用戶
func (r *UserRepository) Update(id string, username string) error {
	query := `UPDATE users SET username = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	result, err := r.db.Exec(query, username, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateRole 更新用戶角色
func (r *UserRepository) UpdateRole(id string, role string) error {
	query := `UPDATE users SET role = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	result, err := r.db.Exec(query, role, id)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete 刪除用戶
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// AddToBlacklist 將 token 加入黑名單，TTL 到期後 Redis 自動清除
func (r *UserRepository) AddToBlacklist(token string, ttl time.Duration) error {
	return r.redis.Set(context.Background(), "blacklist:"+token, "1", ttl).Err()
}
