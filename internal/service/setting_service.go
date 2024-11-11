package service

import (
	"context"
	"log"

	"github.com/gofrs/uuid/v5"

	"github.com/Intiqo/app-platform/internal/domain"
)

type appSettingService struct {
	tr domain.Transactioner
	r  domain.SettingRepository
}

// NewSettingService creates a new setting service
func NewSettingService(tr domain.Transactioner, r domain.SettingRepository) domain.SettingService {
	s := &appSettingService{
		tr: tr,

		r: r,
	}
	s.createDefaultSettings()
	return s
}

func (s *appSettingService) FindByID(id uuid.UUID) (result domain.Setting, err error) {
	return s.r.FindByID(context.TODO(), id)
}

func (s *appSettingService) Filter(in domain.FilterSettingsByCriteriaInput, options domain.QueryOptions) (result []domain.Setting, total int64, err error) {
	return s.r.Filter(context.TODO(), in, options)
}

// Creates default settings in the system
func (s *appSettingService) createDefaultSettings() {
	settingsToCreate := make([]*domain.Setting, 0)
	existingSettings, _, err := s.r.Filter(context.TODO(), domain.FilterSettingsByCriteriaInput{}, domain.QueryOptions{})
	if err != nil {
		log.Fatalf("Error getting existing settings: %v", err)
	}

	// Check if default settings already exist
	for k, v := range domain.SettingDefaultValues {
		settingFound := false
		for _, s := range existingSettings {
			if s.Key == k {
				settingFound = true
				break
			}
		}
		if !settingFound {
			settingsToCreate = append(settingsToCreate, &domain.Setting{
				Key:   k,
				Value: v,
			})
		}
	}

	// Return if no settings to create
	if len(settingsToCreate) == 0 {
		return
	}

	// Create default settings
	log.Println("Creating default settings...")
	err = s.r.CreateMultiple(context.TODO(), settingsToCreate)
	if err != nil {
		log.Fatalf("Error creating default settings: %v", err)
	}
}
