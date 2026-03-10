package config

import "os"

type Config struct {
	Port    string
	MongoDB MongoDBConfig
}

type MongoDBConfig struct {
	URL    string
	DBName string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8083"),
		MongoDB: MongoDBConfig{
			URL:    getEnv("MONGODB_URL", "mongodb://localhost:27017"),
			DBName: getEnv("MONGODB_DB", "notedb"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
