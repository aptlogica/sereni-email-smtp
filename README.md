# sereni-email-smtp - Production-Grade SMTP Email Service

> Enterprise-grade SMTP email service and open source email server for mission-critical applications. A comprehensive transactional email service and email delivery service providing reliable email sending, advanced queueing, retry mechanisms, and seamless cloud-native integration.

[![Version](https://img.shields.io/badge/Version-1.0.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Quality Gate Status](https://sonar.aptlogica.com/api/project_badges/measure?project=aptlogica_sereni-email-smtp_12345678&metric=alert_status&token=sqb_152d71a0f9a3621514372a3e4c87460e3059bbc2)](https://sonar.aptlogica.com/dashboard?id=aptlogica_sereni-email-smtp_12345678)

## Overview

**sereni-email-smtp** is a production-ready SMTP email provider and email microservice engineered for reliability, security, and observability in enterprise environments. This comprehensive backend email service and developer email API features advanced queueing mechanisms, intelligent retry logic, comprehensive metrics, and seamless integration with modern cloud-native infrastructure as an email sending service, an email service for backend systems, and an open source smtp provider for transactional workflows. Complete email backend solution for backend applications.

## Key Features

- **Enterprise SMTP Service**: High-throughput email delivery with advanced error handling
- **Intelligent Queueing**: Redis-backed queue system with priority handling and dead letter queues
- **Retry Logic**: Configurable retry strategies with exponential backoff and circuit breakers
- **Comprehensive Monitoring**: Detailed metrics, logging, and observability dashboards
- **Security First**: Secure credential management, TLS enforcement, and audit logging
- **SMTP Mail Service**: Complete email integration service with SMTP email toolkit support
- **Cloud-Native Ready**: Kubernetes deployment with horizontal scaling capabilities

## Architecture
- Go 1.23+, idiomatic design
- Modular, testable codebase

## Installation
```sh
go get github.com/aptlogica/sereni-email-smtp
```

## Configuration
See `.env.example` for environment variables and configuration options.

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/aptlogica/sereni-email-smtp/pkg/client"
    "github.com/aptlogica/sereni-email-smtp/pkg/config"
    "github.com/aptlogica/sereni-email-smtp/pkg/types"
)

func main() {
    // Initialize configuration
    cfg := config.New()
    cfg.SMTPHost = "smtp.gmail.com"
    cfg.SMTPPort = 587
    cfg.Username = "your-email@gmail.com"
    cfg.Password = "your-app-password"
    
    // Create email client
    client, err := client.New(cfg)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer client.Close()
    
    // Send email
    email := &types.Email{
        To:      []string{"recipient@example.com"},
        Subject: "Welcome to SereniBase",
        Body:    "Hello and welcome to our platform!",
        IsHTML:  false,
    }
    
    ctx := context.Background()
    if err := client.Send(ctx, email); err != nil {
        log.Fatal("Failed to send email:", err)
    }
    
    log.Println("Email sent successfully")
}
```

## Development

### Local Setup
```bash
# Clone the repository
git clone https://github.com/aptlogica/sereni-email-smtp.git
cd sereni-email-smtp

# Install dependencies
go mod download

# Set up environment
cp .env.example .env
# Configure your SMTP settings in .env

# Start development server
go run ./cmd/server
```

### Environment Configuration
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
REDIS_URL=redis://localhost:6379
PORT=8080
LOG_LEVEL=debug
```

### Docker Development
```bash
# Start with dependencies
docker-compose up -d redis

# Run the service
go run ./cmd/server
```

## Testing
- Run `go test ./...` to execute unit tests

## Security
See [SECURITY.md](SECURITY.md) for reporting vulnerabilities.

## License
MIT License. Copyright (c) 2026 Aptlogica Technologies.


