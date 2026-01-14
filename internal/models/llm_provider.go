package models

import "time"

type LLMProviderConfig struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name"`
	Provider     string    `db:"provider" json:"provider"`
	BaseURL      string    `db:"base_url" json:"base_url,omitempty"`
	APIKey       string    `db:"api_key" json:"api_key,omitempty"`
	DefaultModel string    `db:"default_model" json:"default_model,omitempty"`
	IsDefault    bool      `db:"is_default" json:"is_default"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
