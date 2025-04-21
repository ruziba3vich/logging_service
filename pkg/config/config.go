package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		DbCfg       *DbConfig
		ConsumerCfg *ConsumerConfig
	}

	DbConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}

	ConsumerConfig struct {
		Brokers          string
		GroupID          string
		Topic            string
		AutoOffsetReset  string
		EnableAutoCommit bool
		MaxPoolInterval  int
		SessionTimeOut   int
	}
)

func New() *Config {
	return &Config{
		DbCfg: &DbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "9000"),
			User:     getEnv("DB_USER", "default"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "default"),
		},
		ConsumerCfg: &ConsumerConfig{
			Brokers:          getEnv("BROKERS", "localhost:9092"),
			GroupID:          getEnv("GROUP_ID", "logging_consumer"),
			Topic:            getEnv("KAFKA_TOPIC", "logs"),
			AutoOffsetReset:  getEnv("AUTO_OFFSET_RESET", "earliest"),
			EnableAutoCommit: getEnvBool("ENABLE_AUTO_COMMIT", false),
			MaxPoolInterval:  getEnvInt("MAX_POOL_INTERVAL", 300000),
			SessionTimeOut:   getEnvInt("SESSION_TIME_OUT", 30000),
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

func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true"
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		v, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return v
	}
	return fallback
}
