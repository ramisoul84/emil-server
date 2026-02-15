package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ramisoul84/emil-server/internal/domain"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type analyticsRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewAnalyticsRepository(db *sqlx.DB) *analyticsRepository {
	return &analyticsRepository{db, logger.Get()}
}

func (r *analyticsRepository) SaveVisit(ctx context.Context, data *domain.Data) error {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "analytics_repository",
			"method":     "save_visit",
			"request_id": ctx.Value("request_id").(string),
		},
	)
	logger.Info().Msg("Store visit in DB")

	query := `
            INSERT INTO visits (
				session_id, user_id, ip, country, city, os,
				start_time, duration, active_duration, actions_count
				)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

	_, err := r.db.ExecContext(ctx, query,
		data.SessionID,
		data.UserID,
		data.IP,
		data.Country,
		data.City,
		data.OS,
		data.StartTime,
		data.Duration,
		data.ActiveDuration,
		data.ActionsCount,
	)

	if err != nil {
		logger.Error().Err(err).Msg("failed to save visit")
		return fmt.Errorf("failed to save visit: %w", err)
	}

	logger.Info().Msg("visit saved successfully")

	return nil
}

func (r *analyticsRepository) ListVisits(ctx context.Context, limit, offset int) ([]*domain.Data, error) {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "analytics_repository",
			"method":     "list_visits",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("Get list of visit")

	if limit <= 0 {
		limit = 20 // default limit
	}
	if limit > 100 {
		limit = 100 // max limit
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT *
		FROM visits 
		ORDER BY start_time DESC 
		LIMIT $1 OFFSET $2
	`

	var data []*domain.Data
	err := r.db.SelectContext(ctx, &data, query, limit, offset)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to list visits")
		return nil, fmt.Errorf("failed to list visits: %w", err)
	}

	logger.Info().Msg("got list successfully")

	return data, nil
}

func (r *analyticsRepository) GetVisitsStats(ctx context.Context) (*domain.Stats, error) {
	logger := r.logger.WithFields(
		map[string]any{
			"layer":      "analytics_repository",
			"method":     "visits_stats",
			"request_id": ctx.Value("request_id").(string),
		},
	)

	logger.Info().Msg("Get visits stats")

	query := `
		SELECT 
			COUNT(*) as total_visits,
			COUNT(DISTINCT user_id) as unique_users,
			COALESCE(AVG(duration), 0) as avg_duration,
			COALESCE(AVG(active_duration), 0) as avg_active_duration,
			COALESCE(AVG(actions_count), 0) as avg_actions
		FROM visits;
	`

	var stats domain.Stats
	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to get visit stats")
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil
}
