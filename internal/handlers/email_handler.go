// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aptlogica/sereni-email-smtp/internal/email"

	"github.com/gin-gonic/gin"
)

var handlerHeaderInjectionPattern = regexp.MustCompile(`[\r\n\x00]`)

func stripCRLF(s string) string {
	return handlerHeaderInjectionPattern.ReplaceAllString(strings.TrimSpace(s), "")
}

type EmailHandler struct {
	Service *email.EmailService
}

func NewEmailHandler(service *email.EmailService) *EmailHandler {
	return &EmailHandler{Service: service}
}

// SendEmail sends a transactional email
// @Summary Send a transactional email
// @Description Send a single email to a recipient
// @Tags email
// @Accept json
// @Produce json
// @Param request body email.EmailRequest true "Email Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/email/send [post]
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req email.EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.SendTransactionalEmail(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send email: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email sent successfully",
	})
}

// SendBulkEmail sends bulk emails
// @Summary Send bulk emails
// @Description Send emails to multiple recipients
// @Tags email
// @Accept json
// @Produce json
// @Param request body email.BulkEmailRequest true "Bulk Email Request"
// @Success 200 {object} email.EmailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/email/send-bulk [post]
func (h *EmailHandler) SendBulkEmail(c *gin.Context) {
	var req email.BulkEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	failedEmails, err := h.Service.SendBulkEmail(req.Recipients, req.Subject, req.Body, req.IsHTML)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send bulk emails: " + err.Error(),
		})
		return
	}

	response := email.EmailResponse{
		Success: len(failedEmails) == 0,
		Message: "Bulk emails processed",
	}

	if len(failedEmails) > 0 {
		response.FailedEmails = failedEmails
		response.Message += fmt.Sprintf(" with %d failed emails", len(failedEmails))
	}

	c.JSON(http.StatusOK, response)
}

// GenerateOTP generates and sends an OTP
// @Summary Generate and send OTP
// @Description Generate a One-Time Password and send it to the specified email
// @Tags otp
// @Accept json
// @Produce json
// @Param request body email.OTPRequest true "OTP Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/email/otp/generate [post]
func (h *EmailHandler) GenerateOTP(c *gin.Context) {
	var req email.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use default expiry of 5 minutes if not specified
	expiry := req.Expiry
	if expiry <= 0 {
		expiry = 5
	}

	// Generate OTP and store, but send email asynchronously
	go h.Service.GenerateAndSendOTP(req.To, expiry)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP sent successfully",
	})
}

// VerifyOTP verifies an OTP
// @Summary Verify OTP
// @Description Verify a One-Time Password for an email
// @Tags otp
// @Accept json
// @Produce json
// @Param request body email.OTPVerificationRequest true "OTP Verification Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/email/otp/verify [post]
func (h *EmailHandler) VerifyOTP(c *gin.Context) {
	var req email.OTPVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValid := h.Service.VerifyOTP(req.Email, req.OTP)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid or expired OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP verified successfully",
	})
}

// HealthCheck checks the service health
// @Summary Health Check
// @Description Check if the service is up and running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *EmailHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "email-service",
	})
}
