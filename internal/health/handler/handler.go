package handler

import (
	"net/http"
	"time"

	response "gin/internal/response"
	"gin/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// Health performs a health check
// @Summary      Health check
// @Description  Check the health status of the API and database connectivity
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=HealthResponse}
// @Failure      503  {object}  response.Response{data=HealthResponse}
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]string),
	}

	// Check database connectivity
	sqlDB, err := h.db.DB()
	if err != nil {
		health.Status = "unhealthy"
		health.Checks["database"] = "error: failed to get database connection"
		response.SendResponse(c, health, "Health check failed", http.StatusServiceUnavailable)
		return
	}

	// Ping database with timeout
	ctx, cancel := utils.GetContextWithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		health.Status = "unhealthy"
		health.Checks["database"] = "error: " + err.Error()
		response.SendResponse(c, health, "Health check failed", http.StatusServiceUnavailable)
		return
	}

	health.Checks["database"] = "ok"

	// Determine HTTP status code
	statusCode := http.StatusOK
	if health.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	response.SendResponse(c, health, "Health check successful", statusCode)
}
