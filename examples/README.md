# Sereni Email SMTP Examples

This directory contains practical examples demonstrating how to use the Sereni Email SMTP Provider for various email scenarios.

## Examples Overview

| Example | Description | Complexity |
|---------|-------------|------------|
| [Basic Email](./basic-email/) | Send simple text and HTML emails | Beginner |
| [Template Email](./template-email/) | Use email templates with variables | Intermediate |
| [Bulk Email](./bulk-email/) | Send emails to multiple recipients | Intermediate |
| [Email with Attachments](./attachments/) | Send emails with file attachments | Intermediate |
| [Multi-Provider](./multi-provider/) | Configure multiple SMTP providers | Advanced |

## Quick Start

Choose an example that matches your use case and follow the README in each directory.

### Basic Usage

```go
go run basic-email/main.go
```

### With Templates

```go
go run template-email/main.go
```

## Prerequisites

- Go 1.21+
- SMTP server credentials (Gmail, SendGrid, Mailgun, etc.)

## Common Configuration

Most examples use these environment variables:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@gmail.com
```

### Gmail Setup

1. Enable 2-Factor Authentication
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Use the app password instead of your regular password

### SendGrid Setup

```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

### Mailgun Setup

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=your-mailgun-smtp-username
SMTP_PASSWORD=your-mailgun-smtp-password
```