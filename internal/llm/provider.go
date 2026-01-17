package llm

import (
	"agent-coach/internal/models"
	"context"
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

type CompletionRequest struct {
	Model        string        `json:"model"`
	Messages     []Message     `json:"messages"`
	SystemPrompt string        `json:"system_prompt"`
	Tools        []models.Tool `json:"tools"`
}

type CompletionResponse struct {
	Content      string            `json:"content"`
	ToolCalls    []models.ToolCall `json:"tool_calls"`
	Model        string            `json:"model"`
	Usage        int               `json:"usage"`
	FinishReason string            `json:"finish_reason"`
}

type ProviderType string

const (
	ProviderTypeOpenRouter ProviderType = "openrouter"
	ProviderTypeOllama     ProviderType = "ollama"
)

var availableProviderTypes = []ProviderType{
	ProviderTypeOpenRouter,
	ProviderTypeOllama,
}

type Provider interface {
	Name() string
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	IsAvailable() bool
	Type() ProviderType
}
