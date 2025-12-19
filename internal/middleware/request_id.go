package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDKey = "request_id"

// RequestIDMiddleware generates a unique request ID for each request
// and adds it to the response headers and context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already present in header
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			// Generate a new ULID for request ID
			requestID = ulid.Make().String()
		}

		// Set request ID in context for use in handlers and logging
		c.Set(RequestIDKey, requestID)

		// Add request ID to response header
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from the Gin context
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}
