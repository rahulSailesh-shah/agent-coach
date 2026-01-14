package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agent-coach/internal/models"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	task.CreatedAt = time.Now()

	// entity := r.toEntity(task)

	query := `
		INSERT INTO tasks (id, goal_id, parent_id, title, description, due_date, status, 
			priority, difficulty_rating, estimated_minutes, actual_minutes, struggle_notes, 
			created_at, completed_at)
		VALUES (:id, :goal_id, :parent_id, :title, :description, :due_date, :status,
			:priority, :difficulty_rating, :estimated_minutes, :actual_minutes, :struggle_notes,
			:created_at, :completed_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	query := `
		SELECT id, goal_id, parent_id, title, description, due_date, status, priority,
			difficulty_rating, estimated_minutes, actual_minutes, struggle_notes, created_at, completed_at
		FROM tasks WHERE id = ?
	`

	var entity models.Task
	err := r.db.GetContext(ctx, &entity, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *TaskRepository) GetByGoalID(ctx context.Context, goalID string) ([]*models.Task, error) {
	query := `
		SELECT id, goal_id, parent_id, title, description, due_date, status, priority,
			difficulty_rating, estimated_minutes, actual_minutes, struggle_notes, created_at, completed_at
		FROM tasks WHERE goal_id = ? ORDER BY priority DESC, due_date ASC
	`

	var entities []models.Task
	if err := r.db.SelectContext(ctx, &entities, query, goalID); err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, len(entities))
	for i, entity := range entities {
		tasks[i] = &entity
	}

	return tasks, nil
}

func (r *TaskRepository) GetPendingByGoalID(ctx context.Context, goalID string) ([]*models.Task, error) {
	query := `
		SELECT id, goal_id, parent_id, title, description, due_date, status, priority,
			difficulty_rating, estimated_minutes, actual_minutes, struggle_notes, created_at, completed_at
		FROM tasks WHERE goal_id = ? AND status IN ('pending', 'in_progress') 
		ORDER BY priority DESC, due_date ASC
	`

	var entities []models.Task
	if err := r.db.SelectContext(ctx, &entities, query, goalID); err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, len(entities))
	for i, entity := range entities {
		tasks[i] = &entity
	}

	return tasks, nil
}

func (r *TaskRepository) GetDueToday(ctx context.Context) ([]*models.Task, error) {
	query := `
		SELECT id, goal_id, parent_id, title, description, due_date, status, priority,
			difficulty_rating, estimated_minutes, actual_minutes, struggle_notes, created_at, completed_at
		FROM tasks 
		WHERE date(due_date) = date('now') AND status IN ('pending', 'in_progress')
		ORDER BY priority DESC
	`

	var entities []models.Task
	if err := r.db.SelectContext(ctx, &entities, query); err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, len(entities))
	for i, entity := range entities {
		tasks[i] = &entity
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks SET 
			title = :title, description = :description, due_date = :due_date,
			status = :status, priority = :priority, difficulty_rating = :difficulty_rating,
			estimated_minutes = :estimated_minutes, actual_minutes = :actual_minutes,
			struggle_notes = :struggle_notes, completed_at = :completed_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, task)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TaskRepository) MarkComplete(ctx context.Context, id string, actualMinutes *int) error {
	now := time.Now()
	query := `UPDATE tasks SET status = 'completed', completed_at = ?, actual_minutes = ? WHERE id = ?`

	var minutes int
	if actualMinutes != nil {
		minutes = *actualMinutes
	}

	result, err := r.db.ExecContext(ctx, query, now, minutes, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TaskRepository) LogStruggle(ctx context.Context, id string, notes string) error {
	query := `UPDATE tasks SET struggle_notes = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, notes, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
