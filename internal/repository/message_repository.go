package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type messageRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewMessageRepository(db *sqlx.DB, log logger.Logger) *messageRepository {
	return &messageRepository{db, log}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "create_message",
	})

	query := `
		INSERT INTO messages (id, name, email, text, time, unread, city, country)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.Name,
		message.Email,
		message.Text,
		message.Time,
		message.Unread,
		message.City,
		message.Country,
	)

	if err != nil {
		log.WithError(err).Error("failed to create message")
		return fmt.Errorf("failed to create message: %w", err)
	}

	log.Info("message created successfully")

	return nil
}

func (r *messageRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "get_message",
	})

	query := `
		SELECT *
		FROM messaages 
		WHERE id = $1
	`
	var message domain.Message
	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("message not found")
			return nil, domain.ErrNotFound
		}
		log.WithError(err).Error("failed to get message")
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil

}

func (r *messageRepository) Update(ctx context.Context, id uuid.UUID) error {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "update_message",
	})

	query := `
		UPDATE messages
		SET unread = false
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.WithError(err).Error("failed to update message")
		return fmt.Errorf("failed to update message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Info("messsage not found for update")
		return domain.ErrNotFound
	}

	log.Info("message updated successfully")

	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "delete_message",
	})

	query := `
		DELETE FROM messages 
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.WithError(err).Error("failed to delete message")
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Info("messsage not found for deletion")
		return domain.ErrNotFound
	}

	log.Info("message deleted successfully")
	return nil
}

func (r *messageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Message, int, error) {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "list_messages",
	})

	// Query for messages with pagination
	query := `
		SELECT *
		FROM messages
		ORDER BY time DESC
		LIMIT $1 OFFSET $2
	`

	var messages []*domain.Message
	err := r.db.SelectContext(ctx, &messages, query, limit, offset)
	if err != nil {
		log.WithError(err).Error("failed to list messages")
		return nil, 0, fmt.Errorf("failed to list messages: %w", err)
	}

	// Query for total count
	countQuery := `SELECT COUNT(*) FROM messages`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		log.WithError(err).Error("failed to get total count")
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	log.WithFields(map[string]any{
		"count": len(messages),
		"total": total,
	}).Info("messages listed successfully")

	return messages, total, nil
}
