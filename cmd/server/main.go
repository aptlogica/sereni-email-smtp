// cmd/server/main.go
package main

import (
	"email-service/internal/config"
	"email-service/internal/email"
	"email-service/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize email service
	emailService := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.FromEmail,
	)

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", emailHandler.HealthCheck)

	// Email endpoints
	emailGroup := r.Group("/api/v1/email")
	{
		emailGroup.POST("/send", emailHandler.SendEmail)
		emailGroup.POST("/bulk", emailHandler.SendBulkEmail)
		emailGroup.POST("/otp/generate", emailHandler.GenerateOTP)
		emailGroup.POST("/otp/verify", emailHandler.VerifyOTP)
	}

	log.Printf("Email microservice starting on port %s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
