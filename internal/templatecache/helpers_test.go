package templatecache

import (
	"testing"
)

func TestParseAndCacheTemplate_SuccessAndFail(t *testing.T) {
	tc := NewTemplateCache()
	tmpl, err := ParseAndCacheTemplate(tc, "ok", "hello {{.}}")
	if err != nil || tmpl == nil {
		t.Fatalf("expected parse success: %v", err)
	}
	if _, ok := tc.Get("ok"); !ok {
		t.Fatalf("expected cached template")
	}

	// invalid template
	if _, err := ParseAndCacheTemplate(tc, "bad", "{{%"); err == nil {
		t.Fatalf("expected parse error for invalid template")
	}
}
