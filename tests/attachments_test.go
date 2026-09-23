// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

// sendWithMockSMTP wires up a service whose SMTP DATA payload is captured in
// the returned MockWriter, so the actual message bytes written to the wire
// (not just the sanitized inputs SendEmailFunc would see) can be inspected.
func sendWithMockSMTP(t *testing.T, to []string, subject, body string, isHTML bool, attachments []email.Attachment) (*MockWriter, error) {
	t.Helper()
	service := email.NewEmailService("localhost", 587, "user", "pass", "from@test.com", 5)
	mockWriter := &MockWriter{}
	mockClient := &MockSmtpClient{DataWriter: mockWriter}
	service.Dial = func(addr string) (email.SmtpClient, error) {
		return mockClient, nil
	}

	err := service.SendEmail(to, subject, body, isHTML, attachments)
	return mockWriter, err
}

func TestSendEmail_WithAttachment_BuildsMultipartMessage(t *testing.T) {
	content := []byte("id,name\n1,Ada\n")
	attachment := email.Attachment{
		Filename:      "report.csv",
		ContentType:   "text/csv",
		ContentBase64: base64.StdEncoding.EncodeToString(content),
	}

	mockWriter, err := sendWithMockSMTP(t, []string{"test@example.com"}, "Subject", "Body", false, []email.Attachment{attachment})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	message := string(mockWriter.written)

	if !strings.Contains(message, "Content-Type: multipart/mixed; boundary=") {
		t.Errorf("expected a multipart/mixed envelope, got:\n%s", message)
	}
	if !strings.Contains(message, `Content-Disposition: attachment; filename="report.csv"`) {
		t.Errorf("expected a Content-Disposition header naming the attachment, got:\n%s", message)
	}
	if !strings.Contains(message, "Content-Type: text/csv") {
		t.Errorf("expected the attachment's own Content-Type in its part, got:\n%s", message)
	}
	if !strings.Contains(message, "Content-Transfer-Encoding: base64") {
		t.Errorf("expected base64 transfer encoding on the attachment part, got:\n%s", message)
	}

	// The encoded content must round-trip: strip CRLFs from the wrapped
	// lines and decode back to the original bytes.
	decodable := extractBase64Payload(message, "report.csv")
	decoded, err := base64.StdEncoding.DecodeString(decodable)
	if err != nil {
		t.Fatalf("attachment payload did not decode as base64: %v", err)
	}
	if !bytes.Equal(decoded, content) {
		t.Errorf("decoded attachment content = %q, want %q", decoded, content)
	}
}

func TestSendEmail_WithoutAttachments_KeepsSimpleMessageFormat(t *testing.T) {
	mockWriter, err := sendWithMockSMTP(t, []string{"test@example.com"}, "Subject", "Body", false, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	message := string(mockWriter.written)
	if strings.Contains(message, "multipart") {
		t.Errorf("expected the original flat message format with no attachments, got:\n%s", message)
	}
	if !strings.HasPrefix(message, "From: from@test.com\r\nTo: test@example.com\r\nSubject: Subject\r\n\r\nBody") {
		t.Errorf("unexpected message format:\n%s", message)
	}
}

func TestSendEmail_AttachmentWithInvalidBase64_ReturnsError(t *testing.T) {
	attachment := email.Attachment{
		Filename:      "broken.txt",
		ContentType:   "text/plain",
		ContentBase64: "not-valid-base64!!",
	}

	_, err := sendWithMockSMTP(t, []string{"test@example.com"}, "Subject", "Body", false, []email.Attachment{attachment})
	if err == nil {
		t.Fatal("expected an error for invalid attachment content, got nil")
	}
}

func TestSendEmail_AttachmentFilenameWithControlCharacters_IsSanitized(t *testing.T) {
	malicious := email.Attachment{
		Filename:      "evil.txt\r\nBcc: attacker@evil.com",
		ContentType:   "text/plain",
		ContentBase64: base64.StdEncoding.EncodeToString([]byte("hi")),
	}

	mockWriter, err := sendWithMockSMTP(t, []string{"test@example.com"}, "Subject", "Body", false, []email.Attachment{malicious})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	message := string(mockWriter.written)
	if strings.Contains(message, "\r\nBcc:") {
		t.Errorf("expected the injected header to be stripped from the filename, got:\n%s", message)
	}
}

// extractBase64Payload pulls the base64 lines belonging to the named
// attachment's part out of a raw multipart message, joining them back into
// one unwrapped string.
func extractBase64Payload(message, filename string) string {
	marker := `filename="` + filename + `"`
	idx := strings.Index(message, marker)
	if idx == -1 {
		return ""
	}
	rest := message[idx:]
	// The payload starts after the blank line following the part's headers.
	headerEnd := strings.Index(rest, "\r\n\r\n")
	if headerEnd == -1 {
		return ""
	}
	body := rest[headerEnd+4:]
	// The payload ends at the next boundary line ("--...").
	boundaryIdx := strings.Index(body, "\r\n--")
	if boundaryIdx != -1 {
		body = body[:boundaryIdx]
	}
	return strings.ReplaceAll(strings.ReplaceAll(body, "\r\n", ""), "\n", "")
}
