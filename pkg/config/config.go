package config

import (
	"os"
)

type (
	Config struct {
		DbCfg DbConfig
	}

	DbConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
)

func New() *Config {
	return &Config{
		DbCfg: DbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "9000"),
			User:     getEnv("DB_USER", "default"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "default"),
		},
	}
}

// getEnv returns the fallback value if the given key is not provided in env
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
