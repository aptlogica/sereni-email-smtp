// Copyright (c) 2026 Aptlogica Technologies Private Limited
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"github.com/aptlogica/sereni-email-smtp/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_Comprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test with specific origin
	r := gin.New()
	r.Use(middleware.CORSMiddleware("https://example.com"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Test GET request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin https://example.com, got %s", headers.Get("Access-Control-Allow-Origin"))
	}
	if headers.Get("Access-Control-Allow-Credentials") != "true" {
		t.Errorf("Expected Access-Control-Allow-Credentials true, got %s", headers.Get("Access-Control-Allow-Credentials"))
	}
	if headers.Get("Access-Control-Allow-Headers") == "" {
		t.Error("Expected Access-Control-Allow-Headers to be set")
	}
	if headers.Get("Access-Control-Allow-Methods") != "POST, OPTIONS, GET, PUT" {
		t.Errorf("Expected specific methods, got %s", headers.Get("Access-Control-Allow-Methods"))
	}

	// Test OPTIONS request
	w = httptest.NewRecorder()
	req = httptest.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected 204 for OPTIONS, got %d", w.Code)
	}

	// Test with wildcard origin
	r2 := gin.New()
	r2.Use(middleware.CORSMiddleware("*"))
	r2.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/test", nil)
	r2.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected wildcard origin, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	// Test POST request passes through (but will get 404 since we only have GET route)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)

	// CORS headers should still be set even on 404
	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected CORS headers on POST, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	// Test PUT request passes through (but will get 404 since we only have GET route)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/test", nil)
	r.ServeHTTP(w, req)

	// CORS headers should still be set even on 404
	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected CORS headers on PUT, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}
}
