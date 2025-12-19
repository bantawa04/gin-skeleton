package middlewares

import (
	"time"

	"gin/internal/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests and automatically logs errors
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Get request details
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response details
		statusCode := c.Writer.Status()
		contentLength := c.Writer.Size()

		// Get request ID from context if available
		requestID := GetRequestID(c)

		// Prepare log fields
		fields := map[string]interface{}{
			"method":         method,
			"path":           path,
			"status_code":    statusCode,
			"duration_ms":    duration.Milliseconds(),
			"client_ip":      clientIP,
			"user_agent":     userAgent,
			"content_length": contentLength,
		}

		// Add request ID if available
		if requestID != "" {
			fields["request_id"] = requestID
		}

		// Log based on status code
		if statusCode >= 400 {
			// Log errors (status >= 400) to error log
			logger.LogHTTPError(statusCode, "HTTP Error", fields)
		}
	}
}
