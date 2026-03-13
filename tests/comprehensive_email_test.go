package test

import (
	"email-service/internal/email"
	"errors"
	"testing"
)

func TestGenerateOTP_Comprehensive(t *testing.T) {
	for i := 0; i < 100; i++ {
		otp := email.GenerateOTP()
		if len(otp) != 6 {
			t.Errorf("Expected OTP length 6, got %d", len(otp))
		}
		if !email.ValidateOTPFormat(otp) {
			t.Errorf("Generated OTP %s is not valid format", otp)
		}
	}
}

func TestIsValidEmail_Comprehensive(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"123@456.com",
		"a@b.c",
		"test+tag@example.org",
	}

	for _, emailAddr := range validEmails {
		if !email.IsValidEmail(emailAddr) {
			t.Errorf("Expected %s to be valid", emailAddr)
		}
	}

	invalidEmails := []string{
		"notanemail",
		"@domain.com",
		"user@",
		"use r@domain.com",
		"",
		"user@@domain.com",
		"user@domain.",
	}

	for _, emailAddr := range invalidEmails {
		if email.IsValidEmail(emailAddr) {
			t.Errorf("Expected %s to be invalid", emailAddr)
		}
	}
}

func TestIsValidEmailList_Comprehensive(t *testing.T) {
	validList := []string{"a@b.com", "test@example.org", "user@domain.net"}
	if !email.IsValidEmailList(validList) {
		t.Error("Expected valid email list to be valid")
	}

	invalidList := []string{"a@b.com", "invalid", "test@example.org"}
	if email.IsValidEmailList(invalidList) {
		t.Error("Expected list with invalid email to be invalid")
	}

	emptyList := []string{}
	if !email.IsValidEmailList(emptyList) {
		t.Error("Expected empty list to be valid")
	}
}

func TestValidateOTPFormat_Comprehensive(t *testing.T) {
	validOTPs := []string{"123456", "000000", "999999"}
	for _, otp := range validOTPs {
		if !email.ValidateOTPFormat(otp) {
			t.Errorf("Expected %s to be valid OTP format", otp)
		}
	}

	invalidOTPs := []string{"12345", "1234567", "abc123", "12345a", ""}
	for _, otp := range invalidOTPs {
		if email.ValidateOTPFormat(otp) {
			t.Errorf("Expected %s to be invalid OTP format", otp)
		}
	}
}

func TestNewEmailService_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 15)

	if service.SMTPHost != "localhost" {
		t.Errorf("Expected SMTPHost localhost, got %s", service.SMTPHost)
	}
	if service.SMTPPort != 587 {
		t.Errorf("Expected SMTPPort 587, got %d", service.SMTPPort)
	}
	if service.SMTPUsername != "user" {
		t.Errorf("Expected SMTPUsername user, got %s", service.SMTPUsername)
	}
	if service.SMTPPassword != "pass" {
		t.Errorf("Expected SMTPPassword pass, got %s", service.SMTPPassword)
	}
	if service.FromEmail != "from@test.com" {
		t.Errorf("Expected FromEmail from@test.com, got %s", service.FromEmail)
	}
	if service.BulkBatchSize != 15 {
		t.Errorf("Expected BulkBatchSize 15, got %d", service.BulkBatchSize)
	}
}

func TestSendEmail_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	// Test with SendEmailFunc override
	var lastTo []string
	var lastSubject, lastBody string
	var lastIsHTML bool

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		lastTo = to
		lastSubject = subject
		lastBody = body
		lastIsHTML = isHTML
		return nil
	}

	err := service.SendEmail([]string{"test@example.com"}, "Test Subject", "Test Body", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(lastTo) != 1 || lastTo[0] != "test@example.com" {
		t.Errorf("Expected to be [test@example.com], got %v", lastTo)
	}
	if lastSubject != "Test Subject" {
		t.Errorf("Expected subject Test Subject, got %s", lastSubject)
	}
	if lastBody != "Test Body" {
		t.Errorf("Expected body Test Body, got %s", lastBody)
	}
	if !lastIsHTML {
		t.Error("Expected isHTML to be true")
	}

	// Test with SendEmailFunc returning error
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return errors.New("send failed")
	}

	err = service.SendEmail([]string{"test@example.com"}, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "send failed" {
		t.Errorf("Expected 'send failed', got %s", err.Error())
	}
}

func TestSendBulkEmail_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 2)

	// Test with SendBulkEmailFunc override
	service.SendBulkEmailFunc = func(recipients []string, subject, body string, isHTML bool) ([]string, error) {
		// Simulate some failures
		failed := []string{}
		for _, r := range recipients {
			if !email.IsValidEmail(r) {
				failed = append(failed, r)
			}
		}
		return failed, nil
	}

	recipients := []string{"valid@test.com", "invalid", "another@valid.com"}
	failed, err := service.SendBulkEmail(recipients, "Subject", "Body", true)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(failed) != 1 || failed[0] != "invalid" {
		t.Errorf("Expected failed to be [invalid], got %v", failed)
	}

	// Test with SendBulkEmailFunc returning error
	service.SendBulkEmailFunc = func(recipients []string, subject, body string, isHTML bool) ([]string, error) {
		return nil, errors.New("bulk send failed")
	}

	_, err = service.SendBulkEmail(recipients, "Subject", "Body", false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGenerateAndSendOTP_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	// Test successful OTP generation
	otp, err := service.GenerateAndSendOTP("test@example.com", 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !email.ValidateOTPFormat(otp) {
		t.Errorf("Generated OTP %s is not valid format", otp)
	}

	// Test with send email error
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return errors.New("email send failed")
	}

	_, err = service.GenerateAndSendOTP("test@example.com", 10)
	if err == nil {
		t.Error("Expected error when email send fails, got nil")
	}

	// Test with invalid email
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	_, err = service.GenerateAndSendOTP("invalid", 10)
	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
}

func TestVerifyOTP_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}

	// Generate OTP and verify
	otp, err := service.GenerateAndSendOTP("test@example.com", 1)
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// Verify correct OTP
	if !service.VerifyOTP("test@example.com", otp) {
		t.Error("Expected OTP verification to succeed")
	}

	// Verify wrong OTP
	if service.VerifyOTP("test@example.com", "000000") {
		t.Error("Expected wrong OTP verification to fail")
	}

	// Verify OTP for wrong email
	if service.VerifyOTP("wrong@example.com", otp) {
		t.Error("Expected OTP verification for wrong email to fail")
	}
}

func TestSendTransactionalEmail_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var lastSubject string
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		lastSubject = subject
		return nil
	}

	// Test basic transactional email
	req := &email.EmailRequest{
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
		IsHTML:  true,
	}

	err := service.SendTransactionalEmail(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if lastSubject != "Test Subject" {
		t.Errorf("Expected subject Test Subject, got %s", lastSubject)
	}

	// Test with invalid email
	req.To = []string{"invalid"}
	err = service.SendTransactionalEmail(req)
	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
}

func TestJoin_Comprehensive(t *testing.T) {
	// Test normal join
	items := []string{"a", "b", "c"}
	result := email.Join(items, ", ")
	expected := "a, b, c"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test single item
	single := []string{"alone"}
	result = email.Join(single, ", ")
	if result != "alone" {
		t.Errorf("Expected alone, got %s", result)
	}

	// Test empty slice
	empty := []string{}
	result = email.Join(empty, ", ")
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}
