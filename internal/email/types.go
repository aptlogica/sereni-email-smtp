// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
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
	Attachments  []Attachment           `json:"attachments,omitempty"`
}

// Attachment is one file attached to a transactional email. Content travels
// as base64 (rather than a storage reference) so this service - which has no
// credentials or client for storage-provider - never needs to fetch bytes
// itself; the caller (e.g. trigger-processor) resolves the reference and
// reads the file before building the request.
type Attachment struct {
	Filename      string `json:"filename" binding:"required"`
	ContentType   string `json:"content_type,omitempty"`
	ContentBase64 string `json:"content_base64" binding:"required"`
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
