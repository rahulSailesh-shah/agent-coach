package main

import (
	"context"
	"fmt"

	"agent-coach/internal/llm"
	"agent-coach/internal/models"
	"agent-coach/internal/service"
	"agent-coach/internal/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	db        *storage.DB
	llmRouter *llm.Router
	service   *service.Service
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := storage.NewDB()
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize database: %v", err)
		return
	}
	a.db = db

	a.llmRouter, err = llm.NewRouter(a.db)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize LLM router: %v", err)
		return
	}

	a.service = service.NewService(a.db, a.llmRouter)

	runtime.LogInfo(ctx, "Application started successfully")
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
	runtime.LogInfo(ctx, "Application shutdown complete")
}

// ============================================================================
// Goal Operations
// ============================================================================

func (a *App) CreateGoal(goal models.Goal) (*models.Goal, error) {
	if err := a.service.CreateGoal(a.ctx, &goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}
	return &goal, nil
}

func (a *App) GetGoal(id string) (*models.Goal, error) {
	return a.service.GetGoal(a.ctx, id)
}

func (a *App) UpdateGoal(goal models.Goal) error {
	return a.service.UpdateGoal(a.ctx, &goal)
}

func (a *App) DeleteGoal(id string) error {
	return a.service.DeleteGoal(a.ctx, id)
}

// ============================================================================
// Task Operations
// ============================================================================

func (a *App) CreateTask(task models.Task) (*models.Task, error) {
	if err := a.service.CreateTask(a.ctx, &task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return &task, nil
}

func (a *App) GetTask(id string) (*models.Task, error) {
	return a.service.GetTask(a.ctx, id)
}

func (a *App) GetTasksByGoalID(goalID string) ([]*models.Task, error) {
	return a.service.GetTasksByGoalID(a.ctx, goalID)
}

func (a *App) GetTodaysTasks() ([]*models.Task, error) {
	return a.service.GetTodaysTasks(a.ctx)
}

func (a *App) CompleteTask(id string, actualMinutes *int) error {
	return a.service.CompleteTask(a.ctx, id, actualMinutes)
}

func (a *App) LogStruggle(id string, notes string) error {
	return a.service.LogStruggle(a.ctx, id, notes)
}

// ============================================================================
// Agent-Powered Chat Operations
// ============================================================================

type ChatRequest struct {
	Message string `json:"message"`
	GoalID  string `json:"goalId"`
}

type ChatResponse struct {
	Content   string `json:"content"`
	AgentType string `json:"agentType"`
}

func (a *App) Chat(req ChatRequest) (*ChatResponse, error) {
	output, err := a.service.Chat(a.ctx, req.Message, req.GoalID)
	if err != nil {
		return nil, err
	}

	response := &ChatResponse{
		Content:   output.Response,
		AgentType: string(output.AgentType),
	}

	return response, nil
}

func (a *App) GetConversationHistory(coachID string, limit int) ([]*models.Conversation, error) {
	if limit <= 0 {
		limit = 50
	}
	return a.service.GetConversationHistory(a.ctx, coachID, limit)
}

// ============================================================================
// LLM Configuration Operations
// ============================================================================

func (a *App) GetLLMProviders() ([]llm.ProviderType, error) {
	return a.llmRouter.GetAvailableProviders(a.ctx)
}

func (a *App) GetLLMProviderConfigs() ([]*models.LLMProviderConfig, error) {
	return a.llmRouter.GetProviderConfigs(a.ctx)
}

func (a *App) SaveLLMProviderConfig(config models.LLMProviderConfig) (*models.LLMProviderConfig, error) {
	if err := a.llmRouter.SaveProviderConfig(a.ctx, &config); err != nil {
		return nil, fmt.Errorf("failed to save provider config: %w", err)
	}
	return &config, nil
}

func (a *App) UpdateLLMProviderConfig(config models.LLMProviderConfig) error {
	if err := a.llmRouter.UpdateProviderConfig(a.ctx, &config); err != nil {
		return fmt.Errorf("failed to update provider config: %w", err)
	}
	return nil
}

func (a *App) DeleteLLMProviderConfig(id int) error {
	if err := a.llmRouter.DeleteProviderConfig(a.ctx, id); err != nil {
		return fmt.Errorf("failed to delete provider config: %w", err)
	}
	return nil
}
