// Copyright (c) 2026 Aptlogica Technologies Private Limited
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package email

type EmailRequest struct {
	To           []string               `json:"to" binding:"required"`
	Subject      string                 `json:"subject" binding:"required"`
	Body         string                 `json:"body" binding:"required"`
	IsHTML       bool                   `json:"is_html"`
	Template     string                 `json:"template"`
	TemplateData map[string]interface{} `json:"template_data"`
}

type BulkEmailRequest struct {
	Recipients   []string               `json:"recipients" binding:"required"`
	Subject      string                 `json:"subject" binding:"required"`
	Body         string                 `json:"body" binding:"required"`
	IsHTML       bool                   `json:"is_html"`
	Template     string                 `json:"template"`
	TemplateData map[string]interface{} `json:"template_data"`
}

type OTPRequest struct {
	To     string `json:"to" binding:"required"`
	Expiry int    `json:"expiry"`
}

type EmailResponse struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	FailedEmails []string `json:"failed_emails,omitempty"`
}

type OTPVerificationRequest struct {
	Email string `json:"email" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}
