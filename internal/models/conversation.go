package models

import (
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type AgentType string

const (
	AgentTypePlanner        AgentType = "planner"
	AgentTypeExecutor       AgentType = "executor"
	AgentTypeEvaluator      AgentType = "evaluator"
	AgentTypeAccountability AgentType = "accountability"
)

type Conversation struct {
	ID        string                 `db:"id" json:"id"`
	GoalID    *string                `db:"goal_id" json:"goalId,omitempty"`
	SessionID string                 `db:"session_id" json:"sessionId"`
	Role      Role                   `db:"role" json:"role"`
	Content   string                 `db:"content" json:"content"`
	AgentType AgentType              `db:"agent_type" json:"agentType,omitempty"`
	Metadata  map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt time.Time              `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time              `db:"updated_at" json:"updatedAt"`
}
