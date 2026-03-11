package config

import (
	"os"
)

type Config struct {
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
}

func Load() *Config {
	return &Config{
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5433"),
		PostgresUser: getEnv("POSTGRES_USER", "todos"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "todos"),
		PostgresDB:   getEnv("POSTGRES_DB", "todos"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
