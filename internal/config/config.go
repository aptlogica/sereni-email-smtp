// internal/config/config.go
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	FromEmail     string
	RedisURL      string
	BulkBatchSize int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:          GetEnv("PORT", "8080"),
		SMTPHost:      GetEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      GetEnvAsInt("SMTP_PORT", 587),
		SMTPUsername:  GetEnv("SMTP_USERNAME", ""),
		SMTPPassword:  GetEnv("SMTP_PASSWORD", ""),
		FromEmail:     GetEnv("FROM_EMAIL", ""),
		RedisURL:      GetEnv("REDIS_URL", ""),
		BulkBatchSize: GetEnvAsInt("BULK_BATCH_SIZE", 10),
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
