/*
Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
This file is part of software developed by Aptlogica Technologies Private Limited.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
	http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Websites:
https://www.aptlogica.com
https://www.serenibase.com
Support:
support@aptlogica.com
support@serenibase.com
*/

package main

import (
	"log"

	_ "github.com/aptlogica/sereni-email-smtp/docs" // Import generated docs
	"github.com/aptlogica/sereni-email-smtp/internal/config"
	"github.com/aptlogica/sereni-email-smtp/internal/email"
	"github.com/aptlogica/sereni-email-smtp/internal/handlers"
	"github.com/aptlogica/sereni-email-smtp/pkg/middleware"

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
	cfg := config.LoadConfig() // validateSecrets is called inside LoadConfig

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
