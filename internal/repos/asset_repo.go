package repos

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type AssetRepo struct {
	db *bun.DB
}

type AssetListQuery struct {
	ProjectID string
	Directory string
	Query     string
	TagIDs    []string
	Types     []string
	Shapes    []string

	SizeMin int64
	SizeMax int64

	RatingMin int
	RatingMax int

	MtimeFrom int64
	MtimeTo   int64

	WidthMin  int
	WidthMax  int
	HeightMin int
	HeightMax int

	SortBy    string
	SortOrder string

	Limit  int
	Offset int
}

func NewAssetRepo(db *bun.DB) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) GetByPath(ctx context.Context, path string) (*models.Asset, error) {
	var a models.Asset
	err := r.db.NewSelect().
		Model(&a).
		Where("path = ?", path).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AssetRepo) GetByID(ctx context.Context, id string) (*models.Asset, error) {
	var a models.Asset
	err := r.db.NewSelect().
		Model(&a).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AssetRepo) GetByPaths(ctx context.Context, paths []string) ([]models.Asset, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	var out []models.Asset
	err := r.db.NewSelect().
		Model(&out).
		Where("path IN (?)", bun.In(paths)).
		Scan(ctx)
	return out, err
}

func (r *AssetRepo) FindMissingAssetsBySize(ctx context.Context, size int64) ([]models.Asset, error) {
	var out []models.Asset
	err := r.db.NewSelect().
		Model(&out).
		Where("size = ?", size).
		Where("status = ?", "MISSING").
		Limit(10).
		Scan(ctx)
	return out, err
}

func (r *AssetRepo) GetAllPaths(ctx context.Context) ([]string, error) {
	var paths []string
	err := r.db.NewSelect().
		Model((*models.Asset)(nil)).
		Column("path").
		Scan(ctx, &paths)
	return paths, err
}

func (r *AssetRepo) RelinkAsset(ctx context.Context, id string, newPath string, mtime int64) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("path = ?", newPath).
		Set("mtime = ?", mtime).
		Set("status = ?", "READY").
		Set("last_op_log = ?", utils.FormatNow()+": File moved/renamed and relinked").
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) Create(ctx context.Context, path string, size int64, mtime int64) (*models.Asset, error) {
	now := time.Now().Unix()
	a := models.Asset{
		ID:        utils.NewID(),
		Path:      path,
		Size:      size,
		Mtime:     mtime,
		Scope:     "global",
		Status:    "PENDING",
		Shape:     "unknown",
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := r.db.NewInsert().Model(&a).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssetRepo) GetPendingAssets(ctx context.Context, limit int) ([]models.Asset, error) {
	var out []models.Asset
	err := r.db.NewSelect().
		Model(&out).
		Where("status = ?", "PENDING").
		Limit(limit).
		Scan(ctx)
	return out, err
}

func (r *AssetRepo) GetPendingCount(ctx context.Context) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Asset)(nil)).
		Where("status = ?", "PENDING").
		Count(ctx)
	return count, err
}

func (r *AssetRepo) UpdateMetadata(ctx context.Context, id string, size int64, mtime int64, status string, log string) error {
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("size = ?", size).
		Set("mtime = ?", mtime).
		Set("status = ?", status).
		Set("last_op_log = ?", log).
		Set("updated_at = ?", time.Now().Unix()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := r.validateStatusTransition(ctx, id, status); err != nil {
		return err
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) UpdateStatusWithLog(ctx context.Context, id string, status string, opLog string) error {
	if err := r.validateStatusTransition(ctx, id, status); err != nil {
		return err
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("status = ?", status).
		Set("last_op_log = ?", opLog).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) ClearOpLog(ctx context.Context, id string) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("last_op_log = ?", "").
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) UpdateFingerprint(ctx context.Context, id string, fingerprint string, parentAssetID *string) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("fingerprint = ?", fingerprint).
		Set("parent_asset_id = ?", parentAssetID).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) UpdateMediaMeta(ctx context.Context, id string, mediaMeta string) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("media_meta = ?", mediaMeta).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) UpdateDerivedMetadata(ctx context.Context, id string, mediaMeta string, shape string, suggestedRating *int) error {
	now := time.Now().Unix()
	q := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("media_meta = ?", mediaMeta).
		Set("shape = ?", shape).
		Set("updated_at = ?", now).
		Where("id = ?", id)

	if suggestedRating != nil {
		// Keep first suggestion and never override user-rated assets.
		q = q.Set(
			"suggested_rating = CASE WHEN user_rating IS NULL AND suggested_rating IS NULL THEN ? ELSE suggested_rating END",
			*suggestedRating,
		)
	}

	_, err := q.Exec(ctx)
	return err
}

func (r *AssetRepo) UpdateUserRating(ctx context.Context, id string, userRating *int) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("user_rating = ?", userRating).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Asset)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetRepo) FindActiveAssetsByFingerprint(ctx context.Context, fp string) ([]models.Asset, error) {
	var out []models.Asset
	err := r.db.NewSelect().
		Model(&out).
		Where("fingerprint = ?", fp).
		Where("status = ?", "READY").
		Limit(5).
		Scan(ctx)
	return out, err
}

func (r *AssetRepo) FindByFingerprint(ctx context.Context, fp string) ([]models.Asset, error) {
	var out []models.Asset
	err := r.db.NewSelect().
		Model(&out).
		Where("fingerprint = ?", fp).
		Where("status != ?", "IGNORED"). // Skip ignored assets
		Scan(ctx)
	return out, err
}

func (r *AssetRepo) GetByFingerprintIncludeIgnored(ctx context.Context, fp string) (*models.Asset, error) {
	var a models.Asset
	err := r.db.NewSelect().
		Model(&a).
		Where("fingerprint = ?", fp).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AssetRepo) ListRecentByDirectory(ctx context.Context, dir string, excludeAssetID string, fromMtime int64, toMtime int64, limit int) ([]models.Asset, error) {
	if limit <= 0 || limit > 300 {
		limit = 100
	}
	var out []models.Asset
	prefix := filepath.Clean(dir) + string(filepath.Separator) + "%"
	q := r.db.NewSelect().
		Model(&out).
		Where("path LIKE ?", prefix).
		OrderExpr("mtime DESC, created_at DESC").
		Limit(limit)
	if excludeAssetID != "" {
		q = q.Where("id != ?", excludeAssetID)
	}
	if fromMtime > 0 {
		q = q.Where("mtime >= ?", fromMtime)
	}
	if toMtime > 0 {
		q = q.Where("mtime <= ?", toMtime)
	}
	err := q.Scan(ctx)
	return out, err
}

func (r *AssetRepo) CountByDirectory(ctx context.Context, dir string) (int, error) {
	clean := strings.TrimSpace(dir)
	if clean == "" {
		return 0, nil
	}
	clean = filepath.Clean(clean)
	prefix := clean + string(filepath.Separator) + "%"
	return r.db.NewSelect().
		Model((*models.Asset)(nil)).
		Where("(path = ? OR path LIKE ?)", clean, prefix).
		Count(ctx)
}

func (r *AssetRepo) ListIDsByDirectory(ctx context.Context, dir string, limit int, offset int) ([]string, error) {
	clean := strings.TrimSpace(dir)
	if clean == "" {
		return []string{}, nil
	}
	clean = filepath.Clean(clean)
	if limit <= 0 || limit > 5000 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	prefix := clean + string(filepath.Separator) + "%"
	var ids []string
	err := r.db.NewSelect().
		Model((*models.Asset)(nil)).
		Column("id").
		Where("(path = ? OR path LIKE ?)", clean, prefix).
		OrderExpr("id ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &ids)
	return ids, err
}

func (r *AssetRepo) GetStats(ctx context.Context) (totalSize int64, count int64, err error) {
	err = r.db.NewSelect().
		Model((*models.Asset)(nil)).
		ColumnExpr("COALESCE(SUM(size), 0)").
		ColumnExpr("COUNT(*)").
		Scan(ctx, &totalSize, &count)
	return totalSize, count, err
}

// BatchCreate inserts multiple assets in a single transaction.
// Uses batch size of 500 to stay within SQLite variable limits.
func (r *AssetRepo) BatchCreate(ctx context.Context, assets []models.Asset) error {
	if len(assets) == 0 {
		return nil
	}
	const batchSize = 500
	for i := 0; i < len(assets); i += batchSize {
		end := i + batchSize
		if end > len(assets) {
			end = len(assets)
		}
		batch := assets[i:end]
		_, err := r.db.NewInsert().
			Model(&batch).
			On("CONFLICT (path) DO NOTHING").
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchUpdateStatus updates status for multiple asset IDs in one query.
func (r *AssetRepo) BatchUpdateStatus(ctx context.Context, ids []string, status string) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Asset)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", now).
		Where("id IN (?)", bun.In(ids)).
		Exec(ctx)
	return err
}

func (r *AssetRepo) validateStatusTransition(ctx context.Context, id string, to string) error {
	var from string
	err := r.db.NewSelect().
		Model((*models.Asset)(nil)).
		Column("status").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx, &from)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("asset not found: %s", id)
		}
		return err
	}
	if from == to {
		return nil
	}
	allowed := map[string]map[string]bool{
		"PENDING": {"READY": true, "ERROR": true, "MISSING": true, "IGNORED": true, "INDEXED": true},
		"READY":   {"MISSING": true, "IGNORED": true, "ERROR": true},
		"INDEXED": {"READY": true, "MISSING": true, "IGNORED": true, "ERROR": true},
		"MISSING": {"READY": true, "IGNORED": true, "ERROR": true, "INDEXED": true},
		"IGNORED": {"READY": true, "MISSING": true, "INDEXED": true},
		"ERROR":   {"PENDING": true, "READY": true, "IGNORED": true, "INDEXED": true},
	}
	if toMap, ok := allowed[from]; ok && toMap[to] {
		return nil
	}
	return fmt.Errorf("invalid asset status transition: %s -> %s", from, to)
}

func (r *AssetRepo) ListByQuery(ctx context.Context, req AssetListQuery) ([]models.Asset, int, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// 1) Count with the same filter set.
	var total int
	countQ := r.db.NewSelect().
		TableExpr("assets AS asset").
		ColumnExpr("COUNT(DISTINCT asset.id)")
	countQ = r.applyListFilters(countQ, req)
	if err := countQ.Scan(ctx, &total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []models.Asset{}, 0, nil
	}

	// 2) Fetch ids first for stable pagination on filtered result set.
	var ids []string
	idQ := r.db.NewSelect().
		TableExpr("assets AS asset").
		Column("asset.id")
	idQ = r.applyListFilters(idQ, req).
		GroupExpr("asset.id")
	idQ = r.applyListSort(idQ, req.SortBy, req.SortOrder).
		Limit(limit).
		Offset(offset)
	if err := idQ.Scan(ctx, &ids); err != nil {
		return nil, 0, err
	}
	if len(ids) == 0 {
		return []models.Asset{}, total, nil
	}

	// 3) Fetch full rows and preserve id order.
	var assets []models.Asset
	q := r.db.NewSelect().
		Model(&assets).
		Where("id IN (?)", bun.In(ids))

	caseSQL := "CASE id"
	args := make([]any, 0, len(ids))
	for i, id := range ids {
		caseSQL += " WHEN ? THEN " + fmt.Sprint(i)
		args = append(args, id)
	}
	caseSQL += " ELSE " + fmt.Sprint(len(ids)) + " END"
	q = q.OrderExpr(caseSQL, args...)

	if err := q.Scan(ctx); err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

func (r *AssetRepo) applyListFilters(q *bun.SelectQuery, req AssetListQuery) *bun.SelectQuery {
	if v := strings.TrimSpace(req.ProjectID); v != "" {
		q = q.Where(
			`EXISTS (
				SELECT 1
				FROM project_assets pa
				WHERE pa.asset_id = asset.id
				  AND pa.project_id = ?
			)`,
			v,
		)
	}

	if v := strings.TrimSpace(req.Directory); v != "" {
		clean := filepath.Clean(v)
		p1 := clean + string(filepath.Separator) + "%"
		q = q.Where("(asset.path = ? OR asset.path LIKE ?)", clean, p1)
	}

	if v := strings.TrimSpace(req.Query); v != "" {
		// Search tokens are AND-ed (no cross-token OR).
		terms := strings.Fields(strings.ToLower(v))
		for _, term := range terms {
			if term == "" {
				continue
			}
			pattern := "%" + term + "%"
			q = q.Where(
				`(
					LOWER(asset.path) LIKE ?
					OR EXISTS (
						SELECT 1
						FROM asset_tags at
						JOIN tags t ON t.id = at.tag_id
						WHERE at.asset_id = asset.id
						  AND LOWER(t.name) LIKE ?
					)
					OR EXISTS (
						SELECT 1
						FROM project_assets pa
						JOIN projects p ON p.id = pa.project_id
						WHERE pa.asset_id = asset.id
						  AND LOWER(p.name) LIKE ?
					)
				)`,
				pattern, pattern, pattern,
			)
		}
	}

	if len(req.TagIDs) > 0 {
		q = q.Where(
			`EXISTS (
				SELECT 1
				FROM asset_tags at
				WHERE at.asset_id = asset.id
				  AND at.tag_id IN (?)
			)`,
			bun.In(req.TagIDs),
		)
	}

	if len(req.Types) > 0 {
		var ors []string
		var args []any
		for _, t := range req.Types {
			ext := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(t)), ".")
			if ext == "" {
				continue
			}
			ors = append(ors, "LOWER(asset.path) LIKE ?")
			args = append(args, "%."+ext)
		}
		if len(ors) > 0 {
			q = q.Where("("+strings.Join(ors, " OR ")+")", args...)
		}
	}

	if len(req.Shapes) > 0 {
		cleaned := make([]string, 0, len(req.Shapes))
		for _, s := range req.Shapes {
			v := strings.ToLower(strings.TrimSpace(s))
			if v == "" {
				continue
			}
			cleaned = append(cleaned, v)
		}
		if len(cleaned) > 0 {
			q = q.Where("LOWER(asset.shape) IN (?)", bun.In(cleaned))
		}
	}

	if req.RatingMax == -1 {
		// "Unrated" means user has not provided a rating yet.
		q = q.Where("asset.user_rating IS NULL")
	} else {
		if req.RatingMin > 0 {
			q = q.Where("COALESCE(asset.user_rating, asset.suggested_rating, 0) >= ?", req.RatingMin)
		}
		if req.RatingMax > 0 {
			q = q.Where("COALESCE(asset.user_rating, asset.suggested_rating, 0) <= ?", req.RatingMax)
		}
	}

	if req.SizeMin > 0 {
		q = q.Where("asset.size >= ?", req.SizeMin)
	}
	if req.SizeMax > 0 {
		q = q.Where("asset.size <= ?", req.SizeMax)
	}
	if req.MtimeFrom > 0 {
		q = q.Where("asset.mtime >= ?", req.MtimeFrom)
	}
	if req.MtimeTo > 0 {
		q = q.Where("asset.mtime <= ?", req.MtimeTo)
	}

	if req.WidthMin > 0 {
		q = q.Where("CAST(json_extract(asset.media_meta, '$.width') AS INTEGER) >= ?", req.WidthMin)
	}
	if req.WidthMax > 0 {
		q = q.Where("CAST(json_extract(asset.media_meta, '$.width') AS INTEGER) <= ?", req.WidthMax)
	}
	if req.HeightMin > 0 {
		q = q.Where("CAST(json_extract(asset.media_meta, '$.height') AS INTEGER) >= ?", req.HeightMin)
	}
	if req.HeightMax > 0 {
		q = q.Where("CAST(json_extract(asset.media_meta, '$.height') AS INTEGER) <= ?", req.HeightMax)
	}

	return q
}

func (r *AssetRepo) applyListSort(q *bun.SelectQuery, sortBy, sortOrder string) *bun.SelectQuery {
	dir := "DESC"
	if strings.EqualFold(strings.TrimSpace(sortOrder), "asc") {
		dir = "ASC"
	}

	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "size":
		q = q.OrderExpr("asset.size " + dir)
	case "date", "mtime", "modified_at", "import_time":
		q = q.OrderExpr("asset.mtime " + dir)
	case "created", "created_at":
		q = q.OrderExpr("asset.created_at " + dir)
	case "type":
		q = q.OrderExpr("LOWER(asset.path) " + dir)
	default:
		// Use normalized path for deterministic "name/path" sorting.
		q = q.OrderExpr("LOWER(asset.path) " + dir)
	}
	return q.OrderExpr("asset.id ASC")
}
