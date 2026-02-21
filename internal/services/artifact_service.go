package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/repos"
	"media-assistant-os/internal/utils"
)

type ArtifactService struct {
	projects  *repos.ProjectRepo
	artifacts *repos.ProjectArtifactRepo
}

func NewArtifactService(projects *repos.ProjectRepo, artifacts *repos.ProjectArtifactRepo) *ArtifactService {
	return &ArtifactService{
		projects:  projects,
		artifacts: artifacts,
	}
}

func (s *ArtifactService) List(ctx context.Context, projectID string, kind string, limit int) ([]models.ProjectArtifact, error) {
	if projectID == "" {
		return nil, errors.New("project_id is required")
	}
	out, err := s.artifacts.ListByProject(ctx, projectID, kind, limit)
	if err != nil {
		return nil, err
	}
	for i := range out {
		attachArtifactMeta(&out[i])
	}
	return out, nil
}

func (s *ArtifactService) Get(ctx context.Context, id string) (*models.ProjectArtifact, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	a, err := s.artifacts.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, nil
	}
	attachArtifactMeta(a)
	return a, nil
}

type CreateArtifactRequest struct {
	ProjectID      string          `json:"project_id"`
	Kind           string          `json:"kind"`
	Name           string          `json:"name"`
	Path           string          `json:"path"`
	Content        string          `json:"content"`
	Meta           json.RawMessage `json:"meta"`
	SourcePluginID string          `json:"source_plugin_id"`
}

func (s *ArtifactService) Create(ctx context.Context, req CreateArtifactRequest) (*models.ProjectArtifact, error) {
	if req.ProjectID == "" {
		return nil, errors.New("project_id is required")
	}
	if req.Kind == "" {
		return nil, errors.New("kind is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if s.projects != nil {
		p, err := s.projects.Get(ctx, req.ProjectID)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, errors.New("project not found")
		}
	}

	metaJSON := "{}"
	if len(req.Meta) > 0 {
		if !json.Valid(req.Meta) {
			return nil, errors.New("meta must be valid json")
		}
		metaJSON = string(req.Meta)
	}

	now := time.Now().Unix()
	a := &models.ProjectArtifact{
		ID:             utils.NewID(),
		ProjectID:      req.ProjectID,
		Kind:           req.Kind,
		Name:           req.Name,
		Path:           req.Path,
		Content:        req.Content,
		MetaJSON:       metaJSON,
		SourcePluginID: req.SourcePluginID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.artifacts.Create(ctx, a); err != nil {
		return nil, err
	}
	attachArtifactMeta(a)
	return a, nil
}

type UpdateArtifactRequest struct {
	ID             string          `json:"id"`
	Kind           *string         `json:"kind,omitempty"`
	Name           *string         `json:"name,omitempty"`
	Path           *string         `json:"path,omitempty"`
	Content        *string         `json:"content,omitempty"`
	Meta           json.RawMessage `json:"meta,omitempty"`
	SourcePluginID *string         `json:"source_plugin_id,omitempty"`
}

func (s *ArtifactService) Update(ctx context.Context, req UpdateArtifactRequest) (*models.ProjectArtifact, error) {
	if req.ID == "" {
		return nil, errors.New("id is required")
	}
	a, err := s.artifacts.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("artifact not found")
	}

	if req.Kind != nil {
		if *req.Kind == "" {
			return nil, errors.New("kind cannot be empty")
		}
		a.Kind = *req.Kind
	}
	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		a.Name = *req.Name
	}
	if req.Path != nil {
		a.Path = *req.Path
	}
	if req.Content != nil {
		a.Content = *req.Content
	}
	if req.SourcePluginID != nil {
		a.SourcePluginID = *req.SourcePluginID
	}
	if len(req.Meta) > 0 {
		if !json.Valid(req.Meta) {
			return nil, errors.New("meta must be valid json")
		}
		a.MetaJSON = string(req.Meta)
	}

	a.UpdatedAt = time.Now().Unix()
	if err := s.artifacts.Update(ctx, a); err != nil {
		return nil, err
	}
	attachArtifactMeta(a)
	return a, nil
}

func (s *ArtifactService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.artifacts.Delete(ctx, id)
}

func attachArtifactMeta(a *models.ProjectArtifact) {
	if a == nil {
		return
	}
	if a.MetaJSON == "" {
		a.MetaJSON = "{}"
	}
	if a.MetaJSON == "{}" {
		a.Meta = nil
		return
	}
	var raw json.RawMessage
	if json.Unmarshal([]byte(a.MetaJSON), &raw) == nil {
		a.Meta = raw
	}
}
