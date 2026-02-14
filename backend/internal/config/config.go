package config

import "os"

type Config struct {
	DatabasePath string
	JWTSecret    string
	Port         string
}

func Load() *Config {
	return &Config{
		DatabasePath: getEnv("DATABASE_PATH", "./chronovault.db"),
		JWTSecret:    getEnv("JWT_SECRET", "chronovault-secret-key-change-in-production"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
