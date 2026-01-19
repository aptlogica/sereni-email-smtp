// internal/email/service.go (updated)
package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

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

type OTPEntry struct {
	Code    string
	Email   string
	Created time.Time
	Expiry  time.Time
}

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
	// If a SendEmailFunc is provided (tests), use it
	if es.SendEmailFunc != nil {
		return es.SendEmailFunc(to, subject, body, isHTML)
	}
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Set up the message
	var message string
	if isHTML {
		message = fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
			es.FromEmail,
			Join(to, ", "),
			subject,
			body,
		)
	} else {
		message = fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			es.FromEmail,
			Join(to, ", "),
			subject,
			body,
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
	if err = conn.Mail(es.FromEmail); err != nil {
		return err
	}

	for _, addr := range to {
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

	batchSize := es.BulkBatchSize
	if batchSize <= 0 {
		batchSize = 10
	}
	semaphore := make(chan struct{}, batchSize)

	for _, recipient := range recipients {
		if !IsValidEmail(recipient) {
			failedEmails = append(failedEmails, recipient)
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
	// Process template if provided
	subject := request.Subject
	body := request.Body

	if request.Template != "" {
		if !IsValidEmailList(request.To) {
			return fmt.Errorf("invalid recipient email(s)")
		}
		renderedSubject, renderedBody, err := es.RenderTemplate(request.Template, request.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
		subject = renderedSubject
		body = renderedBody
	}

	return es.SendEmail(request.To, subject, body, request.IsHTML)
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
