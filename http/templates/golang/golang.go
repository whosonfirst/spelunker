// Package golang provides functions for working with ".golang" templates.
package golang

import (
	"context"
	"embed"

	sfomuseum_text "github.com/sfomuseum/go-template/text"
	"text/template"
)

//go:embed *.golang
var FS embed.FS

// LoadTemplates instantiates ".golang" templates.
func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_text.LoadTemplatesMatching(ctx, "*.golang", FS)
}
