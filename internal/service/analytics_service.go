package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type AnalyticsRepository interface {
	SaveVisitor(ctx context.Context, visitor *domain.Visitor) error
	GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error)
}

type analyticsService struct {
	analyticsRepository AnalyticsRepository
	log                 logger.Logger
}

func NewAnalyticsService(analyticsRepository AnalyticsRepository, log logger.Logger) *analyticsService {
	return &analyticsService{analyticsRepository, log}
}

func (s *analyticsService) SaveVisitor(ctx context.Context, visitor *domain.Visitor) error {
	visitor.ID = uuid.New()
	visitor.Time = time.Now()
	return s.analyticsRepository.SaveVisitor(ctx, visitor)
}

func (s *analyticsService) GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error) {
	return s.analyticsRepository.GetVisitors(ctx, limit, offset)
}
