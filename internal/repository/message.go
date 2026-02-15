package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ramisoul84/emil-server/internal/domain"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type messageRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewMessageRepository(db *sqlx.DB) *messageRepository {
	return &messageRepository{db, logger.Get()}
}

func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "message_repository",
			"method":     "create_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)
	logger.Info().Msg("Store message in DB")

	query := `
            INSERT INTO messages (
				user_id, name, email, text,
				time, unread, ip, city, country
				)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

	_, err := r.db.ExecContext(ctx, query,
		message.UserID,
		message.Name,
		message.Email,
		message.Text,
		message.Time,
		message.Unread,
		message.IP,
		message.City,
		message.Country,
	)

	if err != nil {
		logger.Error().Err(err).Msg("failed to save message")
		return fmt.Errorf("failed to save message: %w", err)
	}

	logger.Info().Msg("message saved successfully")

	return nil
}

func (r *messageRepository) Get(ctx context.Context, id int) (*domain.Message, error) {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "message_repository",
			"method":     "get_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)
	logger.Info().Msg("get message from DB")

	query := `
		SELECT *
		FROM messaages 
		WHERE id = $1
	`

	var message domain.Message
	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info().Msg("message not found")
			return nil, domain.ErrNotFound
		}
		logger.Error().Err(err).Msg("failed to get message")
		return nil, domain.ErrInternal
	}

	return &message, nil
}

func (r *messageRepository) Update(ctx context.Context, id int) error {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "message_repository",
			"method":     "update_message_status",
			"request_id": ctx.Value("request_id").(string),
		},
	)
	logger.Info().Msg("update message status")

	query := `
		UPDATE messages
		SET unread = false
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error().Err(err).Msg("ffailed to update message")
		return domain.ErrInternal

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		logger.Info().Msg("message not found")
		return domain.ErrNotFound
	}
	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "message_repository",
			"method":     "delete_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)
	logger.Info().Msg("delete message from DB")

	query := `
		DELETE FROM messages 
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error().Err(err).Msg("ffailed to delete message")
		return domain.ErrInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		logger.Info().Msg("message not found")
		return domain.ErrNotFound
	}

	return nil
}

func (r *messageRepository) List(ctx context.Context, limit, offset int) ([]*domain.Message, int, error) {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "message_repository",
			"method":     "list_messages",
			"request_id": ctx.Value("request_id").(string),
		},
	)
	logger.Info().Msg("list messages")

	query := `
		SELECT *
		FROM messages
		ORDER BY time DESC
		LIMIT $1 OFFSET $2
	`

	var messages []*domain.Message
	err := r.db.SelectContext(ctx, &messages, query, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list messages")
		return nil, 0, domain.ErrInternal
	}

	countQuery := `SELECT COUNT(*) FROM messages`
	var total int
	err = r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get total messages")
		return nil, 0, domain.ErrInternal
	}

	return messages, total, nil
}
