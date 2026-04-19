// Copyright 2026-2030 Aptlogica Technologies Pvt Ltd
// Licensed under the Apache License, Version 2.0
// Websites: https://www.aptlogica.com | https://www.serenibase.com
// Support: support@aptlogica.com | support@serenibase.com

package templatecache

import (
	"html/template"
)

// ParseAndCacheTemplate parses a template and stores it in the cache
func ParseAndCacheTemplate(tc *TemplateCache, name, src string) (*template.Template, error) {
	tmpl, err := template.New(name).Parse(src)
	if err != nil {
		return nil, err
	}
	tc.Set(name, tmpl)
	return tmpl, nil
}
