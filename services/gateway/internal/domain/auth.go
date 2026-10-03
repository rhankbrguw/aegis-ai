package domain

import (
	"time"
)

// APIKey represents a client credential entity.
type APIKey struct {
	ID        string    `json:"id"`
	KeyHash   string    `json:"-"`
	Prefix    string    `json:"prefix"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	RPM       int       `json:"rpm"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateKeyRequest defines input payload for key generation.
type CreateKeyRequest struct {
	TenantID string `json:"tenant_id" validate:"required,min=3"`
	Name     string `json:"name" validate:"required,min=2"`
	RPM      int    `json:"rpm" validate:"gte=0"`
}

// CreateKeyResponse returns the generated secret key (only once).
type CreateKeyResponse struct {
	APIKey APIKey `json:"api_key"`
	Secret string `json:"secret"`
}
