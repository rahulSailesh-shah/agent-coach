package models

import (
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusSkipped    TaskStatus = "skipped"
)

type Task struct {
	ID               string     `db:"id" json:"id"`
	GoalID           string     `db:"goal_id" json:"goalId"`
	ParentID         *string    `db:"parent_id" json:"parentId,omitempty"`
	Title            string     `db:"title" json:"title"`
	Description      string     `db:"description" json:"description,omitempty"`
	DueDate          *time.Time `db:"due_date" json:"dueDate,omitempty"`
	Status           TaskStatus `db:"status" json:"status"`
	Priority         int        `db:"priority" json:"priority"`
	DifficultyRating *int       `db:"difficulty_rating" json:"difficultyRating,omitempty"`
	EstimatedMinutes *int       `db:"estimated_minutes" json:"estimatedMinutes,omitempty"`
	ActualMinutes    *int       `db:"actual_minutes" json:"actualMinutes,omitempty"`
	StruggleNotes    string     `db:"struggle_notes" json:"struggleNotes,omitempty"`
	CompletedAt      *time.Time `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
}
