package main

import (
	"fmt"
	"log"
	"os"

	email "github.com/aptlogica/sereni-email-smtp"
)

func main() {
	fmt.Println("=== Sereni Email SMTP - Basic Email Example ===")

	// Load configuration from environment variables
	config := email.Config{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     587,
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", ""),
		UseTLS:   true,
	}

	// Validate configuration
	if config.Username == "" || config.Password == "" || config.From == "" {
		log.Fatal("Please set SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM environment variables")
	}

	// Initialize email provider
	provider := email.NewProvider(config)

	fmt.Printf("Using SMTP server: %s:%d\n", config.Host, config.Port)
	fmt.Printf("From address: %s\n", config.From)

	// Example 1: Send a simple text email
	fmt.Println("\n1. Sending simple text email...")

	textEmail := email.Email{
		To:      []string{"recipient@example.com"},
		Subject: "Test Email from Sereni SMTP",
		Body:    "Hello! This is a test email sent using Sereni Email SMTP provider.",
		IsHTML:  false,
	}

	err := provider.SendEmail(textEmail)
	if err != nil {
		log.Printf("❌ Failed to send text email: %v", err)
	} else {
		fmt.Println("✅ Text email sent successfully!")
	}

	// Example 2: Send HTML email
	fmt.Println("\n2. Sending HTML email...")

	htmlEmail := email.Email{
		To:      []string{"recipient@example.com"},
		Subject: "HTML Email from Sereni SMTP",
		Body: `
			<html>
			<head>
				<title>Welcome to Sereni</title>
			</head>
			<body>
				<h1>🎉 Welcome to Sereni!</h1>
				<p>This is an <strong>HTML email</strong> sent using the Sereni Email SMTP provider.</p>
				
				<h2>Features:</h2>
				<ul>
					<li>✅ Support for multiple SMTP providers</li>
					<li>✅ HTML and text email support</li>
					<li>✅ Template-based emails</li>
					<li>✅ Attachment support</li>
					<li>✅ Bulk email sending</li>
				</ul>
				
				<p>
					<a href="https://github.com/aptlogica/sereni-email-smtp" 
					   style="background-color: #007cba; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
						View Documentation
					</a>
				</p>
				
				<hr>
				<p><small>Sent with ❤️ using Sereni Email SMTP</small></p>
			</body>
			</html>
		`,
		IsHTML: true,
	}

	err = provider.SendEmail(htmlEmail)
	if err != nil {
		log.Printf("❌ Failed to send HTML email: %v", err)
	} else {
		fmt.Println("✅ HTML email sent successfully!")
	}

	// Example 3: Send email with CC and BCC
	fmt.Println("\n3. Sending email with CC and BCC...")

	ccEmail := email.Email{
		To:      []string{"recipient@example.com"},
		CC:      []string{"cc-recipient@example.com"},
		BCC:     []string{"bcc-recipient@example.com"},
		Subject: "Email with CC and BCC",
		Body:    "This email demonstrates CC and BCC functionality.",
		IsHTML:  false,
	}

	err = provider.SendEmail(ccEmail)
	if err != nil {
		log.Printf("❌ Failed to send CC/BCC email: %v", err)
	} else {
		fmt.Println("✅ CC/BCC email sent successfully!")
	}

	// Example 4: Send email with custom headers
	fmt.Println("\n4. Sending email with custom headers...")

	customEmail := email.Email{
		To:      []string{"recipient@example.com"},
		Subject: "Email with Custom Headers",
		Body:    "This email includes custom headers for tracking and identification.",
		IsHTML:  false,
		Headers: map[string]string{
			"X-Priority":   "1",
			"X-Message-ID": "sereni-email-12345",
			"Reply-To":     config.From,
			"X-Campaign":   "welcome-series",
		},
	}

	err = provider.SendEmail(customEmail)
	if err != nil {
		log.Printf("❌ Failed to send custom headers email: %v", err)
	} else {
		fmt.Println("✅ Custom headers email sent successfully!")
	}

	fmt.Println("\n=== Example completed! ===")
	fmt.Println("\nNext steps:")
	fmt.Println("- Check out template-email example for dynamic content")
	fmt.Println("- See bulk-email example for sending to multiple recipients")
	fmt.Println("- Explore attachments example for file uploads")
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
