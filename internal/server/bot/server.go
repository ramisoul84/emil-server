package bot

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ramisoul84/emil-server/config"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type BotServer struct {
	bot      *tgbotapi.BotAPI
	adminIDs map[int64]bool
	mu       sync.RWMutex
	logger   logger.Logger
}

func NewBotServer(cfg *config.Config) (*BotServer, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	server := &BotServer{
		bot:      bot,
		adminIDs: make(map[int64]bool),
		logger:   logger.Get(),
	}

	// Get updates to capture admin chat IDs when they start the bot
	go server.listenForAdmins(cfg.Bot.AdminUsernames)

	server.logger.Info().Msgf("🤖 Bot started: @%s", bot.Self.UserName)
	return server, nil
}

// Listen for /start command from admins
func (s *BotServer) listenForAdmins(adminUsernames []string) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		if update.Message.Command() == "start" {
			// Check if user is admin
			for _, adminUsername := range adminUsernames {
				if update.Message.From.UserName == adminUsername {
					s.registerAdmin(update.Message.Chat.ID)
					s.sendWelcome(update.Message.Chat.ID)
					break
				}
			}
		}
	}
}

func (s *BotServer) registerAdmin(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adminIDs[chatID] = true
	s.logger.Info().Msgf("✅ Admin registered: %d", chatID)
}

func (s *BotServer) sendWelcome(chatID int64) {
	msg := tgbotapi.NewMessage(chatID,
		"👋 Welcome Admin!\n\n"+
			"You'll receive notifications when someone visits the site.")
	s.bot.Send(msg)
}

// SendNotification sends a message to all registered admins
func (s *BotServer) SendNotification(message string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.adminIDs) == 0 {
		return fmt.Errorf("no admins registered")
	}

	var errors []error
	for chatID := range s.adminIDs {
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "Markdown"

		if _, err := s.bot.Send(msg); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %d: %w", chatID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send to some admins: %v", errors)
	}
	return nil
}

// GetAdminCount returns number of registered admins
func (s *BotServer) GetAdminCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.adminIDs)
}
