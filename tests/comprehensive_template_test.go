// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"errors"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

func TestRenderTemplate_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	// Test existing template (otp_template)
	data := map[string]interface{}{
		"otp":    "123456",
		"expiry": 10,
	}

	subject, body, err := service.RenderTemplate("otp_template", data)
	if err != nil {
		t.Errorf("Expected no error for existing template, got %v", err)
	}
	if subject == "" || body == "" {
		t.Error("Expected non-empty subject and body")
	}

	// Test non-existing template
	_, _, err = service.RenderTemplate("non_existing", data)
	if err == nil {
		t.Error("Expected error for non-existing template")
	}
}

func TestAddTemplate_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	template := email.EmailTemplate{
		Subject:  "Test Subject {{.name}}",
		HTMLBody: "<h1>Hello {{.name}}</h1>",
		TextBody: "Hello {{.name}}",
	}

	service.AddTemplate("test_template", template)

	// Verify template was added
	retrieved, ok := service.GetTemplate("test_template")
	if !ok {
		t.Error("Expected template to be found after adding")
	}
	if retrieved.Subject != template.Subject {
		t.Errorf("Expected subject %s, got %s", template.Subject, retrieved.Subject)
	}
	if retrieved.HTMLBody != template.HTMLBody {
		t.Errorf("Expected HTML body %s, got %s", template.HTMLBody, retrieved.HTMLBody)
	}
	if retrieved.TextBody != template.TextBody {
		t.Errorf("Expected text body %s, got %s", template.TextBody, retrieved.TextBody)
	}
}

func TestGetTemplate_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	// Test getting non-existing template
	_, ok := service.GetTemplate("non_existing")
	if ok {
		t.Error("Expected false for non-existing template")
	}

	// Add and get template
	template := email.EmailTemplate{
		Subject:  "Test",
		HTMLBody: "<p>Test</p>",
		TextBody: "Test",
	}
	service.AddTemplate("get_test", template)

	retrieved, ok := service.GetTemplate("get_test")
	if !ok {
		t.Error("Expected true for existing template")
	}
	if retrieved.Subject != "Test" {
		t.Errorf("Expected subject Test, got %s", retrieved.Subject)
	}
}

func TestSendTemplateEmail_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	var lastTo []string
	var lastSubject, lastBody string
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		lastTo = to
		lastSubject = subject
		lastBody = body
		return nil
	}

	// Test with existing template
	data := map[string]interface{}{
		"otp":    "123456",
		"expiry": 10,
	}

	err := service.SendTemplateEmail([]string{"test@example.com"}, "otp_template", data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(lastTo) != 1 || lastTo[0] != "test@example.com" {
		t.Errorf("Expected to be [test@example.com], got %v", lastTo)
	}
	if lastSubject == "" {
		t.Error("Expected non-empty subject")
	}
	if lastBody == "" {
		t.Error("Expected non-empty body")
	}

	// Test with non-existing template
	err = service.SendTemplateEmail([]string{"test@example.com"}, "non_existing", data)
	if err == nil {
		t.Error("Expected error for non-existing template")
	}

	// Test with send error
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return errors.New("send failed")
	}

	err = service.SendTemplateEmail([]string{"test@example.com"}, "otp_template", data)
	if err == nil {
		t.Error("Expected error when send fails")
	}
}

func TestGetAvailableTemplates_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	// Initially should have default templates
	templates := service.GetAvailableTemplates()
	if len(templates) == 0 {
		t.Error("Expected some default templates")
	}

	// Add custom template and check it's included
	customTemplate := email.EmailTemplate{
		Subject:  "Custom",
		HTMLBody: "<p>Custom</p>",
		TextBody: "Custom",
	}
	service.AddTemplate("custom", customTemplate)

	templates = service.GetAvailableTemplates()
	found := false
	for _, tmpl := range templates {
		if tmpl == "custom" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected custom template to be in available templates list")
	}
}
