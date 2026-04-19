// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

// Security-focused tests for email injection prevention
package test

import (
	"crypto/tls"
	"io"
	"net/smtp"
	"strings"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

// mockSmtpClient for testing
type mockSmtpClient struct {
	startTLSFunc func(*tls.Config) error
	authFunc     func(smtp.Auth) error
	mailFunc     func(string) error
	rcptFunc     func(string) error
	dataFunc     func() (io.WriteCloser, error)
	closeFunc    func() error
	quitFunc     func() error
}

func (m *mockSmtpClient) StartTLS(config *tls.Config) error {
	if m.startTLSFunc != nil {
		return m.startTLSFunc(config)
	}
	return nil
}

func (m *mockSmtpClient) Auth(a smtp.Auth) error {
	if m.authFunc != nil {
		return m.authFunc(a)
	}
	return nil
}

func (m *mockSmtpClient) Mail(from string) error {
	if m.mailFunc != nil {
		return m.mailFunc(from)
	}
	return nil
}

func (m *mockSmtpClient) Rcpt(to string) error {
	if m.rcptFunc != nil {
		return m.rcptFunc(to)
	}
	return nil
}

func (m *mockSmtpClient) Data() (io.WriteCloser, error) {
	if m.dataFunc != nil {
		return m.dataFunc()
	}
	return &mockWriteCloser{}, nil
}

func (m *mockSmtpClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func (m *mockSmtpClient) Quit() error {
	if m.quitFunc != nil {
		return m.quitFunc()
	}
	return nil
}

type mockWriteCloser struct{}

func (m *mockWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriteCloser) Close() error {
	return nil
}

// TestEmailInjection_HeaderInjectionPrevention tests that CRLF characters in headers are sanitized
func TestEmailInjection_HeaderInjectionPrevention(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var capturedMessage string
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		// Capture what would be sent
		capturedMessage = subject + "|" + strings.Join(to, ",") + "|" + body
		return nil
	}

	// Test CRLF injection in subject
	maliciousSubject := "Test\r\nBcc: attacker@evil.com"
	err := service.SendEmail([]string{"victim@example.com"}, maliciousSubject, "Body", false)
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}
	if strings.Contains(capturedMessage, "\r") || strings.Contains(capturedMessage, "\n") {
		t.Errorf("CRLF characters not sanitized from subject: %s", capturedMessage)
	}
	// After sanitization, the subject should be "TestBcc: attacker@evil.com" which is harmless

	// Test CRLF injection in recipient
	maliciousRecipient := "victim@example.com\r\nBcc: attacker@evil.com"
	err = service.SendEmail([]string{maliciousRecipient}, "Subject", "Body", false)
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}
	// After sanitization, net/mail will typically remove the invalid parts or quote them
	// In our case, it should at least remove the CRLF
	if strings.Contains(capturedMessage, "\r") || strings.Contains(capturedMessage, "\n") {
		t.Errorf("CRLF characters not sanitized from recipient: %s", capturedMessage)
	}
	// Note: net/mail.ParseAddress might return "<victim@example.comBcc:attacker@evil.com>"
	// because it strips control characters via sanitizeHeader fallback if parsing fails,
	// or it might just error. Our implementation currently falls back to sanitizeHeader if ParseAddress fails.
}

// TestEmailInjection_BodyControlCharactersPrevention tests that control characters in body are removed
func TestEmailInjection_BodyControlCharactersPrevention(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var capturedBody string
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		capturedBody = body
		return nil
	}

	// Test CRLF and null bytes in body
	maliciousBody := "Click here: http://legit.com\r\n\r\nBcc: attacker@evil.com\r\nContent-Type: text/html\r\n\r\n<script>alert('xss')</script>\x00"
	err := service.SendEmail([]string{"victim@example.com"}, "Subject", maliciousBody, false)
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}

	// Body should have null bytes removed, but CRLF should REMAIN
	if !strings.Contains(capturedBody, "\r\n") {
		t.Errorf("CRLF characters SHOULD remain in body for multi-line support: %q", capturedBody)
	}
	if strings.Contains(capturedBody, "\x00") {
		t.Errorf("Null bytes not sanitized from body: %q", capturedBody)
	}
	// After sanitization, "Bcc:" text may remain but without CRLF it's harmless
}

// TestEmailInjection_FromEmailSanitization tests that FromEmail is sanitized
func TestEmailInjection_FromEmailSanitization(t *testing.T) {
	// Service with malicious FromEmail
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com\r\nBcc: attacker@evil.com", 5)

	mailCalled := false
	var capturedFrom string

	// Create a mock client
	service.Dial = func(addr string) (email.SmtpClient, error) {
		return &mockSmtpClient{
			mailFunc: func(from string) error {
				mailCalled = true
				capturedFrom = from
				return nil
			},
		}, nil
	}

	err := service.SendEmail([]string{"victim@example.com"}, "Subject", "Body", false)
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}

	if !mailCalled {
		t.Error("Mail() was not called")
	}

	if strings.Contains(capturedFrom, "\r") || strings.Contains(capturedFrom, "\n") {
		t.Errorf("FromEmail CRLF not sanitized: %q", capturedFrom)
	}
	// After sanitization, "Bcc:" text may remain but without CRLF it's harmless
}

// TestEmailInjection_EmptyFromEmailAfterSanitization tests rejection of empty sender
func TestEmailInjection_EmptyFromEmailAfterSanitization(t *testing.T) {
	// Service with FromEmail that becomes empty after sanitization
	service := email.NewEmailService("localhost", 587, "user", "pass", "\r\n", 5)

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return &mockSmtpClient{}, nil
	}

	err := service.SendEmail([]string{"victim@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for empty sender after sanitization, got nil")
	}
	if !strings.Contains(err.Error(), "invalid or empty sender") {
		t.Errorf("Expected 'invalid or empty sender' error, got: %v", err)
	}
}

// TestEmailInjection_TransactionalEmailValidation tests SendTransactionalEmail validates all recipients
func TestEmailInjection_TransactionalEmailValidation(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	// Test with invalid recipient (no template)
	req := &email.EmailRequest{
		To:      []string{"invalid-email"},
		Subject: "Test",
		Body:    "Test body",
		IsHTML:  false,
	}
	err := service.SendTransactionalEmail(req)
	if err == nil {
		t.Error("Expected validation error for invalid recipient, got nil")
	}
	if !strings.Contains(err.Error(), "invalid recipient") {
		t.Errorf("Expected 'invalid recipient' error, got: %v", err)
	}

	// Test with invalid recipient (with template)
	req.Template = "test_template"
	err = service.SendTransactionalEmail(req)
	if err == nil {
		t.Error("Expected validation error for invalid recipient with template, got nil")
	}
}

// TestEmailInjection_TransactionalEmailSanitization tests that subject/body are sanitized
func TestEmailInjection_TransactionalEmailSanitization(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var capturedSubject, capturedBody string
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		capturedSubject = subject
		capturedBody = body
		return nil
	}

	req := &email.EmailRequest{
		To:      []string{"victim@example.com"},
		Subject: "Test\r\nBcc: attacker@evil.com",
		Body:    "Body\r\n\r\nBcc: attacker@evil.com\x00",
		IsHTML:  false,
	}

	err := service.SendTransactionalEmail(req)
	if err != nil {
		t.Fatalf("SendTransactionalEmail failed: %v", err)
	}

	if strings.Contains(capturedSubject, "\r") || strings.Contains(capturedSubject, "\n") {
		t.Errorf("Subject CRLF not sanitized: %q", capturedSubject)
	}
	if !strings.Contains(capturedBody, "\r\n") {
		t.Errorf("Body CRLF SHOULD remain for multi-line support: %q", capturedBody)
	}
	if strings.Contains(capturedBody, "\x00") {
		t.Errorf("Body null bytes not sanitized: %q", capturedBody)
	}
}

// TestEmailInjection_GenerateAndSendOTPValidation tests OTP email validation
func TestEmailInjection_GenerateAndSendOTPValidation(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	// Test with invalid email
	_, err := service.GenerateAndSendOTP("not-an-email", 10)
	if err == nil {
		t.Error("Expected validation error for invalid email, got nil")
	}
	if !strings.Contains(err.Error(), "invalid email") {
		t.Errorf("Expected 'invalid email' error, got: %v", err)
	}

	// Test with email containing CRLF (should fail validation)
	_, err = service.GenerateAndSendOTP("test@example.com\r\nBcc: attacker@evil.com", 10)
	if err == nil {
		t.Error("Expected validation error for malicious email, got nil")
	}
}

// TestEmailInjection_MultipleRecipientsValidation tests that all recipients are validated
func TestEmailInjection_MultipleRecipientsValidation(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var capturedRecipients []string
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		capturedRecipients = to
		return nil
	}

	// Test with mixed valid and malicious recipients
	recipients := []string{
		"valid1@example.com",
		"valid2@example.com\r\nBcc: attacker@evil.com",
		"valid3@example.com",
	}

	err := service.SendEmail(recipients, "Subject", "Body", false)
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}

	// All recipients should have CRLF sanitized
	for _, recipient := range capturedRecipients {
		if strings.Contains(recipient, "\r") || strings.Contains(recipient, "\n") {
			t.Errorf("Recipient CRLF not sanitized: %q", recipient)
		}
	}
}
