package main

import (
	"fmt"
	"log"

	"github.com/aptlogica/sereni-email-smtp/internal/config"
	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

const exampleRecipient = "recipient@example.com"

func main() {
	fmt.Println("=== Sereni Email SMTP - Basic Email Example ===")

	// Load configuration from environment variables
	cfg := config.LoadConfig()

	// Validate configuration
	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" || cfg.FromEmail == "" {
		log.Fatal("Please set SMTP_USERNAME, SMTP_PASSWORD, and FROM_EMAIL environment variables")
	}

	// Initialize email service
	service := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.FromEmail,
		cfg.BulkBatchSize,
	)

	fmt.Printf("Using SMTP server: %s:%d\n", cfg.SMTPHost, cfg.SMTPPort)
	fmt.Printf("From address: %s\n", cfg.FromEmail)

	// Example 1: Send a simple text email
	fmt.Println("\n1. Sending simple text email...")

	err := service.SendEmail(
		[]string{exampleRecipient},
		"Test Email from Sereni SMTP",
		"Hello! This is a test email sent using Sereni Email SMTP provider.",
		false,
	)
	if err != nil {
		log.Printf("Failed to send text email: %v", err)
	} else {
		fmt.Println("Text email sent successfully!")
	}

	// Example 2: Send HTML email
	fmt.Println("\n2. Sending HTML email...")

	htmlBody := `
<html>
<head>
	<title>Welcome to Sereni</title>
</head>
<body>
	<h1>Welcome to Sereni</h1>
	<p>This is an <strong>HTML email</strong> sent using the Sereni Email SMTP service.</p>

	<h2>Features:</h2>
	<ul>
		<li>Support for multiple SMTP providers</li>
		<li>HTML and text email support</li>
		<li>Template-based emails</li>
		<li>Bulk email sending</li>
	</ul>

	<p>
		<a href="https://github.com/aptlogica/sereni-email-smtp"
		   style="background-color: #007cba; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
			View Documentation
		</a>
	</p>
</body>
</html>
`

	err = service.SendEmail(
		[]string{exampleRecipient},
		"HTML Email from Sereni SMTP",
		htmlBody,
		true,
	)
	if err != nil {
		log.Printf("Failed to send HTML email: %v", err)
	} else {
		fmt.Println("HTML email sent successfully!")
	}

	// Example 3: Send email using a predefined template
	fmt.Println("\n3. Sending templated email...")

	err = service.SendTemplateEmail(
		[]string{exampleRecipient},
		"welcome",
		map[string]interface{}{
			"name": "Sereni User",
		},
	)
	if err != nil {
		log.Printf("Failed to send templated email: %v", err)
	} else {
		fmt.Println("Templated email sent successfully!")
	}

	// Example 4: Send bulk email
	fmt.Println("\n4. Sending bulk email...")

	failed, err := service.SendBulkEmail(
		[]string{"recipient1@example.com", "recipient2@example.com"},
		"Bulk Email from Sereni SMTP",
		"This is a bulk email sent using Sereni Email SMTP service.",
		false,
	)
	if err != nil {
		log.Printf("Failed to send bulk email: %v", err)
	} else if len(failed) > 0 {
		// Log count only to avoid log injection from user-provided email addresses
		log.Printf("Bulk email completed with %d failures", len(failed))
	} else {
		fmt.Println("Bulk email sent successfully!")
	}

	fmt.Println("\n=== Example completed! ===")
}
