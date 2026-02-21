package fallback

import (
	"context"
	"path/filepath"
	"strings"

	"media-assistant-os/internal/processor"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Name() string {
	return "dummy.fallback"
}

func (p *Parser) CanHandle(ext string) bool {
	// Handle everything else
	return true
}

func (p *Parser) Parse(ctx context.Context, path string) (*processor.Result, error) {
	_ = ctx
	ext := strings.ToLower(filepath.Ext(path))
	return &processor.Result{
		Metadata: &processor.Metadata{
			Format: strings.ToUpper(strings.TrimPrefix(ext, ".")),
			Extra: map[string]any{
				"parser": "dummy",
				"note":   "placeholder for unimplemented format",
			},
		},
	}, nil
}
