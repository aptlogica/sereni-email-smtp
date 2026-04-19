// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"os"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/config"
)

func TestLoadConfig_Comprehensive(t *testing.T) {
	// Clear any existing env vars
	oldVars := map[string]string{
		"HOST":            os.Getenv("HOST"),
		"PORT":            os.Getenv("PORT"),
		"ALLOWED_ORIGIN":  os.Getenv("ALLOWED_ORIGIN"),
		"SMTP_HOST":       os.Getenv("SMTP_HOST"),
		"SMTP_PORT":       os.Getenv("SMTP_PORT"),
		"SMTP_USERNAME":   os.Getenv("SMTP_USERNAME"),
		"SMTP_PASSWORD":   os.Getenv("SMTP_PASSWORD"),
		"FROM_EMAIL":      os.Getenv("FROM_EMAIL"),
		"REDIS_URL":       os.Getenv("REDIS_URL"),
		"BULK_BATCH_SIZE": os.Getenv("BULK_BATCH_SIZE"),
	}

	// Clear all
	for k := range oldVars {
		os.Unsetenv(k)
	}

	// Restore at the end
	defer func() {
		for k, v := range oldVars {
			if v != "" {
				os.Setenv(k, v)
			}
		}
	}()

	// Set required secrets for test
	os.Setenv("SMTP_USERNAME", "testuser")
	os.Setenv("SMTP_PASSWORD", "testpass")
	os.Setenv("FROM_EMAIL", "test@example.com")

	// Test with defaults
	cfg := config.LoadConfig()
	if cfg.Host != "0.0.0.0" {
		t.Errorf("Expected Host 0.0.0.0, got %s", cfg.Host)
	}
	if cfg.Port != "8082" {
		t.Errorf("Expected Port 8082, got %s", cfg.Port)
	}
	if cfg.AllowedOrigin != "*" {
		t.Errorf("Expected AllowedOrigin *, got %s", cfg.AllowedOrigin)
	}
	if cfg.SMTPHost != "smtp.gmail.com" {
		t.Errorf("Expected SMTPHost smtp.gmail.com, got %s", cfg.SMTPHost)
	}
	if cfg.SMTPPort != 587 {
		t.Errorf("Expected SMTPPort 587, got %d", cfg.SMTPPort)
	}
	if cfg.BulkBatchSize != 10 {
		t.Errorf("Expected BulkBatchSize 10, got %d", cfg.BulkBatchSize)
	}

	// Test with custom values
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("PORT", "9999")
	os.Setenv("ALLOWED_ORIGIN", "https://example.com")
	os.Setenv("SMTP_HOST", "custom.smtp.com")
	os.Setenv("SMTP_PORT", "2525")
	os.Setenv("SMTP_USERNAME", "testuser")
	os.Setenv("SMTP_PASSWORD", "testpass")
	os.Setenv("FROM_EMAIL", "test@example.com")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("BULK_BATCH_SIZE", "25")

	cfg = config.LoadConfig()
	if cfg.Host != "127.0.0.1" {
		t.Errorf("Expected Host 127.0.0.1, got %s", cfg.Host)
	}
	if cfg.Port != "9999" {
		t.Errorf("Expected Port 9999, got %s", cfg.Port)
	}
	if cfg.AllowedOrigin != "https://example.com" {
		t.Errorf("Expected AllowedOrigin https://example.com, got %s", cfg.AllowedOrigin)
	}
	if cfg.SMTPHost != "custom.smtp.com" {
		t.Errorf("Expected SMTPHost custom.smtp.com, got %s", cfg.SMTPHost)
	}
	if cfg.SMTPPort != 2525 {
		t.Errorf("Expected SMTPPort 2525, got %d", cfg.SMTPPort)
	}
	if cfg.SMTPUsername != "testuser" {
		t.Errorf("Expected SMTPUsername testuser, got %s", cfg.SMTPUsername)
	}
	if cfg.SMTPPassword != "testpass" {
		t.Errorf("Expected SMTPPassword testpass, got %s", cfg.SMTPPassword)
	}
	if cfg.FromEmail != "test@example.com" {
		t.Errorf("Expected FromEmail test@example.com, got %s", cfg.FromEmail)
	}
	if cfg.RedisURL != "redis://localhost:6379" {
		t.Errorf("Expected RedisURL redis://localhost:6379, got %s", cfg.RedisURL)
	}
	if cfg.BulkBatchSize != 25 {
		t.Errorf("Expected BulkBatchSize 25, got %d", cfg.BulkBatchSize)
	}
}

func TestGetEnv_Comprehensive(t *testing.T) {
	// Test with existing env
	os.Setenv("TEST_GET_ENV", "testvalue")
	defer os.Unsetenv("TEST_GET_ENV")

	result := config.GetEnv("TEST_GET_ENV", "default")
	if result != "testvalue" {
		t.Errorf("Expected testvalue, got %s", result)
	}

	// Test with missing env
	result = config.GetEnv("MISSING_ENV_VAR", "defaultval")
	if result != "defaultval" {
		t.Errorf("Expected defaultval, got %s", result)
	}

	// Test with empty env
	os.Setenv("EMPTY_ENV", "")
	defer os.Unsetenv("EMPTY_ENV")
	result = config.GetEnv("EMPTY_ENV", "emptydefault")
	if result != "emptydefault" {
		t.Errorf("Expected emptydefault, got %s", result)
	}
}

func TestGetEnvAsInt_Comprehensive(t *testing.T) {
	// Test with valid int
	os.Setenv("TEST_INT_ENV", "42")
	defer os.Unsetenv("TEST_INT_ENV")

	result := config.GetEnvAsInt("TEST_INT_ENV", 100)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid int
	os.Setenv("TEST_INVALID_INT", "not-a-number")
	defer os.Unsetenv("TEST_INVALID_INT")
	result = config.GetEnvAsInt("TEST_INVALID_INT", 200)
	if result != 200 {
		t.Errorf("Expected 200, got %d", result)
	}

	// Test with missing env
	result = config.GetEnvAsInt("MISSING_INT_VAR", 300)
	if result != 300 {
		t.Errorf("Expected 300, got %d", result)
	}

	// Test with negative int
	os.Setenv("TEST_NEG_INT", "-123")
	defer os.Unsetenv("TEST_NEG_INT")
	result = config.GetEnvAsInt("TEST_NEG_INT", 400)
	if result != -123 {
		t.Errorf("Expected -123, got %d", result)
	}
}
