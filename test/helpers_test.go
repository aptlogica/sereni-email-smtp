package test

import (
	"testing"

	"email-service/internal/templatecache"
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
