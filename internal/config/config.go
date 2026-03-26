// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host          string
	Port          string
	AllowedOrigin string
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	FromEmail     string
	RedisURL      string
	BulkBatchSize int
}

func LoadConfig() *Config {
	// Load .env file if it exists (optional for Docker deployments)
	_ = godotenv.Load()

	cfg := &Config{
		Host:          GetEnv("HOST", "0.0.0.0"),
		Port:          GetEnv("PORT", "8082"),
		AllowedOrigin: GetEnv("ALLOWED_ORIGIN", "*"),
		SMTPHost:      GetEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      GetEnvAsInt("SMTP_PORT", 587),
		SMTPUsername:  GetEnv("SMTP_USERNAME", ""),
		SMTPPassword:  GetEnv("SMTP_PASSWORD", ""),
		FromEmail:     GetEnv("FROM_EMAIL", ""),
		RedisURL:      GetEnv("REDIS_URL", ""),
		BulkBatchSize: GetEnvAsInt("BULK_BATCH_SIZE", 10),
	}
	validateSecrets(cfg)
	return cfg
}

// validateSecrets checks for required secrets and logs a fatal error if missing
func validateSecrets(cfg *Config) {
	missing := []string{}
	if cfg.SMTPUsername == "" {
		missing = append(missing, "SMTP_USERNAME")
	}
	if cfg.SMTPPassword == "" {
		missing = append(missing, "SMTP_PASSWORD")
	}
	if cfg.FromEmail == "" {
		missing = append(missing, "FROM_EMAIL")
	}
	if len(missing) > 0 {
		panic("Missing required secrets: " + joinStrings(missing, ", "))
	}
}

func joinStrings(arr []string, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result += sep + arr[i]
	}
	return result
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
