package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() {

	fmt.Printf("Port: %s\n", getEnv("PORT", "8080"))
	fmt.Printf("SMTPHost: %s\n", getEnv("SMTP_HOST", "smtp.gmail.com"))
	fmt.Printf("SMTPPort: %s\n", getEnv("SMTP_PORT", "587"))
	fmt.Printf("SMTPUsername: %s\n", getEnv("SMTP_USERNAME", ""))
	fmt.Printf("SMTPPassword: %s\n", getEnv("SMTP_PASSWORD", ""))
	fmt.Printf("FromEmail: %s\n", getEnv("FROM_EMAIL", ""))
	fmt.Printf("RedisURL: %s\n", getEnv("REDIS_URL", ""))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	LoadConfig()
}
