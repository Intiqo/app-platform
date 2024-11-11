package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	// JSONB represents a JSONB type
	JSONB map[string]interface{} // @name JSONB
)

// Define all the base models here
type (
	// Base represents the base model for App
	Base struct {
		ID uuid.UUID `db:"id" json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	} // @name Base

	// Audit represents fields capturing the audit information for App
	Audit struct {
		CreatedAt time.Time  `db:"created_at" json:"-" example:"2020-01-01T00:00:00+05:30"`
		UpdatedAt time.Time  `db:"updated_at" json:"-" example:"2020-01-01T00:00:00+05:30"`
		DeletedAt *time.Time `db:"deleted_at" json:"-" example:"2020-01-01T00:00:00+05:30"`
	} // @name Audit

	// Address defines model for an Address.
	Address struct {
		Street   string `db:"street" json:"street,omitempty" example:"123 Main St"`
		District string `json:"district,omitempty" example:"Mumbai"`
		City     string `db:"city" json:"city,omitempty" example:"New York"`
		State    string `db:"state" json:"state,omitempty" example:"NY"`
		Zipcode  string `db:"zipcode" json:"zipCode,omitempty" example:"10001"`
		Country  string `db:"country" json:"country,omitempty" example:"United States"`
	} // @name Address

	// Dimension defines model for a Dimension.
	Dimension struct {
		MeasurementUnit string  `db:"measurement_unit" json:"measurementUnit,omitempty" example:"acres"`
		Area            float64 `db:"area" json:"area,omitempty" example:"100.00"`
	} // @name Dimension
)

type (
	FilterOp string // @name FilterOp

	// Claims represents the claims in the JWT token
	Claims struct {
		UserID uuid.UUID `json:"userId" swaggerignore:"true"`
		Role   string    `json:"role" swaggerignore:"true"`
	} // @name Claims

	// TokenInfo represents the token information
	TokenInfo struct {
		UserID          uuid.UUID   `json:"-"`
		OrganizationID  uuid.UUID   `json:"-"`
		OrganizationIDs []uuid.UUID `json:"-"`
		Role            string      `json:"-"`
	} // @name TokenInfo

	// SortKey defines the sort key for sorting
	SortKey struct {
		// Field represents a column for the entity you are sorting
		Field string `json:"field" example:"name"`
		// Direction represents the direction of the sort
		Direction string `json:"direction" enums:"asc,desc" example:"asc"`
	} // @name SortKey

	// QueryOptions defines the options for a query
	QueryOptions struct {
		Limit  int64 `json:"limit" example:"10"`
		Offset int64 `json:"offset" example:"0"`

		SortKeys []SortKey `json:"sortKeys"`
	} // @name QueryOptions

)

type (
	// BaseResponse is the base response type
	BaseResponse struct {
		Data interface{} `json:"data"`
	} // @name BaseResponse

	// PaginationResponse is the pagination response type
	PaginationResponse struct {
		BaseResponse
		Total int64 `json:"total" example:"100"`
		Page  int64 `json:"page" example:"1"`
		Size  int64 `json:"size" example:"10"`
	} // @name PaginationResponse

	// ErrorResponse is the error response type
	ErrorResponse struct {
		Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
		Message string `json:"message" example:"Internal Server Error"`
	} // @name ErrorResponse
)

type (
	// Transactioner defines the methods that a transactioner should implement
	Transactioner interface {
		// Begin starts a transaction
		Begin(ctx context.Context) (result context.Context, err error)
		// Commit commits a transaction
		Commit(ctx context.Context) (err error)
		// Rollback rolls back a transaction
		Rollback(ctx context.Context, err error)
	}
)

// Value implements the driver.Valuer interface,
func (j *JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// Scan implements the sql.Scanner interface,
func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal([]byte(value.(string)), &j); err != nil {
		return err
	}
	return nil
}

const (
	MessageVALIDATIONFAILED   string = "Validation failed for some or all of the fields in the request"
	MessageUNAUTHORIZEDACCESS string = "You are not authorized to access this resource"
	MessageFORBIDDENACCESS    string = "You are forbidden from accessing this resource"
)

const (
	SortDirectionAsc  string = "asc"
	SortDirectionDesc string = "desc"
)
