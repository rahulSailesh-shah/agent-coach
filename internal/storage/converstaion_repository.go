package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agent-coach/internal/models"

	"github.com/google/uuid"
)

type ConversationRepository struct {
	db *DB
}

func NewConversationRepository(db *DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) Create(ctx context.Context, conv *models.Conversation) error {
	if conv.ID == "" {
		conv.ID = uuid.New().String()
	}
	conv.CreatedAt = time.Now()

	query := `
		INSERT INTO conversations (id, goal_id, session_id, role, content, agent_type, metadata, created_at)
		VALUES (:id, :goal_id, :session_id, :role, :content, :agent_type, :metadata, :created_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, conv)
	return err
}

func (r *ConversationRepository) GetByID(ctx context.Context, id string) (*models.Conversation, error) {
	query := `
		SELECT id, goal_id, session_id, role, content, agent_type, metadata, created_at
		FROM conversations WHERE id = ?
	`

	var entity models.Conversation
	err := r.db.GetContext(ctx, &entity, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *ConversationRepository) GetBySessionID(ctx context.Context, sessionID string) ([]*models.Conversation, error) {
	query := `
		SELECT id, goal_id, session_id, role, content, agent_type, metadata, created_at
		FROM conversations WHERE session_id = ? ORDER BY created_at ASC
	`

	var entities []models.Conversation
	if err := r.db.SelectContext(ctx, &entities, query, sessionID); err != nil {
		return nil, err
	}

	conversations := make([]*models.Conversation, len(entities))
	for i, entity := range entities {
		conversations[i] = &entity
	}

	return conversations, nil
}

func (r *ConversationRepository) GetByGoalID(ctx context.Context, goalID string, limit int) ([]*models.Conversation, error) {
	query := `
		SELECT id, goal_id, session_id, role, content, agent_type, metadata, created_at
		FROM conversations WHERE goal_id = ? ORDER BY created_at DESC LIMIT ?
	`

	var entities []models.Conversation
	if err := r.db.SelectContext(ctx, &entities, query, goalID, limit); err != nil {
		return nil, err
	}

	conversations := make([]*models.Conversation, len(entities))
	for i, entity := range entities {
		conversations[len(entities)-1-i] = &entity
	}

	return conversations, nil
}

func (r *ConversationRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM conversations WHERE session_id = ?", sessionID)
	return err
}
