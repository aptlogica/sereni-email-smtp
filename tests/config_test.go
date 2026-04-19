// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/config"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_CFG_KEY", "value1")
	defer os.Unsetenv("TEST_CFG_KEY")

	got := config.GetEnv("TEST_CFG_KEY", "default")
	if got != "value1" {
		t.Fatalf("expected value1, got %s", got)
	}

	// missing key returns default
	got2 := config.GetEnv("MISSING_KEY", "def")
	if got2 != "def" {
		t.Fatalf("expected def, got %s", got2)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")
	got := config.GetEnvAsInt("TEST_INT", 5)
	if got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}

	// invalid int uses default
	os.Setenv("TEST_INT_BAD", "bad")
	defer os.Unsetenv("TEST_INT_BAD")
	got2 := config.GetEnvAsInt("TEST_INT_BAD", 7)
	if got2 != 7 {
		t.Fatalf("expected 7, got %d", got2)
	}
}

func TestLoadConfig_LoadsDefaults(t *testing.T) {

	t.Setenv("HOST", "")
	t.Setenv("ALLOWED_ORIGIN", "")
	// Unset values that should be loaded from .env so godotenv.Load can override.
	os.Unsetenv("PORT")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	t.Setenv("SMTP_USERNAME", "testuser")
	t.Setenv("SMTP_PASSWORD", "testpass")
	t.Setenv("FROM_EMAIL", "test@example.com")
	t.Setenv("REDIS_URL", "")
	// Unset so .env can override.
	os.Unsetenv("BULK_BATCH_SIZE")

	// Use temp dir to control .env presence
	tmpDir, err := ioutil.TempDir("", "cfgtest")
	if err != nil {
		t.Fatalf("tmpdir failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	os.Chdir(tmpDir)

	// Without .env LoadConfig should use defaults (LoadConfig doesn't panic)
	cfg := config.LoadConfig()
	if cfg.Host != "0.0.0.0" || cfg.Port != "8082" || cfg.SMTPHost != "smtp.gmail.com" {
		t.Fatalf("defaults not working: %+v", cfg)
	}

	// Write a .env and expect it to override defaults
	env := "PORT=9090\nSMTP_HOST=example.com\nSMTP_PORT=2525\nBULK_BATCH_SIZE=20\n"
	if err := ioutil.WriteFile(filepath.Join(tmpDir, ".env"), []byte(env), 0644); err != nil {
		t.Fatalf("write .env failed: %v", err)
	}
	cfg = config.LoadConfig()
	if cfg.Port != "9090" || cfg.SMTPHost != "example.com" || cfg.SMTPPort != 2525 || cfg.BulkBatchSize != 20 {
		t.Fatalf("unexpected cfg values: %+v", cfg)
	}
}
