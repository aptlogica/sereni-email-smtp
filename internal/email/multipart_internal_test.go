// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package email

import (
	"encoding/base64"
	"errors"
	"mime/multipart"
	"strings"
	"testing"
)

// failAfterWriter succeeds its first n Write calls, then fails every call
// after that with err. Used to reach the mw.CreatePart/Write/Close error
// branches that a real strings.Builder (which never fails) can't trigger.
type failAfterWriter struct {
	n   int
	err error
}

func (f *failAfterWriter) Write(p []byte) (int, error) {
	if f.n <= 0 {
		return 0, f.err
	}
	f.n--
	return len(p), nil
}

var errWriteFailed = errors.New("write failed")

func TestWriteMultipartBody_CreatePartError(t *testing.T) {
	_, err := writeMultipartBody(&failAfterWriter{n: 0, err: errWriteFailed}, false, "body", nil)
	if err == nil || !strings.Contains(err.Error(), "create message body part") {
		t.Fatalf("expected create message body part error, got %v", err)
	}
}

func TestWriteMultipartBody_BodyWriteError(t *testing.T) {
	_, err := writeMultipartBody(&failAfterWriter{n: 1, err: errWriteFailed}, false, "body", nil)
	if err == nil || !strings.Contains(err.Error(), "write message body part") {
		t.Fatalf("expected write message body part error, got %v", err)
	}
}

func TestWriteMultipartBody_CloseError(t *testing.T) {
	_, err := writeMultipartBody(&failAfterWriter{n: 2, err: errWriteFailed}, false, "body", nil)
	if err == nil || !strings.Contains(err.Error(), "close multipart message") {
		t.Fatalf("expected close multipart message error, got %v", err)
	}
}

func TestWriteAttachmentPart_CreatePartError(t *testing.T) {
	mw := multipart.NewWriter(&failAfterWriter{n: 0, err: errWriteFailed})
	att := Attachment{Filename: "report.csv", ContentBase64: base64.StdEncoding.EncodeToString([]byte("data"))}

	err := writeAttachmentPart(mw, att)
	if err == nil || !strings.Contains(err.Error(), "create part") {
		t.Fatalf("expected create part error, got %v", err)
	}
}

func TestWriteAttachmentPart_WriteContentError(t *testing.T) {
	mw := multipart.NewWriter(&failAfterWriter{n: 1, err: errWriteFailed})
	att := Attachment{Filename: "report.csv", ContentBase64: base64.StdEncoding.EncodeToString([]byte("data"))}

	err := writeAttachmentPart(mw, att)
	if err == nil || !strings.Contains(err.Error(), "write content") {
		t.Fatalf("expected write content error, got %v", err)
	}
}

func TestWriteAttachmentPart_DefaultsContentType(t *testing.T) {
	var buf strings.Builder
	mw := multipart.NewWriter(&buf)
	att := Attachment{Filename: "report.csv", ContentBase64: base64.StdEncoding.EncodeToString([]byte("data"))}

	if err := writeAttachmentPart(mw, att); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mw.Close()

	if !strings.Contains(buf.String(), "Content-Type: application/octet-stream") {
		t.Errorf("expected default content type application/octet-stream, got:\n%s", buf.String())
	}
}
