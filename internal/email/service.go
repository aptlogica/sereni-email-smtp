// internal/email/service.go (updated)
package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"sync"
	"time"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	otpStore     map[string]OTPEntry
	mutex        sync.RWMutex
}

type OTPEntry struct {
	Code    string
	Email   string
	Created time.Time
	Expiry  time.Time
}

func NewEmailService(smtpHost string, smtpPort int, smtpUsername, smtpPassword, fromEmail string) *EmailService {
	service := &EmailService{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		FromEmail:    fromEmail,
		otpStore:     make(map[string]OTPEntry),
	}

	// Start cleanup goroutine for expired OTPs
	go service.cleanupExpiredOTPs()

	return service
}

func (es *EmailService) SendEmail(to []string, subject, body string, isHTML bool) error {
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Set up the message
	var message string
	if isHTML {
		message = fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
			es.FromEmail,
			join(to, ", "),
			subject,
			body,
		)
	} else {
		message = fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			es.FromEmail,
			join(to, ", "),
			subject,
			body,
		)
	}

	fmt.Println("Email message:", message) // Debugging line

	// Connect to server
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", es.SMTPHost, es.SMTPPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	// TLS
	if err = conn.StartTLS(&tls.Config{ServerName: es.SMTPHost}); err != nil {
		return err
	}

	fmt.Println("Auth object:", auth) // Debugging line

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

func (es *EmailService) SendBulkEmail(recipients []string, subject, body string, isHTML bool) ([]string, error) {
	var failedEmails []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	batchSize := 10 // Limit concurrent connections
	semaphore := make(chan struct{}, batchSize)

	for _, recipient := range recipients {
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
	otp := es.generateOTP()

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
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", email, otp)
	entry, exists := es.otpStore[key]

	if !exists || time.Now().After(entry.Expiry) {
		return false
	}

	// Delete OTP after verification
	delete(es.otpStore, key)
	return true
}

func (es *EmailService) SendTransactionalEmail(request *EmailRequest) error {
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

	return es.SendEmail(request.To, subject, body, request.IsHTML)
}

func (es *EmailService) generateOTP() string {
	// Generate 6-digit numeric OTP
	// In production, use crypto/rand for better security
	return fmt.Sprintf("%06d", time.Now().Unix()%1000000)
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

func join(elements []string, sep string) string {
	if len(elements) == 0 {
		return ""
	}
	if len(elements) == 1 {
		return elements[0]
	}
	result := elements[0]
	for _, element := range elements[1:] {
		result += sep + element
	}
	return result
}
