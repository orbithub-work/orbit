package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ExternalCommandParser calls an external executable to process files
type ExternalCommandParser struct {
	id         string
	executable string
	extensions []string
}

func NewExternalCommandParser(id, executable string, extensions []string) *ExternalCommandParser {
	return &ExternalCommandParser{
		id:         id,
		executable: executable,
		extensions: extensions,
	}
}

func (p *ExternalCommandParser) Name() string {
	return fmt.Sprintf("external.%s", p.id)
}

func (p *ExternalCommandParser) CanHandle(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	for _, e := range p.extensions {
		if strings.ToLower(strings.TrimPrefix(e, ".")) == ext {
			return true
		}
	}
	return false
}

func (p *ExternalCommandParser) Parse(ctx context.Context, path string) (*Result, error) {
	// Call external executable: <executable> <path>
	// Expect JSON output on stdout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, p.executable, path)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("external parser %s failed: %s", p.id, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run external parser %s: %w", p.id, err)
	}

	var res Result
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, fmt.Errorf("failed to parse output from external parser %s: %w", p.id, err)
	}

	return &res, nil
}

var _ Parser = (*ExternalCommandParser)(nil)
