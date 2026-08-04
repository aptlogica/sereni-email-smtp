# Secure Authentication Example

This example demonstrates **secure** email-based authentication patterns using the Sereni Email SMTP service.

## Security Features Demonstrated

This example shows how to **prevent** common email security vulnerabilities:

1. ✅ **Host Header Injection Prevention** - Uses configuration-based URLs instead of HTTP Host headers
2. ✅ **Trusted Domain Validation** - Only allows emails with URLs from configured trusted domains
3. ✅ **HTTPS Enforcement** - Requires HTTPS for all email links (configurable)
4. ✅ **Input Sanitization** - All email data is automatically sanitized
5. ✅ **Secure Token Generation** - Uses cryptographically secure random tokens

## What This Example Prevents

### Host Header Injection Attack (Prevented ✅)

**Attack Scenario:**
```
Attacker sends: POST /api/v1/auth/password-reset/request
With header: Host: malicious.com

Without protection:
  - Victim receives email from trusted service
  - Email contains: https://malicious.com/reset?token=SECRET
  - Victim clicks link, token leaked to attacker
  - Attacker resets victim's password

With this example:
  - URLs are constructed from BASE_URL config (not Host header)
  - If attacker's URL somehow gets through, it's rejected by trusted domain validation
  - Victim is protected ✅
```

## Setup

### 1. Environment Variables

Create a `.env` file:

```bash
# Application Configuration
BASE_URL=https://example.com
TRUSTED_DOMAINS=example.com,app.example.com
PORT=8080

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@example.com
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Run the Example

```bash
# Load environment variables
export $(cat .env | xargs)

# Run the server
go run main.go
```

## API Endpoints

### Request Password Reset

```bash
curl -X POST http://localhost:8080/api/v1/auth/password-reset/request \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

**What happens:**
1. Server generates secure token
2. Server constructs reset URL using `BASE_URL` from config (NOT from HTTP headers)
3. Email service validates URL against trusted domains
4. Email is sent with secure reset link

### Confirm Password Reset

```bash
curl -X POST http://localhost:8080/api/v1/auth/password-reset/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-reset-token",
    "new_password": "NewSecurePassword123!"
  }'
```

### Send Email Verification

```bash
curl -X POST http://localhost:8080/api/v1/auth/email-verification/send \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

### Verify Email

```bash
curl -X POST http://localhost:8080/api/v1/auth/email-verification/verify \
  -H "Content-Type: application/json" \
  -d '{"token": "your-verification-token"}'
```

## Security Best Practices Implemented

### ✅ 1. Configuration-Based URLs

```go
// ✅ SECURE
baseURL := os.Getenv("BASE_URL") // from config
resetURL := baseURL + "/reset?token=" + token

// ❌ VULNERABLE
host := c.Request.Host // from HTTP request
resetURL := "https://" + host + "/reset?token=" + token
```

### ✅ 2. Trusted Domain Validation

```go
// Configure trusted domains at initialization
emailService.SetTrustedDomains(
    []string{"example.com", "app.example.com"},
    false, // HTTPS only
)

// All URLs in emails are validated against these domains
```

### ✅ 3. Secure Token Generation

```go
// Uses crypto/rand for cryptographically secure tokens
func generateSecureToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}
```

### ✅ 4. HTTPS Enforcement

```go
// Second parameter = false means HTTP is not allowed
emailService.SetTrustedDomains(trustedDomains, false)
```

### ✅ 5. Automatic Input Sanitization

All email data is automatically sanitized:
- Email addresses validated using RFC 5322
- Headers sanitized to prevent injection
- Template data HTML-escaped to prevent XSS
- URLs validated for malicious content

## Testing the Security

### Test 1: Verify Trusted Domain Enforcement

Try to send an email with an untrusted domain:

```go
emailService.SendTemplateEmail(
    []string{"test@example.com"},
    "password_reset",
    map[string]interface{}{
        "reset_url": "https://malicious.com/phishing",
    },
)
// Result: Error - domain not in trusted list ✅
```

### Test 2: Verify HTTPS Enforcement

Try to use HTTP URL:

```go
emailService.SendTemplateEmail(
    []string{"test@example.com"},
    "password_reset",
    map[string]interface{}{
        "reset_url": "http://example.com/reset",
    },
)
// Result: Error - HTTP not allowed ✅
```

### Test 3: Verify XSS Protection

Try to inject malicious content:

```go
emailService.SendTemplateEmail(
    []string{"test@example.com"},
    "welcome",
    map[string]interface{}{
        "name": "<script>alert('XSS')</script>",
    },
)
// Result: Success - content is HTML-escaped ✅
```

## Common Mistakes to Avoid

### ❌ Mistake 1: Using HTTP Host Header

```go
// DON'T DO THIS
host := c.Request.Host
resetURL := "https://" + host + "/reset?token=" + token
```

### ❌ Mistake 2: Not Configuring Trusted Domains

```go
// DON'T DO THIS
emailService := email.NewEmailService(...)
// Missing: emailService.SetTrustedDomains(...)
```

### ❌ Mistake 3: Allowing HTTP in Production

```go
// DON'T DO THIS IN PRODUCTION
emailService.SetTrustedDomains(domains, true) // allows HTTP
```

### ❌ Mistake 4: Hardcoding URLs

```go
// DON'T DO THIS
resetURL := "https://example.com/reset?token=" + token
// Use environment variable instead
```

## Deployment Checklist

Before deploying to production:

- [ ] `BASE_URL` set to production domain
- [ ] `TRUSTED_DOMAINS` configured correctly
- [ ] SMTP credentials stored securely (e.g., AWS Secrets Manager)
- [ ] HTTPS-only mode enabled (`allowHTTP = false`)
- [ ] Rate limiting implemented on auth endpoints
- [ ] Token expiry implemented in database
- [ ] Tokens are single-use only
- [ ] Failed login attempts are logged
- [ ] Email sending errors are logged (but don't leak info to users)

## Additional Security Measures

Consider implementing these additional security measures:

1. **Rate Limiting**: Prevent password reset spam
2. **CAPTCHA**: Prevent automated abuse
3. **Account Lockout**: Lock accounts after multiple failed attempts
4. **Two-Factor Authentication**: Add second factor for sensitive operations
5. **Audit Logging**: Log all authentication events
6. **Token Expiry**: Expire tokens after 15-30 minutes
7. **Single-Use Tokens**: Invalidate tokens after use

## References

- [OWASP: Forgot Password](https://cheatsheetseries.owasp.org/cheatsheets/Forgot_Password_Cheat_Sheet.html)
- [OWASP: Email Security](https://owasp.org/www-community/vulnerabilities/Email_Header_Injection)
- [Main Security Guide](../../docs/SECURITY_GUIDE.md)

## License

Copyright 2026-2030 Aptlogica Technologies Pvt Ltd  
Licensed under the Apache License, Version 2.0
