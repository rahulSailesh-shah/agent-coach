package llm

import (
	"context"
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	RoleTool      Role = "tool"
)

type Message struct {
	Role       Role   `json:"role"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id"`
}

type ToolParam struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Enum        []string `json:"enum"`
}

type Tool struct {
	Type        string               `json:"type"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Parameters  map[string]ToolParam `json:"parameters"`
}

type ToolFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type CompletionRequest struct {
	Model        string    `json:"model"`
	Messages     []Message `json:"messages"`
	SystemPrompt string    `json:"system_prompt"`
	Tools        []Tool    `json:"tools"`
}

type CompletionResponse struct {
	Content      string     `json:"content"`
	ToolCalls    []ToolCall `json:"tool_calls"`
	Model        string     `json:"model"`
	Usage        int        `json:"usage"`
	FinishReason string     `json:"finish_reason"`
}

type Provider interface {
	Name() string
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
	IsAvailable() bool
}

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
