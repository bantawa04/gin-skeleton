package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeMiddleware sanitizes user input to prevent XSS attacks
func SanitizeMiddleware() gin.HandlerFunc {
	// Create a strict HTML sanitizer policy
	policy := bluemonday.StrictPolicy()

	return func(c *gin.Context) {
		// Only sanitize JSON request bodies
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.Next()
			return
		}

		// Read request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body.Close()

		// Skip if body is empty
		if len(bodyBytes) == 0 {
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			c.Next()
			return
		}

		// Parse JSON to sanitize string values
		var data interface{}
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			// If not valid JSON, restore body and continue
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			c.Next()
			return
		}

		// Sanitize the data recursively
		sanitized := sanitizeValue(data, policy)

		// Marshal back to JSON
		sanitizedBytes, err := json.Marshal(sanitized)
		if err != nil {
			// If marshaling fails, use original body
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			c.Next()
			return
		}

		// Set sanitized body
		c.Request.Body = io.NopCloser(bytes.NewReader(sanitizedBytes))
		c.Request.ContentLength = int64(len(sanitizedBytes))

		c.Next()
	}
}

// sanitizeValue recursively sanitizes string values in the data structure
func sanitizeValue(v interface{}, policy *bluemonday.Policy) interface{} {
	switch val := v.(type) {
	case string:
		// Sanitize string values
		return policy.Sanitize(val)
	case map[string]interface{}:
		// Recursively sanitize map values
		result := make(map[string]interface{}, len(val))
		for k, v := range val {
			result[k] = sanitizeValue(v, policy)
		}
		return result
	case []interface{}:
		// Recursively sanitize slice values
		result := make([]interface{}, len(val))
		for i, v := range val {
			result[i] = sanitizeValue(v, policy)
		}
		return result
	default:
		// Return as-is for other types (numbers, booleans, null)
		return val
	}
}
