package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mssola/useragent"
	"github.com/ramisoul84/emil-server/internal/domain"
	"github.com/ramisoul84/emil-server/pkg/location"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type analyticsRepository interface {
	SaveVisit(ctx context.Context, data *domain.Data) error
	ListVisits(ctx context.Context, limit, offset int) ([]*domain.Data, error)
	GetVisitsStats(ctx context.Context) (*domain.Stats, error)
}

type botNotifier interface {
	Notify(ctx context.Context, msg string) error
}

type analyticsService struct {
	repo   analyticsRepository
	bot    botNotifier
	logger logger.Logger
}

func NewAnalyticsService(repo analyticsRepository, bot botNotifier) *analyticsService {
	return &analyticsService{
		repo:   repo,
		bot:    bot,
		logger: logger.Get(),
	}
}

func (s *analyticsService) VisitStart(ctx context.Context, data *domain.VisitStartData) error {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "analytics_service",
			"method":     "visit_start",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling visit")

	ip := ctx.Value("ip").(string)
	country, city := location.GetFullClientInfo(ip)
	os := getOS(data.UserAgent)

	msg := fmt.Sprintf(
		"👁 *New Site Visitor*\n\n"+
			"📍 *IP:* %s\n"+
			"📍 *Country:* %s\n"+
			"📍 *City:* %s\n"+
			"👤 *SESSION:* %s\n"+
			"👤 *User:* %s\n"+
			"📱 *Device OS:* %s\n"+
			"🔗 *Referrer:* %s\n",
		ip,
		country,
		city,
		data.SessionID,
		data.UserID,
		os,
		data.Referrer,
	)

	if err := s.bot.Notify(context.Background(), msg); err != nil {
		logger.Error().Err(err).Msg("Failed to send bot notification")
	}

	return nil
}

func (s *analyticsService) VisitEnd(ctx context.Context, visitData *domain.VisitData) error {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "analytics_service",
			"method":     "visit_end",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling visit end")

	ip := ctx.Value("ip").(string)
	country, city := location.GetFullClientInfo(ip)
	os := getOS(visitData.UserAgent)

	tt, err := time.Parse(time.RFC3339, visitData.StartTime)
	var duration time.Duration
	if err != nil {
		logger.Error().Err(err).Msgf("Failed to parse start time: %s", visitData.StartTime)
		duration = time.Duration(visitData.Duration)

	} else {
		duration = time.Since(tt)
	}

	var data domain.Data
	data.SessionID = visitData.SessionID
	data.UserID = visitData.UserID
	data.IP = ip
	data.Country = country
	data.City = city
	data.OS = os
	data.StartTime = visitData.StartTime
	data.Duration = duration.Seconds()
	data.ActiveDuration = visitData.Duration
	data.ActionsCount = getActionsCount(visitData.Actions)

	msg := fmt.Sprintf(
		"📊 *Session Summary*\n\n"+
			"📍 *IP:* %s\n"+
			"📍 *Country:* %s\n"+
			"📍 *City:* %s\n"+
			"👤 *Seesion:* %s\n"+
			"👤 *User:* %s\n"+
			"📱 *Device OS:* %s\n"+
			"🔗 *Referrer:* %s\n"+
			"🕐 *Session Started:* %s\n"+
			"⌛ *Active Duration:* %f\n"+
			"%s",
		ip,
		country,
		city,
		visitData.SessionID,
		visitData.UserID,
		os,
		visitData.Referrer,
		visitData.StartTime,
		visitData.Duration,
		actionsSummary(visitData.Actions),
	)

	if err := s.bot.Notify(context.Background(), msg); err != nil {
		logger.Error().Err(err).Msg("Failed to send bot notification")
	}

	if err := s.repo.SaveVisit(ctx, &data); err != nil {
		logger.Error().Err(err).Msg("Failed to save visit")
		return domain.ErrInternal
	}

	return nil
}

func (s *analyticsService) ListVisits(ctx context.Context, limit, offset int) ([]*domain.Data, error) {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "analytics_service",
			"method":     "list_visits",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling visits list")

	visits, err := s.repo.ListVisits(ctx, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list visits")
		return nil, domain.ErrInternal
	}

	return visits, nil
}

func (s *analyticsService) VisitStats(ctx context.Context) (*domain.Stats, error) {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":      "analytics_service",
			"method":     "visit_stats",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("➡️  [Service] Handling visit stats")

	stats, err := s.repo.GetVisitsStats(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to  get stats")
		return nil, domain.ErrInternal
	}

	return stats, nil
}

/*
	HELPER FUNCTIONS
*/

func getOS(ua string) string {
	if ua == "" {
		return "Unknown"
	}

	agent := useragent.New(ua)

	os := strings.ToLower(agent.OS())

	switch {
	case strings.Contains(os, "windows"):
		return "Windows"
	case strings.Contains(os, "mac os") || strings.Contains(os, "macos"):
		return "macOS"
	case strings.Contains(os, "linux"):
		return "Linux"
	case strings.Contains(os, "android"):
		return "Android"
	case strings.Contains(os, "iphone") || strings.Contains(ua, "ipad"):
		return "iOS"
	default:
		return "Unknown"
	}
}

func actionsSummary(actions map[string]int) string {
	if len(actions) == 0 {
		return ""
	}

	var message strings.Builder
	message.WriteString("🖱️ *Actions:*\n```\n")

	// Find longest button name for alignment
	maxLen := 0
	for name := range actions {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	// Create formatted table
	for name, count := range actions {
		// Pad button name for alignment
		paddedName := fmt.Sprintf("%-*s", maxLen, name)
		fmt.Fprintf(&message, "%s : %3d\n", paddedName, count)
	}

	message.WriteString("```")

	return message.String()
}

func getActionsCount(actions map[string]int) int {
	var count int

	for _, v := range actions {
		count = count + v
	}

	return count
}
