package models

// State represents the current state of the agent workflow
type State string

const (
	StateOnboarding  State = "onboarding"
	StateGoalSetting State = "goal_setting"
	StatePlanning    State = "planning"
	StateActive      State = "active"
	StateCheckIn     State = "check_in"
	StateReplanning  State = "replanning"
	StateReview      State = "review"
	StateCompleted   State = "completed"
)

// AgentContext contains the context information for an agent execution
type AgentContext struct {
	Goal          *Goal
	Tasks         []*Task
	TodaysTasks   []*Task
	Conversations []*Conversation

	CurrentState    State
	StreakDays      int
	RecentStruggles int
	TasksCompleted  int
}
