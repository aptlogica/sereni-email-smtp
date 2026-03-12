// Copyright (c) 2026 Aptlogica Technologies Private Limited
// SPDX-License-Identifier: MIT
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package templatecache

import (
	"html/template"
	"sync"
)

// TemplateCache caches parsed templates for reuse
// Thread-safe for concurrent access

type TemplateCache struct {
	cache map[string]*template.Template
	mu    sync.RWMutex
}

func NewTemplateCache() *TemplateCache {
	return &TemplateCache{
		cache: make(map[string]*template.Template),
	}
}

func (tc *TemplateCache) Get(name string) (*template.Template, bool) {
	tc.mu.RLock()
	tmpl, ok := tc.cache[name]
	tc.mu.RUnlock()
	return tmpl, ok
}

func (tc *TemplateCache) Set(name string, tmpl *template.Template) {
	tc.mu.Lock()
	tc.cache[name] = tmpl
	tc.mu.Unlock()
}

func (tc *TemplateCache) Preload(templates map[string]string) error {
	for name, src := range templates {
		tmpl, err := template.New(name).Parse(src)
		if err != nil {
			return err
		}
		tc.Set(name, tmpl)
	}
	return nil
}
