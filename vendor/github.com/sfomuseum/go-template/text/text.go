// Package text provides methods for loading text (.txt) templates with default functions
package text

import (
	"context"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/sfomuseum/go-template/funcs"
)

// LoadTemplates loads text templates matching ".txt" from 't_fs' with default functions assigned.
func LoadTemplates(ctx context.Context, t_fs ...fs.FS) (*template.Template, error) {

	funcs := TemplatesFuncMap()
	pattern := "*.txt"

	return LoadTemplatesWithFuncMapMatching(ctx, funcs, pattern, t_fs...)
}

// LoadTemplatesWithFuncMap loads text templates matching ".txt" from 't_fs' with default functions defined by 'funcs'.
func LoadTemplatesWithFuncMap(ctx context.Context, funcs template.FuncMap, t_fs ...fs.FS) (*template.Template, error) {

	pattern := "*.txt"
	return LoadTemplatesWithFuncMapMatching(ctx, funcs, pattern, t_fs...)
}

// LoadTemplatesMatching loads text templates matching 'pattern' from 't_fs' with default functions assigned.
func LoadTemplatesMatching(ctx context.Context, pattern string, t_fs ...fs.FS) (*template.Template, error) {

	funcs := TemplatesFuncMap()
	return LoadTemplatesWithFuncMapMatching(ctx, funcs, pattern, t_fs...)
}

// LoadTemplatesWithFuncMapMatching loads text templates matching 'pattern' from 't_fs' with default functions defined by 'funcs'.
func LoadTemplatesWithFuncMapMatching(ctx context.Context, funcs template.FuncMap, pattern string, t_fs ...fs.FS) (*template.Template, error) {

	t := template.New("text").Funcs(funcs)

	var err error

	for idx, f := range t_fs {

		t, err = t.ParseFS(f, pattern)

		if err != nil {
			return nil, fmt.Errorf("Failed to load templates from FS at offset %d, %w", idx, err)
		}
	}

	return t, nil
}

// TemplatesFuncMap() returns a `template.FuncMap` instance with default functions assigned.
func TemplatesFuncMap() template.FuncMap {

	return template.FuncMap{
		// For example: {{ if (IsAvailable "Account" .) }}
		"IsAvailable": funcs.IsAvailable,
		"Add":         funcs.Add,
	}
}
