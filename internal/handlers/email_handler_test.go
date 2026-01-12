package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"email-service/internal/email"

	"github.com/gin-gonic/gin"
)

type fakeService struct {
	sendErr error
}

func (f *fakeService) SendTransactionalEmail(req *email.EmailRequest) error { return f.sendErr }
func (f *fakeService) SendBulkEmail(recipients []string, subject, body string, isHTML bool) ([]string, error) {
	// treat invalid emails as failed
	var failed []string
	for _, r := range recipients {
		if !email.IsValidEmail(r) {
			failed = append(failed, r)
		}
	}
	return failed, nil
}
func (f *fakeService) GenerateAndSendOTP(to string, expiryMinutes int) (string, error) {
	return "123456", nil
}
func (f *fakeService) VerifyOTP(emailAddr, otp string) bool { return otp == "good" }

func setupRouter(s *fakeService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewEmailHandler(&email.EmailService{})
	// replace service with our fake via type assertion
	h.Service = &email.EmailService{}
	// But tests call handler methods directly with context, so we'll set Service to an email service wrapper
	// Instead, set the methods via embedding - easiest is to assign the actual interface methods via closures
	// For brevity, we'll construct handler and set Service to a minimal service and use function injection where needed
	return r
}

func TestSendEmail_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	// Provide invalid JSON body
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString("{bad json"))
	c.Request = req
	h.SendEmail(c)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", recorder.Code)
	}
}

func TestSendBulkEmail_WithInvalidEmails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	body := email.BulkEmailRequest{Recipients: []string{"a@b.com", "bad"}, Subject: "s", Body: "b"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.SendBulkEmail(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", recorder.Code)
	}
	// decode and assert response contains failed emails
	var resp map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if fe, ok := resp["failed_emails"]; !ok || fe == nil {
		t.Fatalf("expected failed_emails in response")
	}
}

func TestGenerateOTPAndVerify(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	// Generate OTP with valid JSON
	genBody := map[string]interface{}{"to": "a@b.com", "expiry": 1}
	b, _ := json.Marshal(genBody)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.GenerateOTP(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", recorder.Code)
	}

	// Now test SendEmail via handler - transactional
	recorder = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(recorder)
	tr := email.EmailRequest{To: []string{"a@b.com"}, Subject: "s", Body: "b"}
	tb, _ := json.Marshal(tr)
	req = httptest.NewRequest("POST", "/", bytes.NewBuffer(tb))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.SendEmail(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected transactional handler 200 got %d", recorder.Code)
	}

	// Create OTP via service and verify via handler
	otp, err := es.GenerateAndSendOTP("a@b.com", 1)
	if err != nil {
		t.Fatalf("generate otp failed: %v", err)
	}

	// Verify OTP
	verifier := map[string]string{"email": "a@b.com", "otp": otp}
	vb, _ := json.Marshal(verifier)
	recorder = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(recorder)
	req = httptest.NewRequest("POST", "/", bytes.NewBuffer(vb))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.VerifyOTP(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected verify 200 got %d", recorder.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	h := NewEmailHandler(&email.EmailService{})
	h.HealthCheck(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", recorder.Code)
	}
}

func TestSendEmail_TemplateCausesServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	// Make SendEmail return error so SendTransactionalEmail returns error
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return errors.New("send failed") }
	h := NewEmailHandler(es)

	// Provide valid JSON with subject/body so binding succeeds and SendTransactionalEmail calls SendEmail
	reqBody := map[string]interface{}{"to": []string{"a@b.com"}, "subject": "s", "body": "b"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.SendEmail(c)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", recorder.Code)
	}
}

func TestSendBulkEmail_BadJSONReturns400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	req := httptest.NewRequest("POST", "/", bytes.NewBufferString("{badjson"))
	c.Request = req
	h.SendBulkEmail(c)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", recorder.Code)
	}
}

func TestGenerateOTP_DefaultExpiry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	// Omit expiry to rely on default
	reqBody := map[string]interface{}{"to": "a@b.com"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.GenerateOTP(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", recorder.Code)
	}
}

func TestVerifyOTP_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	// Verify with a non-existent OTP
	reqBody := map[string]string{"email": "a@b.com", "otp": "wrong"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.VerifyOTP(c)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", recorder.Code)
	}
	// assert response message contains expected text
	var resp map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err == nil {
		if msg, ok := resp["message"]; !ok || msg == "" {
			t.Fatalf("expected error message in response")
		}
	}
}

func TestSendBulkEmail_AllSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	body := email.BulkEmailRequest{Recipients: []string{"a@b.com", "c@d.com"}, Subject: "s", Body: "b"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.SendBulkEmail(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", recorder.Code)
	}
	var resp email.EmailResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success true for all-success bulk")
	}
	if len(resp.FailedEmails) != 0 {
		t.Fatalf("expected no failed emails, got: %v", resp.FailedEmails)
	}
}

func TestSendEmail_SuccessResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendEmailFunc = func(to []string, subject, body string, isHTML bool) error { return nil }
	h := NewEmailHandler(es)

	tr := email.EmailRequest{To: []string{"a@b.com"}, Subject: "s", Body: "b"}
	tb, _ := json.Marshal(tr)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(tb))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.SendEmail(c)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", recorder.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if success, ok := resp["success"].(bool); !ok || !success {
		t.Fatalf("expected success true in response")
	}
}

func TestSendBulkEmail_ServiceErrorReturns500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	es := email.NewEmailService("h", 25, "u", "p", "from@x", 5)
	es.SendBulkEmailFunc = func(recipients []string, subject, body string, isHTML bool) ([]string, error) {
		return nil, errors.New("bulk fail")
	}
	h := NewEmailHandler(es)

	body := email.BulkEmailRequest{Recipients: []string{"a@b.com"}, Subject: "s", Body: "b"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	h.SendBulkEmail(c)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", recorder.Code)
	}
}
