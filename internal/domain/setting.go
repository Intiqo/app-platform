package domain

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type (
	// Setting defines model for Setting.
	Setting struct {
		Base
		Key   string `db:"key" json:"key,omitempty" example:"app.name"`
		Value string `db:"value" json:"value,omitempty" example:"App"`
		Audit
	} // @name Setting
)

type (
	// FilterSettingsByCriteriaInput defines the input for filtering settings by criteria.
	FilterSettingsByCriteriaInput struct {
		Keys []string `json:"keys,omitempty" example:"app.name"`
	} // @name FilterSettingsByCriteriaInput
)

type (
	// SettingRepository defines the setting repository
	SettingRepository interface {
		// FindByID finds a setting by its ID.
		FindByID(ctx context.Context, id uuid.UUID) (result Setting, err error)
		// Filter filters settings by criteria.
		// limit and offset specified through query options are used for pagination.
		// total is the total number of entities in the database matching the criteria.
		Filter(ctx context.Context, in FilterSettingsByCriteriaInput, opts QueryOptions) (result []Setting, total int64, err error)
		// Create creates a setting.
		Create(ctx context.Context, entity *Setting) (err error)
		// CreateMultiple creates multiple settings.
		CreateMultiple(ctx context.Context, entities []*Setting) (err error)
		// Update updates a setting.
		Update(ctx context.Context, entity *Setting) (err error)
		// UpdateMultiple updates multiple settings.
		UpdateMultiple(ctx context.Context, entities []*Setting) (err error)
		// DeleteByID deletes a setting by its ID.
		DeleteByID(ctx context.Context, id uuid.UUID) (err error)
		// DeleteByIDs deletes settings by their IDs.
		DeleteByIDs(ctx context.Context, ids []uuid.UUID) (err error)
	}

	// SettingService defines the setting service
	SettingService interface {
		// FindByID finds a setting by its ID.
		FindByID(id uuid.UUID) (result Setting, err error)
		// Filter filters settings by criteria.
		// limit and offset specified through query options are used for pagination.
		// total is the total number of entities in the database matching the criteria.
		Filter(in FilterSettingsByCriteriaInput, options QueryOptions) (result []Setting, total int64, err error)
	}
)

const (
	SettingKeyAppName       = "app.name"
	SettingTestPhoneNumbers = "test.phone_numbers"
	SettingTestPhoneCode    = "test.phone_code"
)

var SettingDefaultValues = map[string]string{
	SettingKeyAppName:       "App",
	SettingTestPhoneNumbers: "+911234567890",
	SettingTestPhoneCode:    "123456",
}
