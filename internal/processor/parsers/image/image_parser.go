package image

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"strings"

	"media-assistant-os/internal/processor"
	"media-assistant-os/internal/utils"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Name() string {
	return "native.image"
}

func (p *Parser) CanHandle(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	supported := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"webp": true,
		"bmp":  true,
	}
	return supported[ext]
}

func (p *Parser) Parse(ctx context.Context, path string) (*processor.Result, error) {
	_ = ctx
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 1. Decode Config for basic info
	imgConfig, format, err := image.DecodeConfig(file)
	if err != nil {
		return &processor.Result{Error: err}, nil
	}

	res := &processor.Result{
		Metadata: &processor.Metadata{
			Width:  imgConfig.Width,
			Height: imgConfig.Height,
			Format: strings.ToUpper(format),
			Extra:  make(map[string]any),
		},
	}

	// 2. Extract EXIF if JPEG/TIFF
	_, _ = file.Seek(0, 0)
	if x, err := exif.Decode(file); err == nil {
		if tm, err := x.DateTime(); err == nil {
			res.Metadata.Extra["date_taken"] = tm
		}
		if lat, long, err := x.LatLong(); err == nil {
			res.Metadata.Extra["latitude"] = lat
			res.Metadata.Extra["longitude"] = long
		}
		if model, err := x.Get(exif.Model); err == nil {
			res.Metadata.Extra["camera_model"] = model.String()
		}
		if lens, err := x.Get(exif.LensModel); err == nil {
			res.Metadata.Extra["lens_model"] = lens.String()
		}
	}

	// 3. Generate Thumbnail and Compute pHash
	_, _ = file.Seek(0, 0)
	img, _, err := image.Decode(file)
	if err == nil {
		// Compute dHash
		if hash, err := utils.DHashFromImage(img); err == nil {
			res.Metadata.Extra["phash"] = hash
		}

		thumb := resize.Thumbnail(300, 300, img, resize.Lanczos3)
		buf := new(bytes.Buffer)
		var encodeErr error
		if format == "jpeg" || format == "jpg" {
			encodeErr = jpeg.Encode(buf, thumb, nil)
		} else {
			encodeErr = png.Encode(buf, thumb)
		}
		if encodeErr == nil && buf.Len() > 0 {
			res.Thumbnail = buf.Bytes()
		}
	}

	return res, nil
}

// Ensure interface is implemented
var _ processor.Parser = (*Parser)(nil)
