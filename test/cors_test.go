package test

import (
	"net/http/httptest"
	"testing"

	"email-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_OptionsAndGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.OPTIONS("/", func(c *gin.Context) { c.Status(200) })
	r.GET("/", func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
