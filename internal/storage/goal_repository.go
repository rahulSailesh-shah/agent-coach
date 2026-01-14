package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agent-coach/internal/models"

	"github.com/google/uuid"
)

type GoalRepository struct {
	db *DB
}

func NewGoalRepository(db *DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Create(ctx context.Context, goal *models.Goal) error {
	if goal.ID == "" {
		goal.ID = uuid.New().String()
	}
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	query := `
		INSERT INTO goals (id, title, description, target_date, status, context, created_at, updated_at)
		VALUES (:id, :title, :description, :target_date, :status, :context, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, goal)
	return err
}

func (r *GoalRepository) GetByID(ctx context.Context, id string) (*models.Goal, error) {
	query := `
		SELECT id, title, description, target_date, status, context, created_at, updated_at
		FROM goals WHERE id = ?
	`

	var entity models.Goal
	err := r.db.GetContext(ctx, &entity, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Update updates an existing goal
func (r *GoalRepository) Update(ctx context.Context, goal *models.Goal) error {
	goal.UpdatedAt = time.Now()

	query := `
		UPDATE goals SET 
			title = :title, description = :description, target_date = :target_date,
			status = :status, context = :context, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, goal)
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

func (r *GoalRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM goals WHERE id = ?", id)
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

func (r *GoalRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE goals SET status = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
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
