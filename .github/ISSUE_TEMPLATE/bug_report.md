---
name: Bug report
about: Create a report to help us improve
title: '[BUG] '
labels: 'bug'
assignees: ''
---

## Bug Description

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Configure SMTP settings with '...'
2. Send email with '...'
3. Check email delivery
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Actual behavior**
A clear and concise description of what actually happened.

## Environment

**sereni-email-smtp Version:** [e.g. v1.0.0]
**Go Version:** [e.g. 1.23.1]
**SMTP Provider:** [e.g. Gmail, SendGrid, AWS SES]
**OS:** [e.g. Ubuntu 20.04, Windows 11, macOS 13]

## Configuration

```yaml
# Share relevant SMTP configuration (remove sensitive data)
smtp:
  host: smtp.gmail.com
  port: 587
  # Do not include passwords/keys
```

## Error Output

```
// Paste any error messages or logs here
```

## Email Details (if applicable)

- **To:** [recipient email - use example.com for privacy]
- **Subject:** [email subject]
- **Content Type:** [text/html, text/plain]
- **Attachments:** [yes/no, file types]

**Additional Context**

Add any other context about the problem here, including:
- Email delivery status
- SMTP server responses
- Network connectivity issues

**Possible Solution**
If you have suggestions for fixing the bug, please share them here.