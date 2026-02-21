package repos

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type EventLog struct {
	bun.BaseModel `bun:"table:event_logs"`

	ID        int64  `bun:"id,pk,autoincrement" json:"id"`
	EventType string `bun:"event_type,notnull" json:"event_type"`
	Payload   string `bun:"payload,notnull" json:"payload"`
	CreatedAt int64  `bun:"created_at,notnull" json:"created_at"`
}

type EventLogRepo struct {
	db *bun.DB
}

func NewEventLogRepo(db *bun.DB) *EventLogRepo {
	return &EventLogRepo{db: db}
}

func (r *EventLogRepo) Append(ctx context.Context, eventType string, payload string) (int64, error) {
	item := &EventLog{
		EventType: eventType,
		Payload:   payload,
		CreatedAt: time.Now().Unix(),
	}
	_, err := r.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *EventLogRepo) ListSince(ctx context.Context, sinceID int64, limit int) ([]EventLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	var out []EventLog
	err := r.db.NewSelect().
		Model(&out).
		Where("id > ?", sinceID).
		Order("id ASC").
		Limit(limit).
		Scan(ctx)
	return out, err
}
