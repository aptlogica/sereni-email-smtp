// Copyright (c) 2026 Aptlogica Technologies Private Limited
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aptlogica/sereni-email-smtp/internal/email"
	"github.com/aptlogica/sereni-email-smtp/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewEmailHandler_Comprehensive(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	handler := handlers.NewEmailHandler(service)

	if handler.Service != service {
		t.Error("Expected handler to have the provided service")
	}
}

func TestEmailHandler_SendEmail_Comprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}
	handler := handlers.NewEmailHandler(service)

	// Test successful send
	req := email.EmailRequest{
		To:      []string{"test@example.com"},
		Subject: "Test",
		Body:    "Body",
		IsHTML:  true,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendEmail(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test invalid JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendEmail(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", w.Code)
	}

	// Test service error
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return errors.New("service error")
	}

	body, _ = json.Marshal(req)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendEmail(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for service error, got %d", w.Code)
	}
}

func TestEmailHandler_SendBulkEmail_Comprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	service.SendBulkEmailFunc = func(recipients []string, subject, body string, isHTML bool) ([]string, error) {
		return []string{}, nil
	}
	handler := handlers.NewEmailHandler(service)

	// Test successful bulk send
	req := email.BulkEmailRequest{
		Recipients: []string{"test1@example.com", "test2@example.com"},
		Subject:    "Bulk Test",
		Body:       "Bulk Body",
		IsHTML:     true,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendBulkEmail(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test invalid JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendBulkEmail(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", w.Code)
	}

	// Test with failed emails
	service.SendBulkEmailFunc = func(recipients []string, subject, body string, isHTML bool) ([]string, error) {
		return []string{"failed@example.com"}, nil
	}

	body, _ = json.Marshal(req)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendBulkEmail(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 even with failed emails, got %d", w.Code)
	}

	// Test service error
	service.SendBulkEmailFunc = func(recipients []string, subject, body string, isHTML bool) ([]string, error) {
		return nil, errors.New("bulk service error")
	}

	body, _ = json.Marshal(req)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendBulkEmail(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for service error, got %d", w.Code)
	}
}

func TestEmailHandler_GenerateOTP_Comprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test successful OTP generation
	{
		service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
		service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
			return nil
		}
		handler := handlers.NewEmailHandler(service)

		req := map[string]interface{}{
			"to":     "test@example.com",
			"expiry": 10,
		}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

		handler.GenerateOTP(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	}

	// Test with default expiry (missing expiry field)
	{
		service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
		service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
			return nil
		}
		handler := handlers.NewEmailHandler(service)

		req := map[string]interface{}{
			"to": "test@example.com",
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.GenerateOTP(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 with default expiry, got %d", w.Code)
		}
	}

	// Test invalid JSON
	{
		service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
		service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
			return nil
		}
		handler := handlers.NewEmailHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.GenerateOTP(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid JSON, got %d", w.Code)
		}
	}

	// Test service error still returns 200 (async send)
	{
		service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
		service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
			return errors.New("otp service error")
		}
		handler := handlers.NewEmailHandler(service)

		req := map[string]interface{}{
			"to":     "test@example.com",
			"expiry": 10,
		}
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.GenerateOTP(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 even if send fails asynchronously, got %d", w.Code)
		}
	}
}

func TestEmailHandler_VerifyOTP_Comprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		return nil
	}
	handler := handlers.NewEmailHandler(service)

	// Generate OTP first
	otp, _ := service.GenerateAndSendOTP("test@example.com", 10)

	// Test successful verification
	req := map[string]string{
		"email": "test@example.com",
		"otp":   otp,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.VerifyOTP(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test invalid JSON
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.VerifyOTP(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", w.Code)
	}

	// Test wrong OTP
	req = map[string]string{
		"email": "test@example.com",
		"otp":   "000000",
	}

	body, _ = json.Marshal(req)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.VerifyOTP(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong OTP, got %d", w.Code)
	}
}

func TestEmailHandler_HealthCheck_Comprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	handler := handlers.NewEmailHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler.HealthCheck(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse health check response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %v", response["status"])
	}
}
