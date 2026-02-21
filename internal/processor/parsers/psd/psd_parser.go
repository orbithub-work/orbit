package psd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"media-assistant-os/internal/processor"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Name() string {
	return "external.psd"
}

func (p *Parser) CanHandle(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	return ext == "psd" || ext == "psb"
}

func (p *Parser) Parse(ctx context.Context, path string) (*processor.Result, error) {
	res := &processor.Result{
		Metadata: &processor.Metadata{
			Format: "PSD",
			Extra:  map[string]any{},
		},
	}

	if w, h, err := readImageSize(ctx, path); err == nil {
		res.Metadata.Width = w
		res.Metadata.Height = h
	} else {
		res.Metadata.Extra["size_error"] = err.Error()
	}

	if thumb, err := renderThumbnail(ctx, path); err == nil {
		res.Thumbnail = thumb
	} else {
		res.Metadata.Extra["thumbnail_error"] = err.Error()
	}

	return res, nil
}

func readImageSize(ctx context.Context, path string) (int, int, error) {
	if sipsPath, err := exec.LookPath("sips"); err == nil {
		cmd := exec.CommandContext(ctx, sipsPath, "-g", "pixelWidth", "-g", "pixelHeight", path)
		out, err := cmd.Output()
		if err == nil {
			return parseSipsSize(string(out))
		}
	}

	if magickPath, err := exec.LookPath("magick"); err == nil {
		cmd := exec.CommandContext(ctx, magickPath, "identify", "-format", "%w %h", path+"[0]")
		out, err := cmd.Output()
		if err == nil {
			parts := strings.Fields(string(out))
			if len(parts) >= 2 {
				w, _ := strconv.Atoi(parts[0])
				h, _ := strconv.Atoi(parts[1])
				if w > 0 && h > 0 {
					return w, h, nil
				}
			}
		}
	}

	if convertPath, err := exec.LookPath("convert"); err == nil {
		cmd := exec.CommandContext(ctx, convertPath, path+"[0]", "-format", "%w %h", "info:")
		out, err := cmd.Output()
		if err == nil {
			parts := strings.Fields(string(out))
			if len(parts) >= 2 {
				w, _ := strconv.Atoi(parts[0])
				h, _ := strconv.Atoi(parts[1])
				if w > 0 && h > 0 {
					return w, h, nil
				}
			}
		}
	}

	return 0, 0, errors.New("no available tool to read PSD size")
}

func renderThumbnail(ctx context.Context, path string) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "psd-thumb-*.jpg")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	defer os.Remove(tmpPath)

	if sipsPath, err := exec.LookPath("sips"); err == nil {
		cmd := exec.CommandContext(ctx, sipsPath, "-s", "format", "jpeg", "-Z", "300", path, "--out", tmpPath)
		if err := cmd.Run(); err == nil {
			return os.ReadFile(tmpPath)
		}
	}

	if magickPath, err := exec.LookPath("magick"); err == nil {
		cmd := exec.CommandContext(ctx, magickPath, path+"[0]", "-thumbnail", "300x300", "jpeg:"+tmpPath)
		if err := cmd.Run(); err == nil {
			return os.ReadFile(tmpPath)
		}
	}

	if convertPath, err := exec.LookPath("convert"); err == nil {
		cmd := exec.CommandContext(ctx, convertPath, path+"[0]", "-thumbnail", "300x300", "jpeg:"+tmpPath)
		if err := cmd.Run(); err == nil {
			return os.ReadFile(tmpPath)
		}
	}

	return nil, errors.New("no available tool to render PSD thumbnail")
}

func parseSipsSize(out string) (int, int, error) {
	lines := strings.Split(out, "\n")
	w, h := 0, 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "pixelWidth:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				w, _ = strconv.Atoi(parts[1])
			}
		}
		if strings.HasPrefix(line, "pixelHeight:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				h, _ = strconv.Atoi(parts[1])
			}
		}
	}
	if w > 0 && h > 0 {
		return w, h, nil
	}
	return 0, 0, errors.New("failed to parse size from sips output")
}

var _ processor.Parser = (*Parser)(nil)
