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

// Brand tokens, taken from the values serenilrs.com declares as its
// --sereni-* CSS custom properties. Email clients support neither custom
// properties nor reliable webfonts, so the palette is inlined as literals and
// every face ships with a real fallback stack.
const (
	brandBlack  = "#0d0d0d" // --sereni-black, also the site's theme-color
	brandBlue   = "#2f0df2" // --sereni-blue, primary action
	brandGreen  = "#17ed92" // --sereni-green, used on dark ground only
	brandGround = "#f5f5fa" // neutral biased toward the brand blue
	brandCard   = "#ffffff"
	brandInk    = "#1f1f29"
	brandInkMid = "#4a4a5c"
	brandFaint  = "#7a7a90"

	brandSans = "Inter,-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif"
	brandMono = "'JetBrains Mono',Menlo,Consolas,monospace"

	// The site's own logo is a white-filled SVG built for its dark ground, and
	// Outlook and Gmail both strip inline SVG. A hosted PNG reversed out of a
	// brand-black header block is the only form that renders everywhere.
	brandLogoURL = "https://www.serenilrs.com/icons/icon-192.png"
	brandSiteURL = "https://www.serenilrs.com"

	brandWordmark = "Sereni LRS"
	brandSupport  = "support@serenilrs.com"
	brandLegal    = "\u00a9 2015 - 2030 Sereni LRS. All rights reserved."
	brandTagline  = "Enterprise xAPI Learning Record Store"
)

// shell wraps body content in the branded email frame.
//
// Table-based rather than div-based on purpose: Outlook's Word rendering engine
// ignores max-width and border-radius on divs, so a fixed-width table is the
// only layout that holds up across clients. preheader is the grey text an inbox
// shows beside the subject; leaving it empty lets clients scrape the first
// visible words instead, which reads badly.
func shell(preheader, bodyHTML string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="color-scheme" content="light">
</head>
<body style="margin:0;padding:0;background-color:` + brandGround + `;">
<div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;height:0;width:0;">` + preheader + `</div>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:` + brandGround + `;">
<tr><td align="center" style="padding:24px 12px;">

<table role="presentation" width="600" cellpadding="0" cellspacing="0" border="0" style="width:600px;max-width:600px;background-color:` + brandCard + `;border-radius:12px;overflow:hidden;">
<tr><td style="background-color:` + brandBlack + `;padding:22px 32px;">
<img src="` + brandLogoURL + `" width="34" height="34" alt="` + brandWordmark + `" style="display:inline-block;vertical-align:middle;border:0;outline:none;text-decoration:none;">
<span style="display:inline-block;vertical-align:middle;padding-left:12px;font-family:` + brandSans + `;font-size:18px;font-weight:600;letter-spacing:-0.01em;color:#ffffff;">` + brandWordmark + `</span>
</td></tr>
<tr><td style="padding:36px 32px 38px;font-family:` + brandSans + `;font-size:15px;line-height:1.6;color:` + brandInkMid + `;">
` + bodyHTML + `
</td></tr>
</table>

<table role="presentation" width="600" cellpadding="0" cellspacing="0" border="0" style="width:600px;max-width:600px;">
<tr><td align="center" style="padding:20px 32px 8px;font-family:` + brandSans + `;font-size:12px;line-height:1.6;color:` + brandFaint + `;">
` + brandLegal + `<br>
` + brandTagline + ` &middot; <a href="` + brandSiteURL + `" style="color:` + brandBlue + `;text-decoration:none;">serenilrs.com</a>
</td></tr>
</table>

</td></tr>
</table>
</body>
</html>`
}

// heading renders an h1 at brand ink weight.
func heading(text string) string {
	return `<h1 style="margin:0 0 18px;font-family:` + brandSans + `;font-size:22px;line-height:1.25;font-weight:650;letter-spacing:-0.02em;color:` + brandInk + `;">` + text + `</h1>`
}

// button renders a call to action. bgcolor on the cell as well as the style is
// what makes the fill survive Outlook.
func button(href, label string) string {
	return `<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:28px 0;"><tr>
<td align="center" bgcolor="` + brandBlue + `" style="border-radius:8px;">
<a href="` + href + `" style="display:inline-block;padding:13px 30px;font-family:` + brandSans + `;font-size:15px;font-weight:600;color:#ffffff;text-decoration:none;border-radius:8px;">` + label + `</a>
</td></tr></table>`
}

// support renders the standard closing offer of help. Every email points at the
// support address rather than inviting a reply: transactional mail is sent from
// an unattended noreply sender, so a reply goes nowhere.
func support() string {
	return `<p style="margin:16px 0 0;">If you have any questions or need assistance, please contact our support team &mdash; <a href="mailto:` + brandSupport + `" style="color:` + brandBlue + `;text-decoration:none;font-weight:500;">` + brandSupport + `</a></p>`
}

// signOff closes every message with a named team rather than "The Team".
func signOff() string {
	return `<p style="margin:24px 0 0;">Best regards,<br><strong style="color:` + brandInk + `;">The ` + brandWordmark + ` Team</strong></p>`
}

// Predefined templates
var predefinedTemplates = map[string]EmailTemplate{
	"welcome": {
		Subject: "Welcome to Sereni LRS",
		HTMLBody: shell(
			"Your Sereni LRS account is ready.",
			heading("Welcome, {{.name}}")+
				`<p style="margin:0 0 16px;">Your Sereni LRS account is ready. This is where all your learning activity, insights and progress come together in one place.</p>
<p style="margin:0 0 8px;">A good place to start:</p>
<ul style="margin:0 0 16px;padding-left:20px;">
<li style="margin-bottom:6px;">Complete your profile</li>
<li style="margin-bottom:6px;">Send your first xAPI statements</li>
<li style="margin-bottom:6px;">Explore dashboards and reporting</li>
</ul>
<p style="margin:0;">Once you are in, you can start exploring statements, tracking activity and building reports.</p>`+
				support()+
				signOff(),
		),
		TextBody: "Welcome, {{.name}}! Your Sereni LRS account is ready. This is where all your learning activity, insights and progress come together in one place.\n\nIf you have any questions or need assistance, contact our support team at support@serenilrs.com\n\nBest regards,\nThe Sereni LRS Team",
	},
	"password_reset": {
		Subject: "Reset your Sereni LRS password",
		HTMLBody: shell(
			"Reset your Sereni LRS password. This link expires in 24 hours.",
			heading("Reset your password")+
				`<p style="margin:0 0 16px;">You asked to reset the password for your Sereni LRS account. Use the button below to choose a new one.</p>`+
				button("{{.reset_url}}", "Reset password")+
				`<p style="margin:0 0 16px;">This link expires in 24 hours. If you did not request a reset, you can safely ignore this email and your password stays unchanged.</p>`+
				support()+
				signOff(),
		),
		TextBody: "You asked to reset your Sereni LRS password. Open this link to choose a new one: {{.reset_url}}\n\nThe link expires in 24 hours. If you did not request a reset, ignore this email.\n\nIf you have any questions or need assistance, contact our support team at support@serenilrs.com\n\nBest regards,\nThe Sereni LRS Team",
	},
	"verification": {
		Subject: "Verify your email for Sereni LRS",
		HTMLBody: shell(
			"Confirm your email address to finish setting up Sereni LRS.",
			heading("Verify your email")+
				`<p style="margin:0 0 16px;">Confirm this email address to finish setting up your Sereni LRS account.</p>`+
				button("{{.verification_url}}", "Verify email address")+
				`<p style="margin:0 0 16px;">If you did not create a Sereni LRS account, you can safely ignore this email.</p>`+
				support()+
				signOff(),
		),
		TextBody: "Confirm your email address to finish setting up your Sereni LRS account: {{.verification_url}}\n\nIf you did not create an account, ignore this email.\n\nIf you have any questions or need assistance, contact our support team at support@serenilrs.com\n\nBest regards,\nThe Sereni LRS Team",
	},
	"otp_template": {
		Subject: "Your Sereni LRS verification code",
		HTMLBody: shell(
			"Your Sereni LRS verification code.",
			heading("Your verification code")+
				`<p style="margin:0 0 8px;">Enter this code to continue signing in to Sereni LRS.</p>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="margin:24px 0;"><tr>
<td align="center" bgcolor="`+brandBlack+`" style="border-radius:10px;padding:20px 24px;">
<span style="font-family:`+brandMono+`;font-size:30px;font-weight:700;letter-spacing:8px;color:`+brandGreen+`;">{{.otp}}</span>
</td></tr></table>
<p style="margin:0 0 16px;">The code expires in {{.expiry}} minutes. If you did not request it, you can safely ignore this email.</p>`+
				support()+
				signOff(),
		),
		TextBody: "Your Sereni LRS verification code is {{.otp}}. It expires in {{.expiry}} minutes.\n\nIf you did not request it, ignore this email.\n\nIf you have any questions or need assistance, contact our support team at support@serenilrs.com\n\nBest regards,\nThe Sereni LRS Team",
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

	// SECURITY: Sanitize template data to prevent XSS and URL injection attacks
	es.mutex.RLock()
	trustedConfig := es.TrustedDomainConfig
	es.mutex.RUnlock()

	sanitizedData, err := SanitizeTemplateData(data, trustedConfig)
	if err != nil {
		return "", "", fmt.Errorf("template data validation failed: %w", err)
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
	if err := htmlTmpl.Execute(&htmlBuf, sanitizedData); err != nil {
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
	if err := subjTmpl.Execute(&subjBuf, sanitizedData); err != nil {
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
