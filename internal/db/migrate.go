package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	Version int
	Up      func(ctx context.Context, tx *sql.Tx) error
}

func Migrate(ctx context.Context, d *DB) error {
	if err := ensureMigrationsTable(ctx, d.sql); err != nil {
		return err
	}

	migrations := []Migration{
		{Version: 1, Up: migrateV1},
		{Version: 2, Up: migrateV2},
		{Version: 3, Up: migrateV3},
		{Version: 4, Up: migrateV4},
		{Version: 5, Up: migrateV5},
		{Version: 6, Up: migrateV6},
		{Version: 7, Up: migrateV7},
		{Version: 8, Up: migrateV8},
		{Version: 9, Up: migrateV9},
		{Version: 10, Up: migrateV10},
		{Version: 11, Up: migrateV11},
		{Version: 12, Up: migrateV12},
		{Version: 13, Up: migrateV13},
		{Version: 14, Up: migrateV14},
		{Version: 15, Up: migrateV15},
		{Version: 16, Up: migrateV16},
		{Version: 17, Up: migrateV17},
		{Version: 18, Up: migrateV18},
		{Version: 19, Up: migrateV19},
		{Version: 20, Up: migrateV20},
		{Version: 21, Up: migrateV21},
		{Version: 22, Up: migrateV22},
		{Version: 23, Up: migrateV23},
		{Version: 24, Up: migrateV24},
		{Version: 25, Up: migrateV25},
		{Version: 26, Up: migrateV26},
	}

	applied, err := appliedVersions(ctx, d.sql)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}
		tx, err := d.sql.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if err := m.Up(ctx, tx); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := recordMigration(ctx, tx, m.Version); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, s *sql.DB) error {
	_, err := s.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  applied_at INTEGER NOT NULL
);`)
	return err
}

func appliedVersions(ctx context.Context, s *sql.DB) (map[int]bool, error) {
	rows, err := s.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int]bool{}
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

func recordMigration(ctx context.Context, tx *sql.Tx, version int) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version, applied_at) VALUES(?, ?)`, version, time.Now().Unix())
	return err
}

func migrateV1(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`
CREATE TABLE IF NOT EXISTS projects (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  project_type TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);`,
		`
CREATE TABLE IF NOT EXISTS assets (
  id TEXT PRIMARY KEY,
  path TEXT NOT NULL UNIQUE,
  size INTEGER NOT NULL,
  mtime INTEGER NOT NULL,
  fingerprint TEXT,
  parent_asset_id TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);`,
		`
CREATE TABLE IF NOT EXISTS project_assets (
  project_id TEXT NOT NULL,
  asset_id TEXT NOT NULL,
  alias TEXT,
  tags_json TEXT,
  status TEXT,
  project_metadata_json TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  PRIMARY KEY(project_id, asset_id)
);`,
		`
CREATE TABLE IF NOT EXISTS activity_logs (
  id TEXT PRIMARY KEY,
  level TEXT NOT NULL,
  message TEXT NOT NULL,
  created_at INTEGER NOT NULL
);`,
	}

	for _, q := range stmts {
		if _, err := tx.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("migration v1 failed: %w", err)
		}
	}
	return nil
}

func migrateV2(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE projects ADD COLUMN status TEXT NOT NULL DEFAULT 'active';`,
		`ALTER TABLE projects ADD COLUMN description TEXT NOT NULL DEFAULT '';`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV3(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE projects ADD COLUMN path TEXT NOT NULL DEFAULT '';`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV4(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE assets ADD COLUMN scope TEXT NOT NULL DEFAULT 'global';`,
		`ALTER TABLE assets ADD COLUMN project_id TEXT;`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV5(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE assets ADD COLUMN status TEXT NOT NULL DEFAULT 'PENDING';`,
		`ALTER TABLE assets ADD COLUMN media_meta TEXT NOT NULL DEFAULT '';`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV6(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS asset_lineage (
			id TEXT PRIMARY KEY,
			ancestor_id TEXT NOT NULL,
			descendant_id TEXT NOT NULL,
			relation_type TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			FOREIGN KEY(ancestor_id) REFERENCES assets(id),
			FOREIGN KEY(descendant_id) REFERENCES assets(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_assets_fingerprint ON assets(fingerprint);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_lineage_unique ON asset_lineage(ancestor_id, descendant_id, relation_type);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_lineage_ancestor ON asset_lineage(ancestor_id);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_lineage_descendant ON asset_lineage(descendant_id);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV7(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status);`)
	return err
}

func migrateV8(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS media_tasks (
			id TEXT PRIMARY KEY,
			asset_id TEXT NOT NULL,
			task_type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			priority INTEGER NOT NULL DEFAULT 0,
			retry_count INTEGER NOT NULL DEFAULT 0,
			max_retries INTEGER NOT NULL DEFAULT 3,
			worker_id TEXT,
			started_at DATETIME,
			finished_at DATETIME,
			error_message TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(asset_id) REFERENCES assets(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_media_tasks_status ON media_tasks(status);`,
		`CREATE INDEX IF NOT EXISTS idx_media_tasks_asset_id ON media_tasks(asset_id);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV9(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		// Update assets table
		`ALTER TABLE assets ADD COLUMN last_op_log TEXT;`,
		// Update activity_logs table
		`ALTER TABLE activity_logs ADD COLUMN asset_id TEXT;`,
		`ALTER TABLE activity_logs ADD COLUMN project_id TEXT;`,
		`CREATE INDEX IF NOT EXISTS idx_activity_logs_asset_id ON activity_logs(asset_id);`,
		`CREATE INDEX IF NOT EXISTS idx_activity_logs_project_id ON activity_logs(project_id);`,
	}
	for _, s := range stmts {
		// Using ExecContext and ignoring errors for "duplicate column" because SQLite doesn't have "IF NOT EXISTS" for ADD COLUMN
		_, _ = tx.ExecContext(ctx, s)
	}
	return nil
}

func migrateV10(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE media_tasks ADD COLUMN progress INTEGER NOT NULL DEFAULT 0;`,
	}
	for _, s := range stmts {
		_, _ = tx.ExecContext(ctx, s)
	}
	return nil
}

func migrateV11(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE media_tasks ADD COLUMN next_retry_at DATETIME;`,
		`ALTER TABLE media_tasks ADD COLUMN dead_letter_reason TEXT;`,
		`CREATE INDEX IF NOT EXISTS idx_media_tasks_retry_due ON media_tasks(status, next_retry_at);`,
		`CREATE TABLE IF NOT EXISTS media_task_dlq (
			id TEXT PRIMARY KEY,
			task_id TEXT NOT NULL,
			asset_id TEXT NOT NULL,
			task_type TEXT NOT NULL,
			error_message TEXT NOT NULL,
			retry_count INTEGER NOT NULL,
			max_retries INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_media_task_dlq_task_id ON media_task_dlq(task_id);`,
	}
	for _, s := range stmts {
		_, _ = tx.ExecContext(ctx, s)
	}
	return nil
}

func migrateV12(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS event_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_type TEXT NOT NULL,
			payload TEXT NOT NULL,
			created_at INTEGER NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_event_logs_event_type ON event_logs(event_type);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV13(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS project_artifacts (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			kind TEXT NOT NULL DEFAULT '',
			name TEXT NOT NULL DEFAULT '',
			path TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			meta_json TEXT NOT NULL DEFAULT '{}',
			source_plugin_id TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_project_artifacts_project_id ON project_artifacts(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_artifacts_kind ON project_artifacts(kind);`,
		`CREATE INDEX IF NOT EXISTS idx_project_artifacts_created_at ON project_artifacts(created_at);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV14(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS system_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV15(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			color TEXT,
			icon TEXT,
			parent_id TEXT,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(parent_id) REFERENCES tags(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);`,
		`CREATE INDEX IF NOT EXISTS idx_tags_parent_id ON tags(parent_id);`,
		`CREATE TABLE IF NOT EXISTS asset_tags (
			asset_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			PRIMARY KEY(asset_id, tag_id),
			FOREIGN KEY(asset_id) REFERENCES assets(id) ON DELETE CASCADE,
			FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_tags_asset_id ON asset_tags(asset_id);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_tags_tag_id ON asset_tags(tag_id);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV16(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE media_tasks ADD COLUMN lease_until DATETIME;`,
		`CREATE INDEX IF NOT EXISTS idx_media_tasks_lease_until ON media_tasks(status, lease_until);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_media_tasks_inflight_unique ON media_tasks(asset_id, task_type) WHERE status IN ('pending','processing');`,
	}
	for _, s := range stmts {
		_, _ = tx.ExecContext(ctx, s)
	}
	return nil
}

func migrateV17(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE projects ADD COLUMN last_activity_at INTEGER NOT NULL DEFAULT 0;`,
	}
	for _, s := range stmts {
		_, _ = tx.ExecContext(ctx, s)
	}
	return nil
}

func migrateV18(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS workflow_templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			domain TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			is_system INTEGER NOT NULL DEFAULT 0,
			is_active INTEGER NOT NULL DEFAULT 1,
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_workflow_templates_name ON workflow_templates(name);`,
		`CREATE INDEX IF NOT EXISTS idx_workflow_templates_domain ON workflow_templates(domain);`,
		`CREATE TABLE IF NOT EXISTS workflow_template_steps (
			id TEXT PRIMARY KEY,
			template_id TEXT NOT NULL,
			title TEXT NOT NULL,
			step_type TEXT NOT NULL DEFAULT 'task',
			order_index INTEGER NOT NULL DEFAULT 0,
			is_required INTEGER NOT NULL DEFAULT 1,
			suggested_days INTEGER NOT NULL DEFAULT 0,
			checklist_json TEXT NOT NULL DEFAULT '[]',
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(template_id) REFERENCES workflow_templates(id) ON DELETE CASCADE
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_workflow_template_steps_unique_order ON workflow_template_steps(template_id, order_index);`,
		`CREATE INDEX IF NOT EXISTS idx_workflow_template_steps_template ON workflow_template_steps(template_id);`,
		`CREATE TABLE IF NOT EXISTS project_workflows (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL UNIQUE,
			template_id TEXT NOT NULL,
			name TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'active',
			current_step_index INTEGER NOT NULL DEFAULT 0,
			start_at INTEGER NOT NULL DEFAULT 0,
			target_end_at INTEGER NOT NULL DEFAULT 0,
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(template_id) REFERENCES workflow_templates(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_project_workflows_template_id ON project_workflows(template_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_workflows_status ON project_workflows(status);`,
		`CREATE TABLE IF NOT EXISTS project_roadmap_items (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			workflow_id TEXT NOT NULL,
			parent_id TEXT,
			title TEXT NOT NULL,
			item_type TEXT NOT NULL DEFAULT 'task',
			status TEXT NOT NULL DEFAULT 'todo',
			priority INTEGER NOT NULL DEFAULT 0,
			order_index INTEGER NOT NULL DEFAULT 0,
			start_at INTEGER NOT NULL DEFAULT 0,
			due_at INTEGER NOT NULL DEFAULT 0,
			note TEXT NOT NULL DEFAULT '',
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(workflow_id) REFERENCES project_workflows(id) ON DELETE CASCADE,
			FOREIGN KEY(parent_id) REFERENCES project_roadmap_items(id) ON DELETE SET NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_project_roadmap_items_project ON project_roadmap_items(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_roadmap_items_workflow ON project_roadmap_items(workflow_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_roadmap_items_parent ON project_roadmap_items(parent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_roadmap_items_status ON project_roadmap_items(status);`,
		`CREATE TABLE IF NOT EXISTS project_notes (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			workflow_id TEXT,
			note_type TEXT NOT NULL DEFAULT 'note',
			title TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			source_asset_id TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			is_pinned INTEGER NOT NULL DEFAULT 0,
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(workflow_id) REFERENCES project_workflows(id) ON DELETE SET NULL,
			FOREIGN KEY(source_asset_id) REFERENCES assets(id) ON DELETE SET NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_project_notes_project ON project_notes(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_notes_workflow ON project_notes(workflow_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_notes_type ON project_notes(note_type);`,
		`CREATE INDEX IF NOT EXISTS idx_project_notes_pinned ON project_notes(is_pinned);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV19(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS project_workflow_steps (
			id TEXT PRIMARY KEY,
			workflow_id TEXT NOT NULL,
			template_step_id TEXT,
			title TEXT NOT NULL,
			step_type TEXT NOT NULL DEFAULT 'task',
			order_index INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'todo',
			is_required INTEGER NOT NULL DEFAULT 1,
			started_at INTEGER NOT NULL DEFAULT 0,
			completed_at INTEGER NOT NULL DEFAULT 0,
			kpi_rules_json TEXT NOT NULL DEFAULT '[]',
			kpi_result_json TEXT NOT NULL DEFAULT '{}',
			note TEXT NOT NULL DEFAULT '',
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(workflow_id) REFERENCES project_workflows(id) ON DELETE CASCADE,
			FOREIGN KEY(template_step_id) REFERENCES workflow_template_steps(id) ON DELETE SET NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_project_workflow_steps_unique_order ON project_workflow_steps(workflow_id, order_index);`,
		`CREATE INDEX IF NOT EXISTS idx_project_workflow_steps_workflow ON project_workflow_steps(workflow_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_workflow_steps_status ON project_workflow_steps(status);`,
		`CREATE TABLE IF NOT EXISTS publish_channels (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			platform TEXT NOT NULL,
			account_ref TEXT NOT NULL DEFAULT '',
			display_name TEXT NOT NULL DEFAULT '',
			is_enabled INTEGER NOT NULL DEFAULT 1,
			config_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_publish_channels_project ON publish_channels(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_publish_channels_platform ON publish_channels(platform);`,
		`CREATE TABLE IF NOT EXISTS publish_jobs (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			workflow_id TEXT,
			roadmap_item_id TEXT,
			artifact_id TEXT NOT NULL DEFAULT '',
			channel_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			schedule_at INTEGER NOT NULL DEFAULT 0,
			started_at INTEGER NOT NULL DEFAULT 0,
			finished_at INTEGER NOT NULL DEFAULT 0,
			idempotency_key TEXT NOT NULL,
			payload_json TEXT NOT NULL DEFAULT '{}',
			error_message TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(workflow_id) REFERENCES project_workflows(id) ON DELETE SET NULL,
			FOREIGN KEY(roadmap_item_id) REFERENCES project_roadmap_items(id) ON DELETE SET NULL,
			FOREIGN KEY(channel_id) REFERENCES publish_channels(id) ON DELETE CASCADE
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_publish_jobs_idempotency ON publish_jobs(idempotency_key);`,
		`CREATE INDEX IF NOT EXISTS idx_publish_jobs_project ON publish_jobs(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_publish_jobs_status_schedule ON publish_jobs(status, schedule_at);`,
		`CREATE TABLE IF NOT EXISTS publish_records (
			id TEXT PRIMARY KEY,
			job_id TEXT NOT NULL,
			project_id TEXT NOT NULL,
			platform TEXT NOT NULL,
			account_ref TEXT NOT NULL DEFAULT '',
			external_post_id TEXT NOT NULL DEFAULT '',
			post_url TEXT NOT NULL DEFAULT '',
			published_at INTEGER NOT NULL DEFAULT 0,
			raw_response_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(job_id) REFERENCES publish_jobs(id) ON DELETE CASCADE,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_publish_records_job ON publish_records(job_id);`,
		`CREATE INDEX IF NOT EXISTS idx_publish_records_external_post ON publish_records(external_post_id);`,
		`CREATE TABLE IF NOT EXISTS metrics_events (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			record_id TEXT,
			platform TEXT NOT NULL,
			account_ref TEXT NOT NULL DEFAULT '',
			external_post_id TEXT NOT NULL DEFAULT '',
			metric_type TEXT NOT NULL,
			metric_value REAL NOT NULL DEFAULT 0,
			occurred_at INTEGER NOT NULL,
			source TEXT NOT NULL DEFAULT '',
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(record_id) REFERENCES publish_records(id) ON DELETE SET NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_events_project ON metrics_events(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_events_external_post ON metrics_events(external_post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_events_metric_time ON metrics_events(metric_type, occurred_at);`,
		`CREATE TABLE IF NOT EXISTS metrics_snapshots_daily (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			record_id TEXT,
			platform TEXT NOT NULL,
			external_post_id TEXT NOT NULL DEFAULT '',
			snapshot_date TEXT NOT NULL,
			views REAL NOT NULL DEFAULT 0,
			likes REAL NOT NULL DEFAULT 0,
			comments REAL NOT NULL DEFAULT 0,
			shares REAL NOT NULL DEFAULT 0,
			favorites REAL NOT NULL DEFAULT 0,
			followers_delta REAL NOT NULL DEFAULT 0,
			meta_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(record_id) REFERENCES publish_records(id) ON DELETE SET NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_metrics_snapshots_daily_unique ON metrics_snapshots_daily(project_id, external_post_id, snapshot_date);`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_snapshots_daily_project ON metrics_snapshots_daily(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_snapshots_daily_date ON metrics_snapshots_daily(snapshot_date);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV20(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS lineage_candidates (
			id TEXT PRIMARY KEY,
			ancestor_id TEXT NOT NULL,
			descendant_id TEXT NOT NULL,
			project_id TEXT,
			rule_type TEXT NOT NULL DEFAULT '',
			score REAL NOT NULL DEFAULT 0,
			confidence TEXT NOT NULL DEFAULT 'LOW',
			status TEXT NOT NULL DEFAULT 'suggested',
			reason TEXT NOT NULL DEFAULT '',
			evidence_json TEXT NOT NULL DEFAULT '{}',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(ancestor_id) REFERENCES assets(id) ON DELETE CASCADE,
			FOREIGN KEY(descendant_id) REFERENCES assets(id) ON DELETE CASCADE,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE SET NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_lineage_candidates_unique_rule_pair ON lineage_candidates(ancestor_id, descendant_id, rule_type);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_candidates_ancestor ON lineage_candidates(ancestor_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_candidates_descendant ON lineage_candidates(descendant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_candidates_project ON lineage_candidates(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lineage_candidates_status_score ON lineage_candidates(status, score DESC);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV21(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`ALTER TABLE assets ADD COLUMN shape TEXT NOT NULL DEFAULT 'unknown';`,
		`ALTER TABLE assets ADD COLUMN suggested_rating INTEGER;`,
		`ALTER TABLE assets ADD COLUMN user_rating INTEGER;`,
		`CREATE INDEX IF NOT EXISTS idx_assets_shape ON assets(shape);`,
		`CREATE INDEX IF NOT EXISTS idx_assets_suggested_rating ON assets(suggested_rating);`,
		`CREATE INDEX IF NOT EXISTS idx_assets_user_rating ON assets(user_rating);`,
	}
	for _, s := range stmts {
		// Ignore duplicate-column errors for upgrade idempotency on partially migrated DBs.
		_, _ = tx.ExecContext(ctx, s)
	}

	// Backfill shape for assets with valid width/height metadata.
	_, _ = tx.ExecContext(ctx, `
		UPDATE assets
		SET shape = CASE
			WHEN json_valid(media_meta) = 1
				AND CAST(json_extract(media_meta, '$.width') AS INTEGER) > 0
				AND CAST(json_extract(media_meta, '$.height') AS INTEGER) > 0
			THEN CASE
				WHEN ABS(
					(CAST(json_extract(media_meta, '$.width') AS REAL) /
					 CAST(json_extract(media_meta, '$.height') AS REAL)) - 1.0
				) <= 0.12 THEN 'square'
				WHEN (CAST(json_extract(media_meta, '$.width') AS REAL) /
					  CAST(json_extract(media_meta, '$.height') AS REAL)) >= 1.8 THEN 'panorama'
				WHEN CAST(json_extract(media_meta, '$.width') AS INTEGER) >
					 CAST(json_extract(media_meta, '$.height') AS INTEGER) THEN 'landscape'
				ELSE 'portrait'
			END
			ELSE shape
		END
		WHERE shape = '' OR shape = 'unknown' OR shape IS NULL;
	`)

	return nil
}

func migrateV22(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS project_sources (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			root_path TEXT NOT NULL,
			source_type TEXT NOT NULL DEFAULT 'extra',
			watch_enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_project_sources_unique_project_path ON project_sources(project_id, root_path);`,
		`CREATE INDEX IF NOT EXISTS idx_project_sources_project ON project_sources(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_sources_watch_enabled ON project_sources(project_id, watch_enabled);`,
		`ALTER TABLE project_assets ADD COLUMN role TEXT NOT NULL DEFAULT 'source';`,
		`ALTER TABLE project_assets ADD COLUMN bind_mode TEXT NOT NULL DEFAULT 'auto';`,
		`ALTER TABLE project_assets ADD COLUMN confidence REAL NOT NULL DEFAULT 1.0;`,
		`CREATE INDEX IF NOT EXISTS idx_project_assets_project_role ON project_assets(project_id, role);`,
		`CREATE INDEX IF NOT EXISTS idx_project_assets_project_bind_mode ON project_assets(project_id, bind_mode);`,
	}

	for _, s := range stmts {
		// Ignore duplicate-column errors for upgrade idempotency on partially migrated DBs.
		_, _ = tx.ExecContext(ctx, s)
	}

	// Backfill existing project.path as primary project source.
	_, _ = tx.ExecContext(ctx, `
		INSERT INTO project_sources (id, project_id, root_path, source_type, watch_enabled, created_at, updated_at)
		SELECT lower(hex(randomblob(16))), p.id, p.path, 'primary', 1, strftime('%s','now'), strftime('%s','now')
		FROM projects p
		WHERE p.path IS NOT NULL AND TRIM(p.path) != ''
		ON CONFLICT(project_id, root_path) DO UPDATE SET
			source_type = CASE WHEN project_sources.source_type = 'primary' THEN 'primary' ELSE excluded.source_type END,
			watch_enabled = excluded.watch_enabled,
			updated_at = excluded.updated_at;
	`)

	return nil
}

func migrateV23(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS library_sources (
			id TEXT PRIMARY KEY,
			root_path TEXT NOT NULL UNIQUE,
			watch_enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_library_sources_created_at ON library_sources(created_at);`,
		`CREATE TABLE IF NOT EXISTS project_source_bind_jobs (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			source_id TEXT NOT NULL,
			root_path TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			total_assets INTEGER NOT NULL DEFAULT 0,
			processed_assets INTEGER NOT NULL DEFAULT 0,
			error_message TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			finished_at INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY(source_id) REFERENCES project_sources(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_project_source_bind_jobs_project ON project_source_bind_jobs(project_id, created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_project_source_bind_jobs_status ON project_source_bind_jobs(status, updated_at DESC);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_project_source_bind_jobs_inflight_unique
			ON project_source_bind_jobs(project_id, source_id)
			WHERE status IN ('pending', 'running');`,
		`ALTER TABLE project_assets ADD COLUMN source_id TEXT NOT NULL DEFAULT '';`,
		`CREATE INDEX IF NOT EXISTS idx_project_assets_project_source ON project_assets(project_id, source_id);`,
	}

	for _, s := range stmts {
		_, _ = tx.ExecContext(ctx, s)
	}

	// Backfill library sources from existing project source roots and legacy project.path.
	_, _ = tx.ExecContext(ctx, `
		INSERT INTO library_sources (id, root_path, watch_enabled, created_at, updated_at)
		SELECT lower(hex(randomblob(16))), ps.root_path, ps.watch_enabled, strftime('%s','now'), strftime('%s','now')
		FROM project_sources ps
		WHERE ps.root_path IS NOT NULL AND TRIM(ps.root_path) != ''
		ON CONFLICT(root_path) DO UPDATE SET
			watch_enabled = excluded.watch_enabled,
			updated_at = excluded.updated_at;
	`)
	_, _ = tx.ExecContext(ctx, `
		INSERT INTO library_sources (id, root_path, watch_enabled, created_at, updated_at)
		SELECT lower(hex(randomblob(16))), p.path, 1, strftime('%s','now'), strftime('%s','now')
		FROM projects p
		WHERE p.path IS NOT NULL AND TRIM(p.path) != ''
		ON CONFLICT(root_path) DO UPDATE SET
			updated_at = excluded.updated_at;
	`)

	return nil
}

func migrateV24(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS asset_history_events (
			id TEXT PRIMARY KEY,
			asset_id TEXT NOT NULL,
			project_id TEXT NOT NULL DEFAULT '',
			event_type TEXT NOT NULL,
			source_path TEXT NOT NULL DEFAULT '',
			target_path TEXT NOT NULL DEFAULT '',
			confidence TEXT NOT NULL DEFAULT 'medium',
			is_inferred BOOLEAN NOT NULL DEFAULT 1,
			detail TEXT NOT NULL DEFAULT '',
			occurred_at INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			FOREIGN KEY(asset_id) REFERENCES assets(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_history_asset_time ON asset_history_events(asset_id, occurred_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_history_project_time ON asset_history_events(project_id, occurred_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_asset_history_event_type ON asset_history_events(event_type, occurred_at DESC);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV25(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS search_history (
			id TEXT PRIMARY KEY,
			query_hash TEXT NOT NULL,
			query TEXT NOT NULL,
			filters TEXT NOT NULL DEFAULT '',
			count INTEGER NOT NULL DEFAULT 1,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_search_history_query_hash_unique ON search_history(query_hash);`,
		`CREATE INDEX IF NOT EXISTS idx_search_history_updated ON search_history(updated_at DESC);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func migrateV26(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS plugin_runtimes (
			plugin_id TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			info_json TEXT NOT NULL,
			issued_at INTEGER NOT NULL,
			last_used_at INTEGER NOT NULL,
			expires_at INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_runtimes_token_unique ON plugin_runtimes(token);`,
		`CREATE INDEX IF NOT EXISTS idx_plugin_runtimes_updated_at ON plugin_runtimes(updated_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_plugin_runtimes_expires_at ON plugin_runtimes(expires_at);`,
	}
	for _, s := range stmts {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
