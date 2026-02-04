// cmd/server/main.go
package main

import (
	_ "email-service/docs" // Import generated docs
	"email-service/internal/config"
	"email-service/internal/email"
	"email-service/internal/handlers"
	"email-service/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Serenibase Email Service API
// @version 1.0
// @description API for the Email Service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	cfg := config.LoadConfig()

	// Initialize email service
	emailService := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.FromEmail,
		cfg.BulkBatchSize,
	)

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigin))

	// Health check endpoint
	r.GET("/health", emailHandler.HealthCheck)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Email endpoints
	emailGroup := r.Group("/api/v1/email")
	{
		emailGroup.POST("/send", emailHandler.SendEmail)
		emailGroup.POST("/bulk", emailHandler.SendBulkEmail)
		emailGroup.POST("/otp/generate", emailHandler.GenerateOTP)
		emailGroup.POST("/otp/verify", emailHandler.VerifyOTP)
	}

	addr := cfg.Host + ":" + cfg.Port
	log.Printf("Email microservice starting on %s", addr)
	log.Fatal(r.Run(addr))
}
