// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"strings"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

// TestRenderTemplate_Basic tests basic template rendering
func TestRenderTemplate_Basic(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	subject, body, err := service.RenderTemplate("welcome", map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if subject != "Welcome to Sereni LRS" {
		t.Errorf("Expected subject 'Welcome to Sereni LRS', got %s", subject)
	}

	if !strings.Contains(body, "Alice") {
		t.Error("Expected name 'Alice' to be rendered in body")
	}
}

// TestRenderTemplate_AllPredefinedTemplates tests all predefined templates
func TestRenderTemplate_AllPredefinedTemplates(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	templates := []string{"welcome", "password_reset", "verification", "otp_template"}
	data := map[string]interface{}{
		"name":             "Test User",
		"reset_url":        "https://example.com/reset",
		"verification_url": "https://example.com/verify",
		"otp":              "123456",
		"expiry":           10,
	}

	for _, tmplName := range templates {
		subject, body, err := service.RenderTemplate(tmplName, data)
		if err != nil {
			t.Errorf("Error rendering %s: %v", tmplName, err)
		}

		if subject == "" {
			t.Errorf("Expected non-empty subject for %s", tmplName)
		}

		if body == "" {
			t.Errorf("Expected non-empty body for %s", tmplName)
		}
	}
}

// TestRenderTemplate_NonExistent tests rendering non-existent template
func TestRenderTemplate_NonExistent(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	_, _, err := service.RenderTemplate("does_not_exist", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' in error, got %s", err.Error())
	}
}

// TestRenderTemplate_EmptyData tests rendering with empty data
func TestRenderTemplate_EmptyData(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	subject, body, err := service.RenderTemplate("welcome", map[string]interface{}{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if subject != "Welcome to Sereni LRS" {
		t.Errorf("Expected subject 'Welcome to Sereni LRS', got %s", subject)
	}

	if body == "" {
		t.Error("Expected non-empty body even with empty data")
	}
}

// TestRenderTemplate_SpecialCharacters tests rendering with special characters
func TestRenderTemplate_SpecialCharacters(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	data := map[string]interface{}{
		"name": "John & Jane",
	}

	_, body, err := service.RenderTemplate("welcome", data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(body, "John") && !strings.Contains(body, "Jane") {
		t.Error("Expected name to appear in body")
	}
}

// TestRenderTemplate_Caching tests that templates are cached
func TestRenderTemplate_Caching(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	data := map[string]interface{}{"name": "Test"}

	// Render twice
	_, _, err1 := service.RenderTemplate("welcome", data)
	_, _, err2 := service.RenderTemplate("welcome", data)

	if err1 != nil || err2 != nil {
		t.Errorf("Expected no error, got %v and %v", err1, err2)
	}

	// Both should succeed (first one creates cache, second uses cache)
}

// TestAddTemplate tests adding a custom template
func TestAddTemplate_Custom(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	newTemplate := email.EmailTemplate{
		Name:     "promo",
		Subject:  "Special Offer: {{.discount}}% off!",
		HTMLBody: "<html><p>Get {{.discount}}% off on your next purchase!</p></html>",
		TextBody: "Get {{.discount}}% off on your next purchase!",
	}

	service.AddTemplate("promo", newTemplate)

	// Verify it can be retrieved
	retrieved, exists := service.GetTemplate("promo")
	if !exists {
		t.Error("Expected custom template to be added")
	}

	if retrieved.Subject != "Special Offer: {{.discount}}% off!" {
		t.Errorf("Expected correct subject, got %s", retrieved.Subject)
	}

	// Verify it can be rendered
	subject, _, err := service.RenderTemplate("promo", map[string]interface{}{"discount": 50})
	if err != nil {
		t.Errorf("Expected no error rendering custom template, got %v", err)
	}

	if subject != "Special Offer: 50% off!" {
		t.Errorf("Expected 'Special Offer: 50%% off!', got %s", subject)
	}
}

// TestAddTemplate_Override tests overriding an existing template
func TestAddTemplate_Override(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	// Get original welcome template
	originalWelcome, _ := service.GetTemplate("welcome")
	if originalWelcome.Subject != "Welcome to Sereni LRS" {
		t.Error("Expected original welcome template")
	}

	// Override it
	newWelcome := email.EmailTemplate{
		Name:     "welcome",
		Subject:  "Welcome!",
		HTMLBody: "<html>Welcome!</html>",
		TextBody: "Welcome!",
	}

	service.AddTemplate("welcome", newWelcome)

	// Verify it was overridden
	retrieved, _ := service.GetTemplate("welcome")
	if retrieved.Subject != "Welcome!" {
		t.Errorf("Expected overridden subject 'Welcome!', got %s", retrieved.Subject)
	}
}

// TestGetTemplate_NonExistent tests getting non-existent template
func TestGetTemplate_NonExistent(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	_, exists := service.GetTemplate("non_existent")
	if exists {
		t.Error("Expected non-existent template to return false")
	}
}

// TestGetTemplate_Existing tests getting existing template
func TestGetTemplate_Existing(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	template, exists := service.GetTemplate("password_reset")
	if !exists {
		t.Error("Expected password_reset template to exist")
	}

	if template.Subject != "Reset your Sereni LRS password" {
		t.Errorf("Expected 'Reset your Sereni LRS password', got %s", template.Subject)
	}
}

// TestSendTemplateEmail_Success tests successful template email send
func TestSendTemplateEmail_Success(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	sendCalled := false
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		sendCalled = true
		return nil
	}

	err := service.SendTemplateEmail(
		[]string{"test@example.com"},
		"welcome",
		map[string]interface{}{"name": "Bob"},
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !sendCalled {
		t.Error("Expected SendEmail to be called")
	}
}

// TestSendTemplateEmail_TemplateNotFound tests template not found error
func TestSendTemplateEmail_TemplateNotFound(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	err := service.SendTemplateEmail(
		[]string{"test@example.com"},
		"non_existent",
		map[string]interface{}{},
	)

	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' in error, got %s", err.Error())
	}
}

// TestSendTemplateEmail_SendError tests error during send
func TestSendTemplateEmail_SendError(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return &EmailSendError{Message: "SMTP connection failed"}
	}

	err := service.SendTemplateEmail(
		[]string{"test@example.com"},
		"welcome",
		map[string]interface{}{"name": "Charlie"},
	)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// TestGetAvailableTemplates_NotEmpty tests that available templates list is not empty
func TestGetAvailableTemplates_NotEmpty(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	templates := service.GetAvailableTemplates()

	if len(templates) == 0 {
		t.Error("Expected at least one template available")
	}
}

// TestGetAvailableTemplates_AllPredefined tests that all predefined templates are available
func TestGetAvailableTemplates_AllPredefined(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	templates := service.GetAvailableTemplates()
	expectedCount := 4 // welcome, password_reset, verification, otp_template

	if len(templates) < expectedCount {
		t.Errorf("Expected at least %d templates, got %d", expectedCount, len(templates))
	}
}

// TestGetAvailableTemplates_CustomAdded tests that custom templates are included
func TestGetAvailableTemplates_CustomAdded(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	initialCount := len(service.GetAvailableTemplates())

	customTemplate := email.EmailTemplate{
		Name:     "custom_alert",
		Subject:  "Alert",
		HTMLBody: "<html>Alert</html>",
		TextBody: "Alert",
	}

	service.AddTemplate("custom_alert", customTemplate)

	newCount := len(service.GetAvailableTemplates())

	if newCount <= initialCount {
		t.Errorf("Expected more templates after adding custom one. Before: %d, After: %d", initialCount, newCount)
	}
}

// EmailSendError is a test helper for template send errors
type EmailSendError struct {
	Message string
}

func (e *EmailSendError) Error() string {
	return e.Message
}
