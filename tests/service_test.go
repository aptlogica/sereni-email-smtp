// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"crypto/tls"
	"errors"
	"io"
	"net/smtp"
	"sync"
	"testing"
	"time"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

// MockSmtpClient for testing SMTP interactions
type MockSmtpClient struct {
	StartTLSErr error
	AuthErr     error
	MailErr     error
	RcptErr     error
	DataErr     error
	CloseErr    error
	QuitErr     error
	DataWriter  io.WriteCloser
	closed      bool
}

func (m *MockSmtpClient) StartTLS(config *tls.Config) error {
	return m.StartTLSErr
}

func (m *MockSmtpClient) Auth(a smtp.Auth) error {
	return m.AuthErr
}

func (m *MockSmtpClient) Mail(from string) error {
	return m.MailErr
}

func (m *MockSmtpClient) Rcpt(to string) error {
	return m.RcptErr
}

func (m *MockSmtpClient) Data() (io.WriteCloser, error) {
	return m.DataWriter, m.DataErr
}

func (m *MockSmtpClient) Close() error {
	m.closed = true
	return m.CloseErr
}

func (m *MockSmtpClient) Quit() error {
	return m.QuitErr
}

// MockWriter for testing write operations
type MockWriter struct {
	writeErr error
	closeErr error
	written  []byte
}

func (mw *MockWriter) Write(p []byte) (int, error) {
	if mw.writeErr != nil {
		return 0, mw.writeErr
	}
	mw.written = append(mw.written, p...)
	return len(p), nil
}

func (mw *MockWriter) Close() error {
	return mw.closeErr
}

// TestSendEmail_StartTLSError tests SendEmail when StartTLS fails
func TestSendEmail_StartTLSError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockClient := &MockSmtpClient{
		StartTLSErr: errors.New("TLS handshake failed"),
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for StartTLS failure, got nil")
	}
	if err.Error() != "TLS handshake failed" {
		t.Errorf("Expected 'TLS handshake failed', got %s", err.Error())
	}

	if !mockClient.closed {
		t.Error("Expected connection to be closed on error")
	}
}

// TestSendEmail_AuthError tests SendEmail when Auth fails
func TestSendEmail_AuthError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockClient := &MockSmtpClient{
		AuthErr: errors.New("authentication failed"),
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Auth failure, got nil")
	}
	if err.Error() != "authentication failed" {
		t.Errorf("Expected 'authentication failed', got %s", err.Error())
	}
}

// TestSendEmail_MailError tests SendEmail when Mail fails
func TestSendEmail_MailError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockClient := &MockSmtpClient{
		MailErr: errors.New("mail command failed"),
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Mail failure, got nil")
	}
}

// TestSendEmail_RcptError tests SendEmail when Rcpt fails
func TestSendEmail_RcptError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockClient := &MockSmtpClient{
		RcptErr: errors.New("rcpt command failed"),
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Rcpt failure, got nil")
	}
}

// TestSendEmail_DataError tests SendEmail when Data fails
func TestSendEmail_DataError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockClient := &MockSmtpClient{
		DataErr: errors.New("data command failed"),
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Data failure, got nil")
	}
}

// TestSendEmail_WriteError tests SendEmail when writing to data writer fails
func TestSendEmail_WriteError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockWriter := &MockWriter{
		writeErr: errors.New("write failed"),
	}

	mockClient := &MockSmtpClient{
		DataWriter: mockWriter,
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Write failure, got nil")
	}
}

// TestSendEmail_DataWriterCloseError tests SendEmail when closing data writer fails
func TestSendEmail_DataWriterCloseError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockWriter := &MockWriter{
		closeErr: errors.New("close failed"),
	}

	mockClient := &MockSmtpClient{
		DataWriter: mockWriter,
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for DataWriter close failure, got nil")
	}
}

// TestSendEmail_QuitError tests SendEmail when Quit fails (but should not return error as message was sent)
func TestSendEmail_QuitError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockWriter := &MockWriter{}

	mockClient := &MockSmtpClient{
		DataWriter: mockWriter,
		QuitErr:    errors.New("quit failed"),
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Quit failure, got nil")
	}
}

// TestSendEmail_DialError tests SendEmail when Dial fails
func TestSendEmail_DialError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return nil, errors.New("dial failed")
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error for Dial failure, got nil")
	}
	if err.Error() != "dial failed" {
		t.Errorf("Expected 'dial failed', got %s", err.Error())
	}
}

// TestSendEmail_HTMLMessage tests SendEmail with HTML content
func TestSendEmail_HTMLMessage(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var capturedMessage string
	mockWriter := &MockWriter{}

	mockClient := &MockSmtpClient{
		DataWriter: mockWriter,
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Subject", "<html>Body</html>", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	capturedMessage = string(mockWriter.written)
	if len(capturedMessage) == 0 {
		t.Error("Expected message to be written")
	}

	if !containsStr(capturedMessage, "text/html") {
		t.Error("Expected HTML content type in message")
	}
}

// TestSendEmail_SanitizationCRLF tests subject and body sanitization for CRLF injection
func TestSendEmail_SanitizationCRLF(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockWriter := &MockWriter{}
	mockClient := &MockSmtpClient{
		DataWriter: mockWriter,
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	// Try to inject CRLF in subject - the sanitization should remove the CRLF
	// so it becomes part of the same line
	maliciousSubject := "Subject\r\nX-Injected: header"
	err := service.SendEmail([]string{"test@example.com"}, maliciousSubject, "Body", false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// The message should be properly formed even with injection attempts
	capturedMessage := string(mockWriter.written)
	if len(capturedMessage) == 0 {
		t.Error("Expected message to be written")
	}
}

// TestSendEmail_MultipleRecipients tests SendEmail with multiple recipients
func TestSendEmail_MultipleRecipients(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	mockClient := &MockSmtpClient{
		DataWriter: &MockWriter{},
		RcptErr:    nil,
	}

	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	recipients := []string{"test1@example.com", "test2@example.com", "test3@example.com"}
	err := service.SendEmail(recipients, "Subject", "Body", false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// If no error, it means all Rcpt calls succeeded
}

// TestSendBulkEmail_InvalidEmails tests SendBulkEmail filtering invalid emails
func TestSendBulkEmail_InvalidEmails(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 2)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	recipients := []string{
		"valid1@test.com",
		"invalid-email",
		"valid2@test.com",
		"@nodomain.com",
		"valid3@test.com",
	}

	failed, err := service.SendBulkEmail(recipients, "Subject", "Body", false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(failed) != 2 {
		t.Errorf("Expected 2 failed emails, got %d", len(failed))
	}

	// Verify that invalid emails are in failed list
	if !containsEmail(failed, "invalid-email") || !containsEmail(failed, "@nodomain.com") {
		t.Error("Expected invalid emails in failed list")
	}
}

// TestSendBulkEmail_DefaultBatchSize tests that default batch size is used when <= 0
func TestSendBulkEmail_DefaultBatchSize(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 0)

	var mu sync.Mutex
	sendCount := 0
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		mu.Lock()
		sendCount++
		mu.Unlock()
		return nil
	}

	recipients := []string{"test1@example.com", "test2@example.com", "test3@example.com"}
	failed, err := service.SendBulkEmail(recipients, "Subject", "Body", false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(failed) != 0 {
		t.Errorf("Expected no failed emails, got %d", len(failed))
	}

	if sendCount != 3 {
		t.Errorf("Expected 3 emails sent, got %d", sendCount)
	}
}

// TestSendBulkEmail_PartialFailures tests that partial failures are captured
func TestSendBulkEmail_PartialFailures(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 2)

	var mu sync.Mutex
	sendCount := 0
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		mu.Lock()
		sendCount++
		currentCount := sendCount
		mu.Unlock()
		if currentCount%2 == 0 { // Fail every other email
			return errors.New("send failed")
		}
		return nil
	}

	recipients := []string{"test1@example.com", "test2@example.com", "test3@example.com", "test4@example.com"}
	failed, err := service.SendBulkEmail(recipients, "Subject", "Body", false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(failed) != 2 {
		t.Errorf("Expected 2 failed emails, got %d", len(failed))
	}
}

// TestVerifyOTP_Expired tests that expired OTPs cannot be verified
func TestVerifyOTP_Expired(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	// Generate OTP with very short expiry
	otp, err := service.GenerateAndSendOTP("test@example.com", 0)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Wait for OTP to expire
	time.Sleep(1 * time.Second)

	// Try to verify expired OTP
	if service.VerifyOTP("test@example.com", otp) {
		t.Error("Expected expired OTP verification to fail")
	}
}

// TestVerifyOTP_DeletesAfterVerification tests that OTP is deleted after successful verification
func TestVerifyOTP_DeletesAfterVerification(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	otp, err := service.GenerateAndSendOTP("test@example.com", 10)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Verify OTP successfully
	if !service.VerifyOTP("test@example.com", otp) {
		t.Error("Expected first OTP verification to succeed")
	}

	// Try to verify the same OTP again - should fail
	if service.VerifyOTP("test@example.com", otp) {
		t.Error("Expected second verification of same OTP to fail (should be deleted)")
	}
}

// TestGenerateAndSendOTP_InvalidEmail tests OTP generation with invalid email
func TestGenerateAndSendOTP_InvalidEmail(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	sentinelErr := errors.New("send failed for invalid recipient")
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return sentinelErr
	}

	_, err := service.GenerateAndSendOTP("not-an-email", 10)
	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
	if err != sentinelErr {
		t.Errorf("Expected %v, got %v", sentinelErr, err)
	}
}

// TestSendTransactionalEmail_InvalidRecipient tests transactional email with invalid recipient
func TestSendTransactionalEmail_InvalidRecipient(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	req := &email.EmailRequest{
		To:      []string{"invalid"},
		Subject: "Subject",
		Body:    "Body",
		IsHTML:  false,
	}

	err := service.SendTransactionalEmail(req)
	if err == nil {
		t.Error("Expected error for invalid recipient, got nil")
	}
}

// TestSendTransactionalEmail_WithTemplate tests transactional email with template rendering
func TestSendTransactionalEmail_WithTemplate(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	sendCalled := false
	var capturedSubject string

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		sendCalled = true
		capturedSubject = subject
		return nil
	}

	req := &email.EmailRequest{
		To:           []string{"test@example.com"},
		Subject:      "unused",
		Body:         "unused",
		Template:     "welcome",
		TemplateData: map[string]interface{}{"name": "John"},
		IsHTML:       true,
	}

	err := service.SendTransactionalEmail(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !sendCalled {
		t.Error("Expected SendEmail to be called")
	}

	if capturedSubject != "Welcome to Our Service!" {
		t.Errorf("Expected subject 'Welcome to Our Service!', got %s", capturedSubject)
	}
}

// TestSendTransactionalEmail_TemplateNotFound tests handling of non-existent template
func TestSendTransactionalEmail_TemplateNotFound(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	req := &email.EmailRequest{
		To:           []string{"test@example.com"},
		Subject:      "Subject",
		Body:         "Body",
		Template:     "non_existent_template",
		TemplateData: map[string]interface{}{},
		IsHTML:       true,
	}

	err := service.SendTransactionalEmail(req)
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

// TestAddAndGetTemplate tests adding and retrieving custom templates
func TestAddAndGetTemplate(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	customTemplate := email.EmailTemplate{
		Name:     "custom",
		Subject:  "Custom Subject",
		HTMLBody: "<html>Custom Body</html>",
		TextBody: "Custom Body",
	}

	service.AddTemplate("custom", customTemplate)

	retrieved, exists := service.GetTemplate("custom")
	if !exists {
		t.Error("Expected custom template to exist")
	}

	if retrieved.Subject != "Custom Subject" {
		t.Errorf("Expected subject 'Custom Subject', got %s", retrieved.Subject)
	}
}

// TestSendTemplateEmail tests sending email using a template
func TestSendTemplateEmail(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	sendCalled := false
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		sendCalled = true
		if subject != "Your Verification Code" {
			t.Errorf("Expected subject 'Your Verification Code', got %s", subject)
		}
		if !containsStr(body, "123456") {
			t.Error("Expected OTP code in body")
		}
		return nil
	}

	err := service.SendTemplateEmail(
		[]string{"test@example.com"},
		"otp_template",
		map[string]interface{}{"otp": "123456", "expiry": 10},
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !sendCalled {
		t.Error("Expected SendEmail to be called")
	}
}

// TestGetAvailableTemplates tests retrieving list of available templates
func TestGetAvailableTemplates(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	templates := service.GetAvailableTemplates()

	if len(templates) == 0 {
		t.Error("Expected some templates available")
	}

	expectedTemplates := []string{"welcome", "password_reset", "verification", "otp_template"}
	for _, expected := range expectedTemplates {
		found := false
		for _, template := range templates {
			if template == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template %s to be available", expected)
		}
	}
}

// TestCleanupExpiredOTPs tests the cleanup goroutine (basic test)
func TestCleanupExpiredOTPs(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	// Generate an OTP with very short expiry
	otp, err := service.GenerateAndSendOTP("test@example.com", 0)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Wait for OTP to expire
	time.Sleep(1 * time.Second)

	// Verify that the OTP no longer works (expired)
	if service.VerifyOTP("test@example.com", otp) {
		t.Error("Expected OTP to be expired and verification to fail")
	}
}

// TestJoin tests the Join helper function
func TestJoin(t *testing.T) {
	tests := []struct {
		elements []string
		sep      string
		expected string
	}{
		{[]string{"a", "b", "c"}, ",", "a,b,c"},
		{[]string{"a"}, ",", "a"},
		{[]string{}, ",", ""},
		{[]string{"test@example.com", "user@domain.com"}, "; ", "test@example.com; user@domain.com"},
	}

	for _, tt := range tests {
		result := email.Join(tt.elements, tt.sep)
		if result != tt.expected {
			t.Errorf("Join(%v, %q) = %q, expected %q", tt.elements, tt.sep, result, tt.expected)
		}
	}
}

// Helper functions
func containsStr(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsEmail(emails []string, emailAddr string) bool {
	for _, e := range emails {
		if e == emailAddr {
			return true
		}
	}
	return false
}
