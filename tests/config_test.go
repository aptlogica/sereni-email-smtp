// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"email-service/internal/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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

func TestLoadConfig_PanicsAndLoads(t *testing.T) {
	// Use temp dir to control .env presence
	tmpDir, err := ioutil.TempDir("", "cfgtest")
	if err != nil {
		t.Fatalf("tmpdir failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldwd, _ := os.Getwd()
	defer os.Chdir(oldwd)
	os.Chdir(tmpDir)

	// Without .env LoadConfig should panic (godotenv.Load returns error)
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		_ = config.LoadConfig()
	}()
	if !panicked {
		t.Fatalf("expected panic when .env missing")
	}

	// Write a .env and expect success
	env := "PORT=9090\nSMTP_HOST=example.com\nSMTP_PORT=2525\nBULK_BATCH_SIZE=20\n"
	if err := ioutil.WriteFile(filepath.Join(tmpDir, ".env"), []byte(env), 0644); err != nil {
		t.Fatalf("write .env failed: %v", err)
	}
	cfg := config.LoadConfig()
	if cfg.Port != "9090" || cfg.SMTPHost != "example.com" || cfg.SMTPPort != 2525 || cfg.BulkBatchSize != 20 {
		t.Fatalf("unexpected cfg values: %+v", cfg)
	}
}
