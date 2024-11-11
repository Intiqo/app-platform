package security

import "github.com/gofrs/uuid/v5"

// TokenMetadata represents the metadata in the auth token
type TokenMetadata struct {
	UserID          uuid.UUID   `json:"user_id"`
	OrganizationID  uuid.UUID   `json:"organization_id"`
	OrganizationIDs []uuid.UUID `json:"organization_ids"`
	Role            string      `json:"role"`
}

// Manager defines the interface for a security manager
type Manager interface {
	// GenerateAuthToken generates an auth token for a user.
	GenerateAuthToken(metadata TokenMetadata) (token string, err error)
}
