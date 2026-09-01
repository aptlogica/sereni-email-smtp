// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"strings"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

// The values serenilrs.com declares as its --sereni-* custom properties. If the
// site's palette changes, these are the constants to update.
const (
	brandBlack = "#0d0d0d"
	brandBlue  = "#2f0df2"
	brandGreen = "#17ed92"
)

// Stock palettes from the two email boilerplates these templates were
// originally copied from. Their reappearance means a template was written
// without reference to the brand.
var offBrandColours = []string{"#2c3e50", "#3498db", "#27ae60", "#3869D4", "#3869d4"}

func templateNames() []string {
	return []string{"welcome", "password_reset", "verification", "otp_template"}
}

func getBody(t *testing.T, name string) string {
	t.Helper()
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	tmpl, ok := service.GetTemplate(name)
	if !ok {
		t.Fatalf("template %q not found", name)
	}
	return tmpl.HTMLBody
}

// TestTemplatesCarryBrandIdentity guards the branding fix: every predefined
// template must show the wordmark, the logo and the brand ground.
func TestTemplatesCarryBrandIdentity(t *testing.T) {
	for _, name := range templateNames() {
		t.Run(name, func(t *testing.T) {
			body := getBody(t, name)

			for _, want := range []string{
				"Sereni LRS",
				"icon-192.png",
				brandBlack,
				"serenilrs.com",
				"2015 - 2030",
			} {
				if !strings.Contains(body, want) {
					t.Errorf("template %q is missing %q", name, want)
				}
			}
		})
	}
}

// TestTemplatesUseBrandTypeface checks the font stack rather than a bare family
// name: email clients cannot load webfonts, so a fallback chain is required.
func TestTemplatesUseBrandTypeface(t *testing.T) {
	for _, name := range templateNames() {
		body := getBody(t, name)
		if !strings.Contains(body, "Inter,") {
			t.Errorf("template %q does not use Inter", name)
		}
		if !strings.Contains(body, "Arial,sans-serif") {
			t.Errorf("template %q has no fallback after Inter", name)
		}
	}
}

// TestTemplatesHaveNoOffBrandColours fails if a stock boilerplate colour
// reappears in any template.
func TestTemplatesHaveNoOffBrandColours(t *testing.T) {
	for _, name := range templateNames() {
		body := getBody(t, name)
		for _, bad := range offBrandColours {
			if strings.Contains(body, bad) {
				t.Errorf("template %q uses off-brand colour %s", name, bad)
			}
		}
	}
}

// TestTemplatesNameTheProduct catches the generic copy the templates shipped
// with, which never told a recipient who had emailed them.
func TestTemplatesNameTheProduct(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	for _, name := range templateNames() {
		tmpl, _ := service.GetTemplate(name)

		if !strings.Contains(tmpl.Subject, "Sereni LRS") {
			t.Errorf("subject for %q does not name the product: %q", name, tmpl.Subject)
		}
		for _, generic := range []string{"Our Service", "our service"} {
			if strings.Contains(tmpl.HTMLBody, generic) {
				t.Errorf("template %q still contains generic copy %q", name, generic)
			}
		}
		// "The Team" is the old anonymous sign-off; "The Sereni LRS Team" is the
		// replacement and legitimately contains it, so match the standalone form.
		if strings.Contains(tmpl.HTMLBody, ">The Team<") {
			t.Errorf("template %q still signs off as an anonymous team", name)
		}
		if !strings.Contains(tmpl.TextBody, "Sereni LRS") {
			t.Errorf("text body for %q does not name the product", name)
		}
	}
}

// TestActionTemplatesUseBrandButton checks the two templates with a call to
// action render it on the brand blue, with bgcolor set for Outlook.
func TestActionTemplatesUseBrandButton(t *testing.T) {
	for _, name := range []string{"password_reset", "verification"} {
		body := getBody(t, name)
		if !strings.Contains(body, brandBlue) {
			t.Errorf("template %q does not use the brand blue for its action", name)
		}
		if !strings.Contains(body, `bgcolor="`+brandBlue+`"`) {
			t.Errorf("template %q button has no bgcolor attribute; it will render unfilled in Outlook", name)
		}
	}
}

// TestOTPUsesBrandGreenOnDark pins the one place the brand green is legible:
// on the brand black. On white it fails contrast, so it must not appear there.
func TestOTPUsesBrandGreenOnDark(t *testing.T) {
	body := getBody(t, "otp_template")
	if !strings.Contains(body, brandGreen) {
		t.Error("otp_template does not use the brand green for the code")
	}
	if !strings.Contains(body, `bgcolor="`+brandBlack+`"`) {
		t.Error("otp_template does not place the code on the brand black ground")
	}
}

// TestTemplatesHavePreheader checks each template sets inbox preview text
// instead of letting clients scrape the first visible markup.
// TestTemplatesPointAtSupport guards the closing line. Transactional mail is
// sent from an unattended sender, so "reply to this email" sends the recipient
// nowhere; every template must name the support address instead.
func TestTemplatesPointAtSupport(t *testing.T) {
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)

	for _, name := range templateNames() {
		tmpl, _ := service.GetTemplate(name)

		if !strings.Contains(tmpl.HTMLBody, "support@serenilrs.com") {
			t.Errorf("template %q does not name the support address", name)
		}
		if !strings.Contains(tmpl.HTMLBody, `href="mailto:support@serenilrs.com"`) {
			t.Errorf("template %q support address is not a mailto link", name)
		}
		if !strings.Contains(tmpl.HTMLBody, "please contact our support team") {
			t.Errorf("template %q is missing the support sentence", name)
		}
		if !strings.Contains(tmpl.TextBody, "support@serenilrs.com") {
			t.Errorf("text body for %q does not name the support address", name)
		}
		for _, bad := range []string{"reply to this email", "Just reply", "feel free to reach out"} {
			if strings.Contains(tmpl.HTMLBody, bad) {
				t.Errorf("template %q still invites a reply (%q)", name, bad)
			}
		}
	}
}

func TestTemplatesHavePreheader(t *testing.T) {
	for _, name := range templateNames() {
		body := getBody(t, name)
		if !strings.Contains(body, "max-height:0") {
			t.Errorf("template %q has no hidden preheader block", name)
		}
	}
}
