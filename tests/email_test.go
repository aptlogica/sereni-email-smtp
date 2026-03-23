// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/smtp"
	"testing"
	"time"

	"github.com/aptlogica/sereni-email-smtp/internal/email"
)

func TestGenerateOTPAndValidators(t *testing.T) {
	otp := email.GenerateOTP()
	if !email.ValidateOTPFormat(otp) {
		t.Fatalf("otp format invalid: %s", otp)
	}

	if !email.IsValidEmail("a@b.com") || email.IsValidEmail("bad") {
		t.Fatalf("email validator failed")
	}

	if !email.IsValidEmailList([]string{"x@y.com", "a@b.com"}) {
		t.Fatalf("email list validation failed")
	}
}

func TestTemplateRenderAndManagement(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	subj, body, err := es.RenderTemplate("otp_template", map[string]interface{}{"otp": "123456", "expiry": 10})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if subj == "" || body == "" {
		t.Fatalf("expected non-empty render")
	}

	et := email.EmailTemplate{Subject: "S", HTMLBody: "<b>{{.who}}</b>", TextBody: "{{.who}}"}
	es.AddTemplate("tst", et)
	got, ok := es.GetTemplate("tst")
	if !ok || got.Subject != et.Subject {
		t.Fatalf("template add/get failed")
	}
	avail := es.GetAvailableTemplates()
	if len(avail) == 0 {
		t.Fatalf("expected templates available")
	}
}

func TestSendTemplateEmailAndTransactional(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	var lastTo []string
	var lastSubject string
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		lastTo = to
		lastSubject = subject
		return nil
	}

	if err := es.SendTemplateEmail([]string{"a@b.com"}, "otp_template", map[string]interface{}{"otp": "111111", "expiry": 5}); err != nil {
		t.Fatalf("send template failed: %v", err)
	}
	if lastSubject == "" || lastTo == nil {
		t.Fatalf("expected send invoked")
	}

	// Transactional with template but invalid recipient
	req := &email.EmailRequest{To: []string{"bad"}, Template: "otp_template"}
	if err := es.SendTransactionalEmail(req); err == nil {
		t.Fatalf("expected transactional error for invalid recipient")
	}
}

func TestOTPStoreAndCleanup(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	otp, err := es.GenerateAndSendOTP("u@x", 0) // expiry 0 -> immediate expiry
	if err != nil {
		t.Fatalf("generate otp failed: %v", err)
	}
	if !email.ValidateOTPFormat(otp) {
		t.Fatalf("invalid otp format")
	}

	// Verify should fail due to expiry
	time.Sleep(5 * time.Millisecond)
	if es.VerifyOTP("u@x", otp) {
		t.Fatalf("expected expired otp to not verify")
	}
}

// mockClient implements email.SmtpClient for testing SendEmail via Dial injection
type mockClient struct {
	startTLSErr error
	authErr     error
	mailErr     error
	rcptErr     error
	dataErr     error
	writeErr    error
	closeErr    error
	quitErr     error
	buf         bytes.Buffer
}

func (m *mockClient) StartTLS(config *tls.Config) error { return m.startTLSErr }
func (m *mockClient) Auth(a smtp.Auth) error            { return m.authErr }
func (m *mockClient) Mail(from string) error            { return m.mailErr }
func (m *mockClient) Rcpt(to string) error              { return m.rcptErr }
func (m *mockClient) Data() (io.WriteCloser, error) {
	if m.dataErr != nil {
		return nil, m.dataErr
	}
	return nopWriteCloser{b: &m.buf, writeErr: m.writeErr, closeErr: m.closeErr}, nil
}
func (m *mockClient) Close() error { return m.closeErr }
func (m *mockClient) Quit() error  { return m.quitErr }

type nopWriteCloser struct {
	b        *bytes.Buffer
	writeErr error
	closeErr error
}

func (n nopWriteCloser) Write(p []byte) (int, error) {
	if n.writeErr != nil {
		return 0, n.writeErr
	}
	return n.b.Write(p)
}
func (n nopWriteCloser) Close() error { return n.closeErr }

func TestSendEmail_DialAndClientErrors(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)

	// Dial error
	es.Dial = func(addr string) (email.SmtpClient, error) { return nil, errors.New("dial fail") }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected dial error")
	}

	// StartTLS error
	mock := &mockClient{startTLSErr: errors.New("tls fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected tls error")
	}

	// Auth error
	mock = &mockClient{authErr: errors.New("auth fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected auth error")
	}

	// Mail error
	mock = &mockClient{mailErr: errors.New("mail fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected mail error")
	}

	// Rcpt error
	mock = &mockClient{rcptErr: errors.New("rcpt fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected rcpt error")
	}

	// Data write error
	mock = &mockClient{writeErr: errors.New("write fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected write error")
	}

	// Data() error
	mock = &mockClient{dataErr: errors.New("data fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected data error")
	}

	// writer Close error
	mock = &mockClient{closeErr: errors.New("close fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected writer close error")
	}

	// Quit error
	mock = &mockClient{quitErr: errors.New("quit fail")}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err == nil {
		t.Fatalf("expected quit error")
	}

	// success
	mock = &mockClient{}
	es.Dial = func(addr string) (email.SmtpClient, error) { return mock, nil }
	es.FromEmail = "from@x"
	if err := es.SendEmail([]string{"a@b.com"}, "s", "body", true); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
}

func TestRenderTemplate_NotFoundAndSendTemplateError(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	// non-existent template
	if _, _, err := es.RenderTemplate("no_such", nil); err == nil {
		t.Fatalf("expected error for missing template")
	}

	// ensure SendTemplateEmail surfaces RenderTemplate errors
	if err := es.SendTemplateEmail([]string{"a@b.com"}, "no_such", nil); err == nil {
		t.Fatalf("expected send template to fail for missing template")
	}
}

func TestSendBulkEmail_FailuresAggregated(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 10) // Use higher batch size to avoid race
	// Fail send for a specific address via SendEmailFunc
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		if len(to) > 0 && to[0] == "fail@x" {
			return errors.New("send failed")
		}
		return nil
	}

	recipients := []string{"ok@x", "fail@x", "bad"}
	failed, _ := es.SendBulkEmail(recipients, "s", "b", false)
	// bad is invalid and fail@x should be reported
	foundFail := false
	foundBad := false
	for _, f := range failed {
		if f == "fail@x" {
			foundFail = true
		}
		if f == "bad" {
			foundBad = true
		}
	}
	if !foundFail || !foundBad {
		t.Fatalf("expected aggregated failures, got: %v", failed)
	}
}

func TestSendTransactionalEmail_SuccessAndTemplateFlow(t *testing.T) {
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	called := false
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		called = true
		return nil
	}

	req := &email.EmailRequest{To: []string{"a@b.com"}, Subject: "s", Body: "b"}
	if err := es.SendTransactionalEmail(req); err != nil {
		t.Fatalf("expected transactional send success: %v", err)
	}
	if !called {
		t.Fatalf("expected SendEmailFunc to be called")
	}
}

func TestJoinAndAvailableTemplatesAndPlainSend(t *testing.T) {
	if email.Join([]string{"a", "b"}, ",") != "a,b" {
		t.Fatalf("join failed")
	}
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	// non-HTML send path via SendEmailFunc
	called := false
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error {
		if isHTML {
			t.Fatalf("expected plain text send")
		}
		called = true
		return nil
	}
	if err := es.SendEmail([]string{"a@b.com"}, "s", "b", false); err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if !called {
		t.Fatalf("expected SendEmailFunc called")
	}

	if len(es.GetAvailableTemplates()) == 0 {
		t.Fatalf("expected available templates")
	}
}
