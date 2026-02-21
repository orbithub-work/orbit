package services

import (
	"context"
)

type LicenseType string

const (
	LicenseTypeFree LicenseType = "free"
	LicenseTypePro  LicenseType = "pro"
)

type License struct {
	Type        LicenseType `json:"type"`
	MaxProjects int         `json:"max_projects"`
}

// LicenseService is responsible for managing the user's license.
// In this initial version, it's a mock that always returns a Free license.
type LicenseService struct {
}

func NewLicenseService() *LicenseService {
	return &LicenseService{}
}

// GetLicense returns the current license information.
// TODO: Implement server-side validation and local lease caching.
func (s *LicenseService) GetLicense(ctx context.Context) (*License, error) {
	// For now, always return the Free license.
	return &License{
		Type:        LicenseTypeFree,
		MaxProjects: 3,
	}, nil
}
