package db

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	"media-assistant-os/internal/pkg/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

type DB struct {
	sql *sql.DB
	orm *bun.DB
}

func Open(dataDir string) (*DB, error) {
	dbPath := filepath.Join(dataDir, "db", "media_assistant.db")
	s, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	// Increase connection limits to take advantage of WAL mode concurrency.
	// WAL mode allows multiple readers and one writer simultaneously.
	s.SetMaxOpenConns(25)
	s.SetMaxIdleConns(5)
	s.SetConnMaxLifetime(5 * time.Minute)

	// SQLite performance PRAGMAs
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=-64000",
		"PRAGMA busy_timeout=5000",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA mmap_size=268435456",
	}
	for _, p := range pragmas {
		if _, err := s.Exec(p); err != nil {
			logger.Warn("PRAGMA warning", zap.String("pragma", p), zap.Error(err))
		}
	}

	return &DB{
		sql: s,
		orm: bun.NewDB(s, sqlitedialect.New()),
	}, nil
}

func (d *DB) SQL() *sql.DB {
	return d.sql
}

func (d *DB) ORM() *bun.DB {
	return d.orm
}

func (d *DB) Close() error {
	if d == nil || d.sql == nil {
		return nil
	}
	return d.sql.Close()
}

func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.sql == nil {
		return sql.ErrConnDone
	}
	return d.sql.PingContext(ctx)
}
