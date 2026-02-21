package services

import (
	"encoding/json"
	"math"
	"strings"

	"media-assistant-os/internal/models"
)

type DuplicateLevel string

const (
	DuplicateExact    DuplicateLevel = "EXACT_DUP"
	DuplicateLikely   DuplicateLevel = "LIKELY_DUP"
	DuplicateConflict DuplicateLevel = "CONFLICT"
)

type MediaMeta struct {
	Width    int     `json:"width,omitempty"`
	Height   int     `json:"height,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	Codec    string  `json:"codec,omitempty"`
	Format   string  `json:"format,omitempty"`
	FPS      float64 `json:"fps,omitempty"`
}

type duplicateCandidate struct {
	Parent *models.Asset
	Level  DuplicateLevel
}

// chooseDuplicateParent 从重复候选中选择最佳父资产
func chooseDuplicateParent(current *models.Asset, dupes []models.Asset) duplicateCandidate {
	best := duplicateCandidate{}
	for i := range dupes {
		d := dupes[i]
		if current == nil || d.ID == current.ID {
			continue
		}
		level := classifyDuplicate(current, &d)
		if level == DuplicateConflict {
			continue
		}
		if best.Parent == nil || duplicateRank(level) > duplicateRank(best.Level) {
			copy := d
			best.Parent = &copy
			best.Level = level
			if level == DuplicateExact {
				break
			}
		}
	}
	return best
}

// duplicateRank 返回重复级别的优先级分数
func duplicateRank(level DuplicateLevel) int {
	switch level {
	case DuplicateExact:
		return 2
	case DuplicateLikely:
		return 1
	default:
		return 0
	}
}

// classifyDuplicate 对两个资产进行重复分类
func classifyDuplicate(a *models.Asset, b *models.Asset) DuplicateLevel {
	if a == nil || b == nil {
		return DuplicateLikely
	}
	if a.Size > 0 && b.Size > 0 && a.Size != b.Size {
		return DuplicateConflict
	}

	ta := detectFileType(a.Path)
	tb := detectFileType(b.Path)
	if ta != tb {
		return DuplicateConflict
	}
	if ta == "binary" {
		return DuplicateExact
	}

	ma, oka := parseMediaMeta(a.MediaMeta)
	mb, okb := parseMediaMeta(b.MediaMeta)
	if !oka || !okb {
		return DuplicateLikely
	}

	switch ta {
	case "video":
		return classifyVideoMeta(ma, mb)
	case "image":
		return classifyImageMeta(ma, mb)
	case "audio":
		return classifyAudioMeta(ma, mb)
	default:
		return DuplicateLikely
	}
}

// classifyVideoMeta 对视频媒体元数据进行重复分类
func classifyVideoMeta(a MediaMeta, b MediaMeta) DuplicateLevel {
	if a.Width > 0 && b.Width > 0 && a.Width != b.Width {
		return DuplicateConflict
	}
	if a.Height > 0 && b.Height > 0 && a.Height != b.Height {
		return DuplicateConflict
	}
	if a.Duration > 0 && b.Duration > 0 && math.Abs(a.Duration-b.Duration) > 0.2 {
		return DuplicateConflict
	}
	// Different codec is suspicious but can happen in remux scenarios.
	if a.Codec != "" && b.Codec != "" && !eqFold(a.Codec, b.Codec) {
		return DuplicateLikely
	}
	if a.Width > 0 && b.Width > 0 && a.Height > 0 && b.Height > 0 && a.Duration > 0 && b.Duration > 0 {
		return DuplicateExact
	}
	return DuplicateLikely
}

// classifyImageMeta 对图片媒体元数据进行重复分类
func classifyImageMeta(a MediaMeta, b MediaMeta) DuplicateLevel {
	if a.Width > 0 && b.Width > 0 && a.Width != b.Width {
		return DuplicateConflict
	}
	if a.Height > 0 && b.Height > 0 && a.Height != b.Height {
		return DuplicateConflict
	}
	if a.Width > 0 && b.Width > 0 && a.Height > 0 && b.Height > 0 {
		return DuplicateExact
	}
	return DuplicateLikely
}

// classifyAudioMeta 对音频媒体元数据进行重复分类
func classifyAudioMeta(a MediaMeta, b MediaMeta) DuplicateLevel {
	if a.Duration > 0 && b.Duration > 0 && math.Abs(a.Duration-b.Duration) > 0.2 {
		return DuplicateConflict
	}
	if a.Codec != "" && b.Codec != "" && !eqFold(a.Codec, b.Codec) {
		return DuplicateLikely
	}
	if a.Duration > 0 && b.Duration > 0 {
		return DuplicateExact
	}
	return DuplicateLikely
}

// parseMediaMeta 解析媒体元数据JSON字符串
func parseMediaMeta(raw string) (MediaMeta, bool) {
	if strings.TrimSpace(raw) == "" {
		return MediaMeta{}, false
	}
	var m MediaMeta
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return MediaMeta{}, false
	}
	return m, true
}

// eqFold 比较两个字符串是否相等（忽略大小写和空格）
func eqFold(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// lineageType 根据重复级别返回血缘关系类型
func lineageType(level DuplicateLevel) string {
	switch level {
	case DuplicateExact:
		return "COPY_EXACT"
	case DuplicateLikely:
		return "COPY_LIKELY"
	default:
		return "COPY"
	}
}
