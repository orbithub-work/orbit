package services

import (
	"context"
	"fmt"
)

// ExtensionSlot describes a stable host UI mount point.
type ExtensionSlot struct {
	Key         string `json:"key"`
	Surface     string `json:"surface"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Multiple    bool   `json:"multiple"`
}

type CapabilitySnapshot struct {
	API      map[string]any  `json:"api"`
	Features map[string]bool `json:"features"`
	Limits   map[string]int  `json:"limits"`
	Plugin   map[string]any  `json:"plugin"`
}

type CapabilityService struct {
	licenseService *LicenseService
	pluginService  *PluginService
}

func NewCapabilityService(licenseService *LicenseService, pluginService *PluginService) *CapabilityService {
	return &CapabilityService{
		licenseService: licenseService,
		pluginService:  pluginService,
	}
}

func (s *CapabilityService) Get(ctx context.Context) (*CapabilitySnapshot, error) {
	licenseType := LicenseTypeFree
	maxProjects := 3
	if s.licenseService != nil {
		license, err := s.licenseService.GetLicense(ctx)
		if err != nil {
			return nil, err
		}
		if license != nil {
			licenseType = license.Type
			if license.MaxProjects > 0 {
				maxProjects = license.MaxProjects
			}
		}
	}

	features := map[string]bool{
		"core.assets.index":             true,
		"core.assets.search":            true,
		"core.projects.manage":          true,
		"core.logs.activity":            true,
		"core.plugins.runtime":          true,
		"core.plugins.mounts":           true,
		"analytics.basic":               true,
		"analytics.advanced":            licenseType == LicenseTypePro,
		"artifacts.advanced_workflows":  licenseType == LicenseTypePro,
		"automation.batch_operations":   licenseType == LicenseTypePro,
		"smart_search.advanced_filters": true,
	}

	pluginStats := map[string]any{
		"registered_plugins": 0,
		"registered_mounts":  0,
	}
	if s.pluginService != nil {
		pluginStats = s.pluginService.Stats()
	}

	return &CapabilitySnapshot{
		API: map[string]any{
			"versioning":       "path_stable",
			"compatibility":    "additive",
			"protocol_version": "plugin-runtime.v1",
		},
		Features: features,
		Limits: map[string]int{
			"projects.max": maxProjects,
		},
		Plugin: pluginStats,
	}, nil
}

func (s *CapabilityService) ListExtensionSlots(ctx context.Context) ([]ExtensionSlot, error) {
	_ = ctx
	return []ExtensionSlot{
		{
			Key:         "pool.sidebar.bottom",
			Surface:     "pool",
			DisplayName: "素材库侧栏底部",
			Description: "在素材库侧栏底部挂载插件入口或小组件。",
			Multiple:    true,
		},
		{
			Key:         "pool.toolbar.right",
			Surface:     "pool",
			DisplayName: "素材库工具栏右侧",
			Description: "在素材库顶部工具栏挂载轻量操作按钮。",
			Multiple:    true,
		},
		{
			Key:         "project.panel.right",
			Surface:     "project",
			DisplayName: "项目库右侧面板",
			Description: "在项目详情右侧挂载项目相关插件视图。",
			Multiple:    true,
		},
		{
			Key:         "artifact.panel.right",
			Surface:     "artifact",
			DisplayName: "交付库右侧面板",
			Description: "在交付库挂载发布、同步或校验插件视图。",
			Multiple:    true,
		},
		{
			Key:         "analytics.card.extra",
			Surface:     "analytics",
			DisplayName: "数据看板扩展卡片",
			Description: "在数据看板追加统计卡片或趋势图。",
			Multiple:    true,
		},
		{
			Key:         "settings.plugins.section",
			Surface:     "settings",
			DisplayName: "设置-插件分区",
			Description: "在设置页挂载插件配置项。",
			Multiple:    true,
		},
		{
			Key:         "plugin.page",
			Surface:     "global",
			DisplayName: "插件独立页面",
			Description: "插件通过独立页面承载复杂 UI。",
			Multiple:    true,
		},
	}, nil
}

func (s *CapabilityService) ValidateMounts(ctx context.Context, mounts []PluginMount) error {
	knownSlots, err := s.ListExtensionSlots(ctx)
	if err != nil {
		return err
	}
	allowed := make(map[string]struct{}, len(knownSlots))
	for _, slot := range knownSlots {
		allowed[slot.Key] = struct{}{}
	}
	for _, mount := range mounts {
		if _, ok := allowed[mount.Slot]; ok {
			continue
		}
		return fmt.Errorf("unknown mount slot: %s", mount.Slot)
	}
	return nil
}
