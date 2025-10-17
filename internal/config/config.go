package config

import (
	"os"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	// Server configuration
	Port        string
	Environment string

	// Database configuration
	DatabaseURL string
	RedisURL    string

	// Security
	JWTSecret string

	// Telegram
	TelegramBotToken string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),

		JWTSecret:        os.Getenv("JWT_SECRET"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsProduction проверяет, запущено ли приложение в production режиме
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment проверяет, запущено ли приложение в development режиме
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

