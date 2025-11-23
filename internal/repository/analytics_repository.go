package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type analyticsRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewAnalyticsRepository(db *sqlx.DB, log logger.Logger) *analyticsRepository {
	return &analyticsRepository{db, log}
}

func (r *analyticsRepository) SaveVisitor(ctx context.Context, visitor *domain.Visitor) error {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "save_visitor",
	})

	query := `
            INSERT INTO visitors (id, ip, user_agent, city, country, time)
            VALUES ($1, $2, $3, $4, $5, $6)
		`

	_, err := r.db.ExecContext(ctx, query,
		visitor.ID,
		visitor.IP,
		visitor.UserAgent,
		visitor.City,
		visitor.Country,
		visitor.Time,
	)

	if err != nil {
		log.WithError(err).Error("failed to save visitor")
		return fmt.Errorf("failed to save visitor: %w", err)
	}

	log.Info("visitor saved successfully")

	return nil
}

func (r *analyticsRepository) GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error) {
	log := r.log.WithFields(map[string]any{
		"layer":     "repository",
		"operation": "get_visitors",
	})

	query := `
		SELECT id, ip, user_agent, city, country, time
		FROM visitors
		ORDER BY time DESC
		LIMIT $1 OFFSET $2
	`

	var visitors []*domain.Visitor
	err := r.db.SelectContext(ctx, &visitors, query, limit, offset)
	if err != nil {
		log.WithError(err).Error("failed to get visitors")
		return nil, 0, 0, fmt.Errorf("failed to get visitors: %w", err)
	}

	type CountResult struct {
		Total  int `db:"total_count"`
		Unique int `db:"unique_count"`
	}

	var counts CountResult
	countQuery := `
		SELECT 
			COUNT(*) as total_count,
			COUNT(DISTINCT ip) as unique_count
		FROM visitors
	`

	err = r.db.GetContext(ctx, &counts, countQuery)
	if err != nil {
		log.WithError(err).Error("failed to get counts")
		return nil, 0, 0, fmt.Errorf("failed to get counts: %w", err)
	}

	return visitors, counts.Total, counts.Unique, nil
}
