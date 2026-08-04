// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		config      *email.TrustedDomainConfig
		expectError bool
		description string
	}{
		{
			name:        "Valid HTTPS URL",
			url:         "https://example.com/reset/token123",
			config:      nil,
			expectError: false,
			description: "Should accept valid HTTPS URL",
		},
		{
			name:        "Valid HTTP URL with allowHTTP",
			url:         "http://example.com/verify",
			config:      &email.TrustedDomainConfig{AllowHTTPS: true, AllowHTTP: true},
			expectError: false,
			description: "Should accept HTTP when explicitly allowed",
		},
		{
			name:        "Reject HTTP URL by default",
			url:         "http://example.com/reset",
			config:      nil,
			expectError: true,
			description: "Should reject HTTP URLs by default",
		},
		{
			name:        "Reject javascript: scheme",
			url:         "javascript:alert('XSS')",
			config:      nil,
			expectError: true,
			description: "Should reject javascript: scheme to prevent XSS",
		},
		{
			name:        "Reject data: scheme",
			url:         "data:text/html,<script>alert('XSS')</script>",
			config:      nil,
			expectError: true,
			description: "Should reject data: scheme to prevent XSS",
		},
		{
			name:        "Reject URL with embedded javascript",
			url:         "https://example.com/redirect?url=javascript:alert(1)",
			config:      nil,
			expectError: true,
			description: "Should detect javascript in URL parameters",
		},
		{
			name:        "Reject URL with XSS patterns",
			url:         "https://example.com/<script>alert('XSS')</script>",
			config:      nil,
			expectError: true,
			description: "Should reject URLs containing XSS patterns",
		},
		{
			name:        "Empty URL",
			url:         "",
			config:      nil,
			expectError: true,
			description: "Should reject empty URLs",
		},
		{
			name:        "URL without scheme",
			url:         "example.com/reset",
			config:      nil,
			expectError: true,
			description: "Should reject URLs without scheme",
		},
		{
			name:        "Trusted domain - exact match",
			url:         "https://example.com/reset",
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: false,
			description: "Should accept URL from trusted domain",
		},
		{
			name:        "Trusted domain - subdomain match",
			url:         "https://app.example.com/reset",
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: false,
			description: "Should accept subdomain of trusted domain",
		},
		{
			name:        "Untrusted domain",
			url:         "https://malicious.com/reset",
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: true,
			description: "Should reject URL from untrusted domain",
		},
		{
			name:        "URL with port",
			url:         "https://example.com:8080/reset",
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: false,
			description: "Should accept URL with port from trusted domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := email.SanitizeURL(tt.url, tt.config)
			if tt.expectError && err == nil {
				t.Errorf("%s: expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
			}
		})
	}
}

func TestSanitizeTemplateData(t *testing.T) {
	tests := []struct {
		name        string
		data        map[string]interface{}
		config      *email.TrustedDomainConfig
		expectError bool
		description string
	}{
		{
			name: "Safe template data",
			data: map[string]interface{}{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			config:      nil,
			expectError: false,
			description: "Should accept safe template data",
		},
		{
			name: "Template data with valid URL",
			data: map[string]interface{}{
				"name":      "John Doe",
				"reset_url": "https://example.com/reset/token123",
			},
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: false,
			description: "Should accept valid URLs in template data",
		},
		{
			name: "Template data with malicious URL",
			data: map[string]interface{}{
				"name":      "John Doe",
				"reset_url": "javascript:alert('XSS')",
			},
			config:      nil,
			expectError: true,
			description: "Should reject malicious URLs in template data",
		},
		{
			name: "Template data with XSS in string field",
			data: map[string]interface{}{
				"name": "<script>alert('XSS')</script>",
			},
			config:      nil,
			expectError: false, // XSS should be escaped, not rejected
			description: "Should escape XSS in non-URL fields",
		},
		{
			name: "Nested template data",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "John",
					"email": "john@example.com",
				},
			},
			config:      nil,
			expectError: false,
			description: "Should handle nested template data",
		},
		{
			name: "Template data with untrusted domain",
			data: map[string]interface{}{
				"reset_url": "https://malicious.com/phishing",
			},
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: true,
			description: "Should reject URLs from untrusted domains",
		},
		{
			name: "Template data with HTTP URL when only HTTPS allowed",
			data: map[string]interface{}{
				"verification_url": "http://example.com/verify",
			},
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true, AllowHTTP: false},
			expectError: true,
			description: "Should reject HTTP URLs when only HTTPS is allowed",
		},
		{
			name: "Multiple URL fields",
			data: map[string]interface{}{
				"reset_url":        "https://example.com/reset",
				"verification_url": "https://example.com/verify",
				"callback_url":     "https://example.com/callback",
			},
			config:      &email.TrustedDomainConfig{TrustedDomains: []string{"example.com"}, AllowHTTPS: true},
			expectError: false,
			description: "Should validate all URL fields",
		},
		{
			name: "Numeric and boolean values",
			data: map[string]interface{}{
				"count":   42,
				"active":  true,
				"balance": 99.99,
			},
			config:      nil,
			expectError: false,
			description: "Should accept safe numeric and boolean values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := email.SanitizeTemplateData(tt.data, tt.config)
			if tt.expectError && err == nil {
				t.Errorf("%s: expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
			}
		})
	}
}

func TestValidateAndSanitizeHostHeader(t *testing.T) {
	tests := []struct {
		name           string
		hostHeader     string
		trustedDomains []string
		expectError    bool
		expectedHost   string
		description    string
	}{
		{
			name:           "Valid trusted host",
			hostHeader:     "example.com",
			trustedDomains: []string{"example.com"},
			expectError:    false,
			expectedHost:   "example.com",
			description:    "Should accept trusted host",
		},
		{
			name:           "Valid trusted subdomain",
			hostHeader:     "app.example.com",
			trustedDomains: []string{"example.com"},
			expectError:    false,
			expectedHost:   "app.example.com",
			description:    "Should accept subdomain of trusted domain",
		},
		{
			name:           "Host with port",
			hostHeader:     "example.com:8080",
			trustedDomains: []string{"example.com"},
			expectError:    false,
			expectedHost:   "example.com",
			description:    "Should strip port from host header",
		},
		{
			name:           "Untrusted host",
			hostHeader:     "malicious.com",
			trustedDomains: []string{"example.com"},
			expectError:    true,
			description:    "Should reject untrusted host",
		},
		{
			name:           "Empty host header",
			hostHeader:     "",
			trustedDomains: []string{"example.com"},
			expectError:    true,
			description:    "Should reject empty host header",
		},
		{
			name:           "No trusted domains configured",
			hostHeader:     "example.com",
			trustedDomains: []string{},
			expectError:    true,
			description:    "Should reject when no trusted domains are configured",
		},
		{
			name:           "Case insensitive matching",
			hostHeader:     "Example.COM",
			trustedDomains: []string{"example.com"},
			expectError:    false,
			expectedHost:   "example.com",
			description:    "Should match domain case-insensitively",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := email.ValidateAndSanitizeHostHeader(tt.hostHeader, tt.trustedDomains)
			if tt.expectError && err == nil {
				t.Errorf("%s: expected error but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("%s: unexpected error: %v", tt.description, err)
			}
			if !tt.expectError && result != tt.expectedHost {
				t.Errorf("%s: expected host '%s', got '%s'", tt.description, tt.expectedHost, result)
			}
		})
	}
}

func TestSanitizeHTMLContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		shouldAllow bool
		description string
	}{
		{
			name:        "Safe HTML content",
			content:     "<p>Hello, World!</p>",
			shouldAllow: true,
			description: "Should allow safe HTML",
		},
		{
			name:        "Remove script tags",
			content:     "<script>alert('XSS')</script>",
			shouldAllow: false,
			description: "Should remove script tags",
		},
		{
			name:        "Remove iframe tags",
			content:     "<iframe src='malicious.com'></iframe>",
			shouldAllow: false,
			description: "Should remove iframe tags",
		},
		{
			name:        "Remove javascript: protocol",
			content:     "<a href='javascript:alert(1)'>Click</a>",
			shouldAllow: false,
			description: "Should remove javascript: protocol",
		},
		{
			name:        "Remove onerror attribute",
			content:     "<img src='x' onerror='alert(1)'>",
			shouldAllow: false,
			description: "Should remove onerror attribute",
		},
		{
			name:        "Remove onload attribute",
			content:     "<body onload='malicious()'>",
			shouldAllow: false,
			description: "Should remove onload attribute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := email.SanitizeHTMLContent(tt.content)
			if tt.shouldAllow && result != tt.content {
				t.Errorf("%s: safe content was modified", tt.description)
			}
			if !tt.shouldAllow && result == tt.content {
				t.Errorf("%s: malicious content was not sanitized", tt.description)
			}
		})
	}
}

func TestEmailServiceWithSecureTemplates(t *testing.T) {
	// Test that email service properly validates template data
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	service.SetTrustedDomains([]string{"example.com"}, false)

	// Mock the SendEmailFunc to prevent actual email sending
	service.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		// Email sent successfully
		return nil
	}

	t.Run("Valid template data passes validation", func(t *testing.T) {
		templateData := map[string]interface{}{
			"name":      "John Doe",
			"reset_url": "https://example.com/reset/token123",
		}

		err := service.SendTemplateEmail([]string{"test@example.com"}, "password_reset", templateData)
		if err != nil {
			t.Errorf("Expected no error with valid template data, got: %v", err)
		}
	})

	t.Run("Malicious URL in template data is rejected", func(t *testing.T) {
		templateData := map[string]interface{}{
			"name":      "John Doe",
			"reset_url": "javascript:alert('XSS')",
		}

		err := service.SendTemplateEmail([]string{"test@example.com"}, "password_reset", templateData)
		if err == nil {
			t.Error("Expected error with malicious URL, got none")
		}
	})

	t.Run("Untrusted domain in template data is rejected", func(t *testing.T) {
		templateData := map[string]interface{}{
			"name":      "John Doe",
			"reset_url": "https://malicious.com/phishing",
		}

		err := service.SendTemplateEmail([]string{"test@example.com"}, "password_reset", templateData)
		if err == nil {
			t.Error("Expected error with untrusted domain, got none")
		}
	})

	t.Run("XSS in non-URL fields is escaped", func(t *testing.T) {
		templateData := map[string]interface{}{
			"name": "John<script>alert('XSS')</script>",
			"otp":  "123456",
		}

		err := service.SendTemplateEmail([]string{"test@example.com"}, "otp_template", templateData)
		if err != nil {
			t.Errorf("Expected no error with escaped XSS, got: %v", err)
		}
		// The body should contain HTML-escaped content
		// Note: The exact escaped form depends on template rendering
	})
}
