package services

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FingerprintProfile describes how a file fingerprint was produced.
type FingerprintProfile struct {
	Version      string
	FileType     string
	BlockSize    int64
	SamplePoints int
	FullRead     bool
}

func ComputeSampleFingerprint(path string) (string, int64, int64, error) {
	fp, size, mtime, _, err := ComputeAdaptiveFingerprint(path)
	return fp, size, mtime, err
}

// ComputeAdaptiveFingerprint computes a deterministic fingerprint using
// type-aware and size-aware sampling to balance throughput and accuracy.
func ComputeAdaptiveFingerprint(path string) (string, int64, int64, FingerprintProfile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", 0, 0, FingerprintProfile{}, err
	}
	if info.IsDir() {
		return "", 0, 0, FingerprintProfile{}, errors.New("path is directory")
	}

	size := info.Size()
	mtime := info.ModTime().Unix()

	f, err := os.Open(path)
	if err != nil {
		return "", 0, 0, FingerprintProfile{}, err
	}
	defer f.Close()

	fileType := detectFileType(path)
	profile := chooseFingerprintProfile(fileType, size)

	h := sha256.New()

	// Fingerprint version and media class are part of hash input so future
	// algorithm upgrades don't mix with prior fingerprints.
	_, _ = h.Write([]byte(profile.Version))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(profile.FileType))
	_, _ = h.Write([]byte{0})

	var sizeBuf [8]byte
	binary.LittleEndian.PutUint64(sizeBuf[:], uint64(size))
	_, _ = h.Write(sizeBuf[:])

	if profile.FullRead {
		if _, err := io.Copy(h, f); err != nil {
			return "", 0, 0, FingerprintProfile{}, err
		}
		return hex.EncodeToString(h.Sum(nil)), size, mtime, profile, nil
	}

	offsets := sampleOffsets(size, profile.BlockSize, profile.SamplePoints)
	for _, off := range offsets {
		if err := hashAt(h, f, off, profile.BlockSize); err != nil {
			return "", 0, 0, FingerprintProfile{}, err
		}
	}
	return hex.EncodeToString(h.Sum(nil)), size, mtime, profile, nil
}

func hashAt(h io.Writer, f *os.File, offset int64, n int64) error {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	_, err := io.CopyN(h, f, n)
	if err == io.EOF {
		return nil
	}
	return err
}

func chooseFingerprintProfile(fileType string, size int64) FingerprintProfile {
	const (
		mb = int64(1024 * 1024)
		gb = int64(1024 * 1024 * 1024)
	)
	profile := FingerprintProfile{
		Version:   "sample-v2",
		FileType:  fileType,
		BlockSize: 64 * 1024, // 64KB
	}

	switch fileType {
	case "image":
		if size < 50*mb {
			profile.FullRead = true
			return profile
		}
	case "audio":
		if size < 200*mb {
			profile.FullRead = true
			return profile
		}
	default:
		if size < 200*mb {
			profile.FullRead = true
			return profile
		}
	}

	switch {
	case size < 100*mb:
		profile.SamplePoints = 8
	case size < 2*gb:
		profile.SamplePoints = 16
	case size < 20*gb:
		profile.SamplePoints = 32
	default:
		profile.SamplePoints = 64
	}
	return profile
}

func sampleOffsets(size int64, blockSize int64, points int) []int64 {
	if points <= 0 || size <= 0 {
		return []int64{0}
	}
	last := size - blockSize
	if last < 0 {
		last = 0
	}
	add := func(out []int64, seen map[int64]struct{}, off int64) []int64 {
		if off < 0 {
			off = 0
		}
		if off > last {
			off = last
		}
		if _, ok := seen[off]; ok {
			return out
		}
		seen[off] = struct{}{}
		return append(out, off)
	}

	offsets := make([]int64, 0, points)
	seen := make(map[int64]struct{}, points+8)

	// Hot points: head, center, tail.
	hot := []float64{0.0, 0.01, 0.02, 0.5, 0.98, 0.99, 1.0}
	for _, p := range hot {
		off := int64(float64(last) * p)
		offsets = add(offsets, seen, off)
	}

	remaining := points - len(offsets)
	if remaining <= 0 {
		return offsets[:points]
	}

	// Uniform points over [5%, 95%].
	for i := 0; i < remaining; i++ {
		ratio := 0.05 + (0.90*float64(i))/float64(maxInt(1, remaining-1))
		off := int64(float64(last) * ratio)
		offsets = add(offsets, seen, off)
		if len(offsets) >= points {
			break
		}
	}

	return offsets
}

func detectFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4", ".mov", ".avi", ".mkv", ".webm", ".m4v", ".ts", ".mpeg", ".mpg":
		return "video"
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".heic":
		return "image"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a":
		return "audio"
	default:
		return "binary"
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
