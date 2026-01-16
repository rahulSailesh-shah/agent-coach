package agent

import (
	"context"

	"agent-coach/internal/models"
)

// Intent represents the classification of user intent
type Intent string

const (
	IntentPlanning       Intent = "PLANNING"
	IntentExecution      Intent = "EXECUTION"
	IntentEvaluation     Intent = "EVALUATION"
	IntentAccountability Intent = "ACCOUNTABILITY"
	IntentGeneral        Intent = "GENERAL"
)

type AgentInput struct {
	Message   string
	Context   *models.AgentContext
	GoalID    string
	SessionID string
}

type AgentOutput struct {
	Response  string
	NextState *models.State
	AgentType models.AgentType
}

type Agent interface {
	Type() models.AgentType
	SystemPrompt(ctx *models.AgentContext) string
	Execute(ctx context.Context, input *AgentInput) (*AgentOutput, error)
	AvailableTools() []models.Tool
}

type ClassifiedIntent struct {
	Intent     Intent  `json:"intent"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}
