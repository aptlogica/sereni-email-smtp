// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

// Package main demonstrates secure email sending patterns
// This example shows how to avoid common security vulnerabilities
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
	"github.com/gin-gonic/gin"
)

var (
	emailService *email.EmailService
	baseURL      string
	trustedDomains []string
)

// init loads secure configuration and initializes the email service
func init() {
	// Load base URL from environment variable (NOT from HTTP requests)
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL environment variable must be set (e.g., https://example.com)")
	}

	// Load trusted domains from configuration
	trustedDomainsStr := os.Getenv("TRUSTED_DOMAINS")
	if trustedDomainsStr == "" {
		log.Fatal("TRUSTED_DOMAINS environment variable must be set")
	}
	trustedDomains = strings.Split(trustedDomainsStr, ",")

	// Initialize email service
	emailService = email.NewEmailService(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("FROM_EMAIL"),
		10,
	)

	// Configure trusted domains for URL validation
	emailService.SetTrustedDomains(trustedDomains, false) // HTTPS only

	log.Printf("Email service initialized with base URL: %s", baseURL)
	log.Printf("Trusted domains: %v", trustedDomains)
}

func main() {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/auth/password-reset/request", handlePasswordResetRequest)
		api.POST("/auth/password-reset/confirm", handlePasswordResetConfirm)
		api.POST("/auth/email-verification/send", handleSendVerification)
		api.POST("/auth/email-verification/verify", handleVerifyEmail)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// handlePasswordResetRequest demonstrates SECURE password reset
// ✅ Uses base URL from configuration, NOT from HTTP Host header
func handlePasswordResetRequest(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate email exists (in production, check database)
	// For this example, we'll accept any valid email format

	// Generate secure reset token
	token, err := generateSecureToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Store token in database (not shown here)
	// storeResetToken(req.Email, token, time.Now().Add(24*time.Hour))

	// ✅ SECURE: Construct reset URL from TRUSTED configuration
	// This prevents host header injection attacks
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	// Send password reset email with validated URL
	err = emailService.SendTemplateEmail(
		[]string{req.Email},
		"password_reset",
		map[string]interface{}{
			"reset_url": resetURL, // ✅ Safe - from config, will be validated
		},
	)

	if err != nil {
		log.Printf("Error sending password reset email: %v", err)
		// Don't reveal whether email exists or not
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// PasswordResetConfirm represents a password reset confirmation
type PasswordResetConfirm struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// handlePasswordResetConfirm validates token and resets password
func handlePasswordResetConfirm(c *gin.Context) {
	var req PasswordResetConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate token (check database, expiry, etc.)
	// For this example, we'll assume validation passes

	// Reset password in database
	// updatePassword(email, req.NewPassword)

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

// EmailVerificationRequest represents an email verification request
type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// handleSendVerification sends email verification link
// ✅ Uses base URL from configuration
func handleSendVerification(c *gin.Context) {
	var req EmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate verification token
	token, err := generateSecureToken()
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	// Store token in database
	// storeVerificationToken(req.Email, token)

	// ✅ SECURE: Construct verification URL from TRUSTED configuration
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", baseURL, token)

	// Send verification email
	err = emailService.SendTemplateEmail(
		[]string{req.Email},
		"verification",
		map[string]interface{}{
			"verification_url": verificationURL, // ✅ Safe - from config
		},
	)

	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// EmailVerificationConfirm represents email verification confirmation
type EmailVerificationConfirm struct {
	Token string `json:"token" binding:"required"`
}

// handleVerifyEmail verifies email with token
func handleVerifyEmail(c *gin.Context) {
	var req EmailVerificationConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate token (check database)
	// markEmailVerified(email)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ❌❌❌ VULNERABLE EXAMPLE - DO NOT USE ❌❌❌
// This function demonstrates the WRONG way to handle password reset
// It's included here for educational purposes only
func vulnerablePasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ❌ VULNERABLE: Using Host header from HTTP request
	// An attacker can send: Host: malicious.com
	// The victim will receive an email with a link to malicious.com
	host := c.Request.Host // ❌ UNTRUSTED INPUT!

	token, _ := generateSecureToken()

	// ❌ VULNERABLE: Constructing URL with untrusted host
	resetURL := fmt.Sprintf("https://%s/reset-password?token=%s", host, token)

	// Even though the email service will validate this URL,
	// if trusted domains are not configured properly,
	// this could lead to a security breach

	emailService.SendTemplateEmail(
		[]string{req.Email},
		"password_reset",
		map[string]interface{}{
			"reset_url": resetURL, // ❌ DANGEROUS - contains attacker-controlled data
		},
	)

	// The email service will REJECT this if the host is not in trusted domains
	// But it's better to not construct such URLs in the first place
}

// ❌❌❌ END VULNERABLE EXAMPLE ❌❌❌
