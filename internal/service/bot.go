package service

import (
	"context"

	"github.com/ramisoul84/emil-server/internal/server/bot"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type BotService struct {
	server *bot.BotServer
	logger logger.Logger
}

func NewBotService(server *bot.BotServer) *BotService {
	return &BotService{
		server: server,
		logger: logger.Get(),
	}
}

// NotifyVisit sends a visit notification to all admins
func (s *BotService) Notify(ctx context.Context, msg string) error {
	if err := s.server.SendNotification(msg); err != nil {
		s.logger.Error().Err(err).Msg("Failed to send notification")
		return err
	}

	s.logger.Info().Msg("Visit notification sent to admins")
	return nil
}
