package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ramisoul84/emil-server/internal/domain"
	"github.com/ramisoul84/emil-server/pkg/location"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type messageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	Get(ctx context.Context, id int) (*domain.Message, error)
	Update(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.Message, int, error)
}

type messageService struct {
	repo   messageRepository
	bot    botNotifier
	logger logger.Logger
}

func NewMessageService(repo messageRepository, bot botNotifier) *messageService {
	return &messageService{
		repo:   repo,
		bot:    bot,
		logger: logger.Get(),
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *domain.Message) error {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "message_service",
			"method":     "create_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling create message")

	ip := ctx.Value("ip").(string)
	country, city := location.GetFullClientInfo(ip)

	message.Time = time.Now()
	message.Unread = false
	message.IP = ip
	message.City = city
	message.Country = country

	msg := fmt.Sprintf(
		"📊 *You got a message*\n\n"+
			"📍 *IP:* %s\n"+
			"📍 *Country:* %s\n"+
			"📍 *City:* %s\n"+
			"👤 *User ID:* %s\n"+
			"👤 *Name:* %s\n"+
			"👤 *Email:* %s\n\n"+
			"%s",
		ip,
		country,
		city,
		message.UserID,
		message.Name,
		message.Email,
		message.Text,
	)

	if err := s.bot.Notify(context.Background(), msg); err != nil {
		logger.Error().Err(err).Msg("Failed to send bot notification")
	}

	err := s.repo.Create(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *messageService) GetMessage(ctx context.Context, id int) (*domain.Message, error) {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "message_service",
			"method":     "get_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling get message")

	return s.repo.Get(ctx, id)
}

func (s *messageService) UpdateMessage(ctx context.Context, id int) error {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "message_service",
			"method":     "update_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling update message")

	return s.repo.Update(ctx, id)
}

func (s *messageService) DeleteMessage(ctx context.Context, id int) error {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "message_service",
			"method":     "delete_message",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling delete message")

	return s.repo.Delete(ctx, id)
}

func (s *messageService) ListMessages(ctx context.Context, limit, offset int) ([]*domain.Message, int, error) {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "message_service",
			"method":     "list_messages",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling list messages")

	return s.repo.List(ctx, limit, offset)
}
