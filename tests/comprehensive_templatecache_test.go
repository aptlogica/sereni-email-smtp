// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package test

import (
	"email-service/internal/templatecache"
	"html/template"
	"testing"
)

func TestNewTemplateCache_Comprehensive(t *testing.T) {
	cache := templatecache.NewTemplateCache()
	if cache == nil {
		t.Error("Expected cache to be created")
	}

	// Cache should be empty initially (except for Get test)
	if _, ok := cache.Get("non-existent"); ok {
		t.Error("Expected cache to be empty initially")
	}
}

func TestTemplateCache_Get_Comprehensive(t *testing.T) {
	cache := templatecache.NewTemplateCache()

	// Test getting non-existent template
	tmpl, ok := cache.Get("missing")
	if ok || tmpl != nil {
		t.Error("Expected false and nil for missing template")
	}

	// Add template and get it
	testTemplate := template.Must(template.New("test").Parse("Hello {{.name}}"))
	cache.Set("test", testTemplate)

	retrieved, ok := cache.Get("test")
	if !ok || retrieved == nil {
		t.Error("Expected template to be found after setting")
	}
	if retrieved != testTemplate {
		t.Error("Expected retrieved template to match set template")
	}
}

func TestTemplateCache_Set_Comprehensive(t *testing.T) {
	cache := templatecache.NewTemplateCache()

	// Test setting template
	tmpl1 := template.Must(template.New("first").Parse("First {{.}}"))
	cache.Set("first", tmpl1)

	retrieved, ok := cache.Get("first")
	if !ok || retrieved != tmpl1 {
		t.Error("Expected set template to be retrievable")
	}

	// Test overwriting template
	tmpl2 := template.Must(template.New("first").Parse("Second {{.}}"))
	cache.Set("first", tmpl2)

	retrieved, ok = cache.Get("first")
	if !ok || retrieved != tmpl2 {
		t.Error("Expected overwritten template to be retrievable")
	}
	if retrieved == tmpl1 {
		t.Error("Expected old template to be replaced")
	}

	// Test setting multiple templates
	tmpl3 := template.Must(template.New("third").Parse("Third {{.}}"))
	cache.Set("third", tmpl3)

	// Both should exist
	if _, ok := cache.Get("first"); !ok {
		t.Error("Expected first template to still exist")
	}
	if _, ok := cache.Get("third"); !ok {
		t.Error("Expected third template to exist")
	}
}

func TestTemplateCache_Preload_Comprehensive(t *testing.T) {
	cache := templatecache.NewTemplateCache()

	// Test successful preload
	templates := map[string]string{
		"template1": "Hello {{.name}}",
		"template2": "Hi {{.user}}, your code is {{.code}}",
		"template3": "Simple text without variables",
	}

	err := cache.Preload(templates)
	if err != nil {
		t.Errorf("Expected no error during preload, got %v", err)
	}

	// Verify all templates were loaded
	for name := range templates {
		if _, ok := cache.Get(name); !ok {
			t.Errorf("Expected template %s to be preloaded", name)
		}
	}

	// Test preload with invalid template syntax
	invalidTemplates := map[string]string{
		"valid":   "Hello {{.name}}",
		"invalid": "Bad syntax {{%",
	}

	err = cache.Preload(invalidTemplates)
	if err == nil {
		t.Error("Expected error when preloading invalid template")
	}

	// Test empty preload
	emptyTemplates := map[string]string{}
	err = cache.Preload(emptyTemplates)
	if err != nil {
		t.Errorf("Expected no error for empty preload, got %v", err)
	}
}

func TestParseAndCacheTemplate_Comprehensive(t *testing.T) {
	cache := templatecache.NewTemplateCache()

	// Test successful parse and cache
	tmpl, err := templatecache.ParseAndCacheTemplate(cache, "success", "Hello {{.world}}")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tmpl == nil {
		t.Error("Expected template to be returned")
	}

	// Verify template was cached
	cached, ok := cache.Get("success")
	if !ok {
		t.Error("Expected template to be cached")
	}
	if cached != tmpl {
		t.Error("Expected cached template to match returned template")
	}

	// Test with invalid template syntax
	_, err = templatecache.ParseAndCacheTemplate(cache, "invalid", "Bad {{%")
	if err == nil {
		t.Error("Expected error for invalid template syntax")
	}

	// Verify invalid template was not cached
	if _, ok := cache.Get("invalid"); ok {
		t.Error("Expected invalid template not to be cached")
	}

	// Test overwriting existing template
	tmpl2, err := templatecache.ParseAndCacheTemplate(cache, "success", "New content {{.item}}")
	if err != nil {
		t.Errorf("Expected no error overwriting template, got %v", err)
	}

	cached, ok = cache.Get("success")
	if !ok || cached != tmpl2 {
		t.Error("Expected template to be overwritten in cache")
	}
	if cached == tmpl {
		t.Error("Expected old template to be replaced")
	}

	// Test with empty template content
	tmpl3, err := templatecache.ParseAndCacheTemplate(cache, "empty", "")
	if err != nil {
		t.Errorf("Expected no error for empty template, got %v", err)
	}
	if tmpl3 == nil {
		t.Error("Expected empty template to be valid")
	}

	// Test with template containing only text (no variables)
	tmpl4, err := templatecache.ParseAndCacheTemplate(cache, "plain", "Just plain text")
	if err != nil {
		t.Errorf("Expected no error for plain text template, got %v", err)
	}
	if tmpl4 == nil {
		t.Error("Expected plain text template to be valid")
	}
}
