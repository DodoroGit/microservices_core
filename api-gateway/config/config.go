package config

import "os"

// Config 儲存 API Gateway 所有執行時的設定。
type Config struct {
	Port            string
	JWTSecret       string
	UserServiceURL  string
	AIServiceURL    string
	NoteServiceURL  string
}

// Load 從環境變數讀取設定，若未設定則使用預設值。
func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8081"),
		AIServiceURL:   getEnv("AI_SERVICE_URL", "http://localhost:8082"),
		NoteServiceURL: getEnv("NOTE_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
