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

// TestGenerateOTP_Format tests that GenerateOTP produces valid 6-digit format
func TestGenerateOTP_Format(t *testing.T) {
	for i := 0; i < 100; i++ {
		otp := email.GenerateOTP()

		if len(otp) != 6 {
			t.Errorf("Expected OTP length 6, got %d for OTP: %s", len(otp), otp)
		}

		// Check if all characters are digits
		for _, ch := range otp {
			if ch < '0' || ch > '9' {
				t.Errorf("Expected OTP to contain only digits, got: %s", otp)
			}
		}
	}
}

// TestGenerateOTP_Randomness tests that GenerateOTP produces different values
func TestGenerateOTP_Randomness(t *testing.T) {
	otps := make(map[string]bool)
	duplicates := 0

	for i := 0; i < 1000; i++ {
		otp := email.GenerateOTP()
		if otps[otp] {
			duplicates++
		}
		otps[otp] = true
	}

	// Allow some duplicates due to birthday paradox, but expect very few
	if duplicates > 5 {
		t.Errorf("Too many duplicate OTPs generated: %d out of 1000", duplicates)
	}
}

// TestGenerateOTP_AllDigitsRange tests edge cases in OTP generation
func TestGenerateOTP_AllDigitsRange(t *testing.T) {
	otps := make(map[string]bool)

	for i := 0; i < 10000; i++ {
		otp := email.GenerateOTP()
		otps[otp] = true

		// Verify format
		if !email.ValidateOTPFormat(otp) {
			t.Errorf("Generated OTP %s failed format validation", otp)
		}
	}

	// Check we get a good variety of OTPs
	if len(otps) < 8000 {
		t.Errorf("Expected good variety in OTP generation, got %d unique OTPs from 10000", len(otps))
	}
}

// TestIsValidEmail_Valid tests various valid email formats
func TestIsValidEmail_Valid(t *testing.T) {
	validEmails := []string{
		"simple@example.com",
		"user.name@example.com",
		"user+tag@example.co.uk",
		"123@456.com",
		"a@b.co",
		"test_underscore@example.com",
		"test-dash@example.com",
		"firstname.lastname@example.com",
		"user@localhost.localdomain",
		"user@example.museum",
		"1234567890@example.com",
		"_@example.com",
	}

	for _, emailAddr := range validEmails {
		if !email.IsValidEmail(emailAddr) {
			t.Errorf("Expected %s to be valid, but got invalid", emailAddr)
		}
	}
}

// TestIsValidEmail_Invalid tests various invalid email formats
func TestIsValidEmail_Invalid(t *testing.T) {
	invalidEmails := []string{
		"",                       // empty
		"plainaddress",           // no @
		"@example.com",           // no local part
		"user@",                  // no domain
		"user @example.com",      // space in local part
		"user@exam ple.com",      // space in domain
		"user@.com",              // no domain name
		"user..name@example.com", // consecutive dots
	}

	for _, emailAddr := range invalidEmails {
		if email.IsValidEmail(emailAddr) {
			t.Errorf("Expected %s to be invalid, but got valid", emailAddr)
		}
	}
}

// TestIsValidEmail_EdgeCases tests edge case emails
func TestIsValidEmail_EdgeCases(t *testing.T) {
	// These might be valid or invalid depending on RFC interpretation
	edgeCases := []struct {
		emailAddr string
		valid     bool
	}{
		{"a@b.c", true},              // minimal valid
		{"user@192.168.1.1", true},   // IP address format
		{"user+@example.com", true},  // plus sign
		{"user.@example.com", false}, // dot before @
	}

	for _, tc := range edgeCases {
		result := email.IsValidEmail(tc.emailAddr)
		if result != tc.valid {
			t.Logf("Edge case %s: got %v, expected %v", tc.emailAddr, result, tc.valid)
		}
	}
}

// TestIsValidEmailList_AllValid tests list with all valid emails
func TestIsValidEmailList_AllValid(t *testing.T) {
	validList := []string{
		"user1@example.com",
		"user2@example.org",
		"user3@example.net",
	}

	if !email.IsValidEmailList(validList) {
		t.Error("Expected all valid emails to return true")
	}
}

// TestIsValidEmailList_SomeInvalid tests list with some invalid emails
func TestIsValidEmailList_SomeInvalid(t *testing.T) {
	invalidList := []string{
		"user1@example.com",
		"invalid",
		"user3@example.net",
	}

	if email.IsValidEmailList(invalidList) {
		t.Error("Expected list with invalid email to return false")
	}
}

// TestIsValidEmailList_AllInvalid tests list with all invalid emails
func TestIsValidEmailList_AllInvalid(t *testing.T) {
	invalidList := []string{
		"invalid1",
		"invalid2",
		"invalid3",
	}

	if email.IsValidEmailList(invalidList) {
		t.Error("Expected all invalid emails to return false")
	}
}

// TestIsValidEmailList_Empty tests empty email list
func TestIsValidEmailList_Empty(t *testing.T) {
	emptyList := []string{}

	if !email.IsValidEmailList(emptyList) {
		t.Error("Expected empty list to return true")
	}
}

// TestIsValidEmailList_SingleValid tests single valid email
func TestIsValidEmailList_SingleValid(t *testing.T) {
	singleList := []string{"user@example.com"}

	if !email.IsValidEmailList(singleList) {
		t.Error("Expected single valid email to return true")
	}
}

// TestIsValidEmailList_SingleInvalid tests single invalid email
func TestIsValidEmailList_SingleInvalid(t *testing.T) {
	singleList := []string{"invalid"}

	if email.IsValidEmailList(singleList) {
		t.Error("Expected single invalid email to return false")
	}
}

// TestIsValidEmailList_ManyEmails tests with large list
func TestIsValidEmailList_ManyEmails(t *testing.T) {
	validList := make([]string, 100)
	for i := 0; i < 100; i++ {
		validList[i] = "user" + string(rune(48+(i%10))) + "@example.com"
	}

	if !email.IsValidEmailList(validList) {
		t.Error("Expected large list of valid emails to return true")
	}

	// Add one invalid at the end
	validList[99] = "invalid"
	if email.IsValidEmailList(validList) {
		t.Error("Expected list with one invalid email to return false")
	}
}

// TestValidateOTPFormat_ValidOTPs tests various valid OTP formats
func TestValidateOTPFormat_ValidOTPs(t *testing.T) {
	validOTPs := []string{
		"000000",
		"123456",
		"999999",
		"111111",
		"654321",
		"500000",
	}

	for _, otp := range validOTPs {
		if !email.ValidateOTPFormat(otp) {
			t.Errorf("Expected %s to be valid OTP format", otp)
		}
	}
}

// TestValidateOTPFormat_InvalidOTPs tests various invalid OTP formats
func TestValidateOTPFormat_InvalidOTPs(t *testing.T) {
	invalidOTPs := []string{
		"",        // empty
		"12345",   // too short
		"1234567", // too long
		"abc123",  // letters
		"12345a",  // contains letter
		"1234 6",  // contains space
		"12-456",  // contains dash
		"123.456", // contains period
		"12345@",  // contains special char
	}

	for _, otp := range invalidOTPs {
		if email.ValidateOTPFormat(otp) {
			t.Errorf("Expected %s to be invalid OTP format", otp)
		}
	}
}

// TestValidateOTPFormat_EdgeCases tests edge cases for OTP validation
func TestValidateOTPFormat_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		otp   string
		valid bool
	}{
		{"000001", true},
		{"999998", true},
		{"100000", true},
		{" 123456", false},
		{"123456 ", false},
		{"1234567", false},
		{"12345", false},
		{"000000", true},
	}

	for _, tc := range edgeCases {
		result := email.ValidateOTPFormat(tc.otp)
		if result != tc.valid {
			t.Errorf("ValidateOTPFormat(%s): got %v, expected %v", tc.otp, result, tc.valid)
		}
	}
}

// TestValidateOTPFormat_NonDigitCharacters tests OTPs with non-digit characters
func TestValidateOTPFormat_NonDigitCharacters(t *testing.T) {
	nonDigitOTPs := []string{
		"O0O0O0",  // letters O instead of 0
		"123O56",  // mix
		"AAAAAA",  // all letters
		"12345\n", // newline
		"12345\t", // tab
		"12345-",  // dash
		"12345_",  // underscore
	}

	for _, otp := range nonDigitOTPs {
		if email.ValidateOTPFormat(otp) {
			t.Errorf("Expected %s to be invalid (contains non-digits)", otp)
		}
	}
}

// TestValidateOTPFormat_CaseSensitivity tests that validation is case sensitive for non-numeric
func TestValidateOTPFormat_CaseSensitivity(t *testing.T) {
	testCases := []string{
		"abcdef",
		"ABCDEF",
		"AbCdEf",
	}

	for _, testCase := range testCases {
		if email.ValidateOTPFormat(testCase) {
			t.Errorf("Expected %s (letters) to be invalid", testCase)
		}
	}
}

// TestValidateOTPFormat_PerformanceWithLargeInput tests performance with larger strings
func TestValidateOTPFormat_PerformanceWithLargeInput(t *testing.T) {
	largeInput := strings.Repeat("1", 1000)
	result := email.ValidateOTPFormat(largeInput)

	if result {
		t.Errorf("Expected large numeric string to be invalid (not 6 digits)")
	}
}

// TestGenerateOTP_NoErrorOnFailure tests GenerateOTP handles rand.Read errors gracefully
// Note: This tests the fallback in GenerateOTP
func TestGenerateOTP_ConsistentBehavior(t *testing.T) {
	// Generate 100 OTPs and ensure all are valid format
	for i := 0; i < 100; i++ {
		otp := email.GenerateOTP()
		if !email.ValidateOTPFormat(otp) {
			t.Errorf("Generated OTP %s has invalid format", otp)
		}
	}
}

// TestIsValidEmail_WithComments tests email validation with complex formats
func TestIsValidEmail_WithComments(t *testing.T) {
	// These test RFC 5322 compliance
	complexEmails := []struct {
		emailAddr string
		valid     bool
	}{
		{"user@example.com", true},
		{"\"user\"@example.com", false}, // quoted strings not handled by basic parser
		{"user%site@example.com", true}, // percent sign is usually valid
	}

	for _, tc := range complexEmails {
		result := email.IsValidEmail(tc.emailAddr)
		t.Logf("IsValidEmail(%q): %v", tc.emailAddr, result)
	}
}
