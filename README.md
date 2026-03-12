# Sereni Email SMTP - Production-Ready Email Microservice

> A scalable, production-ready REST API microservice for transactional and bulk email delivery with OTP generation, template support, and SMTP integration. Deploy as a standalone service or integrate into any application ecosystem.

[![Version](https://img.shields.io/badge/Version-1.0.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24.4+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)
[![Quality Gate Status](https://sonar.aptlogica.com/api/project_badges/measure?project=aptlogica_sereni-email-smtp_81bfe731-8d64-4b02-a58c-81c48abd4f9a&metric=alert_status&token=sqb_b3f0b470331e1ff36f199a51f47f818ea7fa056d)](https://sonar.aptlogica.com/dashboard?id=aptlogica_sereni-email-smtp_81bfe731-8d64-4b02-a58c-81c48abd4f9a)

## 📋 Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Use Cases](#use-cases)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Integration Guide](#integration-guide)
- [Architecture](#architecture)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [FAQ](#faq)
- [License](#license)

## Overview

**Sereni Email SMTP** is a high-performance email microservice that provides RESTful APIs for sending transactional emails, bulk campaigns, and OTP-based authentication. Built with Go and designed for scalability, it seamlessly integrates with any SMTP provider (Gmail, SendGrid, Mailgun, AWS SES, etc.) and works with applications in any programming language.

Sereni Email SMTP is designed with these key characteristics:

- **Microservice Architecture**: Deploy independently or as part of a larger system with Docker/Kubernetes support

- **Language-Agnostic**: Works with applications in any language (Node.js, Python, Java, .NET, PHP, Ruby, etc.) through RESTful APIs

- **Concurrent Processing**: Handle multiple email operations simultaneously with automatic batch optimization for bulk sending

- **Production-Ready**: Docker support, Swagger documentation, health checks, comprehensive logging, error handling, and automated testing

### Why Choose Sereni Email SMTP?

- **Zero Configuration Complexity**: Environment-based configuration - no complex setup files or XML configs

- **Battle-Tested Foundation**: Built on Go's robust net/smtp package with TLS encryption and authentication

- **High Performance**: Concurrent goroutines handle thousands of emails efficiently with automatic batch processing

- **Easy Integration**: Simple REST API that works with any language or framework - just make HTTP requests

- **Email Templates**: Pre-built templates for common scenarios (welcome, password reset, OTP, verification) with custom template support

- **OTP Built-In**: Generate and verify one-time passwords for authentication without additional dependencies

## Key Features

✅ **Transactional Email Delivery**
- Send individual emails with HTML or plain text content
- Support for multiple recipients (To, CC, BCC)
- Custom templates with dynamic data injection
- Automatic content-type detection and formatting

✅ **Bulk Email Campaigns**
- Send to thousands of recipients efficiently
- Automatic batch processing with configurable batch sizes
- Concurrent delivery for maximum throughput
- Failed recipient tracking and reporting

✅ **OTP (One-Time Password) System**
- Generate secure 6-digit verification codes
- Configurable expiration times (default: 5 minutes)
- Automatic cleanup of expired OTPs
- Email + OTP validation for authentication

✅ **Email Template Engine**
- Pre-built templates: Welcome, Password Reset, Email Verification, OTP
- Custom template support with variable substitution
- Both HTML and plain text versions
- Template caching for performance

✅ **RESTful API with Swagger**
- Full OpenAPI 3.0 documentation
- Interactive Swagger UI for testing
- Standard HTTP methods and status codes
- JSON request/response format

✅ **SMTP Provider Integration**
- Works with any SMTP provider (Gmail, Outlook, SendGrid, Mailgun, AWS SES, etc.)
- TLS/SSL encryption support
- SMTP authentication with username/password
- Configurable host, port, and from address

✅ **Production Features**
- Health check endpoint for monitoring
- CORS support for cross-origin requests
- Comprehensive error handling and logging
- Docker and Docker Compose ready
- Kubernetes deployment support

✅ **Developer Experience**
- Complete test coverage with Go testing framework
- Makefile for common operations
- Environment variable configuration
- Clear error messages and status codes

## Use Cases

### 1. **User Authentication & Verification**
Perfect for implementing email-based user registration, login, and two-factor authentication. Generate OTPs, verify emails, and send password reset links without building your own email infrastructure.

```
User Registration → API Call → Email Service → Verification Email → User Validates
```

### 2. **E-commerce Transactional Emails**
Send order confirmations, shipping notifications, payment receipts, and account updates in real-time. Integrate with your e-commerce platform regardless of the technology stack.

### 3. **Marketing Campaigns & Newsletters**
Deliver bulk promotional emails, newsletters, and announcements to thousands of subscribers efficiently. Track failed deliveries and retry logic built-in.

### 4. **Application Notifications**
Send automated alerts, reminders, system notifications, and status updates to users. Works as a notification backend for web, mobile, or desktop applications.

### 5. **Multi-Tenant SaaS Platforms**
Provide email capabilities as a shared service across multiple applications or tenants. Centralize email configuration and monitoring in one microservice.

### 6. **Legacy System Modernization**
Replace outdated email systems with a modern REST API without rewriting your entire application. Simply make HTTP calls from your existing codebase.

## Quick Start

### Prerequisites
- **Docker** (v20.0+) - For containerized deployment
- **SMTP Credentials** - From Gmail, SendGrid, Mailgun, AWS SES, or any SMTP provider
- **curl or Postman** (optional) - For API testing

### 30-Second Setup

```bash
# 1. Clone the repository
git clone https://github.com/aptlogica/sereni-email-smtp.git
cd sereni-email-smtp

# 2. Create environment file
cp .env.example .env

# 3. Edit .env with your SMTP credentials
# nano .env (or use your preferred editor)

# 4. Start the service with Docker Compose
docker-compose up -d

# 5. Verify installation
curl http://localhost:8082/health
```

**Service is now available at http://localhost:8082**

Visit **http://localhost:8082/swagger/index.html** for interactive API documentation.

**Next steps:** See [Installation](#installation) for more setup options, or [Usage](#usage) to send your first email.

## Installation

### Option 1: Docker Compose (Recommended)

Easiest way to get started. Great for development and production deployment.

```bash
# Step 1: Clone the repository
git clone https://github.com/aptlogica/sereni-email-smtp.git
cd sereni-email-smtp

# Step 2: Create environment configuration
cp .env.example .env

# Step 3: Edit .env with your SMTP settings
nano .env
# Update: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, FROM_EMAIL

# Step 4: Start the service
docker-compose up -d

# Step 5: Verify it's running
curl http://localhost:8082/health
```

**Result:** Service running at http://localhost:8082 with Swagger docs at /swagger/index.html

### Option 2: Docker

For custom container orchestration or production deployment without compose.

```bash
# Step 1: Build the Docker image
docker build -t sereni-email-smtp:latest .

# Step 2: Run the container with environment variables
docker run -d \
  -p 8082:8082 \
  -e HOST=0.0.0.0 \
  -e PORT=8082 \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=your_email@gmail.com \
  -e SMTP_PASSWORD=your_app_password \
  -e FROM_EMAIL=your_email@gmail.com \
  -e ALLOWED_ORIGIN=* \
  --name email-service \
  sereni-email-smtp:latest

# Step 3: Check logs
docker logs -f email-service

# Step 4: Verify service
curl http://localhost:8082/health
```

**Result:** Containerized service accessible on port 8082

### Option 3: From Source (Developers)

For development, testing, or when you want to modify the code.

```bash
# Step 1: Ensure Go 1.24.4+ is installed
go version

# Step 2: Clone repository
git clone https://github.com/aptlogica/sereni-email-smtp.git
cd sereni-email-smtp

# Step 3: Install dependencies
go mod download

# Step 4: Create and configure .env
cp .env.example .env
nano .env

# Step 5: Generate Swagger documentation
make swagger
# or: swag init -g cmd/server/main.go

# Step 6: Build the application
make build
# or: go build -o bin/email-service cmd/server/main.go

# Step 7: Run the service
make run
# or: go run cmd/server/main.go
```

**Result:** Service compiling from source and running on configured port

## Configuration

### Environment Variables

Create `.env` file in your project root:

```dotenv
# === Server Configuration ===
HOST=0.0.0.0                          # Server bind address (0.0.0.0 for all interfaces)
PORT=8082                             # HTTP server port
ALLOWED_ORIGIN=*                      # CORS allowed origins (* for all, or specific domain)

# === SMTP Configuration ===
SMTP_HOST=smtp.gmail.com              # SMTP server hostname
SMTP_PORT=587                         # SMTP server port (587 for TLS, 465 for SSL)
SMTP_USERNAME=your_email@gmail.com    # SMTP authentication username
SMTP_PASSWORD=your_app_password       # SMTP authentication password
FROM_EMAIL=your_email@gmail.com       # Default sender email address

# === Email Processing ===
BULK_BATCH_SIZE=50                    # Number of emails per batch for bulk sending (optional, default: 50)
```

### Default Values

If `.env` file is not provided, these defaults are used:
- `HOST`: `0.0.0.0`
- `PORT`: `8082`
- `ALLOWED_ORIGIN`: `*`
- `BULK_BATCH_SIZE`: `50`

### Configuration Examples

**For Gmail:**
```dotenv
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_16_char_app_password
FROM_EMAIL=your_email@gmail.com
```

> **Note:** Gmail requires an [App Password](https://support.google.com/accounts/answer/185833) if 2FA is enabled.

**For SendGrid:**
```dotenv
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your_sendgrid_api_key
FROM_EMAIL=verified@yourdomain.com
```

**For AWS SES:**
```dotenv
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your_aws_ses_username
SMTP_PASSWORD=your_aws_ses_password
FROM_EMAIL=verified@yourdomain.com
```

**For Production (Kubernetes ConfigMap/Secret):**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: email-service-config
data:
  HOST: "0.0.0.0"
  PORT: "8082"
  SMTP_HOST: "smtp.gmail.com"
  SMTP_PORT: "587"
---
apiVersion: v1
kind: Secret
metadata:
  name: email-service-secret
type: Opaque
data:
  SMTP_USERNAME: <base64-encoded>
  SMTP_PASSWORD: <base64-encoded>
  FROM_EMAIL: <base64-encoded>
```

## Usage

### Basic Usage

Send a simple email:

```bash
curl -X POST http://localhost:8082/api/v1/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "subject": "Hello from Sereni Email",
    "body": "This is a test email.",
    "is_html": false
  }'
```

### Example 1: Send HTML Email with Template

```bash
curl -X POST http://localhost:8082/api/v1/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "Welcome to Our Platform",
    "template": "welcome",
    "template_data": {
      "name": "John Doe"
    },
    "is_html": true
  }'
```

**Output:**
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

### Example 2: Generate and Send OTP

```bash
# Generate OTP
curl -X POST http://localhost:8082/api/v1/email/otp/generate \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "expiry": 300
  }'
```

**Output:**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "otp": "123456",
  "expiry": "2026-03-11T10:35:00Z"
}
```

### Example 3: Send Bulk Emails

```bash
curl -X POST http://localhost:8082/api/v1/email/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      "user1@example.com",
      "user2@example.com",
      "user3@example.com"
    ],
    "subject": "Monthly Newsletter",
    "body": "<h1>Hello Subscribers!</h1><p>Check out our latest updates...</p>",
    "is_html": true
  }'
```

**Output:**
```json
{
  "success": true,
  "message": "Bulk email sent successfully",
  "failed_emails": []
}
```

### Example 4: Verify OTP

```bash
curl -X POST http://localhost:8082/api/v1/email/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'
```

**Output:**
```json
{
  "success": true,
  "message": "OTP verified successfully"
}
```

## API Documentation

### Interactive API Docs

Once the service is running, visit: **http://localhost:8082/swagger/index.html**

Full interactive API documentation is available with the running service where you can test all endpoints.

### Endpoints

#### Send Transactional Email
```http
POST /api/v1/email/send
```

**Description:** Send a single email to one or more recipients with optional template support.

**Request:**
```bash
curl -X POST http://localhost:8082/api/v1/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "subject": "Test Email",
    "body": "This is the email body",
    "is_html": true,
    "template": "welcome",
    "template_data": {"name": "John"}
  }'
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

**Response (Error - 400):**
```json
{
  "error": "invalid request: missing required field 'to'"
}
```

#### Send Bulk Emails
```http
POST /api/v1/email/bulk
```

**Description:** Send emails to multiple recipients with automatic batch processing.

**Request:**
```bash
curl -X POST http://localhost:8082/api/v1/email/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": ["user1@example.com", "user2@example.com"],
    "subject": "Bulk Email",
    "body": "Email content here",
    "is_html": false
  }'
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Bulk email sent successfully",
  "failed_emails": []
}
```

#### Generate OTP
```http
POST /api/v1/email/otp/generate
```

**Description:** Generate a 6-digit OTP and send it via email.

**Request:**
```bash
curl -X POST http://localhost:8082/api/v1/email/otp/generate \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "expiry": 300
  }'
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "otp": "123456",
  "expiry": "2026-03-11T10:35:00Z"
}
```

#### Verify OTP
```http
POST /api/v1/email/otp/verify
```

**Description:** Verify an OTP code for an email address.

**Request:**
```bash
curl -X POST http://localhost:8082/api/v1/email/otp/verify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "123456"
  }'
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "OTP verified successfully"
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "message": "Invalid or expired OTP"
}
```

#### Health Check
```http
GET /health
```

**Description:** Check service health status.

**Request:**
```bash
curl http://localhost:8082/health
```

**Response (Success - 200):**
```json
{
  "status": "healthy",
  "service": "email-service",
  "timestamp": "2026-03-11T10:30:00Z"
}
```

### Error Codes

| Code | Name | Description |
|------|------|-------------|
| 200 | OK | Request successful |
| 400 | Bad Request | Invalid parameters or missing required fields |
| 401 | Unauthorized | SMTP authentication failed |
| 404 | Not Found | Endpoint not found |
| 429 | Too Many Requests | Rate limit exceeded (if implemented) |
| 500 | Server Error | Internal error or SMTP connection failure |

### Available Email Templates

| Template Name | Use Case | Required Data |
|--------------|----------|---------------|
| `welcome` | User registration | `name` |
| `password_reset` | Password recovery | `reset_url` |
| `verification` | Email verification | `verification_url` |
| `otp_template` | OTP delivery | `otp`, `expiry` (auto-generated) |

## Integration Guide

### JavaScript / Node.js + Express

```javascript
const express = require('express');
const axios = require('axios');

const app = express();
const EMAIL_SERVICE_URL = 'http://localhost:8082';

app.use(express.json());

// Send welcome email after user registration
app.post('/register', async (req, res) => {
  try {
    const { email, name } = req.body;
    
    // Your user registration logic here...
    
    // Send welcome email via microservice
    const emailResponse = await axios.post(`${EMAIL_SERVICE_URL}/api/v1/email/send`, {
      to: [email],
      subject: "Welcome to Our Platform",
      template: "welcome",
      template_data: { name: name },
      is_html: true
    });
    
    if (emailResponse.data.success) {
      res.json({ message: 'Registration successful', email_sent: true });
    } else {
      res.status(500).json({ error: 'Registration successful but email failed' });
    }
  } catch (error) {
    console.error('Error:', error.message);
    res.status(500).json({ error: 'Registration failed' });
  }
});

// Generate OTP for two-factor authentication
app.post('/auth/send-otp', async (req, res) => {
  try {
    const { email } = req.body;
    
    const response = await axios.post(`${EMAIL_SERVICE_URL}/api/v1/email/otp/generate`, {
      to: email,
      expiry: 300 // 5 minutes
    });
    
    if (response.data.success) {
      // In production, don't return the OTP - only confirm it was sent
      res.json({ message: 'OTP sent to your email' });
    } else {
      res.status(400).json({ error: 'Failed to send OTP' });
    }
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

app.listen(3000, () => console.log('Server running on port 3000'));
```

### Python / Flask

```python
from flask import Flask, request, jsonify
import requests
from requests.exceptions import RequestException

app = Flask(__name__)
EMAIL_SERVICE_URL = 'http://localhost:8082'

@app.route('/send-notification', methods=['POST'])
def send_notification():
    """Send a notification email to a user"""
    try:
        data = request.get_json()
        user_email = data.get('email')
        message = data.get('message')
        
        # Call email microservice
        response = requests.post(
            f'{EMAIL_SERVICE_URL}/api/v1/email/send',
            json={
                'to': [user_email],
                'subject': 'New Notification',
                'body': f'<p>{message}</p>',
                'is_html': True
            },
            timeout=10
        )
        
        result = response.json()
        
        if result.get('success'):
            return jsonify({'message': 'Notification sent successfully'}), 200
        else:
            return jsonify({'error': 'Failed to send notification'}), 500
            
    except RequestException as e:
        return jsonify({'error': f'Email service unavailable: {str(e)}'}), 503
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/verify-email', methods=['POST'])
def verify_email():
    """Verify OTP for email verification"""
    try:
        data = request.get_json()
        email = data.get('email')
        otp = data.get('otp')
        
        response = requests.post(
            f'{EMAIL_SERVICE_URL}/api/v1/email/otp/verify',
            json={'email': email, 'otp': otp},
            timeout=5
        )
        
        result = response.json()
        
        if result.get('success'):
            return jsonify({'message': 'Email verified successfully'}), 200
        else:
            return jsonify({'error': 'Invalid or expired OTP'}), 400
            
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(debug=True, port=5000)
```

### Java / Spring Boot

```java
package com.example.integration;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api")
public class EmailIntegrationController {
    
    private final RestTemplate restTemplate = new RestTemplate();
    private static final String EMAIL_SERVICE_URL = "http://localhost:8082";
    
    @PostMapping("/send-password-reset")
    public ResponseEntity<?> sendPasswordReset(@RequestBody Map<String, String> request) {
        try {
            String email = request.get("email");
            String resetUrl = "https://yourapp.com/reset?token=" + generateToken();
            
            // Prepare email request
            Map<String, Object> emailRequest = new HashMap<>();
            emailRequest.put("to", List.of(email));
            emailRequest.put("template", "password_reset");
            
            Map<String, String> templateData = new HashMap<>();
            templateData.put("reset_url", resetUrl);
            emailRequest.put("template_data", templateData);
            emailRequest.put("is_html", true);
            
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            
            HttpEntity<Map<String, Object>> entity = new HttpEntity<>(emailRequest, headers);
            
            ResponseEntity<Map> response = restTemplate.exchange(
                EMAIL_SERVICE_URL + "/api/v1/email/send",
                HttpMethod.POST,
                entity,
                Map.class
            );
            
            if (response.getStatusCode().is2xxSuccessful()) {
                return ResponseEntity.ok(Map.of("message", "Password reset email sent"));
            } else {
                return ResponseEntity.status(500).body(Map.of("error", "Failed to send email"));
            }
            
        } catch (Exception e) {
            return ResponseEntity.status(500).body(Map.of("error", e.getMessage()));
        }
    }
    
    @PostMapping("/send-bulk-newsletter")
    public ResponseEntity<?> sendBulkNewsletter(@RequestBody Map<String, Object> request) {
        try {
            @SuppressWarnings("unchecked")
            List<String> recipients = (List<String>) request.get("recipients");
            String subject = (String) request.get("subject");
            String body = (String) request.get("body");
            
            Map<String, Object> bulkRequest = new HashMap<>();
            bulkRequest.put("recipients", recipients);
            bulkRequest.put("subject", subject);
            bulkRequest.put("body", body);
            bulkRequest.put("is_html", true);
            
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            
            HttpEntity<Map<String, Object>> entity = new HttpEntity<>(bulkRequest, headers);
            
            ResponseEntity<Map> response = restTemplate.exchange(
                EMAIL_SERVICE_URL + "/api/v1/email/bulk",
                HttpMethod.POST,
                entity,
                Map.class
            );
            
            return ResponseEntity.ok(response.getBody());
            
        } catch (Exception e) {
            return ResponseEntity.status(500).body(Map.of("error", e.getMessage()));
        }
    }
    
    private String generateToken() {
        // Your token generation logic
        return java.util.UUID.randomUUID().toString();
    }
}
```

### PHP / Laravel

```php
<?php
namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;

class EmailServiceController extends Controller
{
    private $emailServiceUrl = 'http://localhost:8082';
    
    /**
     * Send welcome email to new user
     */
    public function sendWelcomeEmail(Request $request)
    {
        try {
            $request->validate([
                'email' => 'required|email',
                'name' => 'required|string'
            ]);
            
            $response = Http::timeout(10)->post($this->emailServiceUrl . '/api/v1/email/send', [
                'to' => [$request->email],
                'template' => 'welcome',
                'template_data' => [
                    'name' => $request->name
                ],
                'is_html' => true
            ]);
            
            if ($response->successful() && $response->json('success')) {
                return response()->json([
                    'message' => 'Welcome email sent successfully'
                ], 200);
            } else {
                return response()->json([
                    'error' => 'Failed to send email'
                ], 500);
            }
            
        } catch (\Exception $e) {
            Log::error('Email service error: ' . $e->getMessage());
            return response()->json([
                'error' => 'Email service unavailable'
            ], 503);
        }
    }
    
    /**
     * Generate and send OTP for verification
     */
    public function generateOTP(Request $request)
    {
        try {
            $request->validate(['email' => 'required|email']);
            
            $response = Http::post($this->emailServiceUrl . '/api/v1/email/otp/generate', [
                'to' => $request->email,
                'expiry' => 300 // 5 minutes
            ]);
            
            if ($response->successful()) {
                $data = $response->json();
                
                // Store OTP in session or database for verification
                session(['otp_email' => $request->email]);
                
                return response()->json([
                    'message' => 'OTP sent to your email',
                    'expiry' => $data['expiry']
                ]);
            } else {
                return response()->json(['error' => 'Failed to generate OTP'], 500);
            }
            
        } catch (\Exception $e) {
            return response()->json(['error' => $e->getMessage()], 500);
        }
    }
    
    /**
     * Verify OTP
     */
    public function verifyOTP(Request $request)
    {
        try {
            $request->validate([
                'email' => 'required|email',
                'otp' => 'required|string|size:6'
            ]);
            
            $response = Http::post($this->emailServiceUrl . '/api/v1/email/otp/verify', [
                'email' => $request->email,
                'otp' => $request->otp
            ]);
            
            $result = $response->json();
            
            if ($result['success'] ?? false) {
                return response()->json(['message' => 'OTP verified successfully'], 200);
            } else {
                return response()->json(['error' => 'Invalid or expired OTP'], 400);
            }
            
        } catch (\Exception $e) {
            return response()->json(['error' => $e->getMessage()], 500);
        }
    }
}
```

### C# / ASP.NET Core

```csharp
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace EmailIntegration.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class EmailController : ControllerBase
    {
        private readonly HttpClient _httpClient;
        private readonly ILogger<EmailController> _logger;
        private const string EmailServiceUrl = "http://localhost:8082";

        public EmailController(IHttpClientFactory httpClientFactory, ILogger<EmailController> logger)
        {
            _httpClient = httpClientFactory.CreateClient();
            _logger = logger;
        }

        [HttpPost("send-verification")]
        public async Task<IActionResult> SendVerificationEmail([FromBody] VerificationRequest request)
        {
            try
            {
                var emailRequest = new
                {
                    to = new[] { request.Email },
                    template = "verification",
                    template_data = new Dictionary<string, string>
                    {
                        { "verification_url", $"https://yourapp.com/verify?token={request.Token}" }
                    },
                    is_html = true
                };

                var response = await _httpClient.PostAsJsonAsync(
                    $"{EmailServiceUrl}/api/v1/email/send",
                    emailRequest
                );

                if (response.IsSuccessStatusCode)
                {
                    var result = await response.Content.ReadFromJsonAsync<EmailResponse>();
                    
                    if (result?.Success == true)
                    {
                        return Ok(new { message = "Verification email sent" });
                    }
                }

                _logger.LogError("Failed to send verification email");
                return StatusCode(500, new { error = "Failed to send email" });
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error sending verification email");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpPost("send-bulk-announcement")]
        public async Task<IActionResult> SendBulkAnnouncement([FromBody] BulkEmailRequest request)
        {
            try
            {
                var bulkRequest = new
                {
                    recipients = request.Recipients,
                    subject = request.Subject,
                    body = request.Body,
                    is_html = true
                };

                var response = await _httpClient.PostAsJsonAsync(
                    $"{EmailServiceUrl}/api/v1/email/bulk",
                    bulkRequest
                );

                if (response.IsSuccessStatusCode)
                {
                    var result = await response.Content.ReadFromJsonAsync<BulkEmailResponse>();
                    return Ok(result);
                }

                return StatusCode(500, new { error = "Bulk email failed" });
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error sending bulk email");
                return StatusCode(500, new { error = ex.Message });
            }
        }

        [HttpPost("verify-otp")]
        public async Task<IActionResult> VerifyOTP([FromBody] OTPVerificationRequest request)
        {
            try
            {
                var verifyRequest = new
                {
                    email = request.Email,
                    otp = request.OTP
                };

                var response = await _httpClient.PostAsJsonAsync(
                    $"{EmailServiceUrl}/api/v1/email/otp/verify",
                    verifyRequest
                );

                var result = await response.Content.ReadFromJsonAsync<EmailResponse>();

                if (result?.Success == true)
                {
                    return Ok(new { message = "OTP verified successfully" });
                }
                else
                {
                    return BadRequest(new { error = "Invalid or expired OTP" });
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error verifying OTP");
                return StatusCode(500, new { error = ex.Message });
            }
        }
    }

    // DTOs
    public class VerificationRequest
    {
        public string Email { get; set; }
        public string Token { get; set; }
    }

    public class BulkEmailRequest
    {
        public List<string> Recipients { get; set; }
        public string Subject { get; set; }
        public string Body { get; set; }
    }

    public class OTPVerificationRequest
    {
        public string Email { get; set; }
        public string OTP { get; set; }
    }

    public class EmailResponse
    {
        public bool Success { get; set; }
        public string Message { get; set; }
    }

    public class BulkEmailResponse
    {
        public bool Success { get; set; }
        public string Message { get; set; }
        public List<string> FailedEmails { get; set; }
    }
}
```

### Ruby / Rails

```ruby
# app/services/email_service_client.rb
require 'net/http'
require 'json'

class EmailServiceClient
  EMAIL_SERVICE_URL = 'http://localhost:8082'
  
  def self.send_welcome_email(email, name)
    uri = URI("#{EMAIL_SERVICE_URL}/api/v1/email/send")
    
    request_body = {
      to: [email],
      template: 'welcome',
      template_data: { name: name },
      is_html: true
    }.to_json
    
    response = make_request(uri, request_body)
    
    if response.is_a?(Net::HTTPSuccess)
      result = JSON.parse(response.body)
      result['success']
    else
      Rails.logger.error("Failed to send welcome email: #{response.body}")
      false
    end
  rescue StandardError => e
    Rails.logger.error("Email service error: #{e.message}")
    false
  end
  
  def self.generate_otp(email, expiry = 300)
    uri = URI("#{EMAIL_SERVICE_URL}/api/v1/email/otp/generate")
    
    request_body = {
      to: email,
      expiry: expiry
    }.to_json
    
    response = make_request(uri, request_body)
    
    if response.is_a?(Net::HTTPSuccess)
      JSON.parse(response.body)
    else
      { 'success' => false, 'message' => 'Failed to generate OTP' }
    end
  end
  
  def self.verify_otp(email, otp)
    uri = URI("#{EMAIL_SERVICE_URL}/api/v1/email/otp/verify")
    
    request_body = {
      email: email,
      otp: otp
    }.to_json
    
    response = make_request(uri, request_body)
    
    if response.is_a?(Net::HTTPSuccess)
      result = JSON.parse(response.body)
      result['success']
    else
      false
    end
  end
  
  private
  
  def self.make_request(uri, body)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = (uri.scheme == 'https')
    http.read_timeout = 10
    
    request = Net::HTTP::Post.new(uri.path, {
      'Content-Type' => 'application/json'
    })
    request.body = body
    
    http.request(request)
  end
end

# app/controllers/users_controller.rb
class UsersController < ApplicationController
  def create
    @user = User.new(user_params)
    
    if @user.save
      # Send welcome email via microservice
      if EmailServiceClient.send_welcome_email(@user.email, @user.name)
        render json: { message: 'User created successfully' }, status: :created
      else
        render json: { message: 'User created but email failed' }, status: :created
      end
    else
      render json: { errors: @user.errors }, status: :unprocessable_entity
    end
  end
  
  def send_verification_otp
    email = params[:email]
    result = EmailServiceClient.generate_otp(email)
    
    if result['success']
      render json: { message: 'OTP sent to your email' }
    else
      render json: { error: 'Failed to send OTP' }, status: :internal_server_error
    end
  end
  
  private
  
  def user_params
    params.require(:user).permit(:email, :name, :password)
  end
end
```

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Client Applications                        │
│     (Web, Mobile, Desktop, CLI, Other Services)        │
│   (Node.js, Python, Java, PHP, .NET, Ruby, Go, etc.)  │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/REST (JSON)
                     │
┌────────────────────▼────────────────────────────────────┐
│           Sereni Email SMTP Service (Port 8082)         │
│  ┌────────────────────────────────────────────────┐    │
│  │  HTTP Layer                                    │    │
│  │  • Gin Router                                  │    │
│  │  • CORS Middleware                             │    │
│  │  • Swagger Documentation                       │    │
│  │  • Health Check Endpoint                       │    │
│  └────────────┬───────────────────────────────────┘    │
│               │                                         │
│  ┌────────────▼──────────┐                              │
│  │  API Handlers         │                              │
│  │  • Email Handler      │                              │
│  │  • OTP Handler        │                              │
│  │  • Bulk Email Handler │                              │
│  └────────────┬──────────┘                              │
│               │                                         │
│  ┌────────────▼──────────┐                              │
│  │  Business Logic       │                              │
│  │  • Email Service      │                              │
│  │  • Template Engine    │                              │
│  │  • OTP Generator      │                              │
│  │  • Batch Processor    │                              │
│  └────────────┬──────────┘                              │
│               │                                         │
│  ┌────────────▼──────────┐                              │
│  │  SMTP Client          │                              │
│  │  • TLS/SSL Support    │                              │
│  │  • Authentication     │                              │
│  │  • Connection Pool    │                              │
│  └────────────┬──────────┘                              │
└────────────────┼──────────────────────────────────────┘
                 │ SMTP Protocol (Port 587/465)
    ┌────────────┴─────────────┐
    │                          │
    ▼                          ▼
┌─────────┐              ┌─────────┐
│  Gmail  │              │SendGrid │
│ SMTP    │              │ SMTP    │
└─────────┘              └─────────┘
    │                          │
    ▼                          ▼
┌─────────┐              ┌─────────┐
│AWS SES  │              │Mailgun  │
│ SMTP    │              │ SMTP    │
└─────────┘              └─────────┘
```

### Component Responsibilities

| Component | Responsibility |
|-----------|-----------------|
| **HTTP Layer** | Routes incoming requests, handles CORS, serves Swagger docs, provides health checks |
| **API Handlers** | Validates request payloads, invokes business logic, formats responses, handles errors |
| **Email Service** | Manages email sending logic, template processing, OTP generation/verification, batch processing |
| **Template Engine** | Loads and caches email templates, performs variable substitution, renders HTML/text versions |
| **SMTP Client** | Establishes secure connections, authenticates with SMTP servers, sends email data, handles retries |
| **Configuration** | Loads environment variables, provides defaults, validates settings |

### Design Patterns

**Microservice Pattern**: Deployed as an independent service with a well-defined REST API, allowing any application to consume it without language or framework constraints.

**Service Layer Pattern**: Business logic is separated from HTTP handlers, making the code testable and maintainable.

**Template Pattern**: Pre-defined email templates with variable substitution reduce code duplication and ensure consistent branding.

**Batch Processing Pattern**: Bulk emails are automatically split into configurable batches to optimize throughput and prevent memory issues.

**Repository Pattern**: Configuration and template caching are abstracted, allowing for easy swapping of storage mechanisms.

## Development

### Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go             # Configuration loading
│   ├── email/
│   │   ├── service.go            # Core email service logic
│   │   ├── types.go              # Request/response types
│   │   ├── templates.go          # Email templates
│   │   └── generate_otp.go       # OTP generation logic
│   ├── handlers/
│   │   └── email_handler.go      # HTTP handlers
│   └── templatecache/
│       ├── cache.go              # Template caching
│       └── helpers.go            # Template utilities
├── pkg/
│   └── middleware/
│       └── cors.go               # CORS middleware
├── tests/
│   ├── email_test.go             # Email service tests
│   ├── email_handler_test.go     # Handler tests
│   ├── config_test.go            # Configuration tests
│   └── cache_test.go             # Cache tests
├── docs/
│   ├── swagger.yaml              # OpenAPI specification
│   ├── swagger.json              # OpenAPI JSON
│   └── docs.go                   # Generated Swagger docs
├── docker-compose.yml            # Multi-container setup
├── Dockerfile                    # Docker image definition
├── Makefile                      # Build automation
├── .env.example                  # Environment template
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
└── README.md                     # This file
```

### Development Setup

```bash
# 1. Clone repository
git clone https://github.com/aptlogica/sereni-email-smtp.git
cd sereni-email-smtp

# 2. Install Go 1.24.4+ (if not installed)
# Download from: https://golang.org/dl/

# 3. Install dependencies
go mod download

# 4. Install Swagger CLI (for documentation generation)
go install github.com/swaggo/swag/cmd/swag@latest

# 5. Copy environment file
cp .env.example .env

# 6. Edit configuration
nano .env  # or use your preferred editor

# 7. Generate Swagger documentation
make swagger

# 8. Run tests
make test

# 9. Run application
make run
```

### Running Tests

```bash
# All tests
make test
# or: go test -v ./...

# With coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Specific package
go test -v ./internal/email/

# Verbose output with race detection
go test -v -race ./...
```

### Building for Production

```bash
# Build binary
make build
# or: go build -o bin/email-service cmd/server/main.go

# Build optimized binary (smaller size)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/email-service cmd/server/main.go

# Build Docker image
docker build -t sereni-email-smtp:1.0.0 .

# Build multi-platform Docker images
docker buildx build --platform linux/amd64,linux/arm64 -t sereni-email-smtp:1.0.0 .

# Push to registry
docker tag sereni-email-smtp:1.0.0 your-registry/sereni-email-smtp:1.0.0
docker push your-registry/sereni-email-smtp:1.0.0
```

### Development Commands

```bash
# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Update dependencies
go get -u ./...
go mod tidy

# Generate mock for testing
mockgen -source=internal/email/service.go -destination=tests/mocks/email_service_mock.go

# Run with hot reload (requires air)
air
```

## Troubleshooting

### Common Issues

#### 1. SMTP Authentication Failed

**Error Message:**
```
Failed to send email: 535 5.7.8 Username and Password not accepted
```

**Root Cause:**
Invalid SMTP credentials or Gmail blocking "less secure apps".

**Solution:**
```bash
# For Gmail users:
# 1. Enable 2-Factor Authentication in your Google Account
# 2. Generate an App Password:
#    https://myaccount.google.com/apppasswords
# 3. Update .env file with the 16-character app password

# Update .env
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_16_char_app_password

# Restart service
docker-compose restart
```

#### 2. Connection Refused / Service Not Running

**Error Message:**
```
curl: (7) Failed to connect to localhost port 8082: Connection refused
```

**Root Cause:**
Service is not running or port is blocked.

**Solution:**
```bash
# Step 1: Check if service is running
docker ps
# or
ps aux | grep email-service

# Step 2: Check logs
docker logs email-service
# or
docker-compose logs

# Step 3: Start the service
docker-compose up -d

# Step 4: Verify
curl http://localhost:8082/health
```

#### 3. TLS Handshake Timeout

**Error Message:**
```
Failed to send email: dial tcp: i/o timeout
```

**Root Cause:**
SMTP server unreachable or incorrect port.

**Solution:**
```bash
# Verify SMTP settings in .env
# Common ports: 587 (TLS), 465 (SSL), 25 (unencrypted)

# Test SMTP connection
telnet smtp.gmail.com 587
# or
openssl s_client -connect smtp.gmail.com:587 -starttls smtp

# Update .env with correct port
SMTP_PORT=587

# Restart
docker-compose restart
```

#### 4. OTP Verification Fails

**Error Message:**
```json
{"success": false, "message": "Invalid or expired OTP"}
```

**Root Cause:**
OTP expired (default: 5 minutes) or incorrect OTP entered.

**Solution:**
```bash
# Regenerate OTP with longer expiry (in seconds)
curl -X POST http://localhost:8082/api/v1/email/otp/generate \
  -H "Content-Type: application/json" \
  -d '{"to": "user@example.com", "expiry": 600}'

# Ensure correct email and OTP are used
curl -X POST http://localhost:8082/api/v1/email/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "otp": "123456"}'
```

#### 5. Bulk Email Partially Fails

**Error Message:**
```json
{
  "success": true,
  "message": "Bulk email sent with some failures",
  "failed_emails": ["invalid@domain.com"]
}
```

**Root Cause:**
Some recipient email addresses are invalid or SMTP server rejected them.

**Solution:**
```bash
# Check failed_emails in response
# Validate email addresses before sending
# Retry failed emails individually

# Example: Filter out invalid emails
curl -X POST http://localhost:8082/api/v1/email/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": ["valid1@example.com", "valid2@example.com"],
    "subject": "Newsletter",
    "body": "Content here",
    "is_html": true
  }'
```

#### 6. Docker Build Fails

**Error Message:**
```
ERROR [internal] load metadata for docker.io/library/golang:1.24.4-alpine
```

**Root Cause:**
Docker daemon not running or network issues.

**Solution:**
```bash
# Check Docker status
docker info

# Start Docker daemon (Linux)
sudo systemctl start docker

# Pull base image manually
docker pull golang:1.24.4-alpine

# Rebuild
docker-compose build --no-cache
```

### Health Checks

```bash
# Basic health check
curl http://localhost:8082/health

# Expected response:
# {"status":"healthy","service":"email-service","timestamp":"..."}

# Check Docker container health
docker ps
# or
docker inspect email-service | grep -A 10 Health

# Check service logs
docker-compose logs -f email-service

# Check SMTP connectivity from container
docker exec -it email-service sh
telnet smtp.gmail.com 587
```

### Debugging

```bash
# View real-time logs
docker-compose logs -f

# View specific container logs
docker logs -f email-service

# Enable debug logging (add to .env)
LOG_LEVEL=debug

# Test SMTP settings manually
curl -X POST http://localhost:8082/api/v1/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["your_email@gmail.com"],
    "subject": "Test Email",
    "body": "If you get this, SMTP works!",
    "is_html": false
  }'

# Check environment variables in container
docker exec email-service env | grep SMTP

# Shell into container for debugging
docker exec -it email-service sh
```

### Performance Issues

```bash
# Check system resources
docker stats email-service

# Increase batch size for bulk emails (in .env)
BULK_BATCH_SIZE=100

# Scale service with Docker Compose
docker-compose up -d --scale email-service=3

# Monitor concurrent connections
# Check SMTP provider rate limits and adjust accordingly
```

## Contributing

We welcome contributions! Here's how to contribute:

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork: `git clone https://github.com/your-username/sereni-email-smtp.git`
3. Create feature branch: `git checkout -b feature/amazing-feature`
4. Follow code standards (see below)
5. Write tests for new functionality
6. Update README if adding new features
7. Commit with clear messages: `git commit -m "feat: add bulk email rate limiting"`
8. Push to branch: `git push origin feature/amazing-feature`
9. Create Pull Request on GitHub

### Code Standards

- Follow Go idioms and best practices
- Run formatter before committing: `go fmt ./...`
- Use meaningful variable and function names
- Add comments for complex logic and exported functions
- Add Swagger annotations for new API endpoints
- One feature per commit

### Testing Requirements

- Write unit tests for new features using Go's testing framework
- Ensure all tests pass: `make test` or `go test ./...`
- Maintain >80% code coverage
- Test edge cases and error conditions
- Add integration tests for new endpoints

### Commit Message Format

```
type(scope): description

Types: feat, fix, docs, style, refactor, test, chore
Examples:
- feat(api): add rate limiting to bulk email endpoint
- fix(smtp): handle connection timeout gracefully
- docs: update installation instructions for Windows
- test(otp): add expiry validation tests
- refactor(handler): simplify error response formatting
```

### Pull Request Checklist

- [ ] Tests pass (`make test`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Swagger docs updated (`make swagger`)
- [ ] No unnecessary changes or debug code
- [ ] Clear commit messages with conventional format
- [ ] README updated if adding features
- [ ] Addresses an existing issue or feature request

## FAQ

**Q: Can I use Sereni Email SMTP in production?**
A: Yes! It's designed for production use with Docker support, health checks, error handling, comprehensive logging, and battle-tested SMTP integration. Many applications use it for transactional and bulk email delivery.

**Q: What are system requirements?**
A: Minimal requirements: Docker 20.0+ and any SMTP provider credentials. For source builds: Go 1.24.4+. Works on Linux, macOS, and Windows. Resource usage is minimal (typically <50MB RAM).

**Q: Does it work with languages other than Go?**
A: Absolutely! Since Sereni Email SMTP is a REST API microservice, it works with any programming language or framework that can make HTTP requests (Node.js, Python, Java, PHP, .NET, Ruby, etc.). See [Integration Guide](#integration-guide) for examples.

**Q: How do I scale this for high-traffic scenarios?**
A: Deploy multiple instances behind a load balancer (Nginx, HAProxy, or cloud load balancer). Use Docker Compose scale: `docker-compose up -d --scale email-service=5`. For Kubernetes, increase replica count. Adjust `BULK_BATCH_SIZE` for optimal throughput.

**Q: What are the performance limitations?**
A: Performance depends on your SMTP provider's rate limits. Gmail: ~500 emails/day, SendGrid: thousands/hour, AWS SES: 50k+/day. Bulk processing handles thousands concurrently. Monitor your SMTP provider's throughput limits.

**Q: Which SMTP providers are supported?**
A: Any SMTP provider that supports standard SMTP protocol with TLS authentication. Tested with Gmail, Google Workspace, SendGrid, Mailgun, AWS SES, Mailjet, Postmark, Sendinblue, and more. See [Configuration](#configuration) for examples.

**Q: How secure is this service?**
A: Very secure. Uses TLS encryption for SMTP connections, stores credentials in environment variables (never in code), supports Docker secrets for sensitive data, and has no external data storage. Follow security best practices: use strong credentials, restrict network access, and keep Docker images updated.

**Q: Can I customize email templates?**
A: Yes! Modify templates in `internal/email/templates.go` or create custom templates. Templates support variable substitution with {{.variable_name}} syntax for both HTML and plain text versions.

**Q: How do I report issues or request features?**
A: Open an issue on GitHub with a clear description, steps to reproduce (for bugs), and expected vs actual behavior. For features, explain the use case and benefits. See [Contributing](#contributing) for guidelines.

**Q: Is there a hosted version?**
A: Currently self-hosted only. Deploy on your own infrastructure, cloud providers (AWS, GCP, Azure), or container platforms (Kubernetes, Docker Swarm, Nomad). Managed hosting may be available in the future.

**Q: How do OTP expiration and cleanup work?**
A: OTPs expire after configurable time (default: 5 minutes). A background goroutine automatically cleans expired OTPs every 10 minutes to prevent memory buildup. OTP storage is in-memory, so restarting the service clears all OTPs.

## License

This project is licensed under the **MIT License**.

**Why this license?** The MIT License is permissive and allows you to use, modify, and distribute this software freely, including for commercial purposes. It's widely recognized and compatible with most projects.

Full license text: See [LICENSE](LICENSE) file in repository.

## Support & Community

- **Report Issues:** [GitHub Issues](https://github.com/aptlogica/sereni-email-smtp/issues)
- **Discussions:** [GitHub Discussions](https://github.com/aptlogica/sereni-email-smtp/discussions)
- **SonarQube Quality:** [Quality Dashboard](https://sonar.aptlogica.com/dashboard?id=aptlogica_sereni-email-smtp_81bfe731-8d64-4b02-a58c-81c48abd4f9a)

## Acknowledgments

- Built with [Gin Web Framework](https://github.com/gin-gonic/gin)
- Powered by Go's robust [net/smtp](https://pkg.go.dev/net/smtp) package
- Documentation by [Swagger/OpenAPI](https://swagger.io/)
- Inspired by microservice architecture principles

---

**Made with ❤️ by the Sereni Team**

---

## Additional Resources

### Environment Variables Quick Reference

| Variable | Default | Description |
|----------|---------|-------------|
| HOST | 0.0.0.0 | Server bind address |
| PORT | 8082 | HTTP server port |
| ALLOWED_ORIGIN | * | CORS allowed origins |
| SMTP_HOST | smtp.gmail.com | SMTP server hostname |
| SMTP_PORT | 587 | SMTP server port |
| SMTP_USERNAME | - | SMTP authentication username |
| SMTP_PASSWORD | - | SMTP authentication password |
| FROM_EMAIL | - | Default sender email address |
| BULK_BATCH_SIZE | 50 | Batch size for bulk emails |

### Common Commands

```bash
# Build application
make build

# Run tests with coverage
make test
go test -v -cover ./...

# Generate Swagger docs
make swagger

# Run service
make run

# Docker build
docker-compose build

# Docker run
docker-compose up -d

# View logs
docker-compose logs -f

# Clean up
make clean
docker-compose down
```

### Useful Links

- [Official Go Documentation](https://golang.org/doc/)
- [Gin Framework Docs](https://gin-gonic.com/docs/)
- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [Gmail SMTP Settings](https://support.google.com/mail/answer/7126229)
- [SendGrid SMTP Guide](https://docs.sendgrid.com/for-developers/sending-email/getting-started-smtp)
- [Docker Documentation](https://docs.docker.com/)

---

## Quick Integration Checklist

- [ ] Service deployed and running (`docker-compose up -d`)
- [ ] Health check passing (`curl http://localhost:8082/health`)
- [ ] SMTP credentials configured in `.env`
- [ ] Test email sent successfully
- [ ] API client configured in your application
- [ ] Error handling implemented
- [ ] OTP generation tested (if using authentication)
- [ ] Bulk email tested (if sending campaigns)
- [ ] Monitoring/logging configured
- [ ] Production deployment plan created

---

**Ready for production? Deploy with confidence!** 🚀
