package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"realtime-events/internal/models"
	"realtime-events/internal/services"
)

type EventHandler struct {
	service           *services.EventService
	validationService *services.ValidationService
	logger            *zap.SugaredLogger
}

func NewEventHandler(service *services.EventService, logger *zap.SugaredLogger) *EventHandler {
	return &EventHandler{
		service:           service,
		validationService: services.NewValidationService(),
		logger:            logger,
	}
}

func (h *EventHandler) IngestEvent(c *gin.Context) {
	var req models.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}

	// Additional validation
	if err := h.validationService.ValidateEventRequest(&req); err != nil {
		h.logger.Errorw("Validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_failed", "message": err.Error()})
		return
	}

	projectID, exists := c.Get("project_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ip := getClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	event, err := h.service.ProcessEvent(c.Request.Context(), &req, projectID.(string), ip, userAgent)
	if err != nil {
		h.logger.Errorw("Failed to process event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":   "accepted",
		"event_id": event.ID,
	})
}

func (h *EventHandler) IngestBatchEvents(c *gin.Context) {
	var req models.BatchEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}

	// Validate each event in the batch
	for i, eventReq := range req.Events {
		if err := h.validationService.ValidateEventRequest(&eventReq); err != nil {
			h.logger.Errorw("Batch validation failed", "error", err, "index", i)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "validation_failed",
				"message": fmt.Sprintf("Event at index %d: %s", i, err.Error()),
			})
			return
		}
	}

	projectID, exists := c.Get("project_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ip := getClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	events := make([]string, 0, len(req.Events))
	for _, eventReq := range req.Events {
		event, err := h.service.ProcessEvent(c.Request.Context(), &eventReq, projectID.(string), ip, userAgent)
		if err != nil {
			h.logger.Errorw("Failed to process batch event", "error", err)
			continue // Continue processing other events
		}
		events = append(events, event.ID)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":    "accepted",
		"event_ids": events,
	})
}

func getClientIP(c *gin.Context) net.IP {
	// Check X-Forwarded-For header
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			if ip := net.ParseIP(strings.TrimSpace(ips[0])); ip != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		if ip := net.ParseIP(xri); ip != nil {
			return ip
		}
	}

	// Fallback to remote address
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return net.ParseIP("127.0.0.1")
	}
	return net.ParseIP(ip)
}