package processor

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Metadata represents extracted file information
type Metadata struct {
	Width    int            `json:"width,omitempty"`
	Height   int            `json:"height,omitempty"`
	Duration float64        `json:"duration,omitempty"` // in seconds
	Format   string         `json:"format"`             // e.g. "JPEG", "PSD", "MP4"
	Codec    string         `json:"codec,omitempty"`
	Extra    map[string]any `json:"extra,omitempty"` // for format-specific data (EXIF, Layers, etc.)
}

// Result represents the outcome of a file process
type Result struct {
	Metadata  *Metadata
	Thumbnail []byte
	Error     error
}

// Parser is the interface for all file type handlers
type Parser interface {
	// Name returns a unique name for the parser
	Name() string
	// CanHandle returns true if this parser can process the given extension or mime type
	CanHandle(ext string) bool
	// Parse extracts metadata and potentially generates a thumbnail
	Parse(ctx context.Context, path string) (*Result, error)
}

type parserEntry struct {
	parser   Parser
	priority int // Higher value = higher priority
}

// Manager manages the registration and selection of parsers
type Manager struct {
	parsers []parserEntry
	mu      sync.RWMutex
}

var (
	defaultManager *Manager
	once           sync.Once
)

// GetManager returns the singleton instance of Manager
func GetManager() *Manager {
	once.Do(func() {
		defaultManager = &Manager{
			parsers: []parserEntry{},
		}
	})
	return defaultManager
}

// Register adds a new parser to the manager with a priority
func (m *Manager) Register(p Parser, priority int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add or update
	entry := parserEntry{parser: p, priority: priority}
	for i, existing := range m.parsers {
		if existing.parser.Name() == p.Name() {
			m.parsers[i] = entry
			return
		}
	}
	m.parsers = append(m.parsers, entry)

	// Sort by priority (descending)
	sort.Slice(m.parsers, func(i, j int) bool {
		return m.parsers[i].priority > m.parsers[j].priority
	})
}

// Process finds the first suitable parser and processes the file
func (m *Manager) Process(ctx context.Context, path string, ext string) (*Result, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, entry := range m.parsers {
		if entry.parser.CanHandle(ext) {
			return entry.parser.Parse(ctx, path)
		}
	}

	return nil, fmt.Errorf("no parser found for extension: %s", ext)
}
