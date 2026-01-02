// internal/handlers/email_handler.go (updated)
package handlers

import (
	"email-service/internal/email"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	Service *email.EmailService
}

func NewEmailHandler(service *email.EmailService) *EmailHandler {
	return &EmailHandler{Service: service}
}

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

	otp, err := h.Service.GenerateAndSendOTP(req.To, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate OTP: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP sent successfully",
		"otp":     otp, // In production, don't return OTP in response
	})
}

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

func (h *EmailHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "email-service",
	})
}
