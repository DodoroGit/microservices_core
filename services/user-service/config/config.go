package config

import "os"

// Config 應用配置
type Config struct {
	Port     string
	Database DatabaseConfig
	Redis    RedisConfig
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host string
	Port string
}

// Load 從環境變數載入配置
func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8081"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "admin123"),
			DBName:   getEnv("DB_NAME", "userdb"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
