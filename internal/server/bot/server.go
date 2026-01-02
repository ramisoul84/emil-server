package bot

import (
	"fmt"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/ramisoul/emil-server/config"
)

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	ChatID    int64  `json:"chat_id"`
}

type BotServer struct {
	bot        *tgbotapi.BotAPI
	users      map[int64]*User
	usersMutex sync.RWMutex
	stopChan   chan struct{}
}

type Messsage struct {
	Name  string
	Email string
	Text  string
}

type Visit struct {
	Country string
}

func NewBot(cfg config.Config) (*BotServer, error) {
	if cfg.Bot.Token == "" {
		return nil, fmt.Errorf("bot token is empty")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot API instance: %w", err)
	}

	server := &BotServer{
		bot:      bot,
		users:    make(map[int64]*User),
		stopChan: make(chan struct{}),
	}

	log.Printf("🤖 Bot authorized as: @%s", bot.Self.UserName)
	log.Printf("🆔 Bot ID: %d", bot.Self.ID)

	return server, nil
}

// Start begins listening for updates and manages the bot
func (s *BotServer) Start() error {
	log.Println("🚀 Starting bot server...")

	// Set bot commands
	s.setCommands()

	// Start listening for updates in a goroutine
	go s.handleUpdates()

	log.Println("✅ Bot server is running")
	return nil
}

// Stop gracefully shuts down the bot server
func (s *BotServer) Stop() {
	log.Println("🛑 Stopping bot server...")
	close(s.stopChan)
	log.Println("✅ Bot server stopped")
}

// setCommands configures the bot's command menu
func (s *BotServer) setCommands() {
	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Start the bot",
		},
	)

	if _, err := s.bot.Request(commands); err != nil {
		log.Printf("Warning: Failed to set commands: %v", err)
	}
}

// handleUpdates processes incoming Telegram updates
func (s *BotServer) handleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.Limit = 100
	u.Offset = 0

	updates := s.bot.GetUpdatesChan(u)

	for {
		select {
		case <-s.stopChan:
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			// Process message
			s.processMessage(update.Message)
		}
	}
}

// processMessage handles incoming messages
func (s *BotServer) processMessage(msg *tgbotapi.Message) {
	user := msg.From
	chat := msg.Chat

	// Store/update user information
	s.storeUser(user, chat)

	// Handle commands
	if msg.IsCommand() {
		s.handleCommand(msg)
		return
	}
}

// storeUser stores or updates user information
func (s *BotServer) storeUser(user *tgbotapi.User, chat *tgbotapi.Chat) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()

	_, exists := s.users[user.ID]
	if !exists {
		s.users[user.ID] = &User{
			ID:        user.ID,
			FirstName: user.FirstName,
			Username:  user.UserName,
			ChatID:    chat.ID,
		}
	}
}

// handleCommand processes bot commands
func (s *BotServer) handleCommand(msg *tgbotapi.Message) {
	user := msg.From
	chatID := msg.Chat.ID
	command := msg.Command()

	switch command {
	case "start":
		s.sendWelcomeMessage(chatID, user)
	default:
		s.sendMessage(chatID, "❓ Unknown command.")
	}
}

// sendWelcomeMessage sends welcome message to new users
func (s *BotServer) sendWelcomeMessage(chatID int64, user *tgbotapi.User) {
	message := fmt.Sprintf(
		"👋 Welcome *%s*!\n\n"+
			"I'm your private notification bot.\n\n"+
			"✅ You are now registered to receive messages.",
		user.FirstName,
	)

	s.sendMessage(chatID, message)
}

// LoadUsers returns all registered user IDs
func (s *BotServer) LoadUsers() []*User {
	s.usersMutex.RLock()
	defer s.usersMutex.RUnlock()

	users := make([]*User, 0, len(s.users))

	for _, user := range s.users {
		fmt.Println(user.ChatID, user.FirstName, user.Username)
		if user.Username == "Emsusn" || user.Username == "ramisoul" {
			users = append(users, user)
		}
	}

	return users
}

// GetUser returns a specific user by ID
func (s *BotServer) GetUser(userID int64) (*User, bool) {
	s.usersMutex.RLock()
	defer s.usersMutex.RUnlock()

	user, exists := s.users[userID]
	return user, exists
}

// SendMessage sends a message to a specific user
func (s *BotServer) SendMessage(userID int64, message string) error {
	// First, check if we have the user's chat ID stored
	user, exists := s.GetUser(userID)
	if !exists {
		return fmt.Errorf("user %d not found or hasn't registered with the bot", userID)
	}

	// Send the message using the stored Chat ID
	msg := tgbotapi.NewMessage(user.ChatID, message)
	msg.ParseMode = "Markdown"

	_, err := s.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to user %d: %w", userID, err)
	}

	return nil
}

// BroadcastMessage sends a message to all registered users
func (s *BotServer) BroadcastMessage(message string) (int, []error) {
	users := s.LoadUsers()
	var errors []error
	successCount := 0

	log.Printf("📢 Starting broadcast to %d users", len(users))

	for _, user := range users {
		err := s.SendMessage(user.ID, message)
		if err != nil {
			errors = append(errors, fmt.Errorf("user %d: %v", user.ID, err))
			log.Printf("❌ Failed to send to user %d: %v", user.ID, err)
		} else {
			successCount++
			log.Printf("✅ Sent to user %d", user.ID)
		}

		// Rate limiting to avoid hitting Telegram limits
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("📊 Broadcast complete: %d successful, %d failed", successCount, len(errors))
	return successCount, errors
}

// sendMessage is a helper method to send messages
func (s *BotServer) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("Error sending message to chat %d: %v", chatID, err)
	}
}
