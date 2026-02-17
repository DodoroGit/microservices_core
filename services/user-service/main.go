package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	db          *sql.DB
	redisClient *redis.Client
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func main() {
	// Initialize database
	initDB()
	defer db.Close()

	// Initialize Redis
	initRedis()
	defer redisClient.Close()

	// Create tables
	createTables()

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "user-service",
		})
	})

	// User routes
	router.POST("/users/register", registerUser)
	router.POST("/users/login", loginUser)
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	port := getEnv("PORT", "8081")
	log.Printf("User Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "admin"),
		getEnv("DB_PASSWORD", "admin123"),
		getEnv("DB_NAME", "userdb"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return
		}
		log.Printf("Waiting for database... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("Failed to ping database:", err)
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			getEnv("REDIS_HOST", "localhost"),
			getEnv("REDIS_PORT", "6379"),
		),
	})
	log.Println("Redis client initialized")
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(100) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	log.Println("Database tables created successfully")
}

func registerUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	query := `INSERT INTO users (id, email, username, password) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at`
	err = db.QueryRow(query, user.ID, user.Email, user.Username, user.Password).Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func loginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	query := `SELECT id, email, username, password, created_at, updated_at FROM users WHERE email = $1`
	err := db.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}

func getUsers(c *gin.Context) {
	rows, err := db.Query(`SELECT id, email, username, created_at, updated_at FROM users`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	query := `SELECT id, email, username, created_at, updated_at FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE users SET username = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.Exec(query, req.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
