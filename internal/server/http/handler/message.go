package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ramisoul84/emil-server/internal/domain"
)

type messageService interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetMessage(ctx context.Context, id int) (*domain.Message, error)
	UpdateMessage(ctx context.Context, id int) error
	DeleteMessage(ctx context.Context, id int) error
	ListMessages(ctx context.Context, limit, offset int) ([]*domain.Message, int, error)
}

type messageHandler struct {
	service messageService
}

func NewMessageHandler(service messageService) *messageHandler {
	return &messageHandler{
		service: service,
	}
}

func (h *messageHandler) Create(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)
	ip := c.Locals("ip").(string)

	var data domain.Message
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.WithValue(c.Context(), "ip", ip)
	ctx = context.WithValue(ctx, "request_id", requestId)

	if err := h.service.CreateMessage(ctx, &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message Saved",
	})
}

func (h *messageHandler) Get(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	idString := c.Params("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "limit must be a positive integer between 1 and 100",
		})
	}

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	message, err := h.service.GetMessage(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(message)
}

func (h *messageHandler) Update(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	idString := c.Params("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "limit must be a positive integer between 1 and 100",
		})
	}

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	err = h.service.UpdateMessage(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message Updated",
	})
}

func (h *messageHandler) Delete(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	idString := c.Params("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "limit must be a positive integer between 1 and 100",
		})
	}

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	err = h.service.DeleteMessage(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message Deleted",
	})
}

func (h *messageHandler) List(c *fiber.Ctx) error {
	requestId := c.Locals("request_id").(string)

	limitString := c.Params("limit")
	offsetString := c.Params("offset")

	limit := 10
	if limitString != "" {
		parsedLimit, err := strconv.Atoi(limitString)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "limit must be a positive integer between 1 and 100",
			})
		}
		limit = parsedLimit
	}

	offset := 0
	if offsetString != "" {
		parsedOffset, err := strconv.Atoi(offsetString)
		if err != nil || parsedOffset < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "offset must be a non-negative integer",
			})
		}
		offset = parsedOffset
	}

	ctx := context.WithValue(c.Context(), "request_id", requestId)
	messages, total, err := h.service.ListMessages(ctx, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to list messages",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
		"total":    total,
	})
}
