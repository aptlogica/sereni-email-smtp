// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"regexp"
	"strings"
	"sync"
	"time"
)

// headerInjectionPattern matches CRLF sequences and null bytes that could be used for header injection
var headerInjectionPattern = regexp.MustCompile(`[\r\n\x00]`)

// bodyControlCharsPattern matches null bytes that enable content smuggling
var bodyControlCharsPattern = regexp.MustCompile(`[\x00]`)

// bodyControlCharsPattern matches control characters (CR, LF, null) that enable content smuggling
var bodyControlCharsPattern = regexp.MustCompile(`[\r\n\x00]`)

// sanitizeHeader removes CRLF characters to prevent email header injection
func sanitizeHeader(input string) string {
	return headerInjectionPattern.ReplaceAllString(input, "")
}

// sanitizeBody removes null bytes that could enable content smuggling
// while preserving normal message content (including newlines)
func sanitizeBody(input string) string {
	return bodyControlCharsPattern.ReplaceAllString(input, "")
}

// EmailService provides methods for sending emails and managing OTPs.
type EmailService struct {
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	FromEmail     string
	otpStore      map[string]OTPEntry
	mutex         sync.RWMutex
	BulkBatchSize int
	// SendEmailFunc, if set, will be used instead of the real SMTP send flow.
	SendEmailFunc func(to []string, subject, body string, isHTML bool) error
	// SendBulkEmailFunc allows overriding bulk send behavior in tests.
	SendBulkEmailFunc func(recipients []string, subject, body string, isHTML bool) ([]string, error)
	// Dial allows injecting a dialer for smtp clients (used in tests).
	Dial func(addr string) (SmtpClient, error)
}

// OTPEntry represents a one-time password entry for email verification.
type OTPEntry struct {
	Code    string
	Email   string
	Created time.Time
	Expiry  time.Time
}

// NewEmailService creates a new EmailService instance with the given SMTP configuration and batch size.
func NewEmailService(smtpHost string, smtpPort int, smtpUsername, smtpPassword, fromEmail string, batchSize int) *EmailService {
	service := &EmailService{
		SMTPHost:      smtpHost,
		SMTPPort:      smtpPort,
		SMTPUsername:  smtpUsername,
		SMTPPassword:  smtpPassword,
		FromEmail:     fromEmail,
		otpStore:      make(map[string]OTPEntry),
		BulkBatchSize: batchSize,
	}

	// Default Dial uses smtp.Dial
	service.Dial = func(addr string) (SmtpClient, error) {
		return smtp.Dial(addr)
	}

	// Start cleanup goroutine for expired OTPs
	go service.cleanupExpiredOTPs()

	return service
}

func (es *EmailService) SendEmail(to []string, subject, body string, isHTML bool) error {
	// Sanitize inputs to prevent header injection - ALWAYS do this first
	sanitizedSubject := sanitizeHeader(subject)
	sanitizedBody := sanitizeBody(body)
	sanitizedTo := make([]string, len(to))
	for i, addr := range to {
		sanitizedTo[i] = sanitizeHeader(addr)
	}
	// Sanitize FromEmail before SMTP envelope usage
	sanitizedFrom := sanitizeHeader(es.FromEmail)
	if sanitizedFrom == "" {
		return fmt.Errorf("invalid or empty sender email after sanitization")
	}

	// If a SendEmailFunc is provided (tests), use it with sanitized inputs
	if es.SendEmailFunc != nil {
		return es.SendEmailFunc(sanitizedTo, sanitizedSubject, sanitizedBody, isHTML)
	}

	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Set up the message
	var message string
	if isHTML {
		message = fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
			sanitizedFrom,
			Join(sanitizedTo, ", "),
			sanitizedSubject,
			sanitizedBody,
		)
	} else {
		message = fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			sanitizedFrom,
			Join(sanitizedTo, ", "),
			sanitizedSubject,
			sanitizedBody,
		)
	}

	// Connect to server
	// Use injected Dial for testability
	conn, err := es.Dial(fmt.Sprintf("%s:%d", es.SMTPHost, es.SMTPPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	// TLS
	if err = conn.StartTLS(&tls.Config{ServerName: es.SMTPHost}); err != nil {
		return err
	}

	// Authenticate
	if err = conn.Auth(auth); err != nil {
		return err
	}

	// Send email
	if err = conn.Mail(sanitizedFrom); err != nil {
		return err
	}

	for _, addr := range sanitizedTo {
		if err = conn.Rcpt(addr); err != nil {
			return err
		}
	}

	writer, err := conn.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return conn.Quit()
}

// SmtpClient abstracts the subset of smtp.Client used by SendEmail
type SmtpClient interface {
	StartTLS(config *tls.Config) error
	Auth(a smtp.Auth) error
	Mail(from string) error
	Rcpt(to string) error
	Data() (io.WriteCloser, error)
	Close() error
	Quit() error
}

func (es *EmailService) SendBulkEmail(recipients []string, subject, body string, isHTML bool) ([]string, error) {
	if es.SendBulkEmailFunc != nil {
		return es.SendBulkEmailFunc(recipients, subject, body, isHTML)
	}
	var failedEmails []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	es.mutex.RLock()
	batchSize := es.BulkBatchSize
	es.mutex.RUnlock()
	if batchSize <= 0 {
		batchSize = 10
	}
	semaphore := make(chan struct{}, batchSize)

	for _, recipient := range recipients {
		if !IsValidEmail(recipient) {
			mu.Lock()
			failedEmails = append(failedEmails, recipient)
			mu.Unlock()
			continue
		}
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			err := es.SendEmail([]string{email}, subject, body, isHTML)
			if err != nil {
				mu.Lock()
				failedEmails = append(failedEmails, email)
				mu.Unlock()
			}
		}(recipient)
	}

	wg.Wait()
	return failedEmails, nil
}

func (es *EmailService) GenerateAndSendOTP(to string, expiryMinutes int) (string, error) {
	// Validate email address before storing/sending
	if !IsValidEmail(to) {
		return "", fmt.Errorf("invalid email address")
	}

	// Generate random OTP
	otp := GenerateOTP()

	// Create OTP entry
	expiryTime := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)
	otpEntry := OTPEntry{
		Code:    otp,
		Email:   to,
		Created: time.Now(),
		Expiry:  expiryTime,
	}

	// Store OTP
	es.mutex.Lock()
	es.otpStore[fmt.Sprintf("%s:%s", to, otp)] = otpEntry
	es.mutex.Unlock()

	// Use template for OTP email
	templateData := map[string]interface{}{
		"otp":    otp,
		"expiry": expiryMinutes,
	}

	err := es.SendTemplateEmail([]string{to}, "otp_template", templateData)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (es *EmailService) VerifyOTP(email, otp string) bool {
	key := fmt.Sprintf("%s:%s", email, otp)
	es.mutex.Lock()
	defer es.mutex.Unlock()

	entry, exists := es.otpStore[key]
	if !exists || time.Now().After(entry.Expiry) {
		return false
	}
	delete(es.otpStore, key)
	return true
}

func (es *EmailService) SendTransactionalEmail(request *EmailRequest) error {
	// Validate recipient list regardless of template usage
	if !IsValidEmailList(request.To) {
		return fmt.Errorf("invalid recipient email(s)")
	}

	// Process template if provided
	subject := request.Subject
	body := request.Body

	if request.Template != "" {
		renderedSubject, renderedBody, err := es.RenderTemplate(request.Template, request.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
		subject = renderedSubject
		body = renderedBody
	}

	// Sanitize subject and body before passing to SendEmail
	sanitizedSubject := sanitizeHeader(subject)
	sanitizedBody := sanitizeBody(body)

	return es.SendEmail(request.To, sanitizedSubject, sanitizedBody, request.IsHTML)
}

func (es *EmailService) cleanupExpiredOTPs() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		es.mutex.Lock()
		for key, entry := range es.otpStore {
			if time.Now().After(entry.Expiry) {
				delete(es.otpStore, key)
			}
		}
		es.mutex.Unlock()
	}
}

func Join(elements []string, sep string) string {
	return strings.Join(elements, sep)
}
