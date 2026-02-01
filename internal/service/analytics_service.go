package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type AnalyticsRepository interface {
	SaveVisitor(ctx context.Context, visitor *domain.Visitor) error
	GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error)
	GetCountByUserId(ctx context.Context, userId string) (int, error)
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
	visitor.OS = GetOS(visitor.OS)
	return s.analyticsRepository.SaveVisitor(ctx, visitor)
}

func (s *analyticsService) GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error) {
	return s.analyticsRepository.GetVisitors(ctx, limit, offset)
}

func (s *analyticsService) GetCountByUserId(ctx context.Context, userId string) (int, error) {
	fmt.Println("USER ID", userId)
	count, err := s.analyticsRepository.GetCountByUserId(ctx, userId)
	fmt.Println("COUNT,ERR", count, err)
	if err != nil {
		fmt.Println(err)
	}
	return count, nil
}

func GetOS(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macos"):
		return "macOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		return "iOS"
	default:
		return "Unknown"
	}
}
