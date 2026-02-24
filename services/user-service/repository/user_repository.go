package repository

import (
	"database/sql"
	"fmt"

	"user-service/models"
)

// UserRepositoryInterface 定義 repository 層的契約，讓 service 層依賴 interface 而非具體實作
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(id string, username string) error
	Delete(id string) error
}

// UserRepository 用戶資料訪問層
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 創建用戶 Repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 創建用戶
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (id, email, username, password)
	          VALUES ($1, $2, $3, $4)
	          RETURNING created_at, updated_at`

	err := r.db.QueryRow(query, user.ID, user.Email, user.Username, user.Password).
		Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByEmail 根據 email 查找用戶
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, username, password, created_at, updated_at
	          FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
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
	query := `SELECT id, email, username, created_at, updated_at
	          FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username,
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
	query := `SELECT id, email, username, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
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
