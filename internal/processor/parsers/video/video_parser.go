package video

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"media-assistant-os/internal/processor"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Name() string {
	return "native.video"
}

func (p *Parser) CanHandle(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	supported := map[string]bool{
		"mp4":  true,
		"mov":  true,
		"avi":  true,
		"mkv":  true,
		"webm": true,
		"flv":  true,
	}
	return supported[ext]
}

func (p *Parser) Parse(ctx context.Context, path string) (*processor.Result, error) {
	ext := strings.ToUpper(strings.TrimPrefix(filepath.Ext(path), "."))
	res := &processor.Result{
		Metadata: &processor.Metadata{
			Format: ext,
			Extra:  map[string]any{},
		},
	}

	if meta, err := probeVideo(ctx, path); err == nil {
		res.Metadata.Width = meta.Width
		res.Metadata.Height = meta.Height
		res.Metadata.Duration = meta.Duration
		res.Metadata.Codec = meta.Codec
	} else {
		res.Metadata.Extra["probe_error"] = err.Error()
	}

	if thumb, err := captureThumbnail(ctx, path); err == nil {
		res.Thumbnail = thumb
	} else {
		res.Metadata.Extra["thumbnail_error"] = err.Error()
	}

	return res, nil
}

type videoMeta struct {
	Width    int
	Height   int
	Duration float64
	Codec    string
}

type ffprobeOutput struct {
	Streams []struct {
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		CodecName string `json:"codec_name"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func probeVideo(ctx context.Context, path string) (videoMeta, error) {
	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		return videoMeta{}, errors.New("ffprobe not found")
	}

	cmd := exec.CommandContext(ctx,
		ffprobePath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,codec_name",
		"-show_entries", "format=duration",
		"-of", "json",
		path,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return videoMeta{}, err
	}

	var parsed ffprobeOutput
	if err := json.Unmarshal(out.Bytes(), &parsed); err != nil {
		return videoMeta{}, err
	}

	meta := videoMeta{}
	if len(parsed.Streams) > 0 {
		meta.Width = parsed.Streams[0].Width
		meta.Height = parsed.Streams[0].Height
		meta.Codec = parsed.Streams[0].CodecName
	}

	if parsed.Format.Duration != "" {
		if d, err := strconv.ParseFloat(parsed.Format.Duration, 64); err == nil {
			meta.Duration = d
		}
	}

	return meta, nil
}

func captureThumbnail(ctx context.Context, path string) ([]byte, error) {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, errors.New("ffmpeg not found")
	}

	cmd := exec.CommandContext(ctx,
		ffmpegPath,
		"-ss", "0.5",
		"-i", path,
		"-frames:v", "1",
		"-vf", "scale=300:-1",
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-an",
		"-sn",
		"pipe:1",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	if out.Len() == 0 {
		return nil, errors.New("ffmpeg returned empty thumbnail")
	}
	return out.Bytes(), nil
}

var _ processor.Parser = (*Parser)(nil)
