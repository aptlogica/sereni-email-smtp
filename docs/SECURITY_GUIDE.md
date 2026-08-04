# Email Security Guide

## Overview

This document provides security guidelines for using the Sereni Email SMTP service safely. Following these practices will help prevent common email-related vulnerabilities such as email header injection, cross-site scripting (XSS), and account compromise attacks.

## Table of Contents

1. [Common Email Security Vulnerabilities](#common-email-security-vulnerabilities)
2. [Security Features](#security-features)
3. [Secure Usage Patterns](#secure-usage-patterns)
4. [Configuration](#configuration)
5. [Examples](#examples)
6. [References](#references)

---

## Common Email Security Vulnerabilities

### 1. Email Header Injection

**Description**: Attackers inject newline characters (`\r\n`) into email headers to add additional headers or alter email routing.

**Impact**: 
- Spam relay abuse
- BCC injection to send copies to unintended recipients
- Email spoofing

**Prevention**: This service automatically sanitizes all headers using `sanitizeHeader()` and `stripCRLF()` functions.

### 2. Host Header Injection

**Description**: Using untrusted HTTP Host headers to construct URLs in emails.

**Example Attack**:
```go
// VULNERABLE CODE - DO NOT USE
func sendPasswordReset(w http.ResponseWriter, r *http.Request) {
    host := r.Header.Get("Host")  // ❌ UNTRUSTED INPUT
    token := generateResetToken(email)
    resetURL := "https://" + host + "/reset/" + token  // ❌ DANGEROUS
    
    emailService.SendTemplateEmail(
        []string{email},
        "password_reset",
        map[string]interface{}{
            "reset_url": resetURL,  // ❌ ATTACKER-CONTROLLED
        },
    )
}
```

**Attack Scenario**:
1. Attacker sends request with `Host: malicious.com` header
2. Victim receives email from trusted service
3. Email contains link: `https://malicious.com/reset/SECRET_TOKEN`
4. Victim clicks link, token is leaked to attacker
5. Attacker uses token to reset victim's password

**Prevention**: Use the trusted domain configuration described below.

### 3. Cross-Site Scripting (XSS) in Emails

**Description**: Malicious scripts embedded in email content.

**Impact**: 
- JavaScript execution in email clients that render HTML
- Cookie theft
- Phishing

**Prevention**: This service automatically escapes HTML in template data and sanitizes HTML content.

---

## Security Features

### Automatic Input Sanitization

The service provides multiple layers of defense:

1. **Header Injection Prevention**
   - `sanitizeHeader()` - Removes control characters
   - `stripCRLF()` - Removes carriage return/line feed characters
   - Email address validation using RFC 5322 parsing

2. **URL Validation**
   - `SanitizeURL()` - Validates URL format and scheme
   - Blocks `javascript:`, `data:`, and other dangerous schemes
   - XSS pattern detection in URLs
   - Trusted domain enforcement

3. **Template Data Sanitization**
   - `SanitizeTemplateData()` - Validates all template variables
   - Automatic HTML escaping for non-URL fields
   - URL validation for fields containing `url`, `link`, etc.
   - Recursive sanitization of nested data

4. **HTML Content Sanitization**
   - `SanitizeHTMLContent()` - Removes dangerous HTML tags and attributes
   - Blocks `<script>`, `<iframe>`, `<object>`, `<embed>`
   - Removes event handlers like `onerror`, `onload`, etc.

---

## Secure Usage Patterns

### 1. Configure Trusted Domains

**Always configure trusted domains before sending emails with URLs:**

```go
emailService := email.NewEmailService(
    smtpHost, smtpPort, 
    smtpUser, smtpPass, 
    fromEmail, batchSize,
)

// Configure your trusted domains
emailService.SetTrustedDomains(
    []string{"example.com", "app.example.com"},
    false, // allowHTTP - set to false to require HTTPS
)
```

### 2. Use Configuration Files for URLs (NOT Request Headers)

**✅ SECURE APPROACH:**

```go
// Load base URL from secure configuration
baseURL := config.Get("BASE_URL") // e.g., "https://example.com"

func sendPasswordReset(email string) error {
    token := generateResetToken(email)
    resetURL := baseURL + "/reset/" + token
    
    return emailService.SendTemplateEmail(
        []string{email},
        "password_reset",
        map[string]interface{}{
            "reset_url": resetURL, // Safe - from config
        },
    )
}
```

**❌ VULNERABLE APPROACH:**

```go
// NEVER construct URLs from HTTP headers
func sendPasswordReset(w http.ResponseWriter, r *http.Request) {
    host := r.Header.Get("Host") // ❌ UNTRUSTED
    token := generateResetToken(email)
    resetURL := "https://" + host + "/reset/" + token // ❌ DANGEROUS
    
    // This will be rejected by SanitizeURL if host is not in trusted domains
    emailService.SendTemplateEmail(...)
}
```

### 3. Validate and Sanitize All Template Data

The service automatically sanitizes template data, but you should still validate input at the application level:

```go
func sendWelcomeEmail(userName, email string) error {
    // Application-level validation
    if !isValidUserName(userName) {
        return errors.New("invalid username")
    }
    
    // Template data is automatically sanitized by the service
    return emailService.SendTemplateEmail(
        []string{email},
        "welcome",
        map[string]interface{}{
            "name": userName, // Will be HTML-escaped automatically
        },
    )
}
```

### 4. Use Environment Variables for Configuration

**Example `.env` file:**

```bash
# Application base URL - loaded from secure configuration
BASE_URL=https://example.com

# Trusted domains for email URLs (comma-separated)
TRUSTED_DOMAINS=example.com,app.example.com

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-secure-password
FROM_EMAIL=noreply@example.com
```

**Example configuration loading:**

```go
type AppConfig struct {
    BaseURL        string
    TrustedDomains []string
    // ... other config
}

func LoadConfig() *AppConfig {
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        log.Fatal("BASE_URL must be set")
    }
    
    trustedDomainsStr := os.Getenv("TRUSTED_DOMAINS")
    trustedDomains := strings.Split(trustedDomainsStr, ",")
    
    return &AppConfig{
        BaseURL:        baseURL,
        TrustedDomains: trustedDomains,
    }
}
```

---

## Configuration

### TrustedDomainConfig

Configure URL validation behavior:

```go
type TrustedDomainConfig struct {
    TrustedDomains []string // List of trusted domains
    AllowHTTPS     bool     // Allow HTTPS URLs (default: true)
    AllowHTTP      bool     // Allow HTTP URLs (default: false)
}
```

**Example:**

```go
// Strict HTTPS-only configuration
emailService.SetTrustedDomains(
    []string{"example.com", "secure.example.com"},
    false, // disallow HTTP
)

// Development configuration (less secure, for testing only)
emailService.SetTrustedDomains(
    []string{"localhost", "127.0.0.1"},
    true, // allow HTTP for local development
)
```

### Subdomain Matching

Trusted domains automatically include all subdomains:

```go
// This configuration:
emailService.SetTrustedDomains([]string{"example.com"}, false)

// Will accept:
// ✅ https://example.com/reset
// ✅ https://app.example.com/reset
// ✅ https://auth.api.example.com/reset

// Will reject:
// ❌ https://exampleXcom (different domain)
// ❌ https://malicious.com
// ❌ http://example.com (HTTP not allowed)
```

---

## Examples

### Example 1: Secure Password Reset Flow

```go
package main

import (
    "github.com/aptlogica/sereni-email-smtp/internal/email"
    "github.com/gin-gonic/gin"
)

var (
    emailService *email.EmailService
    baseURL      string
)

func init() {
    // Load from secure configuration
    baseURL = os.Getenv("BASE_URL") // https://example.com
    
    emailService = email.NewEmailService(
        smtpHost, smtpPort,
        smtpUser, smtpPass,
        fromEmail, batchSize,
    )
    
    // Configure trusted domains
    emailService.SetTrustedDomains(
        []string{"example.com"},
        false, // HTTPS only
    )
}

func handlePasswordResetRequest(c *gin.Context) {
    var req struct {
        Email string `json:"email" binding:"required,email"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Generate secure token
    token := generateSecureResetToken(req.Email)
    
    // Construct URL from TRUSTED configuration (NOT from request)
    resetURL := baseURL + "/reset?token=" + token
    
    // Send email - URL will be validated against trusted domains
    err := emailService.SendTemplateEmail(
        []string{req.Email},
        "password_reset",
        map[string]interface{}{
            "reset_url": resetURL, // ✅ Safe - from config
        },
    )
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to send email"})
        return
    }
    
    c.JSON(200, gin.H{"message": "Password reset email sent"})
}
```

### Example 2: Secure Email Verification

```go
func sendEmailVerification(userEmail, userName string) error {
    // Generate verification token
    token := generateVerificationToken(userEmail)
    
    // Use configured base URL (NOT from HTTP request)
    verificationURL := baseURL + "/verify?token=" + token
    
    // Send email with automatic sanitization
    return emailService.SendTemplateEmail(
        []string{userEmail},
        "verification",
        map[string]interface{}{
            "name":             userName,           // Auto HTML-escaped
            "verification_url": verificationURL,    // Auto URL-validated
        },
    )
}
```

### Example 3: Custom Template with Security

```go
func sendCustomNotification(recipient, message string) error {
    // Create custom template
    customTemplate := email.EmailTemplate{
        Name:     "notification",
        Subject:  "Important Notification",
        HTMLBody: `
            <html>
            <body>
                <h1>Notification</h1>
                <p>{{.message}}</p>
                <p><a href="{{.action_url}}">Take Action</a></p>
            </body>
            </html>
        `,
        TextBody: "Notification: {{.message}}",
    }
    
    emailService.AddTemplate("notification", customTemplate)
    
    // Send with sanitized data
    return emailService.SendTemplateEmail(
        []string{recipient},
        "notification",
        map[string]interface{}{
            "message":    message,                          // HTML-escaped
            "action_url": baseURL + "/dashboard",           // URL-validated
        },
    )
}
```

### Example 4: Validating Host Headers (If Required)

If you must use host headers, validate them first:

```go
import "github.com/aptlogica/sereni-email-smtp/internal/email"

func buildURLFromRequest(r *http.Request, path string) (string, error) {
    hostHeader := r.Header.Get("Host")
    
    // Validate host against trusted domains
    validHost, err := email.ValidateAndSanitizeHostHeader(
        hostHeader,
        []string{"example.com", "app.example.com"},
    )
    if err != nil {
        return "", fmt.Errorf("untrusted host: %w", err)
    }
    
    // Use validated host
    return "https://" + validHost + path, nil
}
```

---

## Testing Security

Run the comprehensive security tests:

```bash
go test ./tests -run TestSanitizeURL -v
go test ./tests -run TestSanitizeTemplateData -v
go test ./tests -run TestValidateAndSanitizeHostHeader -v
go test ./tests -run TestEmailServiceWithSecureTemplates -v
```

---

## Security Checklist

Before deploying to production:

- [ ] Configured trusted domains using `SetTrustedDomains()`
- [ ] Base URL loaded from secure configuration (NOT from HTTP headers)
- [ ] HTTPS-only mode enabled (unless development environment)
- [ ] All email templates reviewed for XSS vulnerabilities
- [ ] Template data properly validated at application level
- [ ] Security tests passing
- [ ] Environment variables properly secured
- [ ] SMTP credentials stored securely (e.g., secrets manager)
- [ ] Rate limiting implemented for email sending endpoints
- [ ] Logging configured for security events

---

## References

- [OWASP: Email Header Injection](https://owasp.org/www-community/vulnerabilities/Email_Header_Injection)
- [OWASP: Content Spoofing](https://owasp.org/www-community/attacks/Content_Spoofing)
- [CWE-640: Weak Password Recovery Mechanism](https://cwe.mitre.org/data/definitions/640.html)
- [CWE-79: Cross-site Scripting (XSS)](https://cwe.mitre.org/data/definitions/79.html)
- [CWE-601: URL Redirection to Untrusted Site](https://cwe.mitre.org/data/definitions/601.html)
- [RFC 5322: Internet Message Format](https://tools.ietf.org/html/rfc5322)

---

## Reporting Security Issues

If you discover a security vulnerability, please report it to:
- Email: security@aptlogica.com
- See [SECURITY.md](../SECURITY.md) for details

**Do not disclose security vulnerabilities publicly until they have been addressed.**

---

## Version History

- **v1.0.0** (2026-08-04): Initial security implementation
  - Added URL validation and sanitization
  - Added template data sanitization
  - Added trusted domain configuration
  - Added host header validation
  - Added comprehensive security tests
