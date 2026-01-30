package test

import (
	"html/template"
	"testing"

	"email-service/internal/templatecache"
)

func TestTemplateCache_GetSetPreload(t *testing.T) {
	tc := templatecache.NewTemplateCache()

	// Get on missing should be false
	if _, ok := tc.Get("missing"); ok {
		t.Fatalf("expected missing")
	}

	// Set and Get
	tmpl := template.Must(template.New("a").Parse("hello {{.}}"))
	tc.Set("a", tmpl)
	got, ok := tc.Get("a")
	if !ok || got == nil {
		t.Fatalf("expected template to be present")
	}

	// Preload valid templates
	err := tc.Preload(map[string]string{"b": "hi {{.}}"})
	if err != nil {
		t.Fatalf("unexpected preload error: %v", err)
	}

	if _, ok := tc.Get("b"); !ok {
		t.Fatalf("expected preloaded template")
	}

	// Preload with an invalid template should return error
	err = tc.Preload(map[string]string{"bad": "{{%"})
	if err == nil {
		t.Fatalf("expected preload to fail for invalid template")
	}
}
