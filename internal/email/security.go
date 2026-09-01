// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package email

import (
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

// TrustedDomainConfig holds the configuration for trusted domains
type TrustedDomainConfig struct {
	TrustedDomains []string
	AllowHTTPS     bool
	AllowHTTP      bool
}

// DefaultTrustedDomainConfig returns a secure default configuration
func DefaultTrustedDomainConfig() *TrustedDomainConfig {
	return &TrustedDomainConfig{
		TrustedDomains: []string{},
		AllowHTTPS:     true,
		AllowHTTP:      false, // Disallow HTTP by default for security
	}
}

var (
	// Comprehensive XSS pattern to detect potential script injection
	xssPattern = regexp.MustCompile(`(?i)<script|javascript:|onerror=|onload=|<iframe|<embed|<object`)
	
	// URL validation pattern
	urlPattern = regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=%]+$`)
)

// SanitizeURL validates and sanitizes a URL to prevent injection attacks
// Returns an error if the URL is malicious or invalid
func SanitizeURL(rawURL string, config *TrustedDomainConfig) (string, error) {
	if rawURL == "" {
		return "", errors.New("URL cannot be empty")
	}

	// Remove leading/trailing whitespace
	rawURL = strings.TrimSpace(rawURL)

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Validate scheme
	if parsedURL.Scheme == "" {
		return "", errors.New("URL must have a scheme (http:// or https://)")
	}

	if config == nil {
		config = DefaultTrustedDomainConfig()
	}

	if err := validateURLScheme(parsedURL.Scheme, config); err != nil {
		return "", err
	}

	// Check for javascript: or data: schemes in any part
	if containsMaliciousScheme(rawURL) {
		return "", errors.New("potentially malicious URL detected")
	}

	// Validate hostname is not empty
	if parsedURL.Host == "" {
		return "", errors.New("URL must have a valid hostname")
	}

	// If trusted domains are configured, validate against them
	if err := validateTrustedDomain(parsedURL, config); err != nil {
		return "", err
	}

	// Additional XSS checks
	if xssPattern.MatchString(rawURL) {
		return "", errors.New("URL contains potentially malicious content")
	}

	// Return the sanitized URL
	return parsedURL.String(), nil
}

// validateURLScheme checks that the scheme is http/https and allowed by config
func validateURLScheme(scheme string, config *TrustedDomainConfig) error {
	if scheme == "https" && !config.AllowHTTPS {
		return errors.New("HTTPS URLs are not allowed")
	}
	if scheme == "http" && !config.AllowHTTP {
		return errors.New("HTTP URLs are not allowed (use HTTPS)")
	}
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s (only http/https allowed)", scheme)
	}
	return nil
}

// containsMaliciousScheme reports whether the raw URL contains a javascript: or data: scheme
func containsMaliciousScheme(rawURL string) bool {
	lower := strings.ToLower(rawURL)
	return strings.Contains(lower, "javascript:") || strings.Contains(lower, "data:")
}

// validateTrustedDomain checks the parsed URL's hostname against the configured trusted domains
func validateTrustedDomain(parsedURL *url.URL, config *TrustedDomainConfig) error {
	if len(config.TrustedDomains) == 0 {
		return nil
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	for _, domain := range config.TrustedDomains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		// Allow exact match or subdomain match
		if hostname == domain || strings.HasSuffix(hostname, "."+domain) {
			return nil
		}
	}

	return fmt.Errorf("URL domain '%s' is not in the trusted domains list", parsedURL.Hostname())
}

// templateURLFieldKeys are the key substrings used to identify URL fields in template data
var templateURLFieldKeys = []string{"url", "link", "reset_url", "verification_url", "callback_url", "redirect_url"}

// SanitizeTemplateData validates and sanitizes all data in template data map
// This prevents XSS and injection attacks in email templates
func SanitizeTemplateData(data map[string]interface{}, config *TrustedDomainConfig) (map[string]interface{}, error) {
	if data == nil {
		return data, nil
	}

	sanitized := make(map[string]interface{})

	for key, value := range data {
		sanitizedValue, err := sanitizeTemplateValue(key, value, config)
		if err != nil {
			return nil, err
		}
		sanitized[key] = sanitizedValue
	}

	return sanitized, nil
}

// sanitizeTemplateValue sanitizes a single template data value according to its type,
// treating string values whose key looks like a URL field as URLs to validate
func sanitizeTemplateValue(key string, value interface{}, config *TrustedDomainConfig) (interface{}, error) {
	switch v := value.(type) {
	case string:
		if isTemplateURLField(key) {
			// Validate and sanitize URL
			sanitizedURL, err := SanitizeURL(v, config)
			if err != nil {
				return nil, fmt.Errorf("invalid URL in field '%s': %w", key, err)
			}
			return sanitizedURL, nil
		}
		// For non-URL strings, HTML escape to prevent XSS
		return html.EscapeString(v), nil
	case int, int64, float64, bool:
		// Numeric and boolean types are safe
		return v, nil
	case map[string]interface{}:
		// Recursively sanitize nested maps
		nestedSanitized, err := SanitizeTemplateData(v, config)
		if err != nil {
			return nil, fmt.Errorf("error sanitizing nested field '%s': %w", key, err)
		}
		return nestedSanitized, nil
	default:
		// For other types, convert to string and escape
		return html.EscapeString(fmt.Sprintf("%v", v)), nil
	}
}

// isTemplateURLField checks if this is a URL field based on its key
func isTemplateURLField(key string) bool {
	keyLower := strings.ToLower(key)
	for _, urlKey := range templateURLFieldKeys {
		if strings.Contains(keyLower, urlKey) {
			return true
		}
	}
	return false
}

// ValidateAndSanitizeHostHeader validates a host header to prevent host header injection
// Returns the validated host or an error if it's not in the trusted domains
func ValidateAndSanitizeHostHeader(hostHeader string, trustedDomains []string) (string, error) {
	if hostHeader == "" {
		return "", errors.New("host header cannot be empty")
	}

	// Remove port if present
	host := strings.Split(hostHeader, ":")[0]
	host = strings.ToLower(strings.TrimSpace(host))

	// Check against trusted domains
	if len(trustedDomains) == 0 {
		return "", errors.New("no trusted domains configured - cannot validate host header")
	}

	for _, domain := range trustedDomains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return host, nil
		}
	}

	return "", fmt.Errorf("host header '%s' is not in trusted domains", host)
}

// SanitizeHTMLContent sanitizes HTML content to prevent XSS attacks
// This is a basic sanitization - for production use, consider a dedicated library
func SanitizeHTMLContent(content string) string {
	// Remove potentially dangerous tags and attributes
	dangerous := []string{
		"<script", "</script>",
		"<iframe", "</iframe>",
		"<object", "</object>",
		"<embed", "</embed>",
		"javascript:",
		"onerror=",
		"onload=",
		"onclick=",
		"onmouseover=",
	}

	sanitized := content
	for _, danger := range dangerous {
		// Case-insensitive replacement
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(danger))
		sanitized = re.ReplaceAllString(sanitized, "")
	}

	return sanitized
}
