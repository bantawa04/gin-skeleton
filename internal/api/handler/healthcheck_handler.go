package handler

import (
	response "gin/internal/response"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct{}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// HealthCheck handles GET /health request
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	status := "healthy"

	response.SendResponse(c, gin.H{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
	}, "Health check completed")
}
