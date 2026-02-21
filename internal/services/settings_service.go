package services

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/uptrace/bun"
)

type SystemSettings struct {
	FirstLaunchCompleted bool              `json:"first_launch_completed"`
	Theme                string            `json:"theme"`
	Language             string            `json:"language"`
	Window               WindowSettings    `json:"window"`
	RecentProjects       []string          `json:"recent_projects"`
	Custom               map[string]string `json:"custom"` // Extensible
}

type WindowSettings struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type SystemSettingRow struct {
	bun.BaseModel `bun:"table:system_settings"`
	Key           string `bun:"key,pk"`
	Value         string `bun:"value"`
}

type SettingsService struct {
	db       *bun.DB
	mu       sync.RWMutex
	Settings SystemSettings
}

func NewSettingsService(db *bun.DB) *SettingsService {
	s := &SettingsService{
		db: db,
		Settings: SystemSettings{
			FirstLaunchCompleted: false,
			Theme:                "system",
			Language:             "zh-CN",
			Window: WindowSettings{
				Width:  1200,
				Height: 800,
			},
			Custom: make(map[string]string),
		},
	}
	// Try to load
	_ = s.Load(context.Background())
	return s
}

func (s *SettingsService) Load(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var rows []SystemSettingRow
	if err := s.db.NewSelect().Model(&rows).Scan(ctx); err != nil {
		return err
	}

	// Map DB rows to Struct
	for _, row := range rows {
		switch row.Key {
		case "first_launch_completed":
			s.Settings.FirstLaunchCompleted = row.Value == "true"
		case "theme":
			s.Settings.Theme = row.Value
		case "language":
			s.Settings.Language = row.Value
		case "window":
			_ = json.Unmarshal([]byte(row.Value), &s.Settings.Window)
		case "recent_projects":
			_ = json.Unmarshal([]byte(row.Value), &s.Settings.RecentProjects)
		default:
			s.Settings.Custom[row.Key] = row.Value
		}
	}
	return nil
}

func (s *SettingsService) Save(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prepare rows
	rows := []SystemSettingRow{
		{Key: "first_launch_completed", Value: "false"},
		{Key: "theme", Value: s.Settings.Theme},
		{Key: "language", Value: s.Settings.Language},
	}
	if s.Settings.FirstLaunchCompleted {
		rows[0].Value = "true"
	}

	// Window
	winData, _ := json.Marshal(s.Settings.Window)
	rows = append(rows, SystemSettingRow{Key: "window", Value: string(winData)})

	// RecentProjects
	projData, _ := json.Marshal(s.Settings.RecentProjects)
	rows = append(rows, SystemSettingRow{Key: "recent_projects", Value: string(projData)})

	// Custom
	for k, v := range s.Settings.Custom {
		rows = append(rows, SystemSettingRow{Key: k, Value: v})
	}

	// Upsert
	if _, err := s.db.NewInsert().Model(&rows).On("CONFLICT (key) DO UPDATE SET value = EXCLUDED.value").Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *SettingsService) IsFirstLaunch() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.Settings.FirstLaunchCompleted
}

func (s *SettingsService) SetFirstLaunchCompleted(completed bool) error {
	s.mu.Lock()
	s.Settings.FirstLaunchCompleted = completed
	s.mu.Unlock()
	return s.Save(context.Background())
}

func (s *SettingsService) Update(updater func(*SystemSettings)) error {
	s.mu.Lock()
	updater(&s.Settings)
	s.mu.Unlock()
	return s.Save(context.Background())
}
