package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"media-assistant-os/internal/services"
)

func TestServer_HealthAndIndex(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{
		IndexFile: func(ctx context.Context, path string, projectID string) (any, error) {
			return map[string]any{"path": path, "project_id": projectID}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	c := &http.Client{Timeout: 2 * time.Second}

	resp, err := c.Get(srv.BaseURL() + "/api/health")
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	body, _ := json.Marshal(map[string]any{"path": "C:\\test.txt", "project_id": "p1"})
	resp, err = c.Post(srv.BaseURL()+"/api/assets/index", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("index status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_UIReady(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var readyValue bool
	srv, err := Start(ctx, 0, 1, Deps{
		SetUIReady: func(ready bool) {
			readyValue = ready
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	body, _ := json.Marshal(map[string]any{"ready": true})
	resp, err := http.Post(srv.BaseURL()+"/api/internal/ui/ready", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("ui ready: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ui ready status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	if !readyValue {
		t.Fatalf("ui ready handler not called")
	}
}

func TestServer_UIReady_NoHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	body, _ := json.Marshal(map[string]any{"ready": true})
	resp, err := http.Post(srv.BaseURL()+"/api/internal/ui/ready", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("ui ready: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("ui ready status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_EventsReplay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{
		ListReplayEvents: func(ctx context.Context, sinceID int64, limit int) (any, error) {
			return []map[string]any{
				{"type": "asset_ready", "_event_id": "12"},
			}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/events/replay?since_id=10&limit=1")
	if err != nil {
		t.Fatalf("events replay: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("events replay status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_TraceHeader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := Start(ctx, 0, 1, Deps{})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	req, _ := http.NewRequest(http.MethodGet, srv.BaseURL()+"/api/health", nil)
	req.Header.Set("X-Trace-Id", "trace-test-123")
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("trace request: %v", err)
	}
	if got := resp.Header.Get("X-Trace-Id"); got != "trace-test-123" {
		t.Fatalf("trace header mismatch: %q", got)
	}
	_ = resp.Body.Close()
}

func TestServer_IdempotencyKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := 0
	srv, err := Start(ctx, 0, 1, Deps{
		IndexFile: func(ctx context.Context, path string, projectID string) (any, error) {
			calls++
			return map[string]any{"path": path, "project_id": projectID}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	body, _ := json.Marshal(map[string]any{"path": "C:\\same.txt", "project_id": "p1"})
	req1, _ := http.NewRequest(http.MethodPost, srv.BaseURL()+"/api/assets/index", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Idempotency-Key", "k1")
	client := &http.Client{Timeout: 2 * time.Second}
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("first request: %v", err)
	}
	firstBody, _ := io.ReadAll(resp1.Body)
	_ = resp1.Body.Close()

	req2, _ := http.NewRequest(http.MethodPost, srv.BaseURL()+"/api/assets/index", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Idempotency-Key", "k1")
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("second request: %v", err)
	}
	secondBody, _ := io.ReadAll(resp2.Body)
	_ = resp2.Body.Close()

	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	if string(firstBody) != string(secondBody) {
		t.Fatalf("idempotent response mismatch")
	}
}

func TestServer_Metrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{
		GetMetrics: func(ctx context.Context) (any, error) {
			return map[string]any{"ok": true}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/metrics")
	if err != nil {
		t.Fatalf("metrics: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("metrics status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_ArtifactsCRUD(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type artifact struct {
		ID        string `json:"id"`
		ProjectID string `json:"project_id"`
		Kind      string `json:"kind"`
		Name      string `json:"name"`
	}

	store := map[string]artifact{}

	srv, err := Start(ctx, 0, 1, Deps{
		EnableProFeatures: true,
		CreateArtifact: func(ctx context.Context, req services.CreateArtifactRequest) (any, error) {
			a := artifact{
				ID:        "a1",
				ProjectID: req.ProjectID,
				Kind:      req.Kind,
				Name:      req.Name,
			}
			store[a.ID] = a
			return a, nil
		},
		ListArtifacts: func(ctx context.Context, projectID string, kind string, limit int) (any, error) {
			out := []artifact{}
			for _, a := range store {
				if a.ProjectID != projectID {
					continue
				}
				if kind != "" && a.Kind != kind {
					continue
				}
				out = append(out, a)
			}
			return out, nil
		},
		GetArtifact: func(ctx context.Context, id string) (any, error) {
			a, ok := store[id]
			if !ok {
				return nil, nil
			}
			return a, nil
		},
		UpdateArtifact: func(ctx context.Context, req services.UpdateArtifactRequest) (any, error) {
			a := store[req.ID]
			if req.Kind != nil {
				a.Kind = *req.Kind
			}
			if req.Name != nil {
				a.Name = *req.Name
			}
			store[req.ID] = a
			return a, nil
		},
		DeleteArtifact: func(ctx context.Context, id string) error {
			delete(store, id)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	body, _ := json.Marshal(map[string]any{
		"project_id": "p1",
		"kind":       "script",
		"name":       "v1",
		"meta":       map[string]any{"ok": true},
	})
	resp, err := http.Post(srv.BaseURL()+"/api/artifacts/create", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("create status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp, err = http.Get(srv.BaseURL() + "/api/artifacts/list?project_id=p1&kind=script&limit=50")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	updateBody, _ := json.Marshal(map[string]any{
		"id":   "a1",
		"name": "v2",
	})
	resp, err = http.Post(srv.BaseURL()+"/api/artifacts/update", "application/json", bytes.NewReader(updateBody))
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp, err = http.Get(srv.BaseURL() + "/api/artifacts/get?id=a1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	delBody, _ := json.Marshal(map[string]any{"id": "a1"})
	resp, err = http.Post(srv.BaseURL()+"/api/artifacts/delete", "application/json", bytes.NewReader(delBody))
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_ArtifactsRouteDisabledByDefault(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/artifacts/list?project_id=p1")
	if err != nil {
		t.Fatalf("artifacts list: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 when pro routes disabled, got: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_PluginRuntimeProxy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	plugin := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/ping" {
			http.Error(w, "unexpected path", http.StatusBadRequest)
			return
		}
		if r.URL.RawQuery != "x=1" {
			http.Error(w, "unexpected query", http.StatusBadRequest)
			return
		}
		if r.Header.Get("X-Smart-Archive-Plugin-ID") != "demo.plugin" {
			http.Error(w, "missing plugin header", http.StatusBadRequest)
			return
		}
		if r.Header.Get("Authorization") != "" {
			http.Error(w, "authorization should not be forwarded", http.StatusBadRequest)
			return
		}
		w.Header().Set("X-Plugin-Reply", "ok")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("pong"))
	}))
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	plugin.Listener = ln
	plugin.Start()
	defer plugin.Close()

	srv, err := Start(ctx, 0, 1, Deps{
		ResolvePluginRuntimeEndpoint: func(ctx context.Context, pluginID string) (string, error) {
			if pluginID != "demo.plugin" {
				return "", errors.New("plugin not found")
			}
			return plugin.URL, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	req, _ := http.NewRequest(http.MethodGet, srv.BaseURL()+"/api/plugin-runtime/demo.plugin/v1/ping?x=1", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("proxy request: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("proxy status: %d body=%s", resp.StatusCode, string(body))
	}
	if string(body) != "pong" {
		t.Fatalf("proxy body mismatch: %s", string(body))
	}
	if resp.Header.Get("X-Plugin-Reply") != "ok" {
		t.Fatalf("proxy header missing")
	}
}

func TestServer_PluginRuntimeProxy_RejectsNonLoopback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{
		ResolvePluginRuntimeEndpoint: func(ctx context.Context, pluginID string) (string, error) {
			return "https://example.com/plugin", nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/plugin-runtime/demo.plugin/ping")
	if err != nil {
		t.Fatalf("proxy request: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-loopback endpoint, got: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestServer_ListAssetsQueryParsing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var captured services.ListAssetsRequest
	srv, err := Start(ctx, 0, 1, Deps{
		ListAssets: func(ctx context.Context, req services.ListAssetsRequest) (any, error) {
			captured = req
			return map[string]any{
				"items":      []any{},
				"nextCursor": nil,
				"hasMore":    false,
				"total":      0,
			}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	url := srv.BaseURL() + "/api/assets" +
		"?projectId=p1&directory=/data/assets&search=cover%20draft" +
		"&tagIds=t1,t2&types=jpg,png&shapes=landscape,square&sortBy=name&sortOrder=desc" +
		"&ratingMin=3&ratingMax=5" +
		"&cursor=100&limit=50&sizeMin=10&sizeMax=999" +
		"&mtimeFrom=1700000000&mtimeTo=1709999999" +
		"&widthMin=100&widthMax=4000&heightMin=100&heightMax=3000"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("list assets: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list assets status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	if captured.ProjectID != "p1" {
		t.Fatalf("projectId parse failed: %q", captured.ProjectID)
	}
	if captured.Directory != "/data/assets" {
		t.Fatalf("directory parse failed: %q", captured.Directory)
	}
	if captured.Query != "cover draft" {
		t.Fatalf("query parse failed: %q", captured.Query)
	}
	if len(captured.TagIDs) != 2 || captured.TagIDs[0] != "t1" || captured.TagIDs[1] != "t2" {
		t.Fatalf("tagIds parse failed: %#v", captured.TagIDs)
	}
	if len(captured.Types) != 2 || captured.Types[0] != "jpg" || captured.Types[1] != "png" {
		t.Fatalf("types parse failed: %#v", captured.Types)
	}
	if len(captured.Shapes) != 2 || captured.Shapes[0] != "landscape" || captured.Shapes[1] != "square" {
		t.Fatalf("shapes parse failed: %#v", captured.Shapes)
	}
	if captured.RatingMin != 3 || captured.RatingMax != 5 {
		t.Fatalf("rating parse failed: %d %d", captured.RatingMin, captured.RatingMax)
	}
	if captured.SortBy != "name" || captured.SortOrder != "desc" {
		t.Fatalf("sort parse failed: %q %q", captured.SortBy, captured.SortOrder)
	}
	if captured.Cursor != "100" || captured.Limit != 50 {
		t.Fatalf("pagination parse failed: cursor=%q limit=%d", captured.Cursor, captured.Limit)
	}
	if captured.SizeMin != 10 || captured.SizeMax != 999 {
		t.Fatalf("size parse failed: %d %d", captured.SizeMin, captured.SizeMax)
	}
	if captured.MtimeFrom != 1700000000 || captured.MtimeTo != 1709999999 {
		t.Fatalf("mtime parse failed: %d %d", captured.MtimeFrom, captured.MtimeTo)
	}
	if captured.WidthMin != 100 || captured.WidthMax != 4000 || captured.HeightMin != 100 || captured.HeightMax != 3000 {
		t.Fatalf("dimension parse failed: w=%d-%d h=%d-%d", captured.WidthMin, captured.WidthMax, captured.HeightMin, captured.HeightMax)
	}
}

func TestServer_UpdateAssetMeta_WithUserRating(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metaUpdated := false
	ratingUpdated := false
	var gotRating *int

	srv, err := Start(ctx, 0, 1, Deps{
		UpdateAssetMeta: func(ctx context.Context, id string, mediaMeta string) error {
			metaUpdated = id == "a1" && mediaMeta == `{"w":1}`
			return nil
		},
		SetAssetUserRating: func(ctx context.Context, id string, userRating *int) error {
			ratingUpdated = id == "a1"
			gotRating = userRating
			return nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	body, _ := json.Marshal(map[string]any{
		"id":          "a1",
		"media_meta":  `{"w":1}`,
		"user_rating": 4,
	})
	resp, err := http.Post(srv.BaseURL()+"/api/assets/update", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("update asset: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update asset status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	if !metaUpdated || !ratingUpdated {
		t.Fatalf("expected metadata and rating to be updated")
	}
	if gotRating == nil || *gotRating != 4 {
		t.Fatalf("unexpected rating: %#v", gotRating)
	}
}

func TestServer_ListAssets_InvalidCursor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := Start(ctx, 0, 1, Deps{
		ListAssets: func(ctx context.Context, req services.ListAssetsRequest) (any, error) {
			return nil, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/assets?cursor=bad_cursor")
	if err != nil {
		t.Fatalf("list assets invalid cursor: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		t.Fatalf("expected 400, got %d, body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	_ = resp.Body.Close()
}

func TestServer_ProjectSourcesAndHealthRoutes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var listedProjectID string
	var addedProjectID, addedRoot, addedType string
	var addedWatchEnabled *bool
	var removedProjectID, removedRoot string
	var healthProjectID string

	srv, err := Start(ctx, 0, 1, Deps{
		ListProjectSources: func(ctx context.Context, projectID string) (any, error) {
			listedProjectID = projectID
			return []map[string]any{{"project_id": projectID, "root_path": "/tmp/a"}}, nil
		},
		AddProjectSource: func(ctx context.Context, projectID string, rootPath string, sourceType string, watchEnabled *bool) (any, error) {
			addedProjectID = projectID
			addedRoot = rootPath
			addedType = sourceType
			addedWatchEnabled = watchEnabled
			return map[string]any{"project_id": projectID, "root_path": rootPath}, nil
		},
		RemoveProjectSource: func(ctx context.Context, projectID string, rootPath string) error {
			removedProjectID = projectID
			removedRoot = rootPath
			return nil
		},
		GetProjectHealth: func(ctx context.Context, projectID string) (any, error) {
			healthProjectID = projectID
			return map[string]any{"project_id": projectID, "total_bindings": 0}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/projects/sources?projectId=p1")
	if err != nil {
		t.Fatalf("list sources: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list sources status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if listedProjectID != "p1" {
		t.Fatalf("list project id mismatch: %q", listedProjectID)
	}

	addBody, _ := json.Marshal(map[string]any{
		"projectId":    "p1",
		"rootPath":     "/tmp/source",
		"sourceType":   "primary",
		"watchEnabled": false,
	})
	resp, err = http.Post(srv.BaseURL()+"/api/projects/sources/add", "application/json", bytes.NewReader(addBody))
	if err != nil {
		t.Fatalf("add source: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("add source status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if addedProjectID != "p1" || addedRoot != "/tmp/source" || addedType != "primary" {
		t.Fatalf("add source args mismatch: project=%q root=%q type=%q", addedProjectID, addedRoot, addedType)
	}
	if addedWatchEnabled == nil || *addedWatchEnabled {
		t.Fatalf("watchEnabled mismatch: %#v", addedWatchEnabled)
	}

	removeBody, _ := json.Marshal(map[string]any{
		"projectId": "p1",
		"rootPath":  "/tmp/source",
	})
	resp, err = http.Post(srv.BaseURL()+"/api/projects/sources/remove", "application/json", bytes.NewReader(removeBody))
	if err != nil {
		t.Fatalf("remove source: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("remove source status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if removedProjectID != "p1" || removedRoot != "/tmp/source" {
		t.Fatalf("remove source args mismatch: project=%q root=%q", removedProjectID, removedRoot)
	}

	resp, err = http.Get(srv.BaseURL() + "/api/projects/health?projectId=p1")
	if err != nil {
		t.Fatalf("project health: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("project health status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if healthProjectID != "p1" {
		t.Fatalf("health project id mismatch: %q", healthProjectID)
	}
}

func TestServer_ProjectDirectoryTreeRoutes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var boundProjectID string
	var childrenProjectID, childrenPath string
	var warningsProjectID, warningsPath string

	srv, err := Start(ctx, 0, 1, Deps{
		ListProjectBoundDirectories: func(ctx context.Context, projectID string) (any, error) {
			boundProjectID = projectID
			return []map[string]any{
				{"path": "/tmp/library", "name": "library", "is_root": true},
			}, nil
		},
		ListProjectDirectoryChildren: func(ctx context.Context, projectID string, path string) (any, error) {
			childrenProjectID = projectID
			childrenPath = path
			return []map[string]any{
				{"path": "/tmp/library/a", "name": "a", "is_root": false},
			}, nil
		},
		ListProjectDirectoryWarnings: func(ctx context.Context, projectID string, path string) (any, error) {
			warningsProjectID = projectID
			warningsPath = path
			return []map[string]any{
				{"code": "watch_permission_denied", "path": path},
			}, nil
		},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer srv.Close(context.Background())

	resp, err := http.Get(srv.BaseURL() + "/api/projects/directories/bound?projectId=p42")
	if err != nil {
		t.Fatalf("bound directories: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bound directories status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if boundProjectID != "p42" {
		t.Fatalf("bound project id mismatch: %q", boundProjectID)
	}

	resp, err = http.Get(srv.BaseURL() + "/api/projects/directories/children?projectId=p42&path=/tmp/library")
	if err != nil {
		t.Fatalf("children directories: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("children directories status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if childrenProjectID != "p42" || childrenPath != "/tmp/library" {
		t.Fatalf("children args mismatch: project=%q path=%q", childrenProjectID, childrenPath)
	}

	resp, err = http.Get(srv.BaseURL() + "/api/projects/directories/warnings?projectId=p42&path=/tmp/library")
	if err != nil {
		t.Fatalf("directory warnings: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("directory warnings status: %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
	if warningsProjectID != "p42" || warningsPath != "/tmp/library" {
		t.Fatalf("warnings args mismatch: project=%q path=%q", warningsProjectID, warningsPath)
	}
}
