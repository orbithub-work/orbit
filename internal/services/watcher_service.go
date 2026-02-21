package services

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"media-assistant-os/internal/repos"
)

type watchSession struct {
	Key       string
	ProjectID string
	RootPath  string
	ExpiresAt int64
	AddedDirs map[string]struct{}
}

type WatcherService struct {
	watcher           *fsnotify.Watcher
	assetService      *AssetService
	projectRepo       *repos.ProjectRepo
	projectSourceRepo *repos.ProjectSourceRepo
	stopChan          chan struct{}
	stopOnce          sync.Once
	mu                sync.Mutex
	watchedPaths      map[string]string // path -> projectID
	warnedPermissions map[string]struct{}
	sessionWatches    map[string]*watchSession // key(project::path) -> session
	sessionPathRefs   map[string]int           // watched path -> session ref count
}

func NewWatcherService(assetService *AssetService, projectRepo *repos.ProjectRepo, projectSourceRepo *repos.ProjectSourceRepo) (*WatcherService, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &WatcherService{
		watcher:           watcher,
		assetService:      assetService,
		projectRepo:       projectRepo,
		projectSourceRepo: projectSourceRepo,
		stopChan:          make(chan struct{}),
		watchedPaths:      make(map[string]string),
		warnedPermissions: make(map[string]struct{}),
		sessionWatches:    make(map[string]*watchSession),
		sessionPathRefs:   make(map[string]int),
	}, nil
}

func (s *WatcherService) Start(ctx context.Context) {
	projects, err := s.projectRepo.List(ctx)
	if err == nil {
		for _, p := range projects {
			roots, rerr := s.resolveWatchRoots(ctx, p.ID, p.Path)
			if rerr != nil {
				log.Printf("resolve watch roots failed: project=%s err=%v", p.ID, rerr)
				continue
			}
			for _, root := range roots {
				s.watchRootOnly(p.ID, root)
			}
		}
	}

	go func() {
		defer s.watcher.Close()
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}
				s.handleEvent(event)
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("监控器错误: %v", err)
			case <-s.stopChan:
				return
			}
		}
	}()

	go s.runSessionJanitor()
}

func (s *WatcherService) WatchProject(projectID string, rootPath string) {
	_, _ = s.watchRecursive(projectID, rootPath, "")
}

func (s *WatcherService) StartSessionWatch(ctx context.Context, projectID string, path string, ttlSeconds int) (map[string]any, error) {
	projectID = strings.TrimSpace(projectID)
	path = strings.TrimSpace(path)
	if projectID == "" || path == "" {
		return nil, errors.New("project_id and path are required")
	}
	clean := filepath.Clean(path)
	if abs, err := filepath.Abs(clean); err == nil {
		clean = abs
	}
	info, err := os.Stat(clean)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	if ttlSeconds <= 0 {
		ttlSeconds = 300
	}
	if ttlSeconds > 3600 {
		ttlSeconds = 3600
	}

	key := sessionKey(projectID, clean)
	exp := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix()

	s.mu.Lock()
	if existing, ok := s.sessionWatches[key]; ok {
		existing.ExpiresAt = exp
		added := len(existing.AddedDirs)
		s.mu.Unlock()
		return map[string]any{
			"project_id":    projectID,
			"path":          clean,
			"expires_at":    exp,
			"added_watches": added,
			"mode":          "extended",
		}, nil
	}
	s.sessionWatches[key] = &watchSession{
		Key:       key,
		ProjectID: projectID,
		RootPath:  clean,
		ExpiresAt: exp,
		AddedDirs: make(map[string]struct{}),
	}
	s.mu.Unlock()

	added, err := s.watchRecursive(projectID, clean, key)
	if err != nil {
		_ = s.StopSessionWatch(ctx, projectID, clean)
		return nil, err
	}
	return map[string]any{
		"project_id":    projectID,
		"path":          clean,
		"expires_at":    exp,
		"added_watches": added,
		"mode":          "extended",
	}, nil
}

func (s *WatcherService) StopSessionWatch(_ context.Context, projectID string, path string) error {
	projectID = strings.TrimSpace(projectID)
	path = strings.TrimSpace(path)
	if projectID == "" || path == "" {
		return nil
	}
	key := sessionKey(projectID, filepath.Clean(path))
	return s.stopSessionByKey(key)
}

func (s *WatcherService) watchRootOnly(projectID string, rootPath string) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return
	}
	if abs, err := filepath.Abs(rootPath); err == nil {
		rootPath = abs
	}
	info, err := os.Stat(rootPath)
	if err != nil {
		if isPermissionError(err) {
			s.warnPermissionDenied(projectID, rootPath, err)
		}
		return
	}
	if !info.IsDir() {
		return
	}
	s.addWatch(rootPath, projectID)
}

func (s *WatcherService) watchRecursive(projectID string, rootPath string, sessionKey string) (int, error) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return 0, nil
	}
	if abs, err := filepath.Abs(rootPath); err == nil {
		rootPath = abs
	}
	info, err := os.Stat(rootPath)
	if err != nil {
		if isPermissionError(err) {
			s.warnPermissionDenied(projectID, rootPath, err)
		}
		return 0, err
	}
	if !info.IsDir() {
		return 0, nil
	}

	addedCount := 0
	walkErr := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if isPermissionError(err) {
				s.warnPermissionDenied(projectID, path, err)
				if d != nil && d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			return filepath.SkipDir
		}

		added := s.addWatch(path, projectID)
		if added {
			addedCount++
			if sessionKey != "" {
				s.attachSessionWatchPath(sessionKey, path)
			}
		}
		return nil
	})
	if walkErr != nil {
		if isPermissionError(walkErr) {
			s.warnPermissionDenied(projectID, rootPath, walkErr)
		}
		return addedCount, walkErr
	}
	return addedCount, nil
}

func (s *WatcherService) addWatch(path string, projectID string) bool {
	s.mu.Lock()
	if _, exists := s.watchedPaths[path]; exists {
		s.mu.Unlock()
		return false
	}
	s.mu.Unlock()

	if err := s.watcher.Add(path); err != nil {
		if isPermissionError(err) {
			s.warnPermissionDenied(projectID, path, err)
		}
		log.Printf("Error watching path %s: %v", path, err)
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.watchedPaths[path]; exists {
		_ = s.watcher.Remove(path)
		return false
	}
	s.watchedPaths[path] = projectID
	return true
}

func (s *WatcherService) removeWatch(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.watchedPaths[path]; !exists {
		return
	}
	_ = s.watcher.Remove(path)
	delete(s.watchedPaths, path)
}

func (s *WatcherService) handleEvent(event fsnotify.Event) {
	base := filepath.Base(event.Name)
	if strings.HasPrefix(base, ".") || base == "Thumbs.db" {
		return
	}

	dir := filepath.Dir(event.Name)
	projectID := ""
	s.mu.Lock()
	projectID = s.watchedPaths[dir]
	if projectID == "" {
		for path, pid := range s.watchedPaths {
			if strings.HasPrefix(event.Name, path) {
				projectID = pid
				break
			}
		}
	}
	s.mu.Unlock()

	if projectID == "" {
		return
	}

	_ = s.projectRepo.UpdateActivity(context.Background(), projectID)

	if event.Has(fsnotify.Create) {
		info, err := os.Stat(event.Name)
		if err == nil && info.IsDir() {
			for _, key := range s.matchingSessionKeys(projectID, event.Name) {
				_, _ = s.watchRecursive(projectID, event.Name, key)
			}
		}
	}

	if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
		_, _ = s.assetService.IndexFile(context.Background(), IndexFileRequest{
			Path:      event.Name,
			ProjectID: projectID,
			Trigger:   "watcher",
		})
	} else if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
		s.removeWatch(event.Name)
	}
}

func (s *WatcherService) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopChan)
	})
}

func (s *WatcherService) resolveWatchRoots(ctx context.Context, projectID string, fallbackPath string) ([]string, error) {
	roots := make([]string, 0, 4)
	seen := make(map[string]struct{}, 4)

	if s.projectSourceRepo != nil {
		sources, err := s.projectSourceRepo.ListWatchEnabledByProject(ctx, projectID)
		if err != nil {
			return nil, err
		}
		for _, source := range sources {
			root := strings.TrimSpace(source.RootPath)
			if root == "" {
				continue
			}
			clean := filepath.Clean(root)
			if _, ok := seen[clean]; ok {
				continue
			}
			seen[clean] = struct{}{}
			roots = append(roots, clean)
		}
	}

	if len(roots) == 0 {
		root := strings.TrimSpace(fallbackPath)
		if root != "" {
			clean := filepath.Clean(root)
			if _, ok := seen[clean]; !ok {
				roots = append(roots, clean)
			}
		}
	}
	return roots, nil
}

func (s *WatcherService) warnPermissionDenied(projectID string, path string, err error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}

	key := projectID + "::" + path
	s.mu.Lock()
	if _, ok := s.warnedPermissions[key]; ok {
		s.mu.Unlock()
		return
	}
	s.warnedPermissions[key] = struct{}{}
	s.mu.Unlock()

	msg := "monitor permission denied: " + path
	if s.assetService != nil && s.assetService.activities != nil {
		s.assetService.activities.LogEx(context.Background(), "WARN", msg, "", projectID)
	}
	if s.assetService != nil && s.assetService.eventHub != nil {
		s.assetService.eventHub.Broadcast(map[string]any{
			"type": "system_warning",
			"data": map[string]any{
				"code":       "watch_permission_denied",
				"project_id": projectID,
				"path":       path,
				"message":    msg,
				"error":      err.Error(),
			},
		})
	}
}

func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, fs.ErrPermission) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "operation not permitted") || strings.Contains(msg, "permission denied")
}

func sessionKey(projectID string, path string) string {
	return strings.TrimSpace(projectID) + "::" + filepath.Clean(strings.TrimSpace(path))
}

func (s *WatcherService) attachSessionWatchPath(key string, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessionWatches[key]
	if !ok || sess == nil {
		return
	}
	if _, exists := sess.AddedDirs[path]; exists {
		return
	}
	sess.AddedDirs[path] = struct{}{}
	s.sessionPathRefs[path]++
}

func (s *WatcherService) matchingSessionKeys(projectID string, path string) []string {
	clean := filepath.Clean(strings.TrimSpace(path))
	now := time.Now().Unix()
	keys := make([]string, 0, 2)
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, sess := range s.sessionWatches {
		if sess == nil || sess.ProjectID != projectID {
			continue
		}
		if sess.ExpiresAt > 0 && sess.ExpiresAt < now {
			continue
		}
		if pathWithinRoot(clean, sess.RootPath) || pathWithinRoot(sess.RootPath, clean) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s *WatcherService) runSessionJanitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.expireSessions()
		case <-s.stopChan:
			return
		}
	}
}

func (s *WatcherService) expireSessions() {
	now := time.Now().Unix()
	keys := make([]string, 0, 4)
	s.mu.Lock()
	for key, sess := range s.sessionWatches {
		if sess == nil {
			keys = append(keys, key)
			continue
		}
		if sess.ExpiresAt > 0 && sess.ExpiresAt <= now {
			keys = append(keys, key)
		}
	}
	s.mu.Unlock()

	for _, key := range keys {
		_ = s.stopSessionByKey(key)
	}
}

func (s *WatcherService) stopSessionByKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	s.mu.Lock()
	sess, ok := s.sessionWatches[key]
	if !ok || sess == nil {
		s.mu.Unlock()
		return nil
	}
	delete(s.sessionWatches, key)
	paths := make([]string, 0, len(sess.AddedDirs))
	for path := range sess.AddedDirs {
		paths = append(paths, path)
	}
	s.mu.Unlock()

	for _, path := range paths {
		remove := false
		s.mu.Lock()
		ref := s.sessionPathRefs[path]
		if ref <= 1 {
			delete(s.sessionPathRefs, path)
			remove = true
		} else {
			s.sessionPathRefs[path] = ref - 1
		}
		s.mu.Unlock()
		if remove {
			s.removeWatch(path)
		}
	}
	return nil
}

func pathWithinRoot(path string, root string) bool {
	path = filepath.Clean(strings.TrimSpace(path))
	root = filepath.Clean(strings.TrimSpace(root))
	if path == "" || root == "" {
		return false
	}
	if path == root {
		return true
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}
