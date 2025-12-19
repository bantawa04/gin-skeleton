package utils

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts the user ID from the Gin context
// This should be called after JWT middleware has validated the token
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	if id, ok := userID.(string); ok {
		return id, true
	}

	return "", false
}

// GetUserClaimsFromContext extracts the user claims from the Gin context
func GetUserClaimsFromContext(c *gin.Context) (*JWTClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	if jwtClaims, ok := claims.(*JWTClaims); ok {
		return jwtClaims, true
	}

	return nil, false
}

// RequireUserID extracts user ID from context and returns error if not found
func RequireUserID(c *gin.Context) (string, error) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		return "", errors.New("user ID not found in context. Ensure JWT middleware is applied")
	}
	return userID, nil
}

// GetContextWithTimeout creates a context with timeout from a parent context
// If parent is nil, creates a new background context
func GetContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		return context.WithTimeout(context.Background(), timeout)
	}
	return context.WithTimeout(parent, timeout)
}
