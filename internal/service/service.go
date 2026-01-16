package service

import (
	"context"

	"agent-coach/internal/agent"
	"agent-coach/internal/llm"
	"agent-coach/internal/models"
	"agent-coach/internal/storage"
)

type Service struct {
	db           *storage.DB
	goalRepo     *storage.GoalRepository
	taskRepo     *storage.TaskRepository
	convRepo     *storage.ConversationRepository
	llmRouter    *llm.Router
	orchestrator *agent.Orchestrator
}

func NewService(db *storage.DB, router *llm.Router) *Service {
	return &Service{
		db:           db,
		goalRepo:     storage.NewGoalRepository(db),
		taskRepo:     storage.NewTaskRepository(db),
		convRepo:     storage.NewConversationRepository(db),
		llmRouter:    router,
		orchestrator: agent.NewOrchestrator(db, router),
	}
}

// Goal Operations
func (s *Service) CreateGoal(ctx context.Context, goal *models.Goal) error {
	return s.goalRepo.Create(ctx, goal)
}

func (s *Service) GetGoal(ctx context.Context, id string) (*models.Goal, error) {
	return s.goalRepo.GetByID(ctx, id)
}

func (s *Service) UpdateGoal(ctx context.Context, goal *models.Goal) error {
	return s.goalRepo.Update(ctx, goal)
}

func (s *Service) DeleteGoal(ctx context.Context, id string) error {
	return s.goalRepo.Delete(ctx, id)
}

// Task Operations

func (s *Service) CreateTask(ctx context.Context, task *models.Task) error {
	return s.taskRepo.Create(ctx, task)
}

func (s *Service) GetTask(ctx context.Context, id string) (*models.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *Service) GetTasksByGoalID(ctx context.Context, goalID string) ([]*models.Task, error) {
	return s.taskRepo.GetByGoalID(ctx, goalID)
}

func (s *Service) GetTodaysTasks(ctx context.Context) ([]*models.Task, error) {
	return s.taskRepo.GetDueToday(ctx)
}

func (s *Service) CompleteTask(ctx context.Context, id string, actualMinutes *int) error {
	return s.taskRepo.MarkComplete(ctx, id, actualMinutes)
}

func (s *Service) LogStruggle(ctx context.Context, id string, notes string) error {
	return s.taskRepo.LogStruggle(ctx, id, notes)
}

// Chat Operations
func (s *Service) Chat(ctx context.Context, message string, goalID string) (*agent.AgentOutput, error) {
	return s.orchestrator.ProcessMessage(ctx, message, goalID)
}

func (s *Service) StartOnboarding(ctx context.Context, coachID string, goalTitle string) (*agent.AgentOutput, error) {
	return s.orchestrator.StartOnboarding(ctx, coachID, goalTitle)
}

func (s *Service) GetConversationHistory(ctx context.Context, goalID string, limit int) ([]*models.Conversation, error) {
	return s.convRepo.GetByGoalID(ctx, goalID, limit)
}
