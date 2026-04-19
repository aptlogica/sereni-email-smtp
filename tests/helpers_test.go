// Copyright (c) 2026 Aptlogica Technologies Private Limited
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"testing"

	"github.com/aptlogica/sereni-email-smtp/internal/templatecache"
)

func TestParseAndCacheTemplate_SuccessAndFail(t *testing.T) {
	tc := templatecache.NewTemplateCache()
	tmpl, err := templatecache.ParseAndCacheTemplate(tc, "ok", "hello {{.}}")
	if err != nil || tmpl == nil {
		t.Fatalf("expected parse success: %v", err)
	}
	if _, ok := tc.Get("ok"); !ok {
		t.Fatalf("expected cached template")
	}

	// invalid template
	if _, err := templatecache.ParseAndCacheTemplate(tc, "bad", "{{%"); err == nil {
		t.Fatalf("expected parse error for invalid template")
	}
}
