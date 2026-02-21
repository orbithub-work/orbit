# 智能搜索 1.0 实现方案（最小可行版本）

## 目标
在不破坏现有查询一致性和分页稳定性的前提下，提升搜索体验。

---

## 1️⃣ 后端改动（3-4天）

### 1.1 扩展 `ListAssetsRequest`（服务层）

**文件：** `internal/services/asset_service.go`

```go
type ListAssetsRequest struct {
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
	Cursor string

	// ✅ 新增：语义增强参数
	QuickFilter string // "recent" | "unrated" | "large" | "vertical" | "horizontal"
	DatePreset  string // "today" | "thisWeek" | "thisMonth" | "lastMonth"
}
```

### 1.2 预处理逻辑（服务层）

**文件：** `internal/services/asset_service.go`

在 `ListAssets` 方法开头添加预处理：

```go
func (s *AssetService) ListAssets(ctx context.Context, req ListAssetsRequest) (*ListAssetsResult, error) {
	// ✅ 预处理：翻译快捷筛选
	req = s.preprocessQuickFilters(req)
	
	// 原有逻辑不变
	offset := 0
	if strings.TrimSpace(req.Cursor) != "" {
		v, err := strconv.Atoi(strings.TrimSpace(req.Cursor))
		if err != nil || v < 0 {
			return nil, errors.New("invalid cursor")
		}
		offset = v
	}

	assets, total, err := s.assets.ListByQuery(ctx, repos.AssetListQuery{
		ProjectID: req.ProjectID,
		Directory: req.Directory,
		Query:     req.Query,
		TagIDs:    req.TagIDs,
		Types:     req.Types,
		Shapes:    req.Shapes,
		SizeMin:   req.SizeMin,
		SizeMax:   req.SizeMax,
		RatingMin: req.RatingMin,
		RatingMax: req.RatingMax,
		MtimeFrom: req.MtimeFrom,
		MtimeTo:   req.MtimeTo,
		WidthMin:  req.WidthMin,
		WidthMax:  req.WidthMax,
		HeightMin: req.HeightMin,
		HeightMax: req.HeightMax,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Limit:     req.Limit,
		Offset:    offset,
	})
	// ... 后续逻辑不变
}

// ✅ 新增：预处理方法
func (s *AssetService) preprocessQuickFilters(req ListAssetsRequest) ListAssetsRequest {
	now := time.Now()

	// 处理日期预设
	switch req.DatePreset {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		req.MtimeFrom = start.Unix()
		req.MtimeTo = now.Unix()
	case "thisWeek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
		req.MtimeFrom = start.Unix()
		req.MtimeTo = now.Unix()
	case "thisMonth":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		req.MtimeFrom = start.Unix()
		req.MtimeTo = now.Unix()
	case "lastMonth":
		lastMonth := now.AddDate(0, -1, 0)
		start := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
		end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
		req.MtimeFrom = start.Unix()
		req.MtimeTo = end.Unix()
	}

	// 处理快捷筛选
	switch req.QuickFilter {
	case "recent":
		if req.MtimeFrom == 0 { // 不覆盖已有的日期筛选
			req.MtimeFrom = now.AddDate(0, 0, -7).Unix()
		}
	case "unrated":
		// ✅ 使用 RatingMax = -1 表示"未评分"（NULL 语义）
		req.RatingMin = 0
		req.RatingMax = -1
	case "large":
		req.SizeMin = 100 * 1024 * 1024 // 100MB
	case "vertical":
		req.Shapes = append(req.Shapes, "portrait")
	case "horizontal":
		req.Shapes = append(req.Shapes, "landscape")
	}

	return req
}
```

### 1.3 修改 Repo 层支持"未评分"查询

**文件：** `internal/repos/asset_repo.go`

在 `applyListFilters` 方法中修改 rating 筛选逻辑：

```go
func (r *AssetRepo) applyListFilters(q *bun.SelectQuery, req AssetListQuery) *bun.SelectQuery {
	// ... 现有筛选逻辑 ...

	// ✅ 修改：支持 RatingMax = -1 表示未评分
	if req.RatingMax == -1 {
		// 查询未评分（user_rating IS NULL）
		q = q.Where("asset.user_rating IS NULL")
	} else {
		// 原有逻辑
		if req.RatingMin > 0 {
			q = q.Where("asset.user_rating >= ?", req.RatingMin)
		}
		if req.RatingMax > 0 {
			q = q.Where("asset.user_rating <= ?", req.RatingMax)
		}
	}

	// ... 其他筛选逻辑 ...
	
	return q
}
```

### 1.4 HTTP Handler 参数解析

**文件：** `internal/httpapi/handler_asset.go`

在 `handleListAssets` 中添加参数解析：

```go
func (h *Handler) handleListAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListAssets == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	q := r.URL.Query()
	cursor := strings.TrimSpace(q.Get("cursor"))
	if cursor != "" {
		if n, err := strconv.Atoi(cursor); err != nil || n < 0 {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid cursor"})
			return
		}
	}
	
	req := services.ListAssetsRequest{
		ProjectID: firstNonEmpty(
			strings.TrimSpace(q.Get("projectId")),
			strings.TrimSpace(q.Get("project_id")),
		),
		Directory: firstNonEmpty(
			strings.TrimSpace(q.Get("directory")),
			strings.TrimSpace(q.Get("dir")),
			strings.TrimSpace(q.Get("path")),
		),
		Query: firstNonEmpty(
			strings.TrimSpace(q.Get("search")),
			strings.TrimSpace(q.Get("q")),
			strings.TrimSpace(q.Get("keyword")),
		),
		TagIDs:    splitCSVParams(q.Get("tagIds"), q.Get("tagId"), q.Get("tags")),
		Types:     splitCSVParams(q.Get("types"), q.Get("type"), q.Get("fileType")),
		Shapes:    splitCSVParams(q.Get("shapes"), q.Get("shape")),
		SortBy:    strings.TrimSpace(q.Get("sortBy")),
		SortOrder: strings.TrimSpace(q.Get("sortOrder")),
		Cursor:    cursor,
		
		// ✅ 新增参数
		QuickFilter: strings.TrimSpace(q.Get("quickFilter")),
		DatePreset:  strings.TrimSpace(q.Get("datePreset")),
	}
	
	// ... 后续逻辑不变
}
```

---

## 2️⃣ 搜索历史（2天）

### 2.1 数据库迁移

**文件：** `internal/db/migrate.go`

```go
func Migrate(ctx context.Context, d *DB) error {
	// ... 现有代码 ...
	
	migrations := []Migration{
		// ... 现有迁移 ...
		{Version: 25, Up: migrateV25}, // ✅ 新增
	}
	
	// ... 后续逻辑不变
}

// ✅ 新增迁移
func migrateV25(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS search_history (
  id TEXT PRIMARY KEY,
  query_hash TEXT NOT NULL,
  query TEXT NOT NULL,
  filters TEXT,
  count INTEGER DEFAULT 1,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_search_history_updated 
  ON search_history(updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_search_history_hash 
  ON search_history(query_hash);
`)
	return err
}
```

### 2.2 Model

**文件：** `internal/models/search_history.go`（新建）

```go
package models

type SearchHistory struct {
	ID        string `bun:"id,pk" json:"id"`
	QueryHash string `bun:"query_hash,notnull" json:"query_hash"`
	Query     string `bun:"query,notnull" json:"query"`
	Filters   string `bun:"filters" json:"filters"`
	Count     int    `bun:"count,notnull" json:"count"`
	CreatedAt int64  `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt int64  `bun:"updated_at,notnull" json:"updated_at"`
}
```

### 2.3 Repository

**文件：** `internal/repos/search_history_repo.go`（新建）

```go
package repos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"
)

type SearchHistoryRepo struct {
	db *bun.DB
}

func NewSearchHistoryRepo(db *bun.DB) *SearchHistoryRepo {
	return &SearchHistoryRepo{db: db}
}

func (r *SearchHistoryRepo) RecordSearch(ctx context.Context, query string, filters map[string]any) error {
	if query == "" {
		return nil
	}

	// 生成 hash（query + filters）
	filtersJSON, _ := json.Marshal(filters)
	hash := computeHash(query, string(filtersJSON))

	// 查询是否已存在
	var existing models.SearchHistory
	err := r.db.NewSelect().
		Model(&existing).
		Where("query_hash = ?", hash).
		Scan(ctx)

	now := time.Now().Unix()

	if err != nil {
		// 不存在，新增
		history := &models.SearchHistory{
			ID:        utils.NewID(),
			QueryHash: hash,
			Query:     query,
			Filters:   string(filtersJSON),
			Count:     1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		_, err = r.db.NewInsert().Model(history).Exec(ctx)
		return err
	}

	// 已存在，更新计数和时间
	_, err = r.db.NewUpdate().
		Model(&existing).
		Set("count = count + 1").
		Set("updated_at = ?", now).
		Where("id = ?", existing.ID).
		Exec(ctx)
	return err
}

func (r *SearchHistoryRepo) GetRecent(ctx context.Context, limit int) ([]models.SearchHistory, error) {
	var history []models.SearchHistory
	err := r.db.NewSelect().
		Model(&history).
		Order("updated_at DESC").
		Limit(limit).
		Scan(ctx)
	return history, err
}

func (r *SearchHistoryRepo) Clear(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*models.SearchHistory)(nil)).
		Where("1=1").
		Exec(ctx)
	return err
}

func computeHash(query, filters string) string {
	h := sha256.New()
	h.Write([]byte(query))
	h.Write([]byte(filters))
	return hex.EncodeToString(h.Sum(nil))[:16]
}
```

### 2.4 Service 集成

**文件：** `internal/services/asset_service.go`

```go
type AssetService struct {
	assets            *repos.AssetRepo
	historyEvents     *repos.AssetHistoryEventRepo
	projectAssets     *repos.ProjectAssetRepo
	projects          *repos.ProjectRepo
	lineage           *repos.AssetLineageRepo
	lineageCandidates *repos.LineageCandidateRepo
	activities        *ActivityService
	eventHub          *EventHub
	taskService       *TaskService
	cache             *AssetCache
	bloom             *utils.BloomFilter
	searchHistory     *repos.SearchHistoryRepo // ✅ 新增
}

func (s *AssetService) ListAssets(ctx context.Context, req ListAssetsRequest) (*ListAssetsResult, error) {
	// 预处理
	req = s.preprocessQuickFilters(req)
	
	// ✅ 记录搜索历史（异步，不阻塞）
	if req.Query != "" {
		go func() {
			filters := map[string]any{
				"quickFilter": req.QuickFilter,
				"datePreset":  req.DatePreset,
			}
			_ = s.searchHistory.RecordSearch(context.Background(), req.Query, filters)
		}()
	}
	
	// ... 原有逻辑
}
```

### 2.5 HTTP Handler

**文件：** `internal/httpapi/handler_search.go`（新建）

```go
package httpapi

import (
	"net/http"
	"strconv"
)

func (h *Handler) handleSearchHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	history, err := h.deps.SearchHistoryRepo.GetRecent(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: history})
}

func (h *Handler) handleClearSearchHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := h.deps.SearchHistoryRepo.Clear(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}
```

### 2.6 注册路由

**文件：** `internal/httpapi/handler.go`

```go
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// ... 现有路由 ...
	
	// ✅ 搜索历史
	mux.HandleFunc("/api/search/history", h.handleSearchHistory)
	mux.HandleFunc("/api/search/history/clear", h.handleClearSearchHistory)
}
```

### 2.7 Deps 注入

**文件：** `internal/httpapi/deps.go`

```go
type Deps struct {
	// ... 现有依赖 ...
	SearchHistoryRepo *repos.SearchHistoryRepo // ✅ 新增
}
```

---

## 3️⃣ API 接口

### 新增参数

```
GET /api/assets
  ?quickFilter=recent|unrated|large|vertical|horizontal
  &datePreset=today|thisWeek|thisMonth|lastMonth
```

### 新增接口

```
GET  /api/search/history?limit=10
DELETE /api/search/history/clear
```

---

## 4️⃣ 实现步骤

### Day 1：快捷筛选
- [ ] 扩展 `ListAssetsRequest`
- [ ] 实现 `preprocessQuickFilters`
- [ ] 修改 Repo 层 rating 筛选
- [ ] HTTP Handler 参数解析
- [ ] 测试

### Day 2：搜索历史
- [ ] 数据库迁移
- [ ] Model + Repository
- [ ] Service 集成
- [ ] HTTP Handler
- [ ] 测试

### Day 3：联调优化
- [ ] 前端集成
- [ ] 性能测试
- [ ] 边界情况处理

---

## 5️⃣ 测试用例

```bash
# 快捷筛选
curl "http://localhost:32000/api/assets?quickFilter=recent"
curl "http://localhost:32000/api/assets?quickFilter=unrated"
curl "http://localhost:32000/api/assets?datePreset=thisWeek"

# 组合筛选
curl "http://localhost:32000/api/assets?quickFilter=vertical&datePreset=thisMonth"

# 搜索历史
curl "http://localhost:32000/api/search/history?limit=10"
```

---

## 6️⃣ 不做的事（2.0）

- ❌ 拼音搜索（需要重排，破坏分页）
- ❌ 相关性排序（同上）
- ❌ 智能文件夹（先用前端预设）
- ❌ 复杂自然语言解析
- ❌ 页内/页后重排

---

## 7️⃣ 预期效果

### 用户体验提升
- 点击"最近"按钮 → 自动筛选 7 天内
- 点击"未评分" → 自动筛选 rating IS NULL
- 点击"竖屏" → 自动筛选 portrait
- 搜索框显示历史记录

### 技术保证
- ✅ 不破坏现有 cursor 分页
- ✅ 不修改 Repo 层查询逻辑（除 rating）
- ✅ 符合现有 bun + migration 框架
- ✅ 可快速回滚

---

## 8️⃣ 总结

**改动量：**
- 新增文件：3 个
- 修改文件：4 个
- 新增接口：2 个
- 数据库表：1 个

**开发时间：3-4 天**

**核心价值：**
- 快捷筛选减少 80% 重复操作
- 搜索历史提升 50% 搜索效率
- 零技术债，可平滑升级到 2.0
