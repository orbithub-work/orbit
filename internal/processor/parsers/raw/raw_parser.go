package raw

import (
	"bytes"
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
	return "external.raw"
}

func (p *Parser) CanHandle(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	supported := map[string]bool{
		"cr2": true,
		"cr3": true,
		"nef": true,
		"arw": true,
		"dng": true,
		"raf": true,
		"orf": true,
		"rw2": true,
		"pef": true,
		"srw": true,
		"3fr": true,
	}
	return supported[ext]
}

func (p *Parser) Parse(ctx context.Context, path string) (*processor.Result, error) {
	res := &processor.Result{
		Metadata: &processor.Metadata{
			Format: "RAW",
			Extra:  map[string]any{},
		},
	}

	if w, h, err := readRawSize(ctx, path); err == nil {
		res.Metadata.Width = w
		res.Metadata.Height = h
	} else {
		res.Metadata.Extra["size_error"] = err.Error()
	}

	if thumb, err := renderRawThumbnail(ctx, path); err == nil {
		res.Thumbnail = thumb
	} else {
		res.Metadata.Extra["thumbnail_error"] = err.Error()
	}

	return res, nil
}

func readRawSize(ctx context.Context, path string) (int, int, error) {
	if sipsPath, err := exec.LookPath("sips"); err == nil {
		cmd := exec.CommandContext(ctx, sipsPath, "-g", "pixelWidth", "-g", "pixelHeight", path)
		out, err := cmd.Output()
		if err == nil {
			return parseSipsSize(string(out))
		}
	}

	if dcrawPath, err := exec.LookPath("dcraw"); err == nil {
		cmd := exec.CommandContext(ctx, dcrawPath, "-i", "-v", path)
		out, err := cmd.Output()
		if err == nil {
			return parseDcrawSize(string(out))
		}
	}

	return 0, 0, errors.New("no available tool to read RAW size")
}

func renderRawThumbnail(ctx context.Context, path string) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "raw-thumb-*.jpg")
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

	if dcrawPath, err := exec.LookPath("dcraw"); err == nil {
		cmd := exec.CommandContext(ctx, dcrawPath, "-e", "-c", path)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err == nil && out.Len() > 0 {
			return out.Bytes(), nil
		}
	}

	return nil, errors.New("no available tool to render RAW thumbnail")
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

func parseDcrawSize(out string) (int, int, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Image size:") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				w, _ := strconv.Atoi(parts[2])
				h, _ := strconv.Atoi(parts[4])
				if w > 0 && h > 0 {
					return w, h, nil
				}
			}
		}
	}
	return 0, 0, errors.New("failed to parse size from dcraw output")
}

var _ processor.Parser = (*Parser)(nil)
