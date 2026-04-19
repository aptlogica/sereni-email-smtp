// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package email

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/aptlogica/sereni-email-smtp/internal/templatecache"
)

// EmailTemplate represents a template for emails
type EmailTemplate struct {
	Name     string
	Subject  string
	HTMLBody string
	TextBody string
}

// TemplateData holds the data to be used in templates
type TemplateData struct {
	Data map[string]interface{}
}

// Mutex for synchronizing access to predefined templates
var templateMu sync.RWMutex
var tmplCache = templatecache.NewTemplateCache()

// Predefined templates
var predefinedTemplates = map[string]EmailTemplate{
	"welcome": {
		Subject: "Welcome to Our Service!",
		HTMLBody: `
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
            <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
                <h1 style="color: #2c3e50;">Welcome, {{.name}}!</h1>
                <p>Thank you for joining our service. We're excited to have you on board.</p>
                <p>Here are some next steps you might want to take:</p>
                <ul>
                    <li>Complete your profile</li>
                    <li>Explore our features</li>
                    <li>Connect with other users</li>
                </ul>
                <p>If you have any questions, feel free to reach out to our support team.</p>
                <p>Best regards,<br>The Team</p>
            </div>
        </body>
        </html>
        `,
		TextBody: "Welcome, {{.name}}! Thank you for joining our service. We're excited to have you on board.",
	},
	"password_reset": {
		Subject: "Password Reset Request",
		HTMLBody: `
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
            <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
                <h1 style="color: #2c3e50;">Password Reset</h1>
                <p>You have requested to reset your password. Click the button below to reset it:</p>
                <div style="text-align: center; margin: 30px 0;">
                    <a href="{{.reset_url}}" style="background-color: #3498db; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">Reset Password</a>
                </div>
                <p>If you didn't request this, please ignore this email.</p>
                <p>This link will expire in 24 hours.</p>
                <p>Best regards,<br>The Team</p>
            </div>
        </body>
        </html>
        `,
		TextBody: "You have requested to reset your password. Visit this link to reset it: {{.reset_url}}",
	},
	"verification": {
		Subject: "Email Verification Required",
		HTMLBody: `
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
            <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
                <h1 style="color: #2c3e50;">Verify Your Email</h1>
                <p>Please click the button below to verify your email address:</p>
                <div style="text-align: center; margin: 30px 0;">
                    <a href="{{.verification_url}}" style="background-color: #27ae60; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">Verify Email</a>
                </div>
                <p>If you didn't create an account with us, please ignore this email.</p>
                <p>Best regards,<br>The Team</p>
            </div>
        </body>
        </html>
        `,
		TextBody: "Please verify your email by visiting: {{.verification_url}}",
	},
	"otp_template": {
		Subject: "Your Verification Code",
		HTMLBody: `
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
            <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
                <h1 style="color: #2c3e50;">Verification Code</h1>
                <p>Your verification code is:</p>
                <div style="text-align: center; margin: 30px 0;">
                    <span style="font-size: 24px; font-weight: bold; background-color: #f8f9fa; padding: 10px 20px; border-radius: 4px; letter-spacing: 3px;">{{.otp}}</span>
                </div>
                <p>This code will expire in {{.expiry}} minutes.</p>
                <p>If you didn't request this code, please ignore this email.</p>
                <p>Best regards,<br>The Team</p>
            </div>
        </body>
        </html>
        `,
		TextBody: "Your verification code is: {{.otp}} (expires in {{.expiry}} minutes)",
	},
}

// RenderTemplate renders a template with the provided data
func (es *EmailService) RenderTemplate(templateName string, data map[string]interface{}) (string, string, error) {
	templateMu.RLock()
	template, exists := predefinedTemplates[templateName]
	templateMu.RUnlock()
	if !exists {
		return "", "", errors.New("template not found: " + templateName)
	}

	// Get parsed HTML template from cache or parse and store
	htmlTmpl, ok := tmplCache.Get(templateName + ":html")
	if !ok {
		var err error
		htmlTmpl, err = templatecache.ParseAndCacheTemplate(tmplCache, templateName+":html", template.HTMLBody)
		if err != nil {
			return "", "", fmt.Errorf("error parsing HTML body: %w", err)
		}
	}
	var htmlBuf bytes.Buffer
	if err := htmlTmpl.Execute(&htmlBuf, data); err != nil {
		return "", "", fmt.Errorf("error executing HTML body: %w", err)
	}

	// Get parsed subject template from cache or parse and store
	subjTmpl, ok := tmplCache.Get(templateName + ":subject")
	if !ok {
		var err error
		subjTmpl, err = templatecache.ParseAndCacheTemplate(tmplCache, templateName+":subject", template.Subject)
		if err != nil {
			return "", "", fmt.Errorf("error parsing subject: %w", err)
		}
	}
	var subjBuf bytes.Buffer
	if err := subjTmpl.Execute(&subjBuf, data); err != nil {
		return "", "", fmt.Errorf("error executing subject: %w", err)
	}

	return subjBuf.String(), htmlBuf.String(), nil
}

// AddTemplate adds a new template to the service
func (es *EmailService) AddTemplate(name string, template EmailTemplate) {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	predefinedTemplates[name] = template
}

// GetTemplate returns a template by name
func (es *EmailService) GetTemplate(name string) (EmailTemplate, bool) {
	template, exists := predefinedTemplates[name]
	return template, exists
}

// SendTemplateEmail sends an email using a predefined template
func (es *EmailService) SendTemplateEmail(to []string, templateName string, templateData map[string]interface{}) error {
	subject, htmlBody, err := es.RenderTemplate(templateName, templateData)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return es.SendEmail(to, subject, htmlBody, true)
}

// GetAvailableTemplates returns a list of available template names
func (es *EmailService) GetAvailableTemplates() []string {
	var templates []string
	templateMu.RLock()
	for name := range predefinedTemplates {
		templates = append(templates, name)
	}
	templateMu.RUnlock()
	return templates
}
