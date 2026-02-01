package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/location"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type trackReq struct {
	UserId string `json:"user_id"`
}

type AnalyticsService interface {
	SaveVisitor(ctx context.Context, visitor *domain.Visitor) error
	GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error)
}

type analyticsHandler struct {
	service    AnalyticsService
	botService BotService
	log        logger.Logger
}

func NewAnalyticsHandler(service AnalyticsService, botService BotService, log logger.Logger) *analyticsHandler {
	return &analyticsHandler{service, botService, log}
}

func (h *analyticsHandler) TrackVisitor(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "track_visitor",
	})

	var req trackReq
	if err := c.Bind(&req); err != nil {
		log.Info("Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	ip := c.RealIP()
	userAgent := c.Request().UserAgent()

	var visitor domain.Visitor

	visitor.UserID = req.UserId
	visitor.IP = ip
	visitor.OS = userAgent

	info, err := location.GetFullClientInfo(c)

	if err != nil {
		visitor.Country = "Unknown"
		visitor.City = "Unknown"
	} else {
		visitor.City = info.City
		visitor.Country = info.Country

	}

	if err := h.service.SaveVisitor(c.Request().Context(), &visitor); err != nil {
		log.WithError(err).Error("Failed to create message")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create message",
		})
	}

	message := fmt.Sprintf(
		"👋 a new visitor!\n"+
			"Time:%s.\n"+
			"User_ID:%s.\n"+
			"IP:%s.\n"+
			"OS:%s.\n"+
			"Country:%s.\n"+
			"City:%s",
		visitor.Time,
		visitor.UserID,
		visitor.IP,
		visitor.OS,
		visitor.Country,
		visitor.City,
	)

	h.botService.BroadcastMessage(message)

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Message created successfully",
	})
}

func (h *analyticsHandler) GetVisitors(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "get_visitors",
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

	visitors, total, unique, err := h.service.GetVisitors(c.Request().Context(), limit, offset)
	if err != nil {
		log.WithError(err).Error("Failed to get visitors")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visitors",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"visitors": visitors,
		"total":    total,
		"unique":   unique,
	})
}
