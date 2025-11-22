package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	Update(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Message, int, error)
}

type messageService struct {
	messageRepository MessageRepository
	log               logger.Logger
}

func NewMessageService(messageRepository MessageRepository, log logger.Logger) *messageService {
	return &messageService{messageRepository, log}
}

func (s *messageService) CreateMessage(ctx context.Context, req *domain.CreateMessageRequest) error {
	log := s.log.WithFields(map[string]any{
		"layer":     "service",
		"operation": "create_message",
	})

	message := &domain.Message{
		ID:      uuid.New(),
		Name:    req.Name,
		Email:   req.Email,
		Text:    req.Text,
		Time:    time.Now(),
		Unread:  true,
		Country: "SYRIA",
	}

	if err := s.messageRepository.Create(ctx, message); err != nil {
		log.WithError(err).Error("failed to create message")
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (s *messageService) GetMessageByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	log := s.log.WithFields(map[string]any{
		"layer":     "service",
		"operation": "get_message_by_id",
	})

	message, err := s.messageRepository.Get(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to get message")
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

func (s *messageService) MarkMessageAsRead(ctx context.Context, id uuid.UUID) error {
	log := s.log.WithFields(map[string]any{
		"layer":     "service",
		"operation": "mark_message_as_read",
	})

	if err := s.messageRepository.Update(ctx, id); err != nil {
		log.WithError(err).Error("Failed to mark message as read")
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

func (s *messageService) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	log := s.log.WithFields(map[string]any{
		"layer":     "service",
		"operation": "delete_message",
	})

	if err := s.messageRepository.Delete(ctx, id); err != nil {
		log.WithError(err).Error("failed to delete message")
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *messageService) GetMessagesList(ctx context.Context, limit, offset int) ([]*domain.Message, int, error) {
	log := s.log.WithFields(map[string]any{
		"layer":     "service",
		"operation": "get_messages_list",
	})

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	messages, total, err := s.messageRepository.List(ctx, limit, offset)
	if err != nil {
		log.WithError(err).Error("failed to get messages list")
		return nil, 0, fmt.Errorf("failed to get messages list: %w", err)
	}

	return messages, total, nil
}
