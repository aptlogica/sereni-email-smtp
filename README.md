<h1 align="center">sereni-email-smtp - Production-Grade SMTP Email Service</h1>

<p align="center">Enterprise-grade SMTP email service and open source email server for mission-critical applications. A comprehensive transactional email service and email delivery service providing reliable email sending, advanced queueing, retry mechanisms, and seamless cloud-native integration.</p>

<p align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/Version-1.0.0-blue?style=for-the-badge" alt="Version"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.26.2-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"></a>
<a href="https://en.wikipedia.org/wiki/Simple_Mail_Transfer_Protocol"><img src="https://img.shields.io/badge/SMTP-Protocol-EA4335?style=for-the-badge&logo=gmail&logoColor=white" alt="SMTP"></a>
<a href="https://gin-gonic.com/"><img src="https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge&logo=gin&logoColor=white" alt="Gin"></a>
<a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker"></a>
<a href="https://swagger.io/"><img src="https://img.shields.io/badge/Swagger-Documented-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" alt="Swagger"></a>
</p>

<p align="center">
<a href="https://github.com/aptlogica/sereni-email-smtp/actions/workflows/ci.yml"><img src="https://github.com/aptlogica/sereni-email-smtp/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
<a href="https://github.com/aptlogica/sereni-email-smtp/actions/workflows/codeql.yml"><img src="https://github.com/aptlogica/sereni-email-smtp/actions/workflows/codeql.yml/badge.svg" alt="CodeQL"></a>
<a href="https://sonarcloud.io/dashboard?id=aptlogica_sereni-email-smtp"><img src="https://sonarcloud.io/api/project_badges/measure?project=aptlogica_sereni-email-smtp&metric=alert_status" alt="Quality Gate"></a>
<a href="https://sonarcloud.io/dashboard?id=aptlogica_sereni-email-smtp"><img src="https://sonarcloud.io/api/project_badges/measure?project=aptlogica_sereni-email-smtp&metric=coverage" alt="Coverage"></a>
<a href="https://sonarcloud.io/dashboard?id=aptlogica_sereni-email-smtp"><img src="https://sonarcloud.io/api/project_badges/measure?project=aptlogica_sereni-email-smtp&metric=security_rating" alt="Security"></a>
</p>

<p align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License: Apache 2.0"></a>
</p>


## Overview

**sereni-email-smtp** is an open-source, self-hosted SMTP email service and developer email API built for backend applications. It enables reliable delivery of transactional emails such as verification emails, password resets, alerts, and notifications using queueing, retry logic, and observability features. Designed for scalability and flexibility, it integrates easily with REST APIs, microservices, and popular SMTP providers like Gmail, SendGrid, and Mailgun.

sereni-email-smtp runs on <code>:8082</code> as part of the SereniBase backend platform. SereniBase uses this service for user registration emails, password resets, and workflow notifications. See <a href="https://github.com/aptlogica/sereni-base">sereni-base</a> to deploy the full stack.</p>
## Key Features

- **Enterprise SMTP Service**: High-throughput email delivery with advanced error handling
- **Intelligent Queueing**: Redis-backed queue system with priority handling and dead letter queues
- **Retry Logic**: Configurable retry strategies with exponential backoff and circuit breakers
- **Comprehensive Monitoring**: Detailed metrics, logging, and observability dashboards
- **Security First**: Secure credential management, TLS enforcement, and audit logging
- **SMTP Mail Service**: Complete email integration service with SMTP email toolkit support
- **Cloud-Native Ready**: Kubernetes deployment with horizontal scaling capabilities

## Architecture
- Go 1.26.2, idiomatic design
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

# Queue & retry settings
# Maximum times a message will be retried before moving to dead-letter queue
MAX_RETRY_ATTEMPTS=3
# Base backoff in seconds used for exponential backoff between retries
RETRY_BACKOFF_SECONDS=30
# Number of consecutive failures before a circuit breaker trips
CIRCUIT_BREAKER_THRESHOLD=5
# Dead letter queue name (Redis key / stream / list depending on implementation)
DEAD_LETTER_QUEUE=sereni:email:dlq

# Priority queue names
EMAIL_PRIORITY_HIGH=sereni:email:high
EMAIL_PRIORITY_NORMAL=sereni:email:normal
```

### Docker Development

Use the example `docker-compose.yml` below to run Redis and the email service together. The `email` service is configured to connect to Redis using `REDIS_URL=redis://redis:6379` and exposes port `8082`.

```yaml
version: "3.8"
services:
    redis:
        image: redis:7-alpine
        ports:
            - "6379:6379"
        volumes:
            - redis-data:/data

    email:
        build: .
        environment:
            - REDIS_URL=redis://redis:6379
            - PORT=8082
            - LOG_LEVEL=debug
        ports:
            - "8082:8082"
        depends_on:
            - redis

volumes:
    redis-data:
```

Start the stack:

```bash
docker-compose up -d
# Run the service locally (or build the image to run the `email` service)
go run ./cmd/server
```

## Testing
- Run `go test ./...` to execute unit tests

## Security
See [SECURITY.md](SECURITY.md) for reporting vulnerabilities.

## License
Apache License 2.0. Copyright (c) 2026 Aptlogica Technologies.

## API Documentation

See [docs/openapi.yaml](docs/openapi.yaml) for the OpenAPI 3.0 specification of the SMTP Email Service API endpoints.

You can use tools like Swagger UI or Redoc to visualize and interact with the API.


