package models

import (
	"time"
)

type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusPaused    GoalStatus = "paused"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusAbandoned GoalStatus = "abandoned"
)

type Goal struct {
	ID          string                 `db:"id" json:"id"`
	Title       string                 `db:"title" json:"title"`
	Description string                 `db:"description" json:"description,omitempty"`
	TargetDate  *time.Time             `db:"target_date" json:"targetDate,omitempty"`
	Status      GoalStatus             `db:"status" json:"status"`
	Context     map[string]interface{} `db:"context" json:"context,omitempty"`
	CreatedAt   time.Time              `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time              `db:"updated_at" json:"updatedAt"`
}
