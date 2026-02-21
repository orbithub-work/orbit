package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/processor"
	"media-assistant-os/internal/repos"
	"media-assistant-os/internal/utils"
)

type registeredPlugin struct {
	PluginInfo
	Token string
}

type PluginService struct {
	runtimeRepo *repos.PluginRuntimeRepo
	plugins     map[string]registeredPlugin
	tokenIndex  map[string]string // token -> pluginID
	tokenTTL    time.Duration
	onlineTTL   time.Duration
	mu          sync.Mutex
}

func NewPluginService(runtimeRepo *repos.PluginRuntimeRepo) *PluginService {
	const defaultTokenTTL = 90 * 24 * time.Hour
	const defaultOnlineTTL = 2 * time.Minute
	return &PluginService{
		runtimeRepo: runtimeRepo,
		plugins:     map[string]registeredPlugin{},
		tokenIndex:  map[string]string{},
		tokenTTL:    defaultTokenTTL,
		onlineTTL:   defaultOnlineTTL,
	}
}

// Restore rebuilds in-memory token and plugin indices from persisted runtime rows.
func (s *PluginService) Restore(ctx context.Context) error {
	if s.runtimeRepo == nil {
		return nil
	}
	if _, err := s.runtimeRepo.DeleteExpired(ctx, time.Now()); err != nil {
		return err
	}
	rows, err := s.runtimeRepo.List(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	nextPlugins := make(map[string]registeredPlugin, len(rows))
	nextTokenIndex := make(map[string]string, len(rows))
	for _, row := range rows {
		plugin, ok := decodeRuntimeRow(row)
		if !ok {
			continue
		}
		if !plugin.ExpiresAt.IsZero() && now.After(plugin.ExpiresAt) {
			continue
		}
		plugin.Online = now.Sub(plugin.LastUsedAt) <= s.onlineTTL
		nextPlugins[plugin.PluginID] = registeredPlugin{
			PluginInfo: plugin,
			Token:      row.Token,
		}
		nextTokenIndex[row.Token] = plugin.PluginID
	}

	s.mu.Lock()
	s.plugins = nextPlugins
	s.tokenIndex = nextTokenIndex
	s.mu.Unlock()
	return nil
}

func (s *PluginService) Register(ctx context.Context, req PluginRegistrationRequest) (*PluginRegistrationResponse, error) {
	req.PluginID = strings.TrimSpace(req.PluginID)
	req.Name = strings.TrimSpace(req.Name)
	req.Version = strings.TrimSpace(req.Version)
	req.Description = strings.TrimSpace(req.Description)
	req.Executable = strings.TrimSpace(req.Executable)
	req.Endpoint = strings.TrimSpace(req.Endpoint)
	req.ProtocolVersion = strings.TrimSpace(req.ProtocolVersion)
	req.Permissions = normalizeNonEmptyStrings(req.Permissions)
	req.TaskTypes = normalizeNonEmptyStrings(req.TaskTypes)
	req.Extensions = normalizeNonEmptyStrings(req.Extensions)
	req.Mounts = normalizePluginMounts(req.Mounts, req.UI)

	if req.PluginID == "" {
		return nil, errors.New("plugin_id is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Mode == "" {
		return nil, errors.New("mode is required")
	}
	if req.Mode != PluginModeLocalProcess && req.Mode != PluginModeNetworkWorker && req.Mode != PluginModeFrontend {
		return nil, fmt.Errorf("unsupported mode: %s", req.Mode)
	}
	if req.Mode == PluginModeLocalProcess {
		if req.Executable == "" {
			return nil, errors.New("executable is required for local_process mode")
		}
		if req.Endpoint != "" {
			return nil, errors.New("endpoint is not allowed for local_process mode")
		}
	}
	if req.Mode == PluginModeNetworkWorker {
		if req.Endpoint == "" {
			return nil, errors.New("endpoint is required for network_service mode")
		}
		if req.Executable != "" {
			return nil, errors.New("executable is not allowed for network_service mode")
		}
	}
	if req.Mode == PluginModeFrontend {
		if req.Endpoint != "" {
			return nil, errors.New("endpoint is not allowed for frontend mode")
		}
		if req.Executable != "" {
			return nil, errors.New("executable is not allowed for frontend mode")
		}
		if len(req.Mounts) == 0 && (req.UI == nil || strings.TrimSpace(req.UI.Entry) == "") {
			return nil, errors.New("mounts or ui.entry is required for frontend mode")
		}
	}

	if req.Mode != PluginModeFrontend && len(req.TaskTypes) == 0 && len(req.Capabilities) == 0 {
		return nil, errors.New("task_types or capabilities is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// local_process + extensions => register as parser directly in process manager
	if req.Mode == PluginModeLocalProcess && req.Executable != "" && len(req.Extensions) > 0 {
		extParser := processor.NewExternalCommandParser(req.PluginID, req.Executable, req.Extensions)
		// Register with high priority (100) so external plugins can override native ones
		processor.GetManager().Register(extParser, 100)
	}

	now := time.Now()
	res := PluginRegistrationResponse{
		PluginID: req.PluginID,
		Token:    utils.NewID(),
		Mode:     string(req.Mode),
	}
	if old, ok := s.plugins[req.PluginID]; ok && old.Token != "" {
		delete(s.tokenIndex, old.Token)
	}

	registered := registeredPlugin{
		Token: res.Token,
		PluginInfo: PluginInfo{
			PluginID:        req.PluginID,
			ID:              req.PluginID,
			Name:            req.Name,
			Version:         req.Version,
			Description:     req.Description,
			Mode:            req.Mode,
			Endpoint:        req.Endpoint,
			Executable:      req.Executable,
			Extensions:      req.Extensions,
			TaskTypes:       req.TaskTypes,
			Capabilities:    req.Capabilities,
			ProtocolVersion: req.ProtocolVersion,
			Permissions:     req.Permissions,
			UI:              clonePluginUI(req.UI),
			Mounts:          req.Mounts,
			IssuedAt:        now,
			LastUsedAt:      now,
			ExpiresAt:       now.Add(s.tokenTTL),
			Online:          true,
			RegisteredAt:    now.Unix(),
			LastHeartbeat:   now.Unix(),
		},
	}
	if err := s.persistRegisteredPlugin(ctx, registered); err != nil {
		return nil, err
	}
	s.plugins[req.PluginID] = registered
	s.tokenIndex[res.Token] = req.PluginID
	return &res, nil
}

func (s *PluginService) ValidateToken(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	pluginID, ok := s.tokenIndex[token]
	if !ok {
		return false
	}
	plugin, ok := s.plugins[pluginID]
	if !ok {
		delete(s.tokenIndex, token)
		return false
	}
	if time.Now().After(plugin.ExpiresAt) {
		delete(s.tokenIndex, token)
		delete(s.plugins, pluginID)
		if s.runtimeRepo != nil {
			_ = s.runtimeRepo.Delete(context.Background(), pluginID)
		}
		return false
	}
	plugin.LastUsedAt = time.Now()
	plugin.LastHeartbeat = plugin.LastUsedAt.Unix()
	plugin.Online = true
	plugin.ID = plugin.PluginID
	ensureLegacyPluginUI(&plugin.PluginInfo)
	_ = s.persistRegisteredPlugin(context.Background(), plugin)
	s.plugins[pluginID] = plugin
	return true
}

func (s *PluginService) Heartbeat(ctx context.Context, req PluginHeartbeatRequest) (*PluginInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pluginID := strings.TrimSpace(req.PluginID)
	token := strings.TrimSpace(req.Token)

	if pluginID == "" && token == "" {
		return nil, errors.New("plugin_id or token is required")
	}
	if pluginID == "" {
		mapped, ok := s.tokenIndex[token]
		if !ok {
			return nil, errors.New("plugin not found for token")
		}
		pluginID = mapped
	}

	plugin, ok := s.plugins[pluginID]
	if !ok {
		return nil, errors.New("plugin not found")
	}
	if token != "" && token != plugin.Token {
		return nil, errors.New("token mismatch")
	}

	now := time.Now()
	plugin.LastUsedAt = now
	plugin.LastHeartbeat = now.Unix()
	plugin.Online = true
	plugin.ID = plugin.PluginID
	ensureLegacyPluginUI(&plugin.PluginInfo)
	if err := s.persistRegisteredPlugin(ctx, plugin); err != nil {
		return nil, err
	}
	s.plugins[pluginID] = plugin

	info := plugin.PluginInfo
	return &info, nil
}

func (s *PluginService) List(ctx context.Context) ([]PluginInfo, error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	out := make([]PluginInfo, 0, len(s.plugins))
	for id, plugin := range s.plugins {
		plugin.Online = now.Sub(plugin.LastUsedAt) <= s.onlineTTL
		plugin.ID = plugin.PluginID
		if plugin.RegisteredAt == 0 {
			plugin.RegisteredAt = plugin.IssuedAt.Unix()
		}
		if plugin.LastHeartbeat == 0 {
			plugin.LastHeartbeat = plugin.LastUsedAt.Unix()
		}
		ensureLegacyPluginUI(&plugin.PluginInfo)
		s.plugins[id] = plugin
		out = append(out, plugin.PluginInfo)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].PluginID < out[j].PluginID
	})
	return out, nil
}

func (s *PluginService) ListTaskTypes(ctx context.Context) ([]string, error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	seen := map[string]struct{}{}
	for _, plugin := range s.plugins {
		for _, taskType := range plugin.TaskTypes {
			taskType = strings.TrimSpace(taskType)
			if taskType != "" {
				seen[taskType] = struct{}{}
			}
		}
		for _, cap := range plugin.Capabilities {
			for _, taskType := range cap.TaskTypes {
				taskType = strings.TrimSpace(taskType)
				if taskType != "" {
					seen[taskType] = struct{}{}
				}
			}
		}
	}

	out := make([]string, 0, len(seen))
	for taskType := range seen {
		out = append(out, taskType)
	}
	sort.Strings(out)
	return out, nil
}

func (s *PluginService) ListMountedSlots(ctx context.Context) ([]PluginMountedSlot, error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	out := make([]PluginMountedSlot, 0, 32)
	for id, plugin := range s.plugins {
		plugin.Online = now.Sub(plugin.LastUsedAt) <= s.onlineTTL
		s.plugins[id] = plugin
		for _, mount := range plugin.Mounts {
			if strings.TrimSpace(mount.Slot) == "" {
				continue
			}
			out = append(out, PluginMountedSlot{
				PluginID:   plugin.PluginID,
				PluginName: plugin.Name,
				Online:     plugin.Online,
				Mount:      mount,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Mount.Slot == out[j].Mount.Slot {
			if out[i].Mount.Order == out[j].Mount.Order {
				return out[i].PluginID < out[j].PluginID
			}
			return out[i].Mount.Order < out[j].Mount.Order
		}
		return out[i].Mount.Slot < out[j].Mount.Slot
	})
	return out, nil
}

func (s *PluginService) Stats() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	expiringSoon := 0
	online := 0
	mounts := 0
	for id, plugin := range s.plugins {
		if plugin.ExpiresAt.Sub(now) <= 7*24*time.Hour {
			expiringSoon++
		}
		plugin.Online = now.Sub(plugin.LastUsedAt) <= s.onlineTTL
		if plugin.Online {
			online++
		}
		mounts += len(plugin.Mounts)
		s.plugins[id] = plugin
	}
	return map[string]any{
		"registered_plugins": len(s.plugins),
		"active_tokens":      len(s.tokenIndex),
		"online_plugins":     online,
		"expiring_7d":        expiringSoon,
		"registered_mounts":  mounts,
	}
}

// ResolveRuntimeEndpoint returns the runtime endpoint for a network-service plugin.
func (s *PluginService) ResolveRuntimeEndpoint(ctx context.Context, pluginID string) (string, error) {
	_ = ctx
	pluginID = strings.TrimSpace(pluginID)
	if pluginID == "" {
		return "", errors.New("plugin_id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	plugin, ok := s.plugins[pluginID]
	if !ok {
		return "", errors.New("plugin not found")
	}
	if plugin.Mode != PluginModeNetworkWorker {
		return "", errors.New("plugin mode does not support runtime proxy")
	}
	endpoint := strings.TrimSpace(plugin.Endpoint)
	if endpoint == "" {
		return "", errors.New("plugin endpoint is empty")
	}
	return endpoint, nil
}

func decodeRuntimeRow(row models.PluginRuntime) (PluginInfo, bool) {
	infoRaw := strings.TrimSpace(row.InfoJSON)
	if infoRaw == "" {
		return PluginInfo{}, false
	}
	var info PluginInfo
	if err := json.Unmarshal([]byte(infoRaw), &info); err != nil {
		return PluginInfo{}, false
	}
	info.PluginID = strings.TrimSpace(info.PluginID)
	if info.PluginID == "" {
		info.PluginID = strings.TrimSpace(row.PluginID)
	}
	if info.PluginID == "" {
		return PluginInfo{}, false
	}
	info.ID = info.PluginID
	info.Name = strings.TrimSpace(info.Name)
	info.Version = strings.TrimSpace(info.Version)
	info.Description = strings.TrimSpace(info.Description)
	info.ProtocolVersion = strings.TrimSpace(info.ProtocolVersion)
	info.Permissions = normalizeNonEmptyStrings(info.Permissions)
	info.TaskTypes = normalizeNonEmptyStrings(info.TaskTypes)
	info.Extensions = normalizeNonEmptyStrings(info.Extensions)
	info.Mounts = normalizePluginMounts(info.Mounts, info.UI)
	info.UI = clonePluginUI(info.UI)
	info.IssuedAt = unixToTime(row.IssuedAt, info.IssuedAt)
	info.LastUsedAt = unixToTime(row.LastUsedAt, info.LastUsedAt)
	info.ExpiresAt = unixToTime(row.ExpiresAt, info.ExpiresAt)
	if info.RegisteredAt == 0 && !info.IssuedAt.IsZero() {
		info.RegisteredAt = info.IssuedAt.Unix()
	}
	if info.LastHeartbeat == 0 && !info.LastUsedAt.IsZero() {
		info.LastHeartbeat = info.LastUsedAt.Unix()
	}
	ensureLegacyPluginUI(&info)
	return info, true
}

func unixToTime(ts int64, fallback time.Time) time.Time {
	if ts > 0 {
		return time.Unix(ts, 0)
	}
	return fallback
}

func (s *PluginService) persistRegisteredPlugin(ctx context.Context, plugin registeredPlugin) error {
	if s.runtimeRepo == nil {
		return nil
	}
	info := plugin.PluginInfo
	info.ID = info.PluginID
	info.Permissions = normalizeNonEmptyStrings(info.Permissions)
	info.TaskTypes = normalizeNonEmptyStrings(info.TaskTypes)
	info.Extensions = normalizeNonEmptyStrings(info.Extensions)
	info.Mounts = normalizePluginMounts(info.Mounts, info.UI)
	info.UI = clonePluginUI(info.UI)
	ensureLegacyPluginUI(&info)
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}
	row := models.PluginRuntime{
		PluginID:   info.PluginID,
		Token:      plugin.Token,
		InfoJSON:   string(payload),
		IssuedAt:   info.IssuedAt.Unix(),
		LastUsedAt: info.LastUsedAt.Unix(),
		ExpiresAt:  info.ExpiresAt.Unix(),
	}
	return s.runtimeRepo.Upsert(ctx, row)
}

func normalizeNonEmptyStrings(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		v := strings.TrimSpace(item)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func normalizePluginMounts(mounts []PluginMount, ui *PluginUIConfig) []PluginMount {
	out := make([]PluginMount, 0, len(mounts)+1)
	for _, mount := range mounts {
		slot := strings.TrimSpace(mount.Slot)
		entry := strings.TrimSpace(mount.Entry)
		if slot == "" || entry == "" {
			continue
		}
		mount.Slot = slot
		mount.Entry = entry
		mount.Title = strings.TrimSpace(mount.Title)
		mount.Icon = strings.TrimSpace(mount.Icon)
		mount.Location = strings.TrimSpace(mount.Location)
		out = append(out, mount)
	}

	// Backward compatibility: if only legacy ui is provided, auto-map to a generic slot.
	if len(out) == 0 && ui != nil && strings.TrimSpace(ui.Entry) != "" {
		out = append(out, PluginMount{
			Slot:     "plugin.page",
			Entry:    strings.TrimSpace(ui.Entry),
			Width:    ui.Width,
			Height:   ui.Height,
			Location: strings.TrimSpace(ui.Location),
		})
	}
	return out
}

func clonePluginUI(ui *PluginUIConfig) *PluginUIConfig {
	if ui == nil {
		return nil
	}
	out := *ui
	out.Entry = strings.TrimSpace(out.Entry)
	out.Location = strings.TrimSpace(out.Location)
	return &out
}

func ensureLegacyPluginUI(info *PluginInfo) {
	if info == nil {
		return
	}
	if strings.TrimSpace(info.ID) == "" {
		info.ID = info.PluginID
	}
	if info.RegisteredAt == 0 && !info.IssuedAt.IsZero() {
		info.RegisteredAt = info.IssuedAt.Unix()
	}
	if info.LastHeartbeat == 0 && !info.LastUsedAt.IsZero() {
		info.LastHeartbeat = info.LastUsedAt.Unix()
	}
	if info.UI != nil {
		return
	}
	for _, mount := range info.Mounts {
		if strings.TrimSpace(mount.Entry) == "" {
			continue
		}
		info.UI = &PluginUIConfig{
			Entry:    mount.Entry,
			Width:    mount.Width,
			Height:   mount.Height,
			Location: mount.Location,
		}
		return
	}
}
