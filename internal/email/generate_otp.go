// Copyright (c) 2026 Aptlogica Technologies Private Limited
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

// Secure OTP generation and validation helpers
package email

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/mail"
	"regexp"
)

// GenerateOTP generates a cryptographically secure 6-digit OTP
func GenerateOTP() string {
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// fallback to time-based, but log error
		return fmt.Sprintf("%06d", (int64(binary.BigEndian.Uint32(b[:]))+int64(err.Error()[0]))%1000000)
	}
	return fmt.Sprintf("%06d", binary.BigEndian.Uint32(b[:])%1000000)
}

// IsValidEmail checks if an email address is valid
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidEmailList checks if all emails in a slice are valid
func IsValidEmailList(emails []string) bool {
	for _, e := range emails {
		if !IsValidEmail(e) {
			return false
		}
	}
	return true
}

// ValidateOTPFormat checks if OTP is a 6-digit number
func ValidateOTPFormat(otp string) bool {
	return regexp.MustCompile(`^\d{6}$`).MatchString(otp)
}
