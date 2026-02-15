package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
	)
}

func ObservabilityMiddleware(logger logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		start := time.Now()

		requestId := generateRequestID()

		reqLogger := logger.WithFields(
			map[string]any{
				"request_id": requestId,
				"method":     c.Method(),
				"path":       c.Path(),
				"ip":         c.IP(),
			},
		)

		c.Locals("request_id", requestId)

		reqLogger.Debug().Msg("📥 HTTP Request Started")

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()

		httpRequestsTotal.WithLabelValues(c.Method(), c.Path(), fmt.Sprintf("%d", status)).Inc()
		httpRequestDuration.WithLabelValues(c.Method(), c.Path(), fmt.Sprintf("%d", status)).Observe(duration)

		logEvent := reqLogger.Info().
			Int("status", status).
			Int("bytes", len(c.Response().Body())).
			Float64("duration_seconds", duration)

		if err != nil {
			logEvent = logEvent.Err(err)
		}

		logEvent.Msg("📤 HTTP Request Completed")

		return err

	}
}

func generateRequestID() string {
	return uuid.New().String()
}
