package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/ramisoul84/emil-server/internal/domain"
)

type analyticsService interface {
	VisitStart(ctx context.Context, data *domain.VisitStartData) error
	VisitEnd(ctx context.Context, data *domain.VisitData) error
	ListVisits(ctx context.Context, limit, offset int) ([]*domain.Data, error)
	VisitStats(ctx context.Context) (*domain.Stats, error)
}

type analyticsHandler struct {
	service analyticsService
}

func NewAnalyticsHandler(service analyticsService) *analyticsHandler {
	return &analyticsHandler{
		service: service,
	}
}

func (h *analyticsHandler) VisitStart(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	var data domain.VisitStartData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.WithValue(c.Context(), "ip", c.IP())
	ctx = context.WithValue(ctx, "request_id", requestId)

	if err := h.service.VisitStart(ctx, &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process visit",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Visit tracked",
	})
}

func (h *analyticsHandler) VisitEnd(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	var data domain.VisitData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.WithValue(c.Context(), "ip", c.IP())
	ctx = context.WithValue(ctx, "request_id", requestId)

	if err := h.service.VisitEnd(ctx, &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to handle visit end",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Visit end tracked",
	})
}

func (h *analyticsHandler) List(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	limit := 20
	offset := 0

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	list, err := h.service.ListVisits(ctx, limit, offset)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get visits list",
		})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

func (h *analyticsHandler) Stats(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	stats, err := h.service.VisitStats(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get stats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}
