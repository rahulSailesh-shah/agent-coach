package tool

import (
	"agent-coach/internal/models"
	"agent-coach/internal/storage"
	"context"
	"time"
)

type ToolExecutor struct {
	taskRepo *storage.TaskRepository
	goalRepo *storage.GoalRepository
}

func NewToolExecutor(db *storage.DB) *ToolExecutor {
	return &ToolExecutor{
		taskRepo: storage.NewTaskRepository(db),
		goalRepo: storage.NewGoalRepository(db),
	}
}

func (e *ToolExecutor) ExecuteTool(ctx context.Context, toolName string, args map[string]any, agentCtx *models.AgentContext) (map[string]any, error) {
	switch toolName {
	case "create_task":
		return e.executeCreateTask(ctx, args, agentCtx)
	case "mark_complete":
		return e.executeMarkComplete(ctx, args)
	case "log_struggle":
		return e.executeLogStruggle(ctx, args)
	default:
		return args, nil
	}
}

func (e *ToolExecutor) executeCreateTask(ctx context.Context, args map[string]any, agentCtx *models.AgentContext) (map[string]any, error) {
	task := &models.Task{
		GoalID: agentCtx.Goal.ID,
		Title:  args["title"].(string),
		Status: models.TaskStatusPending,
	}

	if desc, ok := args["description"].(string); ok {
		task.Description = desc
	}
	if dueStr, ok := args["due_date"].(string); ok {
		if due, err := time.Parse("2006-01-02", dueStr); err == nil {
			task.DueDate = &due
		}
	}
	if mins, ok := args["estimated_minutes"].(float64); ok {
		m := int(mins)
		task.EstimatedMinutes = &m
	}
	if diff, ok := args["difficulty"].(float64); ok {
		d := int(diff)
		task.DifficultyRating = &d
	}
	if prio, ok := args["priority"].(float64); ok {
		task.Priority = int(prio)
	}

	if err := e.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return map[string]any{
		"task_id": task.ID,
		"message": "Task created successfully",
	}, nil
}

func (e *ToolExecutor) executeMarkComplete(ctx context.Context, args map[string]any) (map[string]any, error) {
	taskID := args["task_id"].(string)
	var mins *int
	if m, ok := args["actual_minutes"].(float64); ok {
		mi := int(m)
		mins = &mi
	}

	if err := e.taskRepo.MarkComplete(ctx, taskID, mins); err != nil {
		return nil, err
	}

	return map[string]any{
		"message": "Task marked as complete",
	}, nil
}

func (e *ToolExecutor) executeLogStruggle(ctx context.Context, args map[string]any) (map[string]any, error) {
	taskID := args["task_id"].(string)
	notes := args["notes"].(string)

	if err := e.taskRepo.LogStruggle(ctx, taskID, notes); err != nil {
		return nil, err
	}

	return map[string]any{
		"message": "Struggle logged",
	}, nil
}
