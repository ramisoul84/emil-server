package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/location"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *domain.CreateMessageRequest) error
	GetMessageByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	MarkMessageAsRead(ctx context.Context, id uuid.UUID) error
	DeleteMessage(ctx context.Context, id uuid.UUID) error
	GetMessagesList(ctx context.Context, limit, offset int) ([]*domain.Message, int, error)
}

type messageHandler struct {
	messageService MessageService
	log            logger.Logger
}

func NewMessageHandler(messageService MessageService, log logger.Logger) *messageHandler {
	return &messageHandler{messageService, log}
}

func (h *messageHandler) Create(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "create_message",
	})

	var req domain.CreateMessageRequest
	if err := c.Bind(&req); err != nil {
		log.Info("Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if req.Name == "" || req.Email == "" || req.Text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Name and Email and Text are required",
		})
	}

	info, err := location.GetFullClientInfo(c)

	if err != nil {
		req.Country = "Unknown"
		req.City = "Unknown"
	} else {
		req.City = info.City
		req.Country = info.Country
	}

	if err := h.messageService.CreateMessage(c.Request().Context(), &req); err != nil {
		log.WithError(err).Error("Failed to create message")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create message",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Message created successfully",
	})

}

func (h *messageHandler) Get(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "get_message",
	})

	idString := c.Param("id")
	if strings.TrimSpace(idString) == "" {
		log.Info("message ID is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "message ID is required",
		})
	}

	messageID, err := uuid.Parse(idString)
	if err != nil {
		log.Info("Invalid message ID format")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid message ID format",
		})
	}

	message, err := h.messageService.GetMessageByID(c.Request().Context(), messageID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Message not found",
			})
		}
		log.WithError(err).Error("Failed to get message")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get message",
		})
	}

	return c.JSON(http.StatusOK, message)
}

func (h *messageHandler) Update(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "update_message",
	})

	idString := c.Param("id")
	if strings.TrimSpace(idString) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Message ID is required",
		})
	}

	messageID, err := uuid.Parse(idString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid message ID format",
		})
	}

	if err := h.messageService.MarkMessageAsRead(c.Request().Context(), messageID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Message not found",
			})
		}
		log.WithError(err).Error("Failed to mark message as read")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed  to mark message as read",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Message marked as read successfully",
	})
}

func (h *messageHandler) Delete(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "delete_message",
	})

	idString := c.Param("id")
	if strings.TrimSpace(idString) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Message ID is required",
		})
	}

	messageID, err := uuid.Parse(idString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid message ID format",
		})
	}

	if err := h.messageService.DeleteMessage(c.Request().Context(), messageID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Message not found",
			})
		}
		log.WithError(err).Error("Failed to delete message")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete message",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Message deleted successfully",
	})
}

func (h *messageHandler) List(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "list_messages",
	})

	limitString := c.QueryParam("limit")
	offsetString := c.QueryParam("offset")

	limit := 10
	if limitString != "" {
		parsedLimit, err := strconv.Atoi(limitString)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": "limit must be a positive integer between 1 and 100",
			})
		}
		limit = parsedLimit
	}

	offset := 0
	if offsetString != "" {
		parsedOffset, err := strconv.Atoi(offsetString)
		if err != nil || parsedOffset < 0 {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": "offset must be a non-negative integer",
			})
		}
		offset = parsedOffset
	}

	messages, total, err := h.messageService.GetMessagesList(c.Request().Context(), limit, offset)
	if err != nil {
		log.WithError(err).Error("Failed to list messages")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list messages",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"messages": messages,
		"total":    total,
	})
}
